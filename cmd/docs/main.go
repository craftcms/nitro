package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/craftcms/nitro/command/nitro"
	"github.com/spf13/cobra/doc"
)

func main() {
	cli := nitro.NewCommand()

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err := doc.GenMarkdownTree(cli, filepath.Join(path, "docs")); err != nil {
		log.Fatal(err)
	}
}
