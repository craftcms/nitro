package host

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/executor"
	"github.com/craftcms/nitro/internal/validate"
)

func Command(e executor.Executor) *cli.Command {
	return &cli.Command{
		Name:  "add-host",
		Usage: "Add a virtual host",
		Action: func(c *cli.Context) error {
			return run(c, e)
		},
		Before: func(c *cli.Context) error {
			if c.Args().First() == "" {
				// TODO validate the domain name with validate.Domain(d)
				return errors.New("you must pass a domain name")
			}

			if err := validate.PHPVersion(c.String("php-version")); err != nil {
				return err
			}

			if err := validate.Path(c.String("path")); err != nil {
				return err
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "php-version",
				Usage:       "Provide version of PHP",
				DefaultText: "7.4",
			},
			&cli.StringFlag{
				Name:     "path",
				Usage:    "The path to the directory to mount",
				Required: true,
			},
		},
	}
}

func run(c *cli.Context, e executor.Executor) error {
	machine := c.String("machine")
	host := c.Args().First()
	php := c.String("php-version")

	if host == "" {
		return errors.New("missing param host")
	}

	if php == "" {
		fmt.Println("missing php-version")
		php = "7.4"
	}

	args := []string{"multipass", "exec", machine, "--", "sudo", "bash", "/opt/nitro/nginx/add-site.sh", host, php}

	return e.Exec(e.Path(), args, os.Environ())
}
