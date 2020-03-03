package bootstrap

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/executor"
)

func Command(e executor.Executor) *cli.Command {
	return &cli.Command{
		Name:  "bootstrap",
		Usage: "Bootstrap machine",
		Action: func(c *cli.Context) error {
			return run(c, e)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "php-version",
				Usage:       "Provide version of PHP",
				Value:       "7.4",
				DefaultText: "7.4",
			},
			&cli.StringFlag{
				Name:        "database",
				Usage:       "Provide version of PHP",
				Value:       "mariadb",
				DefaultText: "mariadb",
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
