package redis

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/executor"
)

func Command(e executor.Executor) *cli.Command {
	return &cli.Command{
		Name:        "redis",
		Usage:       "Enter a redis shell",
		Description: "Redis is installed by default on the machine, this command will drop you immediately into a shell inside the machine to run commands.",
		Action: func(c *cli.Context) error {
			return run(c, e)
		},
	}
}

func run(c *cli.Context, e executor.Executor) error {
	machine := c.String("machine")

	args := []string{"multipass", "exec", machine, "--", "redis-cli"}

	return e.Exec(e.Path(), args, os.Environ())
}
