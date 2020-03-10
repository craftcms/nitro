package command

import (
	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal"
)

// IP will look for a specific machine IP address by name
func IP(r internal.Runner) *cli.Command {
	return &cli.Command{
		Name:  "ip",
		Usage: "Show machine IP address",
		Action: func(c *cli.Context) error {
			return ipAction(c, r)
		},
	}
}

func ipAction(c *cli.Context, r internal.Runner) error {
	return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "bash", "/opt/nitro/ip.sh"})
}
