package xdebug

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/executor"
)

func CommandOn(e executor.Executor) *cli.Command {
	return &cli.Command{
		Name:        "xon",
		Usage:       "Enable Xdebug",
		Description: "Enable Xdebug for machine",
		Action: func(c *cli.Context) error {
			return e.Exec(e.Path(), []string{"multipass", "exec", c.String("machine"), "--", "sudo", "bash", "/opt/nitro/php/enable-xdebug.sh"}, os.Environ())
		},
	}
}

func CommandOff(e executor.Executor) *cli.Command {
	return &cli.Command{
		Name:        "xoff",
		Usage:       "Disable Xdebug",
		Description: "Disable Xdebug on machine",
		Action: func(c *cli.Context) error {
			return e.Exec(e.Path(), []string{"multipass", "exec", c.String("machine"), "--", "sudo", "bash", "/opt/nitro/php/disable-xdebug.sh"}, os.Environ())
		},
	}
}
