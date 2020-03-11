package command

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/validate"
)

var (
	ErrAttachNoHostArgProvided = errors.New("you must pass a domain name")
	ErrAttachNoPathArgProvided = errors.New("")
)

// Attach will mount a directory to a machine
func Attach(r Runner) *cli.Command {
	return &cli.Command{
		Name:   "attach",
		Usage:  "Add directory to machine",
		Before: attachBeforeAction,
		Action: func(c *cli.Context) error {
			return attachAction(c, r)
		},
		After: attachAfterAction,
	}
}

func attachBeforeAction(c *cli.Context) error {
	if host := c.Args().First(); host == "" {
		// TODO validate the host name with validate.Host(h)
		return ErrAttachNoHostArgProvided
	}

	if path := c.Args().Get(1); path == "" {
		return ErrAttachNoPathArgProvided
	}

	if err := validate.Path(c.Args().Get(1)); err != nil {
		return err
	}

	return nil
}

func attachAction(c *cli.Context, r Runner) error {
	return r.Run([]string{"mount", c.Args().Get(1), c.String("machine") + ":/home/ubuntu/sites/" + c.Args().First()})
}

func attachAfterAction(c *cli.Context) error {
	fmt.Println("attached", c.Args().Get(1), "to", "/home/ubuntu/sites/"+c.Args().First())

	fmt.Println("\nto edit your hosts file, run the following as sudo:")
	fmt.Println("\nsudo", c.App.Name, "--machine", c.String("machine"), "hosts", c.Args().First())

	return nil
}
