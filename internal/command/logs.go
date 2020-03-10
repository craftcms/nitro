package command

import (
	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal"
)

// Logs show system logs for a machine
func Logs(r internal.Runner) *cli.Command {
	return &cli.Command{
		Name:  "logs",
		Usage: "Show machine logs",
		Action: func(c *cli.Context) error {
			return cli.ShowSubcommandHelp(c)
		},
		Subcommands: []*cli.Command{
			{
				Name:        "nginx",
				Description: "Show logs from nginx",
				Action: func(c *cli.Context) error {
					return logsNginx(c, r)
				},
			},
		},
	}
}

func logsNginx(c *cli.Context, r internal.Runner) error {
	return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "bash", "/opt/nitro/nginx/tail-logs.sh"})
}
