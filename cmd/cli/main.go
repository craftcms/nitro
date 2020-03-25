package main

import (
	"bufio"
	"log"
	"os"

	"github.com/mitchellh/cli"

	"github.com/craftcms/nitro/command"
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

	c := cli.NewCLI("nitro", "1.0.0")
	c.Args = os.Args[1:]

	baseCommand := &command.CoreCommand{
		UI:     &ui,
		Runner: r,
	}

	c.Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &command.InitCommand{CoreCommand: baseCommand}, nil
		},
	}

	status, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(status)
}
