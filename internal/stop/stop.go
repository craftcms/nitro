package stop

import (
	"errors"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "Stop the machine",
		Action: func(c *cli.Context) error {
			machine := c.String("machine")
			return errors.New(machine)
		},
	}
}
