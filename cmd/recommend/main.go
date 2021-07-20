package main

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/client"
)

func main() {
	// create the docker client
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	info, err := docker.Info(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(info.MemTotal / 1000000000)
}
