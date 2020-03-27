package command

import (
	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
)

type CoreCommand struct {
	UI     cli.Ui
	Runner ShellRunner
	Config *viper.Viper
}
