package main

import (
	"fmt"
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

	dir := filepath.Join(path, "docs")
	// file the file
	// writer := bufio.NewWriter()
	// doc.GenMarkdownCustom(cli, )

	if err := doc.GenMarkdownTree(cli, dir); err != nil {
		log.Fatal(err)
	}

	fmt.Println("docs output to", dir)
}
