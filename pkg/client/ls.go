package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Start is used to start all containers related to a specific environment
func (cli *Client) LS(ctx context.Context, name string, args []string) error {
	// get all the containers using a filter, we only want to start containers which
	// have the label com.craftcms.nitro.environment=name
	filter := filters.NewArgs()
	filter.Add("label", "com.craftcms.nitro.environment="+name)
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to get a list of the containers, %w", err)
	}

	// if there are no containers, were done
	if len(containers) == 0 {
		fmt.Println("There are no container running for the", name, "environment")

		return nil
	}

	fmt.Println("Listing containers for", name)

	// list each container for for the environment
	for _, container := range containers {
		fmt.Println(
			fmt.Sprintf("  ==> container: %q \ttype: web\t \tid: %s", strings.TrimLeft(container.Names[0], "/"), container.ID),
		)
	}

	return nil
}
