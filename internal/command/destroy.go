package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Destroy will completely destroy a machine
func Destroy(r Runner) *cli.Command {
	return &cli.Command{
		Name:        "destroy",
		Usage:       "Destroy machine",
		Description: "By default, when deleting a machine it is soft-deleted which means it can be recovered. This command will destroy the machine making it unrecoverable.",
		Action: func(c *cli.Context) error {
			return destroyAction(c, r)
		},
		After: func(c *cli.Context) error {
			if c.Bool("permanent") {
				fmt.Println("permanently deleted machine", c.String("machine"))
				return nil
			}

			fmt.Println("deleted machine", c.String("machine"))
			fmt.Println("")
			fmt.Println("If this was an accident, this machine can still be recovered, run the following command to restore the machine:")
			fmt.Println("  $ multipass recover", c.String("machine"))
			return nil
		},
		Flags: []cli.Flag{
			permanentDeleteFlag,
		},
	}
}

func destroyAction(c *cli.Context, r Runner) error {
	if c.Bool("permanent") {
		return r.Run([]string{"delete", c.String("machine"), "--purge"})
	}

	return r.Run([]string{"delete", c.String("machine")})
}
