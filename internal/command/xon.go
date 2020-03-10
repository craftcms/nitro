package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

// Xon will enable xdebug on a machine
func XOn() *cli.Command {
	return &cli.Command{
		Name:        "xon",
		Usage:       "Enable Xdebug",
		Description: "Enable Xdebug for machine",
		Action:      xOnAction,
	}
}

func xOnAction(c *cli.Context) error {
	machine := c.String("machine")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "bash", "/opt/nitro/php/enable-xdebug.sh")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
