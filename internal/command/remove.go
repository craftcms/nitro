package command

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	ErrRemoveNoHostArgProvided = errors.New("no host was specified for removal")
)

// Remove will remove a host from a machine
func Remove(r Runner) *cli.Command {
	return &cli.Command{
		Name:   "remove",
		Usage:  "Remove virtual host",
		Before: removeBeforeAction,
		Action: func(c *cli.Context) error {
			return removeAction(c, r)
		},
		After: removeAfterAction,
	}
}

func removeBeforeAction(c *cli.Context) error {
	if c.Args().First() == "" {
		return ErrRemoveNoHostArgProvided
	}

	return nil
}

func removeAction(c *cli.Context, r Runner) error {
	machine := c.String("machine")
	host := c.Args().First()

	return r.Run([]string{"exec", "--name", machine, "--", "sudo", "bash", "/opt/nitro/nginx/remove-host.sh", host})
}

func removeAfterAction(c *cli.Context) error {
	fmt.Println("removed host", c.Args().First())

	return nil
}
