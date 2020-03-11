package command

import (
	"github.com/urfave/cli/v2"
)

// SSH SSH will login a user to a specific machine
func SSH(r Runner) *cli.Command {
	return &cli.Command{
		Name:  "ssh",
		Usage: "SSH into machine",
		Action: func(c *cli.Context) error {
			return sshAction(c, r)
		},
	}
}

func sshAction(c *cli.Context, r Runner) error {
	r.UseSyscall(true)

	return r.Run([]string{"shell", c.String("machine")})
}
