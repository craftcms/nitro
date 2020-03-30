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
	s.StringVar(&c.flagName, "name", "", "name of the machine")
	s.BoolVar(&c.flagDryRun, "dry-run", false, "skip executing the command")
	return s
}
