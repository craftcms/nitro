package command

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var CloudInit = `#cloud-config
packages:
  - redis
  - jq
  - apt-transport-https
  - ca-certificates
  - curl
  - gnupg-agent
  - software-properties-common
runcmd:
  - sudo add-apt-repository -y ppa:nginx/stable
  - sudo add-apt-repository -y ppa:ondrej/php
  - curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
  - sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
  - sudo apt-get update -y
  - sudo apt install -y nginx docker-ce docker-ce-cli containerd.io
  - sudo usermod -aG docker ubuntu
  - wget -q -O - https://packages.blackfire.io/gpg.key | sudo apt-key add -
  - echo "deb http://packages.blackfire.io/debian any main" | sudo tee /etc/apt/sources.list.d/blackfire.list
`

type InitCommand struct {
	*CoreCommand

	// flags
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
	s := flag.NewFlagSet("init", 128)

	s.StringVar(&c.flagName, "name", "", "name of the machine")
	s.IntVar(&c.flagCpus, "cpu", 0, "Number of CPUs to allocate to machine")
	s.StringVar(&c.flagMemory, "memory", "", "Amount of memory to allocate to machine")
	s.StringVar(&c.flagDisk, "disk", "", "Amount of disk space to allocate to machine")
	s.BoolVar(&c.flagSkipInstall, "skip-install", false, "Skip installing software on machine")
	s.BoolVar(&c.flagDryRun, "dry-run", false, "skip executing the command")

	return s
}

func (c *InitCommand) Run(args []string) int {
	if err := c.Flags().Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if c.Config.IsSet("name") {
		c.flagName = c.Config.GetString("name")
	}
	if c.Config.IsSet("cpus") {
		c.flagCpus = c.Config.GetInt("cpus")
	}
	if c.Config.IsSet("memory") {
		c.flagMemory = c.Config.GetString("memory")
	}
	if c.Config.IsSet("disk") {
		c.flagDisk = c.Config.GetString("disk")
	}

	// set defaults if the flag is not set
	if err := c.setDefaultOptions(); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	commands := []string{
		"launch",
		"--name",
		c.flagName,
		"--cpus",
		strconv.Itoa(c.flagCpus),
		"--mem",
		c.flagMemory,
		"--disk",
		c.flagDisk,
		"--cloud-init",
		"-",
	}

	c.UI.Info("Setting up machine, this may take a while...")
	c.UI.Info("----")
	c.UI.Info("name: " + c.flagName)
	c.UI.Info("cpus: " + strconv.Itoa(c.flagCpus))
	c.UI.Info("memory: " + c.flagMemory)
	c.UI.Info("disk: " + c.flagDisk)
	c.UI.Info("----")

	// pass the cloud init file to the machine
	if err := c.Runner.SetInput(CloudInit); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// we are just looking
	if c.flagDryRun {
		fmt.Println(commands)
		return 0
	}

	if err := c.Runner.Run(commands); err != nil {
		return 1
	}

	if c.flagSkipInstall {
		c.UI.Info("Skipping software installation")
		c.UI.Info("To install software on machine, run the following command:")
		c.UI.Info("nitro -name " + c.flagName + "install")
		return 0
	}

	c.UI.Info("Installing software on machine...")

	return 0
}

func (c *InitCommand) setDefaultOptions() error {
	if c.flagName == "" {
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

	return nil
}
