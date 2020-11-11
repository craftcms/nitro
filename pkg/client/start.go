package client

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Start is used to start all containers related to a specific environment
func (cli *Client) Start(ctx context.Context, name string, args []string) error {
	// get all the containers using a filter, we only want to start containers which
	// have the label com.craftcms.nitro.environment=name
	filter := filters.NewArgs()
	filter.Add("label", "com.craftcms.nitro.environment="+name)
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to get a list of the containers, %w", err)
	}

	// if there are no containers, were done
	if len(containers) == 0 {
		fmt.Println("There are no containers to start for the", name, "environment")

		return nil
	}

	fmt.Println("Starting environment for", name)

	// start each environment container
	for _, c := range containers {
		fmt.Println("  ==> starting container for", getContainerName(c))

		if err := cli.docker.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
			return fmt.Errorf("unable to start container %s: %w", c.Names[0], err)
		}
	}

	fmt.Println("Development environment for", name, "started")

	return nil
}

func getContainerName(c types.Container) string {
	if c.Labels["com.craftcms.nitro.host"] != "" {
		return "site"
	}

	if c.Labels["com.craftcms.nitro.proxy"] != "" {
		return "proxy"
	}

	return ""
}
