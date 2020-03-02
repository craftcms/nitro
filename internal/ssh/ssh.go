package ssh

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/executor"
)

// Command SSH will login a user to a specific machine
func Command(e executor.Executor) *cli.Command {
	return &cli.Command{
		Name:  "ssh",
		Usage: "SSH into a machine",
		Action: func(c *cli.Context) error {
			return run(c.String("machine"), e)
		},
	}
}

func run(m string, e executor.Executor) error {
	return e.Exec(e.Path(), []string{"multipass", "shell", m}, os.Environ())
}
