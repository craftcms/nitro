package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (cli *Client) Stop(ctx context.Context, name string, args []string) error {
	fmt.Println("Shutting down development environment for", name)

	// get all the containers using a filter, we only want to stop nitro related containers
	// get all of the sites
	siteFilters := filters.NewArgs()
	siteFilters.Add("label", "com.craftcms.nitro.machine="+name)
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: siteFilters})
	if err != nil {
		return fmt.Errorf("unable to get a list of the containers, %w", err)
	}

	// stop each site container
	for _, container := range containers {
		fmt.Println("  ==> stopping container for", container.Labels["com.craftcms.nitro.host"])

		if err := cli.docker.ContainerStop(ctx, container.ID, nil); err != nil {
			return fmt.Errorf("unable to stop container %s: %w", container.Names[0], err)
		}
	}

	// get all the proxy container using a filter
	proxyFilter := filters.NewArgs()
	proxyFilter.Add("label", "com.craftcms.nitro.proxy="+name)
	proxyContainers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: proxyFilter})
	if err != nil {
		return fmt.Errorf("unable to find the proxy container, %w", err)
	}

	// stop each site container
	for _, container := range proxyContainers {
		fmt.Println("  ==> stopping proxy", strings.TrimLeft(container.Names[0], "/"))

		if err := cli.docker.ContainerStop(ctx, container.ID, nil); err != nil {
			return fmt.Errorf("unable to stop container %s: %w", container.Names[0], err)
		}
	}

	fmt.Println("Development environment for", name, "shutdown")

	return nil
}
