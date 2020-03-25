package command

import (
	"flag"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
)

type InitCommand struct {
	UI     cli.Ui
	runner ShellRunner

	// flags
	flagName        string
	flagCpus        int
	flagMemory      string
	flagDisk        string
	flagSkipInstall bool
}

func (c *InitCommand) Synopsis() string {
	return "create new machine"
}

func (c *InitCommand) Help() string {
	helpText := `
Usage: nitro init [options]
  This command starts a nitro virtual machine and will provision the machine with the requested specifications.
  
  Create a new virtual machine and override the default system specifications:
      $ nitro init -name=diesel -cpu=4 -memory=4G -disk=40GB
  
  Create a new virtual machine and with the defaults and skip bootstrapping the machine with the default installations:
      $ nitro init -name=diesel -skip-install
`
	return strings.TrimSpace(helpText)
}

func (c *InitCommand) Flags() *flag.FlagSet {
	set := flag.NewFlagSet("init", 128)

	set.StringVar(&c.flagName, "name", "", "name of the machine")
	set.IntVar(&c.flagCpus, "cpu", 0, "Number of CPUs to allocate to machine")
	set.StringVar(&c.flagMemory, "memory", "", "Amount of memory to allocate to machine")
	set.StringVar(&c.flagDisk, "disk", "", "Amount of disk space to allocate to machine")
	set.BoolVar(&c.flagSkipInstall, "skip-install", false, "Skip installing software on machine")

	return set
}

func (c *InitCommand) Run(args []string) int {
	if err := c.Flags().Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// check each flag and try to get the values from the config file
	if c.flagName == "" {
		// TODO check if the config file has the option
		c.flagName = "nitro-dev"
	}

	if c.flagCpus == 0 {
		c.flagCpus = 2
	}

	if c.flagMemory == "" {
		c.flagMemory = "2G"
	}

	if c.flagDisk == "" {
		c.flagDisk = "20G"
	}

	mpArgs := []string{
		"multipass",
		"launch",
		"--name",
		c.flagName,
		"--cpus",
		strconv.Itoa(c.flagCpus),
		"--memory",
		c.flagMemory,
		"--disk",
		c.flagDisk,
	}
	if err := c.runner.Run(mpArgs); err != nil {
		return 1
	}

	// otherwise prompt the user for the questions
	return 0
}
