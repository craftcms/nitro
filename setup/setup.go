package setup

import (
	"path/filepath"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/terminal"
)

// FirstTime is used when there is no configuration file found in a users
// home/.nitro directory. We do not prompt for input such as memory, cpu,
// disk space in version 2 as that is defined and managed at the docker
// level. If anything fails, we return an error.
func FirstTime(home, env string, output terminal.Outputer) error {
	// TODO(jasonmccallister) consider prompts for which type(s) of database?
	c := config.Config{
		File: filepath.Join(home, ".nitro", env+".yaml"),
	}

	// add a default mysql database
	c.Databases = append(c.Databases, config.Database{
		Engine:  "mysql",
		Version: "8.0",
		Port:    "3306",
	})

	// add a default postgres database
	c.Databases = append(c.Databases, config.Database{
		Engine:  "postgres",
		Version: "12",
		Port:    "5432",
	})

	// save the file
	if err := c.Save(); err != nil {
		return err
	}

	return nil
}
