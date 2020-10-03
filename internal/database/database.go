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
func DetermineEngine(file *os.File) (string, error) {
	engine := ""
	line := 1

	s := bufio.NewScanner(file)

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
func HasCreateStatement(file *os.File) (bool, error) {
	s := bufio.NewScanner(file)
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
