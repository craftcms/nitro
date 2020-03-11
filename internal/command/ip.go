package command

import (
	"github.com/urfave/cli/v2"
)

// IP will look for a specific machine IP address by name
func IP(r Runner) *cli.Command {
	return &cli.Command{
		Name:  "ip",
		Usage: "Show machine IP address",
		Action: func(c *cli.Context) error {
			return ipAction(c, r)
		},
	}
}

func ipAction(c *cli.Context, r Runner) error {
	return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "bash", "/opt/nitro/ip.sh"})
}
