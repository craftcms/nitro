package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

// XOff will disable xdebug on a machine
func XOff() *cli.Command {
	return &cli.Command{
		Name:        "xoff",
		Usage:       "Disable Xdebug",
		Description: "Disable Xdebug on machine",
		Action:      xOffAction,
	}
}

func xOffAction(c *cli.Context) error {
	machine := c.String("machine")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "bash", "/opt/nitro/php/disable-xdebug.sh")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
