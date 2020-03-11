package command

import (
	"github.com/urfave/cli/v2"
)

// Redis executes a shell for redis on a machine
func Redis(r Runner) *cli.Command {
	return &cli.Command{
		Name:        "redis",
		Usage:       "Enter redis shell",
		Description: "Redis is installed by default on the machine, this command will drop you immediately into a shell inside the machine to run commands.",
		Action: func(c *cli.Context) error {
			return redisAction(c, r)
		},
	}
}

func redisAction(c *cli.Context, r Runner) error {
	r.UseSyscall(true)
	return r.Run([]string{"multipass", "exec", c.String("machine"), "--", "redis-cli"})
}
