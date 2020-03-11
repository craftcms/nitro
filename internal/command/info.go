package command

import "github.com/urfave/cli/v2"

// Info will display system information on a machine
func Info(r Runner) *cli.Command {
	return &cli.Command{
		Name:  "info",
		Usage: "Show information about machine",
		Action: func(c *cli.Context) error {
			return infoAction(c, r)
		},
	}
}

func infoAction(c *cli.Context, r Runner) error {
	return r.Run([]string{"info", c.String("machine")})
}
