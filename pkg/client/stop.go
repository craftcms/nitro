package client

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Stop is used to stop all containers related to a specific environment
func (cli *Client) Stop(ctx context.Context, name string, args []string) error {
	fmt.Println("Stopping down environment for", name)

	// get all the containers using a filter, we only want to stop containers which
	// have the label com.craftcms.nitro.machine=name
	siteFilters := filters.NewArgs()
	siteFilters.Add("label", "com.craftcms.nitro.environment="+name)
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: siteFilters})
	if err != nil {
		return fmt.Errorf("unable to get a list of the containers, %w", err)
	}

	// stop each environment container
	for _, container := range containers {
		// TODO maket this more dynamic
		fmt.Println("  ==> stopping container for", container.Labels["com.craftcms.nitro.host"])

		if err := cli.docker.ContainerStop(ctx, container.ID, nil); err != nil {
			return fmt.Errorf("unable to stop container %s: %w", container.Names[0], err)
		}
	}

	fmt.Println("Development environment for", name, "shutdown")

	return nil
}
