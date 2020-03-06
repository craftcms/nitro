package sql

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/executor"
)

// SQL Command creates a shell command into the provided database shell as a root user.
func Command(e executor.Executor) *cli.Command {
	return &cli.Command{
		Name:  "sql",
		Usage: "Enter database shell",
		Description: `Bootstrapping a machine will install mariadb by default, this command allows you to enter a SQL shell as a root user. 
	 If you chose postgres as the database, you can pass a flag --postgres to enter a postgres shell`,
		Action: func(c *cli.Context) error {
			return handle(c, e)
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "postgres",
				Usage:       "Enter a postgres shell",
				Value:       false,
				DefaultText: "false",
			},
		},
	}
}

func handle(c *cli.Context, e executor.Executor) error {
	machine := c.String("machine")

	if c.Bool("postgres") {
		args := []string{"multipass", "exec", machine, "--", "sudo", "-u", "postgres", "psql"}

		return e.Exec(e.Path(), args, os.Environ())
	}

	args := []string{"multipass", "exec", machine, "--", "sudo", "mysql", "-u", "root"}

	return e.Exec(e.Path(), args, os.Environ())
}
