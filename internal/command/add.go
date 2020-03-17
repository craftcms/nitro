package command

import (
	"errors"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/validate"
)

func Add(r Runner) *cli.Command {
	cwd, _ := os.Getwd()
	pathFlag.DefaultText = cwd

	return &cli.Command{
		Name:   "add",
		Usage:  "Add virtual host",
		Before: addBeforeAction,
		Action: func(c *cli.Context) error {
			return addAction(c, r)
		},
		After: addAfterAction,
		Flags: []cli.Flag{
			phpVersionFlag,
			pathFlag,
			publicDirFlag,
		},
	}
}

func addBeforeAction(c *cli.Context) error {
	if host := c.Args().First(); host == "" {
		// TODO validate the domain name with validate.Domain(d)
		return errors.New("you must pass a domain name")
	}

	if err := validate.PHPVersion(c.String("php-version")); err != nil {
		return err
	}

	if c.String("path") == "" {
		return nil
	}

	if err := validate.Path(c.String("path")); err != nil {
		return err
	}

	return nil
}

func addAction(c *cli.Context, r Runner) error {
	machine := c.String("machine")
	host := c.Args().First()
	php := c.String("php-version")
	dir := c.String("public-dir")

	return r.Run([]string{"exec", machine, "--", "sudo", "bash", "/opt/nitro/nginx/add-site.sh", host, php, dir})
}

func addAfterAction(c *cli.Context) error {
	host := c.Args().First()

	var path string
	if c.String("path") == "" {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return err
		}

		return c.App.RunContext(c.Context, []string{c.App.Name, "--machine", c.String("machine"), "attach", host, path})
	}

	return c.App.RunContext(c.Context, []string{c.App.Name, "--machine", c.String("machine"), "attach", host, c.String("path")})
}
