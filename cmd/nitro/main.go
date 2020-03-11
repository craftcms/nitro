package main

import (
	"log"
	"os"

	"github.com/craftcms/nitro/internal/app"
	"github.com/craftcms/nitro/internal/command"
)

func run(args []string) {
	if err := app.NewApp(command.NewRunner("multipass")).Run(args); err != nil {
		log.Fatal(err)
	}
}

func main() {
	run(os.Args)
}
