package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/validate"
)

// Attach will mount a directory to a machine
func Attach() *cli.Command {
	return &cli.Command{
		Name:   "attach",
		Usage:  "Add directory to machine",
		Before: attachBeforeAction,
		Action: attachAction,
	}
}

func attachBeforeAction(c *cli.Context) error {
	if host := c.Args().First(); host == "" {
		// TODO validate the domain name with validate.Domain(d)
		return errors.New("you must pass a domain name")
	}

	if path := c.Args().Get(1); path == "" {
		// TODO validate the domain name with validate.Domain(d)
		return errors.New("you must provide a path to mount")
	}

	if err := validate.Path(c.Args().Get(1)); err != nil {
		return err
	}

	return nil
}

func attachAction(c *cli.Context) error {
	machine := c.String("machine")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	host := c.Args().First()
	path := c.Args().Get(1)

	cmd := exec.Command(
		multipass,
		"mount",
		path,
		machine+":/home/ubuntu/sites/"+host,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
