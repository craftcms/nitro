package databasecontainer

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var (
	// DatabaseImage is used for determining the engine and version
	DatabaseImage = "docker.io/library/%s:%s"
)

// StartOrCreate is used to find a specific database and start the container. If there is no container for the database,
// it will create a new volume and container for the database.
func StartOrCreate(ctx context.Context, docker client.CommonAPIClient, networkID string, db config.Database) (string, string, error) {
	// create the filters for the database
	filter := filters.NewArgs()
	filter.Add("label", labels.Nitro)
	filter.Add("label", labels.DatabaseEngine+"="+db.Engine)
	filter.Add("label", labels.DatabaseVersion+"="+db.Version)
	filter.Add("label", labels.Type+"=database")

	hostname, err := db.GetHostname()
	if err != nil {
		return "", "", err
	}

	// set the container database compatibility
	if db.Engine == "mariadb" || db.Engine == "mysql" {
		filter.Add("label", labels.DatabaseCompatibility+"=mysql")
	} else {
		filter.Add("label", labels.DatabaseCompatibility+"=postgres")
	}

	// get the containers for the database
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return "", "", fmt.Errorf("error getting a list of containers")
	}

	// if there is a container, we should start it and return
	if len(containers) == 1 {
		// check if the container is running
		if containers[0].State != "running" {
			// start the container
			if err := docker.ContainerStart(ctx, containers[0].ID, types.ContainerStartOptions{}); err != nil {
				return "", "", err
			}
		}

		return containers[0].ID, hostname, nil
	}

	// create the database labels for the new container
	lbls := map[string]string{
		labels.Nitro:           "true",
		labels.DatabaseEngine:  db.Engine,
		labels.DatabaseVersion: db.Version,
		labels.Type:            "database",
	}

	// if the database is mysql or mariadb, mark them as
	// mysql compatible (used for importing backups)
	if db.Engine == "mariadb" || db.Engine == "mysql" {
		lbls[labels.DatabaseCompatibility] = "mysql"
	}

	// if the database is postgres, mark it as compatible
	// with postgres. This is not needed but a place holder
	// if cockroachdb is ever supported by craft.
	if db.Engine == "postgres" {
		lbls[labels.DatabaseCompatibility] = "postgres"
	}

	// create the volume
	volume, err := docker.VolumeCreate(ctx, volumetypes.VolumeCreateBody{Driver: "local", Name: hostname, Labels: lbls})
	if err != nil {
		return "", "", fmt.Errorf("unable to create the volume, %w", err)
	}

	// determine the image name
	image := fmt.Sprintf(DatabaseImage, db.Engine, db.Version)

	// set mounts and environment based on the database type
	target := "/var/lib/mysql"
	var envs []string
	if strings.Contains(image, "postgres") {
		target = "/var/lib/postgresql/data"
		envs = []string{"POSTGRES_USER=nitro", "POSTGRES_DB=nitro", "POSTGRES_PASSWORD=nitro"}
	} else {
		envs = []string{"MYSQL_ROOT_PASSWORD=nitro", "MYSQL_DATABASE=nitro", "MYSQL_USER=nitro", "MYSQL_PASSWORD=nitro"}
	}

	// pull the image
	rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
	if err != nil {
		return "", "", fmt.Errorf("unable to pull image %s, %w", image, err)
	}

	// read the output to pull the image
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(rdr); err != nil {
		return "", "", fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
	}

	// set the port for the database
	port, err := nat.NewPort("tcp", db.Port)
	if err != nil {
		return "", "", fmt.Errorf("unable to create the port, %w", err)
	}

	containerConfig := &container.Config{
		Image:  image,
		Labels: lbls,
		ExposedPorts: nat.PortSet{
			port: struct{}{},
		},
		Env: envs,
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volume.Name,
				Target: target,
			},
		},
		PortBindings: map[nat.Port][]nat.PortBinding{
			port: {
				{
					HostIP:   "127.0.0.1",
					HostPort: db.Port,
				},
			},
		},
	}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"nitro-network": {
				NetworkID: networkID,
			},
		},
	}

	// create the container for the database
	resp, err := docker.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, hostname)
	if err != nil {
		return "", "", fmt.Errorf("unable to create the container, %w", err)
	}

	// start the container
	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", "", fmt.Errorf("unable to start the container, %w", err)
	}

	// if the container is mysql compatible
	if db.Engine == "mysql" || db.Engine == "mariadb" {
		cmds := []string{"mysql", "-uroot", "-pnitro", fmt.Sprintf(`-e GRANT ALL PRIVILEGES ON *.* TO '%s'@'%s' WITH GRANT OPTION;`, "nitro", "%")}

		// create the exec
		exec, err := docker.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
			AttachStdout: true,
			AttachStderr: true,
			Tty:          false,
			Cmd:          cmds,
		})
		if err != nil {
			return "", "", err
		}

		// attach to the container
		resp, err := docker.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{
			Tty: false,
		})
		if err != nil {
			return "", "", err
		}
		defer resp.Close()

		// start the exec
		if err := docker.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{}); err != nil {
			return "", "", fmt.Errorf("unable to start the container, %w", err)
		}

		// wait for the container exec to complete
		for {
			resp, err := docker.ContainerExecInspect(ctx, exec.ID)
			if err != nil {
				return "", "", err
			}

			if !resp.Running {
				break
			}
		}
	}

	return resp.ID, hostname, nil
}
