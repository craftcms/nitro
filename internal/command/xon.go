package command

import (
	"github.com/urfave/cli/v2"
)

// Xon will enable xdebug on a machine
func XOn(r Runner) *cli.Command {
	return &cli.Command{
		Name:        "xon",
		Usage:       "Enable Xdebug",
		Description: "Enable Xdebug for machine",
		Action: func(c *cli.Context) error {
			return xOnAction(c, r)
		},
	}
}

func xOnAction(c *cli.Context, r Runner) error {
	return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "bash", "/opt/nitro/php/enable-xdebug.sh"})
}
