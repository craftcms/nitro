package main

import (
	"flag"
	"log"

	"github.com/craftcms/nitro/command"
	"github.com/craftcms/nitro/internal/nitro"
)

func main() {
	ssh := flag.Bool("ssh", false, "ssh into machine")
	flag.Parse()
	runner := command.NewMultipassRunner("multipass")

	var commands []nitro.Command
	if *ssh {
		commands = append(commands, nitro.SSH("somename")...)
	} else {
		commands = nitro.Init("somename", "4", "4G", "20G", "7.4", "mysql", "5.7")
	}

	if err := nitro.Run(runner, commands); err != nil {
		log.Fatal(err)
	}
}
