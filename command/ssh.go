package command

import (
	"flag"
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

func (c *SSHCommand) Flags() *flag.FlagSet {
	s := flag.NewFlagSet("ssh", 0)
	s.StringVar(&c.flagName, "name", "", "name of the machine")
	s.BoolVar(&c.flagDryRun, "dry-run", false, "skip executing the command")
	return s
}

func (c *SSHCommand) Run(args []string) int {
	if err := c.Flags().Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if c.flagName == "" {
		if c.Config.IsSet("name") {
			c.flagName = c.Config.GetString("name")
		}
	}

	commands := []string{"shell", c.flagName}

	if c.flagDryRun {
		c.UI.Info(strings.Join(commands, " "))
		return 0
	}

	c.Runner.UseSyscall(true)

	if err := c.Runner.Run(commands); err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	return 0
}
