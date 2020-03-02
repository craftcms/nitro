package destroy

import (
	"errors"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        "destroy",
		Usage:       "Destroy a machine",
		Description: "By default, when deleting a machine it is soft-deleted which means it can be recovered. This command will destroy the machine making it unrecoverable.",
		Action: func(c *cli.Context) error {
			return errors.New("not implemented")
		},
	}
}
