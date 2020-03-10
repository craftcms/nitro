package command

import (
	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal"
)

// SSH SSH will login a user to a specific machine
func SSH(runner internal.Runner) *cli.Command {
	return &cli.Command{
		Name:  "ssh",
		Usage: "SSH into machine",
		Action: func(c *cli.Context) error {
			return sshAction(c, runner)
		},
	}
}

func sshAction(c *cli.Context, runner internal.Runner) error {
	runner.UseSyscall(true)

	return runner.Run([]string{"shell", c.String("machine")})
}
