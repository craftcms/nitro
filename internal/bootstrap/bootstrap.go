package bootstrap

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/nitro/internal/executor"
)

func Command(e executor.Executor) *cli.Command {
	return &cli.Command{
		Name:  "bootstrap",
		Usage: "Bootstrap the installation of a new machine",
		Action: func(c *cli.Context) error {
			return run(c, e)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "php-version",
				Usage:       "Provide version of PHP",
				DefaultText: "7.4",
			},
		},
	}
}

func run(c *cli.Context, e executor.Executor) error {
	machine := c.String("machine")
	php := c.String("php-version")
	database := c.String("database")

	args := []string{"multipass", "exec", machine, "--", "sudo", "bash", "/opt/nitro/bootstrap.sh", php, database}

	return e.Exec(e.Path(), args, os.Environ())
}
