package command

import (
	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal"
)

// SSH SSH will login a user to a specific machine
func SSH(r internal.Runner) *cli.Command {
	return &cli.Command{
		Name:  "ssh",
		Usage: "SSH into machine",
		Action: func(c *cli.Context) error {
			return sshAction(c, r)
		},
	}
}

func sshAction(c *cli.Context, r internal.Runner) error {
	r.UseSyscall(true)

	return r.Run([]string{"shell", c.String("machine")})
}
