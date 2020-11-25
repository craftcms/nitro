package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/craftcms/nitro/internal/database"
)

func main() {
	fileArg := flag.String("file", "../iraas-staging.dump", "The file to compress")
	flag.Parse()

	path := filepath.Clean(*fileArg)

	engine, err := database.DetermineEngine(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("path:", path)
	fmt.Println("engine:", engine)
}
