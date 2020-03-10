package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

// SSH SSH will login a user to a specific machine
func SSH() *cli.Command {
	return &cli.Command{
		Name:   "ssh",
		Usage:  "SSH into machine",
		Action: sshAction,
	}
}

func sshAction(c *cli.Context) error {
	machine := c.String("machine")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	cmd := exec.Command(
		multipass,
		"shell",
		machine,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
