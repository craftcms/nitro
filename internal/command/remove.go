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
func Remove() *cli.Command {
	return &cli.Command{
		Name:   "remove",
		Before: removeBeforeAction,
		Action: removeAction,
		After:  removeAfterAction,
	}
}

func removeBeforeAction(c *cli.Context) error {
	if c.Args().First() == "" {
		return ErrRemoveNoHostArgProvided
	}

	return nil
}

func removeAction(c *cli.Context) error {
	return nil
}

func removeAfterAction(c *cli.Context) error {
	fmt.Println("removed host", c.Args().First())

	return nil
}
