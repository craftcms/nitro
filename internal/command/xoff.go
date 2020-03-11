package command

import (
	"github.com/urfave/cli/v2"
)

// XOff will disable xdebug on a machine
func XOff(r Runner) *cli.Command {
	return &cli.Command{
		Name:        "xoff",
		Usage:       "Disable Xdebug",
		Description: "Disable Xdebug on machine",
		Action: func(c *cli.Context) error {
			return xOffAction(c, r)
		},
	}
}

func xOffAction(c *cli.Context, r Runner) error {
	return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "bash", "/opt/nitro/php/disable-xdebug.sh"})
}
