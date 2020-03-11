package command

import (
	"github.com/urfave/cli/v2"
)

func Start(r Runner) *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start machine",
		Action: func(c *cli.Context) error {
			return r.Run([]string{"start", c.String("machine")})
		},
	}
}
