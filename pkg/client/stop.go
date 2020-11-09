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

	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		return err
	}

	for _, container := range containers {
		fmt.Println(container.Names)
		if err := cli.docker.ContainerStop(ctx, container.ID, nil); err != nil {
			return fmt.Errorf("unable to stop container, %w", err)
		}
	}

	return fmt.Errorf("not implemented")
}
