package command

import (
	"github.com/urfave/cli/v2"
)

func Update(r Runner) *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update machine",
		Action: func(c *cli.Context) error {
			return updateAction(c, r)
		},
	}
}

func updateAction(c *cli.Context, r Runner) error {
	machine := c.String("machine")

	r.UseSyscall(true)

	return r.Run([]string{"exec", machine, "--", "sudo", "bash", "/opt/nitro/update.sh"})
}
