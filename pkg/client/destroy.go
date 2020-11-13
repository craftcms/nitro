package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Destroy is used to completed remove all containers, networks, and volumes for an environment.
// is it a destructive action and will prompt a user for verification and perform a database
// backup before removing the resources.
func (cli *Client) Destroy(ctx context.Context, env string, args []string) error {
	filter := filters.NewArgs()
	filter.Add("label", "com.craftcms.nitro.environment="+env)

	// get all related containers
	cli.out.Info("Checking for containers")
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return fmt.Errorf("unable to list the containers, %w", err)
	}

	// make sure there are containers
	if len(containers) == 0 {
		cli.out.Error("  ==> no containers found for environment", env)
	} else {
		cli.out.Info("  ==> found", len(containers), "containers for environment", env)
	}

	// get all related volumes
	cli.out.Info("Checking for volumes")
	volumes, err := cli.docker.VolumeList(ctx, filter)
	if err != nil {
		cli.out.Error(" ==> error listing volumes for the environment")
	}

	// make sure there are volumes
	if len(volumes.Volumes) == 0 {
		cli.out.Error("  ==> no volumes found for the environment")
	} else {
		cli.out.Info("  ==> found", len(volumes.Volumes), "volumes for environment", env)
	}

	// get all related networks
	networks, err := cli.docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
	if err != nil {
		cli.out.Error(" ==> error listing networks for the environment")
	}

	// make sure there are networks
	if len(networks) == 0 {
		cli.out.Error("  ==> no networks found for the environment")
	}

	timeout := time.Duration(5000) * time.Millisecond

	// stop all of the container
	cli.out.Info("Removing containers")
	for _, c := range containers {
		name := strings.TrimLeft(c.Names[0], "/")

		// only perform a backup if the container is for databases
		if c.Labels["com.craftcms.nitro.database"] != "" {
			cli.out.Error(" ==> removing databases is not yet supported")
			break

			// TODO(jasonmccallister) implement backups
			cli.out.Info("Backing up database")
			time.Sleep(time.Second * 2)
			cli.out.Info("  ==> database backup for container", strings.TrimLeft(c.Names[0], "/"), "completed")
		}

		// stop the container
		cli.out.Info("  ==> stopping container", name)
		if err := cli.docker.ContainerStop(ctx, c.ID, &timeout); err != nil {
			return fmt.Errorf("unable to stop the container, %w", err)
		}

		cli.out.Info("  ==> container", name, "stopped")

		// remove the container
		cli.out.Info("  ==> removing container", name)
		if err := cli.docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			return fmt.Errorf("unable to remove the container, %w", err)
		}

		cli.out.Info("  ==> container", name, "removed")
	}

	// get all the volumes
	cli.out.Info("Removing volumes")
	for _, v := range volumes.Volumes {
		cli.out.Info("  ==> removing volume", v.Name)

		// remove the volume
		if err := cli.docker.VolumeRemove(ctx, v.Name, true); err != nil {
			cli.out.Error("  ==> unable to remove volume"+v.Name+",", "you may need to manually remove the volume")
			break
		}

		cli.out.Info("  ==> volume", v.Name, "removed")
	}

	// get all the networks
	cli.out.Info("Removing network")
	for _, n := range networks {
		cli.out.Info("  ==> removing network", n.Name)

		if err := cli.docker.NetworkRemove(ctx, n.ID); err != nil {
			cli.out.Error("  ==> unable to remove network", n.Name, "you may need to manually remove the network")
		}

		cli.out.Info("  ==> network", n.Name, "removed")
	}

	cli.out.Info("Development environment for", env, "removed")

	return nil
}
