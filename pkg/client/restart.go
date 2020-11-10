package client

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Restart is used to restart all containers related to a specific environment
func (cli *Client) Restart(ctx context.Context, name string, args []string) error {
	fmt.Println("Restarting all containers for environment", name)

	// get all the containers using a filter, we only want to restart containers which
	// have the label com.craftcms.nitro.environment=name
	siteFilters := filters.NewArgs()
	siteFilters.Add("label", "com.craftcms.nitro.environment="+name)
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: siteFilters})
	if err != nil {
		return fmt.Errorf("unable to get a list of the containers, %w", err)
	}

	// restart each container for the environment
	for _, container := range containers {
		// TODO maket this more dynamic
		fmt.Println("  ==> restarting container for", container.Labels["com.craftcms.nitro.host"])

		if err := cli.docker.ContainerRestart(ctx, container.ID, time.Duration(time.Second*5)); err != nil {
			return fmt.Errorf("unable to restart container %s: %w", container.Names[0], err)
		}
	}

	fmt.Println("Development environment for", name, "restarted")

	return nil
}
