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
	for _, c := range containers {
		var containerType string
		if c.Labels["com.craftcms.nitro.host"] != "" {
			containerType = "web"
		}
		if c.Labels["com.craftcms.nitro.proxy"] != "" {
			containerType = "proxy"
		}

		fmt.Println(
			fmt.Sprintf("  ==> \t%q \ttype: %s \tmounts: %d \tid: %s", strings.TrimLeft(c.Names[0], "/"), containerType, len(c.Mounts), c.ID),
		)
	}

	return nil
}
