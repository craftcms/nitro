package password

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/nitro/internal/action"
)

func Command(e action.CommandLineExecutor) *cli.Command {
	return &cli.Command{
		Name:        "password",
		Usage:       "Get the database password for the user nitro",
		Description: "Regardless of the database engine, there is one password for the non-root user. This password is unique to each machine and generated on startup.",
		Action: func(c *cli.Context) error {
			return run(c, e)
		},
	}
}

func run(c *cli.Context, e action.CommandLineExecutor) error {
	return e.Exec(e.Path(), []string{"multipass", "exec", c.String("machine"), "--", "cat", "/home/ubuntu/.db_password"}, os.Environ())
}
