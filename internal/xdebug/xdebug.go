package xdebug

import (
	"errors"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        "xdebug",
		Usage:       "Enable Xdebug",
		Description: "Calling Xdebug will default to enabling the extension, if the flag --disable is passed it will be disabled",
		Action: func(c *cli.Context) error {
			return errors.New("not implemented")
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "disable",
				Usage:       "Disable Xdebug",
				Value:       false,
				DefaultText: "false",
			},
		},
	}
}
