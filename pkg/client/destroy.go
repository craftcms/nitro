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
	filter.Add("label", EnvironmentLabel+"="+env)

	// get all related containers
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return fmt.Errorf("unable to list the containers, %w", err)
	}

	// make sure there are containers
	if len(containers) == 0 {
		cli.SubError("no containers found for environment", env)
	}

	// get all related volumes
	volumes, err := cli.docker.VolumeList(ctx, filter)
	if err != nil {
		cli.SubError("error listing volumes for the environment")
	}

	// make sure there are volumes
	if len(volumes.Volumes) == 0 {
		cli.SubError("no volumes found for the environment")
	}

	// get all related networks
	networks, err := cli.docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
	if err != nil {
		cli.SubError("error listing networks for the environment")
	}

	// make sure there are networks
	if len(networks) == 0 {
		cli.SubError("no networks found for the environment")
	}

	// stop all of the container
	if len(containers) > 0 {
		timeout := time.Duration(5000) * time.Millisecond

		cli.Info("Removing Containers...")

		for _, c := range containers {
			name := strings.TrimLeft(c.Names[0], "/")

			// only perform a backup if the container is for databases
			if c.Labels["com.craftcms.nitro.todo"] != "" {
				cli.SubError("removing databases is not yet supported")
				break

				// TODO(jasonmccallister) implement backups
				fmt.Println("Backing up database")
				time.Sleep(time.Second * 2)
				cli.SubInfo("database backup for container", strings.TrimLeft(c.Names[0], "/"), "completed")
			}

			// stop the container
			cli.InfoPending("removing", name)
			if err := cli.docker.ContainerStop(ctx, c.ID, &timeout); err != nil {
				return fmt.Errorf("unable to stop the container, %w", err)
			}

			// remove the container
			if err := cli.docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				return fmt.Errorf("unable to remove the container, %w", err)
			}

			cli.InfoDone()
		}
	}

	// get all the volumes
	if len(volumes.Volumes) > 0 {
		cli.Info("Removing Volumes...")

		for _, v := range volumes.Volumes {
			cli.InfoPending("removing", v.Name)

			// remove the volume
			if err := cli.docker.VolumeRemove(ctx, v.Name, true); err != nil {
				cli.SubError("unable to remove volume"+v.Name+",", "you may need to manually remove the volume")
				break
			}

			cli.InfoDone()
		}
	}

	// get all the networks
	if len(networks) > 0 {
		cli.Info("Removing Networks...")

		for _, n := range networks {
			cli.InfoPending("removing", n.Name)

			if err := cli.docker.NetworkRemove(ctx, n.ID); err != nil {
				cli.SubInfo("unable to remove network", n.Name, "you may need to manually remove the network")
			}

			cli.InfoDone()
		}
	}

	cli.Info("Environment", env, "destroyed ✨")

	return nil
}
