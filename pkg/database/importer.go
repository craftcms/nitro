package database

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"syscall"

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
	Compressed      bool
	CompressionType string
	Engine          string
	Version         string
	Hostname        string
	Port            string
	DatabaseName    string
	File            string
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

	// generate the commands to execute
	var createCommand, importCommand []string
	switch opts.Engine {
	case "postgres":
		createCommand = []string{fmt.Sprintf("--host=%s", opts.Hostname), "--port=" + opts.Port, "--username=nitro", fmt.Sprintf(`-c CREATE DATABASE %s;`, opts.DatabaseName)}
		importCommand = []string{fmt.Sprintf("--host=%s", opts.Hostname), "--port=" + opts.Port, "--username=nitro", opts.DatabaseName, "--file=" + opts.File}
	default:
		createCommand = []string{"--user=nitro", fmt.Sprintf("--host=%s", opts.Hostname), "-pnitro", fmt.Sprintf(`-e CREATE DATABASE IF NOT EXISTS %s;`, opts.DatabaseName)}
		// https://dev.mysql.com/doc/refman/8.0/en/mysql-command-options.html
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
	if err := importer.exec(tool, importCommand); err != nil {
		return err
	}

	return nil
}

func (importer *importer) exec(tool string, commands []string) error {
	c := exec.Command(tool, commands...)

	c.Stderr = ioutil.Discard
	c.Stdout = ioutil.Discard

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
