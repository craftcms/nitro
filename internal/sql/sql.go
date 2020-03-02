package sql

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/executor"
)

// SQL Command creates a shell command into the provided database shell as a root user.
func Command(executor executor.Executor) *cli.Command {
	return &cli.Command{
		Name:  "sql",
		Usage: "Enter a shell for the database",
		Description: `Bootstrapping a machine will install mariadb by default, this command allows you to enter a SQL shell as a root user. 
	 If you chose postgres as the database, you can pass a flag --pgsql to enter a postgres shell`,
		Action: func(c *cli.Context) error {
			return run(c, executor)
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "pgsql",
				Usage:       "Enter a postgres shell",
				Value:       false,
				DefaultText: "false",
			},
		},
	}
}

func run(c *cli.Context, e executor.Executor) error {
	machine := c.String("machine")
	pgsql := c.Bool("pgsql")

	if pgsql {
		args := []string{"multipass", "exec", machine, "--", "sudo", "-u", "postgres", "psql"}

		return e.Exec(e.Path(), args, os.Environ())
	}

	args := []string{"multipass", "exec", machine, "--", "sudo", "mysql", "-u", "root"}

	return e.Exec(e.Path(), args, os.Environ())
}
