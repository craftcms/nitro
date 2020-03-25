package command

import "github.com/mitchellh/cli"

type CoreCommand struct {
	UI     cli.Ui
	Runner ShellRunner
}
