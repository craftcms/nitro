package command

import (
	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal"
)

func Stop(r internal.Runner) *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "Stop machine",
		Action: func(c *cli.Context) error {
			return stopAction(c, r)
		},
	}
}

func stopAction(c *cli.Context, r internal.Runner) error {
	return r.Run([]string{"stop", c.String("machine")})
}
