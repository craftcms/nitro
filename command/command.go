package command

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/phpdev/action"
	"github.com/pixelandtonic/phpdev/validate"
)

func Initialize() *cli.Command {
	return &cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "Initialize a new machine",
		Action: func(c *cli.Context) error {
			if err := action.Initialize(c); err != nil {
				return err
			}

			return nil
		},
	}
}

func Update() *cli.Command {
	return &cli.Command{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "Update a machine with the latest software",
		Action: func(c *cli.Context) error {
			if err := action.Update(c); err != nil {
				return err
			}
			return nil
		},
	}
}

func SSH() *cli.Command {
	return &cli.Command{
		Name:    "ssh",
		Aliases: []string{"s"},
		Usage:   "SSH into a machine as administrator",
		Action: func(c *cli.Context) error {
			if err := action.SSH(c); err != nil {
				return err
			}
			return nil
		},
	}
}

func Install() *cli.Command {
	return &cli.Command{
		Name:        "install",
		Usage:       "Install common tools such as PHP, web servers, and databases",
		Description: "Install offers common options for installing packages on a machine",
		Action: func(c *cli.Context) error {
			return cli.ShowSubcommandHelp(c)
		},
		Subcommands: []*cli.Command{
			installPHP(),
			installNginx(),
			installMariaDB(),
			installRedis(),
			installPostgres(),
		},
	}
}

func Stop() *cli.Command {
	return &cli.Command{
		Name:        "stop",
		Usage:       "Stop a machine",
		Description: "Stop a machine when not in use",
		Action: func(c *cli.Context) error {
			if err := action.Stop(c); err != nil {
				return err
			}
			return nil
		},
	}
}

func Delete() *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Usage:       "Delete a machine",
		Description: "Delete a machine when no longer needed",
		Action: func(c *cli.Context) error {
			rdr := bufio.NewReader(os.Stdin)
			fmt.Print("WARNING: Are you sure you wish to perform this task (y/N)? ")

			answer, err := rdr.ReadString(' ')
			if err != nil {
				return err
			}

			if strings.Contains(answer, "n") {
				fmt.Println("Skipping the deletion of", c.String("machine"))
				return nil
			}

			if err := action.Delete(c); err != nil {
				return err
			}
			return nil
		},
	}
}

func Create() *cli.Command {
	return &cli.Command{
		Name:    "create",
		Aliases: []string{"c", "up"},
		Usage:   "Create a new machine for development",
		Action: func(context *cli.Context) error {
			// get the PHP version
			if err := validate.PHPVersionFlag(context.String("php")); err != nil {
				return err
			}

			// get the database to install
			if err := validate.DatabaseFlag(context.String("database")); err != nil {
				return err
			}

			// move a phpinfo file into the root site
			// get the machine ip
			// print the help to set domain name to ip address
			// ? add phpdev --machine=phpdev hosts

			return errors.New("not yet implemented")
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "php",
				Aliases:     []string{"p"},
				Usage:       "The PHP version to install",
				Value:       "7.4",
				DefaultText: "7.4",
			},
			&cli.StringFlag{
				Name:        "database",
				Aliases:     []string{"d", "db"},
				Usage:       "The database to install",
				Required:    false,
				Value:       "mariadb or postgres",
				DefaultText: "mariadb",
			},
		},
	}
}

func Mount() *cli.Command {
	return &cli.Command{
		Name:        "mount",
		Aliases:     []string{"m", "mnt"},
		Usage:       "Mount a folder to a machine",
		Description: "Mount a folder to use as a site in the machine",
		Action: func(c *cli.Context) error {

			// check if the path exists
			if _, err := os.Stat(c.String("path")); os.IsNotExist(err) {
				return err
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Aliases:     []string{"p"},
				Usage:       "The path to the folder to mount",
				DefaultText: "dummy",
				Required:    true,
			},
			&cli.StringFlag{
				Name:     "domain",
				Aliases:  []string{"d"},
				Usage:    "The domain name to mount into the machine",
				Required: true,
			},
		},
	}
}

func installPHP() *cli.Command {
	return &cli.Command{
		Name:    "php",
		Aliases: []string{"p"},
		Usage:   "Install PHP on a machine",
		Action: func(c *cli.Context) error {
			if err := action.InstallPHP(c); err != nil {
				return err
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "version",
				Aliases:     []string{"v"},
				Usage:       "Select which version of PHP to install",
				Value:       "7.4",
				DefaultText: "7.4",
			},
		},
	}
}

func installNginx() *cli.Command {
	return &cli.Command{
		Name:    "nginx",
		Aliases: []string{"n"},
		Usage:   "Install nginx on a machine",
		Action: func(c *cli.Context) error {
			if err := action.InstallNginx(c); err != nil {
				return err
			}

			return nil
		},
	}
}

func installMariaDB() *cli.Command {
	return &cli.Command{
		Name:    "maria",
		Aliases: []string{"m", "mariadb"},
		Usage:   "Install MariaDb Server on a machine",
		Action: func(c *cli.Context) error {
			if err := action.InstallMariaDB(c); err != nil {
				return err
			}

			return nil
		},
	}
}

func installRedis() *cli.Command {
	return &cli.Command{
		Name:    "redis",
		Aliases: []string{"r"},
		Usage:   "Install Redis on a machine",
		Action: func(c *cli.Context) error {
			if err := action.InstallRedis(c); err != nil {
				return err
			}

			return nil
		},
	}
}

func installPostgres() *cli.Command {
	return &cli.Command{
		Name:    "postgres",
		Aliases: []string{"postgresql", "pgsql"},
		Usage:   "Install PostgreSQL on a machine",
		Action: func(c *cli.Context) error {
			if err := action.InstallPostgres(c); err != nil {
				return err
			}

			return nil
		},
	}
}
