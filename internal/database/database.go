package database

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

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
	s := bufio.NewScanner(f)
	line := 1
	for s.Scan() {
		// check if its postgres
		if strings.Contains(s.Text(), "PostgreSQL") || strings.Contains(s.Text(), "pg_dump") {
			engine = "postgres"
			break
		}

		// check if its mysql
		if strings.Contains(s.Text(), "MySQL") || strings.Contains(s.Text(), "mysqldump") || strings.Contains(s.Text(), "mariadb") {
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
		return "", errors.New("unknown database engine detected from file")
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

	creates := false

	s := bufio.NewScanner(f)
	line := 1
	for s.Scan() {

		// check if its mysql
		if strings.Contains(s.Text(), "CREATE DATABASE") {
			creates = true
			break
		}

		if line >= 50 || creates != false {
			break
		}

		line++
	}

	return creates, nil
}
