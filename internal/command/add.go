package command

import (
	"errors"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal"
	"github.com/craftcms/nitro/internal/validate"
)

func Add(r internal.Runner) *cli.Command {
	return &cli.Command{
		Name:   "add",
		Usage:  "Add virtual host",
		Before: addBeforeAction,
		Action: func(c *cli.Context) error {
			return addAction(c, r)
		},
		After: addAfterAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "php-version",
				Usage:       "Provide version of PHP",
				Value:       "7.4",
				DefaultText: "7.4",
			},
			&cli.StringFlag{
				Name:        "public-dir",
				Usage:       "The public directory for the server",
				Value:       "web",
				DefaultText: "web",
			},
		},
	}
}

func addBeforeAction(c *cli.Context) error {
	if host := c.Args().First(); host == "" {
		// TODO validate the domain name with validate.Domain(d)
		return errors.New("you must pass a domain name")
	}

	if path := c.Args().Get(1); path == "" {
		// TODO validate the domain name with validate.Domain(d)
		return errors.New("you must provide a path to mount")
	}

	if err := validate.PHPVersion(c.String("php-version")); err != nil {
		return err
	}

	if err := validate.Path(c.Args().Get(1)); err != nil {
		return err
	}

	return nil
}

func addAction(c *cli.Context, r internal.Runner) error {
	machine := c.String("machine")
	host := c.Args().First()
	php := c.String("php-version")
	dir := c.String("public-dir")

	return r.Run([]string{"exec", machine, "--", "sudo", "bash", "/opt/nitro/nginx/add-site.sh", host, php, dir})
}

func addAfterAction(c *cli.Context) error {
	return c.App.RunContext(c.Context, []string{c.App.Name, "--machine", c.String("machine"), "attach", c.Args().First(), c.Args().Get(1)})
}
