package command

import (
	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal"
)

// SQL creates a shell command into the provided database shell as a root user.
func SQL(r internal.Runner) *cli.Command {
	return &cli.Command{
		Name:  "sql",
		Usage: "Enter database shell",
		Description: `Bootstrapping a machine will install mariadb by default, this command allows you to enter a SQL shell as a root user. 
	 If you chose postgres as the database, you can pass a flag --postgres to enter a postgres shell`,
		Action: func(c *cli.Context) error {
			return sqlAction(c, r)
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

func sqlAction(c *cli.Context, r internal.Runner) error {
	machine := c.String("machine")

	r.UseSyscall(true)

	if c.Bool("postgres") {
		return r.Run([]string{"exec", machine, "--", "sudo", "-u", "postgres", "psql"})
	}

	return r.Run([]string{"exec", machine, "--", "sudo", "mysql", "-u", "root"})
}
