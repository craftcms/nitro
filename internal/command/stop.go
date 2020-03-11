package command

import (
	"github.com/urfave/cli/v2"
)

func Stop(r Runner) *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "Stop machine",
		Action: func(c *cli.Context) error {
			return stopAction(c, r)
		},
	}
}

func stopAction(c *cli.Context, r Runner) error {
	return r.Run([]string{"stop", c.String("machine")})
}
