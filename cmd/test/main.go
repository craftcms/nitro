package main

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	// create the docker client
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	results, err := docker.ImageSearch(ctx, "mariadb", types.ImageSearchOptions{Limit: 10})
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range results {
		fmt.Println(r.Name)
	}
}
