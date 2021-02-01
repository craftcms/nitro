package database

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/craftcms/nitro/pkg/filetype"
	"github.com/craftcms/nitro/pkg/pathexists"
)

var (
	MySQLImportCommand    = "mysql"
	PostgresImportCommand = "psql"
)

// Importer is an interface that is designed to import a database backup
// into a database server.
type Importer interface {
	Import(opts *ImportOptions, finder func(engine, version string) (string, error)) error
}

// ImportOptions are used to create a new importer.
// It contains all of the information needed to run an import.
type ImportOptions struct {
	Engine       string
	Version      string
	Hostname     string
	Port         string
	DatabaseName string
	File         string
	Compressed   bool
}

type importer struct{}

// NewImporter takes options and returns a new
// database importer.
func NewImporter() *importer {
	return &importer{}
}

// Import performs the import operation for a database.
func (importer *importer) Import(opts *ImportOptions, find func(engine, version string) (string, error)) error {
	// ensure there are options
	if opts == nil {
		return fmt.Errorf("no options were provider")
	}

	// validate all of the options
	if err := Validate(opts); err != nil {
		return err
	}

	// check to verify the path exists and is a file
	if !pathexists.IsFile(opts.File) {
		return fmt.Errorf("unable to file the file %s", opts.File)
	}

	// find the import tool
	tool, err := find(opts.Engine, opts.Version)
	if err != nil {
		return err
	}

	// check if the file is compressed
	if opts.Compressed {
		// extract the contents into a new file
		extracted, err := importer.extract(opts.File)
		if err != nil {
			return err
		}

		opts.File = extracted
	}

	// generate the commands to execute
	var createCommand, importCommand []string
	switch opts.Engine {
	case "postgres":
		createCommand = []string{fmt.Sprintf("--host=%s", opts.Hostname), "--port=" + opts.Port, "--username=nitro", fmt.Sprintf(`-c CREATE DATABASE %s;`, opts.DatabaseName)}
		importCommand = []string{fmt.Sprintf("--host=%s", opts.Hostname), "--port=" + opts.Port, "--username=nitro", opts.DatabaseName, "--file=" + opts.File}
	default:
		createCommand = []string{"--user=nitro", fmt.Sprintf("--host=%s", opts.Hostname), "-pnitro", fmt.Sprintf(`-e CREATE DATABASE IF NOT EXISTS %s;`, opts.DatabaseName)}
		// https: //dev.mysql.com/doc/refman/8.0/en/mysql-command-options.html
		importCommand = []string{"--user=nitro", fmt.Sprintf("--host=%s", opts.Hostname), "-pnitro", opts.DatabaseName, fmt.Sprintf(`-e source %s`, opts.File)}
	}

	// if there is a create command, lets create the database
	if createCommand != nil {
		if err := importer.exec(tool, createCommand); err != nil {
			// do not exit on error with the crate command - the error could be "Database already exists"
			fmt.Println(err)
		}
	}

	// import the database
	return importer.exec(tool, importCommand)
}

func (importer *importer) extract(path string) (string, error) {
	// create the temp file to store the data
	temp, err := ioutil.TempFile(os.TempDir(), "nitro-db-compressed")
	if err != nil {
		return "", err
	}
	defer temp.Close()

	// open the file
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// detect the kind
	kind, err := filetype.Determine(path)
	if err != nil {
		return "", err
	}
	switch kind {
	case "zip":
		// create a new gzip reader for the uploading path/file
		r, err := zip.OpenReader(path)
		if err != nil {
			return "", err
		}
		defer r.Close()

		// look at all the files
		for _, f := range r.File {
			if strings.HasSuffix(f.Name, ".sql") && !strings.Contains(f.Name, "MACOSX") {
				// open the file
				rc, err := f.Open()
				if err != nil {
					return "", err
				}
				defer rc.Close()

				buf := new(bytes.Buffer)
				if _, err := buf.ReadFrom(rc); err != nil && !errors.Is(err, io.EOF) {
					return "", err
				}

				// write to the temp file
				if _, err := temp.Write(buf.Bytes()); err != nil {
					return "", err
				}

				return temp.Name(), nil
			}
		}
	case "tar":
		// open the compressed file
		f, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer f.Close()

		// read the file
		reader, err := gzip.NewReader(f)
		if err != nil {
			return "", err
		}

		// copy the content into the new temp file
		if _, err := io.Copy(temp, reader); err != nil {
			return "", err
		}

		return temp.Name(), nil
	}

	return "", fmt.Errorf("unsupported compressed file type %q provided", kind)
}

func (importer *importer) exec(tool string, commands []string) error {
	c := exec.Command(tool, commands...)

	c.Stderr = os.Stderr
	c.Stdout = os.Stdout

	if err := c.Start(); err != nil {
		return fmt.Errorf("unable to start the command: %w", err)
	}

	if err := c.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return fmt.Errorf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			return err
		}
	}

	return nil
}

// Validate takes import options and returns an
// error if the options are missing details
// we need to run the import.
func Validate(opts *ImportOptions) error {
	if opts.Engine == "" {
		return fmt.Errorf("import options is missing the engine")
	}

	// we are not using the version yet
	// if opts.Version == "" {
	// 	return fmt.Errorf("import options is missing the version")
	// }

	if opts.Port == "" {
		return fmt.Errorf("import options is missing the port")
	}

	if opts.Hostname == "" {
		return fmt.Errorf("import options is missing the hostname")
	}

	return nil
}

// DefaultImportToolFinder is a tool that is used to find the executable path
// to the import tool such as mysql or psql. It is a func that is provided
// to the Importer.Import func. It will return the path to the executable
// or an error if the command is not found
func DefaultImportToolFinder(engine, version string) (string, error) {
	switch engine {
	case "postgres":
		t, err := exec.LookPath(PostgresImportCommand)
		if err != nil {
			return "", fmt.Errorf("unable to find the `%q` import tool", PostgresImportCommand)
		}

		return t, nil
	case "mysql":
		t, err := exec.LookPath(MySQLImportCommand)
		if err != nil {
			return "", fmt.Errorf("unable to find the `%q` import tool", MySQLImportCommand)
		}

		return t, nil
	}

	return "", fmt.Errorf("unknown engine %q and version %q options provided", engine, version)
}
