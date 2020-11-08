package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/craftcms/nitro/pkg/client"
)

func main() {
	name := flag.String("machine", "nitro-dev", "the name of the machine")
	stop := flag.Bool("stop", false, "stop the containers")
	flag.Parse()

	ctx := context.Background()

	cli, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	if *stop {
		if err := cli.Stop(ctx, *name, os.Args); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	} else {
		if err := cli.Init(ctx, *name, os.Args); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}
}
