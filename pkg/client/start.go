package client

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Start is used to start all containers related to a specific environment
func (cli *Client) Start(ctx context.Context, name string, args []string) error {
	fmt.Println("Starting environment for", name)

	// get all the containers using a filter, we only want to start containers which
	// have the label com.craftcms.nitro.environment=name
	siteFilters := filters.NewArgs()
	siteFilters.Add("label", "com.craftcms.nitro.environment="+name)
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: siteFilters})
	if err != nil {
		return fmt.Errorf("unable to get a list of the containers, %w", err)
	}

	// syaty each environment container
	for _, container := range containers {
		// TODO maket this more dynamic
		fmt.Println("  ==> starting container for", container.Labels["com.craftcms.nitro.host"])

		if err := cli.docker.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
			return fmt.Errorf("unable to start container %s: %w", container.Names[0], err)
		}
	}

	fmt.Println("Development environment for", name, "started")

	return nil
}
