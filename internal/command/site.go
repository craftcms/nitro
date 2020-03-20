package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/validate"
)

func Site(r Runner) *cli.Command {
	cwd, _ := os.Getwd()
	pathFlag.DefaultText = cwd

	return &cli.Command{
		Name:   "site",
		Usage:  "Add virtual host",
		Before: siteBeforeAction,
		Action: func(c *cli.Context) error {
			return siteAction(c, r)
		},
		After: siteAfterAction,
		Flags: []cli.Flag{
			phpVersionFlag,
			pathFlag,
			publicDirFlag,
			removeFlag,
		},
	}
}

func siteBeforeAction(c *cli.Context) error {
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

func siteAction(c *cli.Context, r Runner) error {
	host := c.Args().First()
	machine := c.String("machine")
	php := c.String("php-version")
	dir := c.String("public-dir")

	if c.Bool("remove") {
		fmt.Println("in remove")
		return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "bash", "/opt/nitro/nginx/remove-site.sh", c.Args().First()})
	}

	return r.Run([]string{"exec", machine, "--", "sudo", "bash", "/opt/nitro/nginx/add-site.sh", host, php, dir})
}

func siteAfterAction(c *cli.Context) error {
	host := c.Args().First()

	if c.Bool("remove") {
		fmt.Println("removed site", host)
		return nil
	}

	// get the current working directory if no path is provided
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
