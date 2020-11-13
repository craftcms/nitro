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
	fmt.Println("Checking for containers")
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return fmt.Errorf("unable to list the containers, %w", err)
	}

	// make sure there are containers
	if len(containers) == 0 {
		fmt.Println("  ==> no containers found for environment", env)
	}

	fmt.Println("  ==> found", len(containers), "containers for environment", env)

	// get all related volumes
	fmt.Println("Checking for volumes")
	volumes, err := cli.docker.VolumeList(ctx, filter)
	if err != nil {
		fmt.Println(" ==> error listing volumes for the environment")
	}

	// make sure there are volumes
	if len(volumes.Volumes) == 0 {
		fmt.Println("  ==> no volumes found for the environment")
	}

	// get all related networks
	networks, err := cli.docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
	if err != nil {
		fmt.Println(" ==> error listing networks for the environment")
	}

	// make sure there are networks
	if len(networks) == 0 {
		fmt.Println("  ==> no networks found for the environment")
	}

	timeout := time.Duration(5000) * time.Millisecond

	// stop all of the container
	fmt.Println("Removing containers")
	for _, c := range containers {
		name := strings.TrimLeft(c.Names[0], "/")

		// only perform a backup if the container is for databases
		if c.Labels["com.craftcms.nitro.database"] != "" {
			fmt.Println(" ==> removing databases is not yet supported")
			break

			// TODO(jasonmccallister) implement backups
			fmt.Println("Backing up database")
			time.Sleep(time.Second * 2)
			fmt.Println("  ==> database backup for container", strings.TrimLeft(c.Names[0], "/"), "completed")
		}

		// stop the container
		fmt.Println("  ==> stopping container", name)
		if err := cli.docker.ContainerStop(ctx, c.ID, &timeout); err != nil {
			return fmt.Errorf("unable to stop the container, %w", err)
		}

		fmt.Println("  ==> container", name, "stopped")

		// remove the container
		fmt.Println("  ==> removing container", name)
		if err := cli.docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			return fmt.Errorf("unable to remove the container, %w", err)
		}

		fmt.Println("  ==> container", name, "removed")
	}

	// get all the volumes
	fmt.Println("Removing volumes")
	for _, v := range volumes.Volumes {
		fmt.Println("  ==> removing volume", v.Name)

		// remove the volume
		if err := cli.docker.VolumeRemove(ctx, v.Name, true); err != nil {
			fmt.Println("  ==> unable to remove volume"+v.Name+",", "you may need to manually remove the volume")
			break
		}

		fmt.Println("  ==> volume", v.Name, "removed")
	}

	// get all the networks
	fmt.Println("Removing network")
	for _, n := range networks {
		fmt.Println("  ==> removing network", n.Name)

		if err := cli.docker.NetworkRemove(ctx, n.ID); err != nil {
			fmt.Println("  ==> unable to remove network", n.Name, "you may need to manually remove the network")
		}

		fmt.Println("  ==> network", n.Name, "removed")
	}

	fmt.Println("Development environment for", env, "removed")

	return nil
}
