package command

import (
	"log"
	"strings"
)

type SSHCommand struct {
	*CoreCommand
}

func (c *SSHCommand) Synopsis() string {
	return "SSH into machine"
}

func (c *SSHCommand) Help() string {
	return strings.TrimSpace(`
Usage: nitro ssh [options]
  This command allows you to SSH into a virtual machine.
  
  SSH to a virtual machine:
      $ nitro ssh -name diesel
  
  SSH to a virtual machine using the config file:
      $ nitro ssh
`)
}

func (c *SSHCommand) Run(args []string) int {
	if err := c.Flags().Parse(args); err != nil {
		c.UI.Error(err.Error())
		log.Fatal("in the parse")
		return 1
	}

	if err := c.Runner.Run([]string{"shell", c.flagName}); err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	return 0
}
