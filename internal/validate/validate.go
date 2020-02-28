package validate

import (
	"errors"
	"fmt"
	"os"
)

// Database takes a string that represents a type of database to install and returns an error if it is a database that
// is not supported.
func Database(v string) error {
	switch v {
	case "mariadb":
		return nil
	case "postgres":
		return nil
	}

	return errors.New(fmt.Sprintf("the database %q is not supported", v))
}

// Path will check is a fali
func Path(p string) error {
	f, err := os.Stat(p)
	if err != nil {

		return err
	}

	if f.IsDir() {
		return nil
	}

	return errors.New("you must provide a path to a valid directory")
}

// PHPVersion takes a string that represents a PHP version to install and returns and error if that PHP version
// is not yet supported.
func PHPVersion(v string) error {
	switch v {
	case "7.4":
		return nil
	case "7.3":
		return nil
	case "7.2":
		return nil
	}

	return errors.New(fmt.Sprintf("the PHP version %q is not valid", v))
}
