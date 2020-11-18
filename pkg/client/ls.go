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
	filter.Add("label", EnvironmentLabel+"="+name)
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to get a list of the containers, %w", err)
	}

	// if there are no containers, were done
	if len(containers) == 0 {
		cli.Error("There are no containers running for the", name, "environment")

		return nil
	}

	cli.Info(fmt.Sprintf("Listing containers for %s...", name))

	// list each container for for the environment
	for _, c := range containers {
		containerType := "web"
		if c.Labels["com.craftcms.nitro.proxy"] != "" {
			containerType = "proxy"
		}

		if c.Labels[DatabaseEngineLabel] != "" {
			containerType = "database"
		}

		n := strings.TrimLeft(c.Names[0], "/")

		fmt.Println("  ==> type:", containerType, "\thostname:", n)
		fmt.Println("      aliases:", "\t\texamplealias.demo,", "anotheralias.test")
		fmt.Println("      ip:", c.NetworkSettings.Networks["nitro-dev"].IPAddress, "\timage:", c.Image)
		if c.Mounts[0].Source != "" {
			fmt.Println("      mount:", c.Mounts[0].Source)
		}
		fmt.Println("      uptime:", c.Status)
		fmt.Println("      ---")
	}

	return nil
}
