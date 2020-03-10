package command

import (
	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal"
)

// Redis executes a shell for redis on a machine
func Redis(r internal.Runner) *cli.Command {
	return &cli.Command{
		Name:        "redis",
		Usage:       "Enter redis shell",
		Description: "Redis is installed by default on the machine, this command will drop you immediately into a shell inside the machine to run commands.",
		Action: func(c *cli.Context) error {
			return redisAction(c, r)
		},
	}
}

func redisAction(c *cli.Context, r internal.Runner) error {
	r.UseSyscall(true)
	return r.Run([]string{"multipass", "exec", c.String("machine"), "--", "redis-cli"})
}
