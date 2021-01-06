package main

import (
	"log"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/craftcms/nitro/internal/cmd"
)

func main() {
	nitro := cmd.New()

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err := doc.GenMarkdownTree(nitro, path+"/docs"); err != nil {
		log.Fatal(err)
	}
}
