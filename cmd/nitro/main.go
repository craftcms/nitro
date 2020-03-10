package main

import (
	"log"
	"os"

	"github.com/craftcms/nitro/internal"
	"github.com/craftcms/nitro/internal/app"
)

func run(args []string) {
	if err := app.NewApp(internal.NewRunner("multipass")).Run(args); err != nil {
		log.Fatal(err)
	}
}

func main() {
	run(os.Args)
}
