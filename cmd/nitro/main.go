package main

import (
	"os"

	"github.com/craftcms/nitro/command/nitro"
)

func main() {
	// execute the nitro root command
	if err := nitro.NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
