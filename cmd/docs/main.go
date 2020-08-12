package main

import (
	"log"

	"github.com/spf13/cobra/doc"

	"github.com/craftcms/nitro/internal/cmd"
)

func main() {
	nitro := cmd.New()
	if err := doc.GenMarkdownTree(nitro, "/tmp"); err != nil {
		log.Fatal(err)
	}
}
