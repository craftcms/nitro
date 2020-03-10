package command

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/executor"
)

// Redis executes a shell for redis on a machine
func Redis(e executor.Executor) *cli.Command {
	return &cli.Command{
		Name:        "redis",
		Usage:       "Enter redis shell",
		Description: "Redis is installed by default on the machine, this command will drop you immediately into a shell inside the machine to run commands.",
		Action: func(c *cli.Context) error {
			return redisAction(c, e)
		},
	}
}

func redisAction(c *cli.Context, e executor.Executor) error {
	machine := c.String("machine")

	args := []string{"multipass", "exec", machine, "--", "redis-cli"}

	return e.Exec(e.Path(), args, os.Environ())
}
