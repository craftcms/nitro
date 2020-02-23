package command

import (
	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/dev/action"
)

var (
	machineFlag = cli.StringFlag{
		Name:        "machine, m",
		Value:       "dev",
		Usage:       "Provide a machine name",
		DefaultText: "dev",
	}
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
		Flags: []cli.Flag{&machineFlag},
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
		Flags: []cli.Flag{&machineFlag},
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
		Flags: []cli.Flag{&machineFlag},
	}
}

func Install() *cli.Command {
	return &cli.Command{
		Name:        "install",
		Usage:       "",
		Description: "Install offers common options for installing packages on a machine",
		Action: func(c *cli.Context) error {
			return nil
		},
		Flags: []cli.Flag{&machineFlag},
		Subcommands: []*cli.Command{
			installPHP(),
			installNginx(),
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
			&machineFlag,
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
		Flags: []cli.Flag{
			&machineFlag,
		},
	}
}
