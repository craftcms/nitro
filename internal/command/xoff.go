package command

import (
	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal"
)

// XOff will disable xdebug on a machine
func XOff(r internal.Runner) *cli.Command {
	return &cli.Command{
		Name:        "xoff",
		Usage:       "Disable Xdebug",
		Description: "Disable Xdebug on machine",
		Action: func(c *cli.Context) error {
			return xOffAction(c, r)
		},
	}
}

func xOffAction(c *cli.Context, r internal.Runner) error {
	return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "bash", "/opt/nitro/php/disable-xdebug.sh"})
}
