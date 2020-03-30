package command

import (
	"flag"

	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
)

type CoreCommand struct {
	UI     cli.Ui
	Runner ShellRunner
	Config *viper.Viper

	// global flags
	flagName   string
	flagDryRun bool
}

func (c *CoreCommand) Flags() *flag.FlagSet {
	s := flag.NewFlagSet("init", 128)
	// define the global flags
	s.StringVar(&c.flagName, "name", "", "name of the machine")
	s.BoolVar(&c.flagDryRun, "dry-run", false, "skip executing the command")

	// set the defaults from a config file
	// TODO this check is here temporarily for a test (ssh command run)
	if c.Config != nil {
		if c.Config.IsSet("name") && c.flagName == "" {
			c.flagName = c.Config.GetString("name")
		}
	}

	return s
}
