package command

import (
	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/dev/actions"
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
			if err := actions.Initialize(c); err != nil {
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
			if err := actions.Update(c); err != nil {
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
		Aliases: []string{"ssh"},
		Usage:   "SSH into a machine as administrator",
		Action: func(c *cli.Context) error {
			if err := actions.SSH(c); err != nil {
				return err
			}
			return nil
		},
		Flags: []cli.Flag{&machineFlag},
	}
}

func Prepare() *cli.Command {
	return &cli.Command{
		Name:    "prepare",
		Aliases: []string{"p"},
		Usage:   "Prepare the machine for PHP development",
		Action: func(c *cli.Context) error {
			if err := actions.Prepare(c); err != nil {
				return err
			}
			return nil
		},
		Flags: []cli.Flag{
			&machineFlag,
			&cli.StringFlag{
				Name:        "php",
				Value:       "7.4",
				Usage:       "Select a PHP version to install",
				DefaultText: "7.4",
			},
		},
	}
}

func Build() *cli.Command {
	return &cli.Command{
		Name:    "build",
		Aliases: []string{"b"},
		Usage:   "Build the machine for PHP development",
		Action: func(c *cli.Context) error {
			if err := actions.Build(c); err != nil {
				return err
			}
			return nil
		},
		Flags: []cli.Flag{
			&machineFlag,
		},
	}
}
