package ssh

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/nitro/internal/action"
)

// Command SSH will login a user to a specific machine
func Command(e action.CommandLineExecutor) *cli.Command {
	return &cli.Command{
		Name:  "ssh",
		Usage: "SSH into a machine as administrator",
		Action: func(c *cli.Context) error {
			return run(c.String("machine"), e)
		},
	}
}

func run(m string, e action.CommandLineExecutor) error {
	return e.Exec(e.Path(), []string{"multipass", "shell", m}, os.Environ())
}
