package validate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/craftcms/nitro/internal/config"
)

func DatabaseEngine(v string) error {
	switch v {
	case "mysql":
		return nil
	case "postgres":
		return nil
	}
	return errors.New("Unsupported database engine: " + v)
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

func DatabaseConfig(databases []config.Database) error {
	ports := map[string]string{}
	versions := map[string]string{}

	for _, database := range databases {
		if ports[database.Port] != "" {
			return errors.New(fmt.Sprintf("Duplicate port %s assigned to %s version %s", database.Port, database.Engine, database.Version))
		}

		if versions[database.Engine] == database.Version {
			return errors.New(fmt.Sprintf("Duplicate engine %s and version %s", database.Engine, database.Version))
		}

		ports[database.Port] = database.Port
		versions[database.Engine] = database.Version
	}

	return nil
}

// DatabaseName is used to validate a give database name to ensure its valid
func DatabaseName(s string) error {
	// if the string is empty
	if s == "" {
		return errors.New("no name was provided")
	}

	// cant be longer than 65
	if len(s) > 64 {
		return errors.New("length of the name must be less than 64 chars")
	}

	// check if the first character is an int
	if f, err := strconv.Atoi(string(s[0])); err == nil && f != 0 {
		return errors.New("name cannot start with a number")
	}

	// check the string for any special chars
	if strings.ContainsAny(s, " $-") {
		return errors.New("invalid name, can only contain letters, numbers, and underscores")
	}

	// check for pg_
	if strings.HasPrefix(s, "pg_") {
		return errors.New("name cannot contain pg_")
	}

	return nil
}
