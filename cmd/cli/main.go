package main

import (
	"bufio"
	"log"
	"os"

	"github.com/mitchellh/cli"

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

	r := command.NewMultipassRunner("multipass")

	c := cli.NewCLI("nitro", Version)
	c.Args = os.Args[1:]

	coreCommand := &command.CoreCommand{
		UI:     &ui,
		Runner: r,
	}

	c.Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &command.InitCommand{CoreCommand: coreCommand}, nil
		},
		"install": func() (cli.Command, error) {
			return &command.InstallCommand{CoreCommand: coreCommand}, nil
		},
	}

	status, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(status)
}
