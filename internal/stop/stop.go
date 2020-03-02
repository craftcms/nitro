package stop

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "Stop the machine",
		Action: func(c *cli.Context) error {
			return run(c)
		},
	}
}

func run(c *cli.Context) error {
	machine := c.String("machine")

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	cmd := exec.Command(multipass, "stop", machine)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}