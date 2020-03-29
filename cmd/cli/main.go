package main

import (
	"bufio"
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/command"
)

var (
	Version = "1.0.0"
)

func main() {
	ui := cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Reader:      bufio.NewReader(os.Stdin),
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	v := viper.New()
	v.SetConfigName("nitro")
	v.SetConfigType("yaml")
	v.AddConfigPath("$HOME/.nitro")
	if err := v.ReadInConfig(); err != nil {
		log.Println(err)
	}

	r := command.NewMultipassRunner("multipass")

	c := cli.NewCLI("nitro", Version)
	c.Args = os.Args[1:]

	cmd := coreCommand(ui, r, v)

	c.Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &command.InitCommand{CoreCommand: cmd}, nil
		},
		"install": func() (cli.Command, error) {
			return &command.InstallCommand{CoreCommand: cmd}, nil
		},
	}

	status, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(status)
}

func coreCommand(ui cli.ColoredUi, r command.ShellRunner, v *viper.Viper) *command.CoreCommand {
	return &command.CoreCommand{
		UI:     &ui,
		Runner: r,
		Config: v,
	}
}
