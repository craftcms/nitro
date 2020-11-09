package client

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (cli *Client) Stop(ctx context.Context, name string, args []string) error {
	filter := filters.Args{}
	filter.ExactMatch("label", name)

	fmt.Println("Starting shutdown for", name)

	// get all the containers using a filter, we only want to stop nitro related containers
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to get a list of the containers, %w", err)
	}

	// stop each container we found
	for _, container := range containers {
		if err := cli.docker.ContainerStop(ctx, container.ID, nil); err != nil {
			return fmt.Errorf("unable to stop container %s: %w", container.Names[0], err)
		}
	}

	fmt.Println("  ==> stopped all nitro related containers")

	return nil
}
