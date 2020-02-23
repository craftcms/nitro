package command

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"bufio"

	"github.com/pixelandtonic/dev/action"
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
			reader := bufio.NewReader(os.Stdin)

			fmt.Println("What is your name?")
			confirm, _ := reader.ReadString("\n")


			if err := action.Delete(c); err != nil {
				return err
			}
			return nil
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
		Aliases: []string{"m"},
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
