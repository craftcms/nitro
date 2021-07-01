package setup

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/portavail"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	mysqlDefaultPort    = 3306
	postgresDefaultPort = 5432
)

// FirstTime is used when there is no configuration file found in a users
// home/.nitro directory. We do not prompt for input such as memory, cpu,
// disk space in version 2 as that is defined and managed at the docker
// level. If anything fails, we return an error.
func FirstTime(home string, reader io.Reader, output terminal.Outputer) error {
	c := config.Config{File: filepath.Join(home, config.DirectoryName, config.FileName)}

	output.Info("Setting up Nitro…")

	// if this is running on Apple Silicon, we need to prompt for mariadb instead until this issue is resolved: https://docs.docker.com/docker-for-mac/apple-m1/
	switch runtime.GOARCH == "arm64" || runtime.GOARCH == "arm" {
	case true:
		if runtime.GOOS == "darwin" {
			output.Info("Apple computers with new silicon do not work with MySQL images.")
		} else {
			output.Info("ARM computers do not work with MySQL images.")
		}

		mariadb, err := output.Confirm("Would you like to use MariaDB?", true, "")
		if err != nil {
			return err
		}

		if mariadb {
			// prompt for the version
			opts := []string{"10.5", "10.4", "10.3", "10.2", "10.1", "10"}
			selected, err := output.Select(os.Stdin, "Select MariaDB version: ", opts)
			if err != nil {
				return err
			}

			version := opts[selected]

			// check if the port is available
			var port string
			for {
				if err := portavail.Check("", strconv.Itoa(mysqlDefaultPort)); err != nil {
					mysqlDefaultPort = mysqlDefaultPort + 1
					continue
				}

				port = strconv.Itoa(mysqlDefaultPort)

				break
			}

			// add a default mariadb database
			c.Databases = append(c.Databases, config.Database{
				Engine:  "mariadb",
				Version: version,
				Port:    port,
			})
		}
	default:
		mysql, err := output.Confirm("Would you like to use MySQL?", true, "")
		if err != nil {
			return err
		}

		if mysql {
			// prompt for the version
			opts := []string{"8.0", "5.7", "5.6"}
			selected, err := output.Select(os.Stdin, "Select MySQL version: ", opts)
			if err != nil {
				return err
			}

			version := opts[selected]

			// check if the port is available
			var port string
			for {
				if err := portavail.Check("", strconv.Itoa(mysqlDefaultPort)); err != nil {
					mysqlDefaultPort = mysqlDefaultPort + 1
					continue
				}

				port = strconv.Itoa(mysqlDefaultPort)

				break
			}

			// add a default mysql database
			c.Databases = append(c.Databases, config.Database{
				Engine:  "mysql",
				Version: version,
				Port:    port,
			})
		}
	}

	postgres, err := output.Confirm("Would you like to use PostgreSQL?", true, "")
	if err != nil {
		return err
	}

	if postgres {
		// prompt for the version
		opts := []string{"13", "12", "11", "10", "9"}
		selected, err := output.Select(os.Stdin, "Select PostgreSQL version: ", opts)
		if err != nil {
			return err
		}

		version := opts[selected]

		// check if the port is available
		var port string
		for {
			if err := portavail.Check("", strconv.Itoa(postgresDefaultPort)); err != nil {
				postgresDefaultPort = postgresDefaultPort + 1
				continue
			}

			port = strconv.Itoa(postgresDefaultPort)

			break
		}

		// add a default postgres database
		c.Databases = append(c.Databases, config.Database{
			Engine:  "postgres",
			Version: version,
			Port:    port,
		})
	}

	redis, err := output.Confirm("Would you like to use Redis?", true, "")
	if err != nil {
		return err
	}

	if redis {
		output.Pending("adding redis service")

		c.Services.Redis = true

		output.Done()
	}

	// save the file
	if err := c.Save(); err != nil {
		return err
	}

	return nil
}
