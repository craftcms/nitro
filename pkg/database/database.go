package database

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/pkg/filetype"
	"github.com/docker/docker/pkg/archive"
)

// ErrUnknownDatabaseEngine is returned when we are unable to determine the engine type from a database backup file.
var ErrUnknownDatabaseEngine = fmt.Errorf("unknown database engine detected from file")

// DetermineEngine takes a file and will check if the
// content of the file is for mysql or postgres db
// imports. It will return the engine "mysql" or
// "postgres" if it can determine the engine.
// If it cannot, it will return an error.
func DetermineEngine(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	engine := ""
	line := 1

	s := bufio.NewScanner(f)
	for s.Scan() {
		txt := s.Text()

		// check if its postgres
		if strings.Contains(txt, "PostgreSQL") || strings.Contains(txt, "pg_dump") {
			engine = "postgres"
			break
		}

		// check if its mysql
		if strings.Contains(txt, "MySQL") || strings.Contains(txt, "mysqldump") || strings.Contains(txt, "mariadb") || strings.Contains(txt, "MariaDB") || strings.Contains(txt, "ENGINE=InnoDB") {
			engine = "mysql"
			break
		}

		if line >= 50 || engine != "" {
			break
		}

		line++
	}

	// final check for empty engine
	if engine == "" {
		return "", ErrUnknownDatabaseEngine
	}

	return engine, nil
}

// HasCreateStatement takes a file and will determine
// if the file will create a database during import.
// If it creates a database, it will return true
// otherwise it will return false.
func HasCreateStatement(file string) (bool, error) {
	f, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	line := 0
	for s.Scan() {
		// check if the line has a create database statement
		if strings.Contains(s.Text(), "CREATE DATABASE") {
			return true, nil
		}

		if line >= 100 {
			break
		}

		line++
	}

	return false, nil
}

// PrepareArchiveFromPath takes a path to a file, which is presumed to
// be a database backup and will determine if the file is already a format
// the Docker API can use. If the file is a zip or tar file, it will open
// the sql file in the archive and write the contents to a temporary file
// and then use the docker archive.Generate functionality to prepare for
// copying the file to the API.
func PrepareArchiveFromPath(path string) (io.Reader, string, error) {
	// get the filename from the path directory
	_, name := filepath.Split(path)

	// if this is a support type docker can already use
	if archive.IsArchivePath(path) {
		// read the file
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, "", err
		}

		r, err := archive.Generate(name, string(b))
		return r, name, err
	}

	// determine the kind of file
	kind, err := filetype.Determine(path)
	if err != nil {
		return nil, "", err
	}

	switch kind {
	case "zip":
		// create a new zip reader
		r, err := zip.OpenReader(path)
		if err != nil {
			return nil, "", err
		}
		defer r.Close()

		// read each of the files
		for _, file := range r.File {
			if strings.HasSuffix(file.Name, ".sql") {
				// create the temp file
				temp, err := ioutil.TempFile(os.TempDir(), "nitro-import-zip-")
				if err != nil {
					return nil, "", err
				}
				defer temp.Close()

				// read of the zip file
				rc, err := file.Open()
				if err != nil {
					return nil, "", err
				}
				defer rc.Close()

				buf := new(bytes.Buffer)
				if _, err := buf.ReadFrom(rc); err != nil && !errors.Is(err, io.EOF) {
					return nil, "", err
				}

				// write to the temp file
				if _, err := temp.Write(buf.Bytes()); err != nil {
					return nil, "", err
				}

				// read from the temp file
				b, err := ioutil.ReadFile(temp.Name())
				if err != nil {
					return nil, "", err
				}

				reader, err := archive.Generate(file.Name, string(b))
				if err != nil {
					return nil, "", err
				}

				return reader, file.Name, err
			}
		}

		// we did not find a sql file, so we need to return an error
		return nil, "", fmt.Errorf("unable to find a .sql file in the zip")
	case "tar":
		f, err := os.Open(path)
		if err != nil {
			return nil, "", err
		}
		defer f.Close()

		r, err := gzip.NewReader(f)
		if err != nil {
			return nil, "", err
		}

		// create the temp file
		temp, err := ioutil.TempFile(os.TempDir(), "nitro-import-gzip-")
		if err != nil {
			return nil, "", err
		}
		defer temp.Close()

		if _, err := io.Copy(temp, r); err != nil {
			return nil, "", err
		}

		// read from the temp file
		b, err := ioutil.ReadFile(temp.Name())
		if err != nil {
			return nil, "", err
		}

		reader, err := archive.Generate(name, string(b))
		if err != nil {
			return nil, "", err
		}

		return reader, name, err
	}

	// if we are here, its a plain file so just read the file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	// generate the reader
	reader, err := archive.Generate(name, string(b))
	if err != nil {
		return nil, "", err
	}

	return reader, name, err
}
