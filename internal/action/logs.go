package action

import (
	"os"

	"github.com/urfave/cli/v2"
)

func LogsNginx(c *cli.Context, e CommandLineExecutor) error {
	machine := c.String("machine")

	args := []string{"multipass", "exec", machine, "--", "sudo", "bash", "/opt/nitro/nginx/tail-logs.sh"}

	return e.Exec(e.Path(), args, os.Environ())
}
