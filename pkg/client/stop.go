package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Stop is used to stop all containers related to a specific environment
func (cli *Client) Stop(ctx context.Context, name string, args []string) error {
	// get all the containers using a filter, we only want to stop containers which
	// have the label com.craftcms.nitro.machine=name
	filter := filters.NewArgs()
	filter.Add("label", EnvironmentLabel+"="+name)
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to get a list of the containers, %w", err)
	}

	// if there are no containers, were done
	if len(containers) == 0 {
		cli.Error("No containers are running in the", name, "environment")

		return nil
	}

	cli.Info("Stopping", name)

	// stop each environment container
	for _, c := range containers {
		n := strings.TrimLeft(c.Names[0], "/")

		cli.InfoPending("stopping", n)

		if err := cli.docker.ContainerStop(ctx, c.ID, nil); err != nil {
			errWrap := fmt.Errorf("unable to stop container %s: %w", n, err)

			cli.SubError(errWrap.Error())

			return errWrap
		}

		cli.InfoDone()
	}

	cli.Info(name, "shutdown 😴")

	return nil
}
