package command

import (
	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal"
)

func Update(r internal.Runner) *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update machine",
		Action: func(c *cli.Context) error {
			return updateAction(c, r)
		},
	}
}

func updateAction(c *cli.Context, r internal.Runner) error {
	machine := c.String("machine")

	r.UseSyscall(true)

	return r.Run([]string{"exec", machine, "--", "sudo", "bash", "/opt/nitro/update.sh"})
}
