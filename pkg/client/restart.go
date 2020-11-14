package client

import (
	"context"
	"fmt"
	"strings"
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
		cli.out.Error("There are no containers for the", name, "environment")

		return nil
	}

	cli.out.Info("Restarting containers for environment", name)

	// set a timeout, consider making this a flag
	timeout := time.Duration(5000) * time.Millisecond

	// restart each container for the environment
	for _, c := range containers {
		n := strings.TrimLeft(c.Names[0], "/")

		cli.out.Info("  ==> restarting", n)

		if err := cli.docker.ContainerRestart(ctx, c.ID, &timeout); err != nil {
			return fmt.Errorf("unable to restart container %s: %w", n, err)
		}
	}

	cli.out.Info("Development environment for", name, "restarted")

	return nil
}
