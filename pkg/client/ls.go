package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// LS is used to return a list of containers related to a specific environment
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
		cli.out.Error("There are no containers running for the", name, "environment")

		return nil
	}

	cli.out.Info("Listing containers for", name)

	// list each container for for the environment
	for _, c := range containers {
		var containerType string
		if c.Labels["com.craftcms.nitro.host"] != "" {
			containerType = "web"
		}
		if c.Labels["com.craftcms.nitro.proxy"] != "" {
			containerType = "proxy"
		}

		n := strings.TrimLeft(c.Names[0], "/")

		cli.out.Info("  ==> type:", containerType, "\thostname:", n)
		cli.out.Info("      ip:", c.NetworkSettings.Networks["nitro-dev"].IPAddress)
		cli.out.Info("      uptime:", c.Status)
	}

	return nil
}
