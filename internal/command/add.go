package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/validate"
)

func Add() *cli.Command {
	return &cli.Command{
		Name:   "add",
		Usage:  "Add virtual host",
		Before: addBeforeAction,
		Action: addAction,
		After:  addAfterAction,
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

func addAction(c *cli.Context) error {
	machine := c.String("machine")
	host := c.Args().First()
	php := c.String("php-version")
	dir := c.String("public-dir")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	if host == "" {
		return errors.New("missing param host")
	}

	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "bash", "/opt/nitro/nginx/add-site.sh", host, php, dir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func addAfterAction(c *cli.Context) error {
	return c.App.RunContext(c.Context, []string{c.App.Name, "--machine", c.String("machine"), "attach", c.Args().First(), c.Args().Get(1)})
}
