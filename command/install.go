package command

import (
	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/nitro/action"
	"github.com/pixelandtonic/nitro/validate"
)

func installPHP() *cli.Command {
	return &cli.Command{
		Name:    "php",
		Aliases: []string{"p"},
		Usage:   "Install PHP on a machine",
		Action: func(c *cli.Context) error {
			if err := validate.PHPVersion(c.String("version")); err != nil {
				return err
			}

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
