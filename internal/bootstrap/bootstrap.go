package bootstrap

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/validate"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "bootstrap",
		Usage: "Bootstrap machine",
		Before: func(c *cli.Context) error {
			return beforeAction(c)
		},
		Action: func(c *cli.Context) error {
			return handle(c)
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

func beforeAction(c *cli.Context) error {
	if err := validate.PHPVersion(c.String("php-version")); err != nil {
		return err
	}

	if err := validate.Database(c.String("database")); err != nil {
		return err
	}

	return nil
}

func handle(c *cli.Context) error {
	machine := c.String("machine")
	php := c.String("php-version")
	database := c.String("database")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	// create the machine
	cmd := exec.Command(
		multipass,
		"exec",
		machine,
		"--",
		"sudo",
		"bash",
		"/opt/nitro/bootstrap.sh",
		php,
		database,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
