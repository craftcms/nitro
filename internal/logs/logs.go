package logs

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/nitro/internal/action"
)

func Command(e action.CommandLineExecutor) *cli.Command {
	return &cli.Command{
		Name:  "logs",
		Usage: "Show logs for a machine",
		Action: func(c *cli.Context) error {
			return cli.ShowSubcommandHelp(c)
		},
		Subcommands: []*cli.Command{
			{
				Name:        "nginx",
				Description: "Show logs from nginx",
				Action: func(c *cli.Context) error {
					return logsNginx(c, e)
				},
			},
		},
	}
}

func logsNginx(c *cli.Context, e action.CommandLineExecutor) error {
	machine := c.String("machine")

	args := []string{"multipass", "exec", machine, "--", "sudo", "bash", "/opt/nitro/nginx/tail-logs.sh"}

	return e.Exec(e.Path(), args, os.Environ())
}
