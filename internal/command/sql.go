package command

import (
	"fmt"
	"os/exec"

	"github.com/urfave/cli/v2"
)

// SQL creates a shell command into the provided database shell as a root user.
func SQL() *cli.Command {
	return &cli.Command{
		Name:  "sql",
		Usage: "Enter database shell",
		Description: `Bootstrapping a machine will install mariadb by default, this command allows you to enter a SQL shell as a root user. 
	 If you chose postgres as the database, you can pass a flag --postgres to enter a postgres shell`,
		Action: sqlAction,
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

func sqlAction(c *cli.Context) error {
	machine := c.String("machine")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	cmd := exec.Command(multipass)
	if c.Bool("postgres") {
		cmd.Args = []string{"exec", "--name", machine, "--", "sudo", "-u", "postgres", "psql"}
	} else {
		cmd.Args = []string{"exec", "--name", machine, "--", "sudo", "mysql", "-u", "root"}
	}

	return cmd.Run()
}
