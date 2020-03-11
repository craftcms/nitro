package command

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	ErrDetachNoPathArgProvided = errors.New("missing the host to remove")
)

func Detach(r Runner) *cli.Command {
	return &cli.Command{
		Name: "detach",
		Before: func(c *cli.Context) error {
			return detachBeforeAction(c)
		},
		Action: func(c *cli.Context) error {
			return detachAction(c, r)
		},
	}
}

func detachBeforeAction(c *cli.Context) error {
	if c.Args().First() == "" {
		return ErrDetachNoPathArgProvided
	}

	return nil
}

func detachAction(c *cli.Context, r Runner) error {
	fmt.Println(fmt.Sprintf("removed mount %v from machine", c.Args().First()))
	return r.Run([]string{"umount", c.String("machine") + ":/home/ubuntu/sites/" + c.Args().First()})
}
