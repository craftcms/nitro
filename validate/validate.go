package validate

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/craftcms/nitro/config"
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

func Domain(v string) error {
	msg := "you must provide a valid domain, without a TLD and only lowercase"

	if strings.Contains(v, " ") {
		return errors.New(msg)
	}

	if !strings.Contains(v, ".") {
		return errors.New(msg)
	}

	lower := strings.ToLower(v)
	if lower != v {
		return errors.New(msg)
	}

	return nil
}

// path will check is a fali
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

func Memory(v string) error {
	if !strings.HasSuffix(v, "G") {
		return errors.New("memory must end with a G")
	}

	return nil
}

func DiskSize(v string) error {
	if !strings.HasSuffix(v, "G") {
		return errors.New("disk must end with a G")
	}

	return nil
}

func DatabaseEngine(v string) error {
	switch v {
	case "mysql":
		return nil
	case "postgres":
		return nil
	}
	return errors.New("unsupported database engine: " + v)
}

func DatabaseEngineAndVersion(e, v string) error {
	if err := DatabaseEngine(e); err != nil {
		return err
	}

	if e == "mysql" {
		switch v {
		case "8.0":
			return nil
		case "8":
			return nil
		case "5.8":
			return nil
		case "5.7":
			return nil
		case "5.6":
			return nil
		case "5":
			return nil
		}
	}

	if e == "postgres" {
		switch v {
		case "12.2":
			return nil
		case "12":
			return nil
		case "11.7":
			return nil
		case "11":
			return nil
		case "10.12":
			return nil
		case "10":
			return nil
		case "9.6":
			return nil
		case "9.5":
			return nil
		case "9":
			return nil
		}
	}

	return errors.New("unsupported version of " + e + ": " + v)
}

func MachineName(v string) error {
	if v == "" {
		return errors.New("machine name cannot be empty")
	}
	if strings.Contains(v, " ") {
		return errors.New("machine name cannot contain spaces")
	}

	return nil
}

func DatabaseConfig(databases []config.Database) error {
	ports := map[string]string{}
	versions := map[string]string{}

	for _, database := range databases {
		if ports[database.Port] != "" {
			return errors.New(fmt.Sprintf("duplicate port %s assigned to %s version %s", database.Port, database.Engine, database.Version))
		}

		if versions[database.Engine] == database.Version {
			return errors.New(fmt.Sprintf("duplicate engine %s and version %s", database.Engine, database.Version))
		}

		ports[database.Port] = database.Port
		versions[database.Engine] = database.Version
	}

	return nil
}
