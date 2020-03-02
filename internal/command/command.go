package command

import (
	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/action"
)

func Stop() *cli.Command {
	return &cli.Command{
		Name:        "stop",
		Usage:       "Stop a machine",
		Description: "Stop a machine when not in use (this does not delete the machine)",
		Action: func(c *cli.Context) error {
			return action.Stop(c)
		},
	}
}

