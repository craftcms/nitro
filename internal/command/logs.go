package command

import (
	"github.com/urfave/cli/v2"
)

// Logs show system logs for a machine
func Logs(r Runner) *cli.Command {
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
			{
				Name:        "xdebug",
				Description: "Show logs from xdebug",
				Action: func(c *cli.Context) error {
					return logsXdebug(c, r)
				},
			},
		},
	}
}

func logsNginx(c *cli.Context, r Runner) error {
	return r.Run([]string{"exec", c.String("machine"), "--", "sudo", "bash", "/opt/nitro/nginx/tail-logs.sh"})
}

func logsXdebug(c *cli.Context, r Runner) error {
	return r.Run([]string{"exec", c.String("machine"), "--", "tail", "-f", "/var/log/nitro/xdebug.log"})
}
