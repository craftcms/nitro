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
	// get all the containers using a filter, we only want to restart containers which
	// have the label com.craftcms.nitro.environment=name
	filter := filters.NewArgs()
	filter.Add("label", "com.craftcms.nitro.environment="+name)
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to get a list of the containers, %w", err)
	}

	// if there are no containers, were done
	if len(containers) == 0 {
		fmt.Println("There are no container for the", name, "environment")

		return nil
	}

	fmt.Println("Restarting containers for environment", name)

	// set a timeout, consider making this a flag
	timeout := time.Duration(5000) * time.Millisecond

	// restart each container for the environment
	for _, container := range containers {
		// TODO make this more dynamic
		fmt.Println("  ==> restarting container for", container.Labels["com.craftcms.nitro.host"])

		if err := cli.docker.ContainerRestart(ctx, container.ID, &timeout); err != nil {
			return fmt.Errorf("unable to restart container %s: %w", container.Names[0], err)
		}
	}

	fmt.Println("Development environment for", name, "restarted")

	return nil
}
