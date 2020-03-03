package xdebug

import (
	"errors"

	"github.com/urfave/cli/v2"
)

func CommandOn() *cli.Command {
	return &cli.Command{
		Name:        "xon",
		Usage:       "Enable Xdebug",
		Description: "Enable Xdebug for machine",
		Action: func(c *cli.Context) error {
			return errors.New("not implemented")
		},
	}
}

func CommandOff() *cli.Command {
	return &cli.Command{
		Name:        "xoff",
		Usage:       "Disable Xdebug",
		Description: "Disable Xdebug on machine",
		Action: func(c *cli.Context) error {
			return errors.New("not implemented")
		},
	}
}
