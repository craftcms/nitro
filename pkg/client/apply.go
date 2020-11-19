package client

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/internal/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/go-connections/nat"
	"github.com/mitchellh/go-homedir"
)

// Apply is used to create a
func (cli *Client) Apply(ctx context.Context, env string, cfg config.Config) error {
	// get the network for the environment
	var networkID string

	// create a filter for the network
	filter := filters.NewArgs()
	filter.Add("label", EnvironmentLabel+"="+env)

	fmt.Println(fmt.Sprintf("Looking for %s network", env))

	// find networks
	networks, err := cli.docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to list docker networks\n%w", err)
	}

	for _, n := range networks {
		if n.Name == env {
			networkID = n.ID
		}
	}

	// if the network is not found
	if networkID == "" {
		return ErrNoNetwork
	}

	cli.SubInfo("using network", networkID)

	// get the users home dir
	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("unable to get the users home directory, %w", err)
	}

	cli.Info("Checking for databases")
	for _, db := range cfg.Databases {
		// add filters to check for the container
		filter.Add("label", DatabaseEngineLabel+"="+db.Engine)
		filter.Add("label", DatabaseVersionLabel+"="+db.Version)

		containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
		if err != nil {
			return fmt.Errorf("error getting a list of containers")
		}

		// if there are no containers, create a volume, container, and start the container
		var containerID string
		var startContainer bool
		switch len(containers) {
		case 1:
			cli.SubInfo("using existing container for", db.Name())

			// set the container id
			containerID = containers[0].ID

			// check if the container is running
			if containers[0].State != "running" {
				startContainer = true
			}
		default:
			cli.SubInfo("creating volume for", db.Name())

			// create the labels
			labels := map[string]string{
				EnvironmentLabel:     env,
				DatabaseEngineLabel:  db.Engine,
				DatabaseVersionLabel: db.Version,
			}

			// create the volume
			volResp, err := cli.docker.VolumeCreate(ctx, volumetypes.VolumesCreateBody{
				Driver: "local",
				Name:   db.Name(),
				Labels: labels,
			})
			if err != nil {
				return fmt.Errorf("unable to create the volume, %w", err)
			}

			// determine the image name
			image := fmt.Sprintf("docker.io/library/%s:%s", db.Engine, db.Version)

			target := "/var/lib/mysql"
			var envs []string
			if strings.Contains(image, "postgres") {
				target = "/var/lib/postgresql/data"
				envs = []string{"POSTGRES_USER=nitro", "POSTGRES_DB=nitro", "POSTGRES_PASSWORD=nitro"}
			} else {
				envs = []string{"MYSQL_ROOT_PASSWORD=nitro", "MYSQL_DATABASE=nitro", "MYSQL_USER=nitro", "MYSQL_PASSWORD=nitro"}
			}

			// pull the image
			cli.SubInfo("pulling image", image)
			rdr, err := cli.docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
			if err != nil {
				return fmt.Errorf("unable to pull image %s, %w", image, err)
			}

			buf := &bytes.Buffer{}
			if _, err := buf.ReadFrom(rdr); err != nil {
				return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
			}

			port, err := nat.NewPort("tcp", db.Port)
			if err != nil {
				return fmt.Errorf("unable to create the port, %w", err)
			}

			// create the container
			cli.SubInfo("creating container for", db.Name())
			conResp, err := cli.docker.ContainerCreate(
				ctx,
				&container.Config{
					Image:  image,
					Labels: labels,
					ExposedPorts: nat.PortSet{
						port: struct{}{},
					},
					Env: envs,
				},
				&container.HostConfig{
					Mounts: []mount.Mount{
						{
							Type:   mount.TypeVolume,
							Source: volResp.Name,
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
				},
				&network.NetworkingConfig{
					EndpointsConfig: map[string]*network.EndpointSettings{
						env: {
							NetworkID: networkID,
						},
					},
				},
				db.Name(),
			)
			if err != nil {
				return fmt.Errorf("unable to create the container, %w", err)
			}

			containerID = conResp.ID
			startContainer = true
			// set the container id to start
		}

		// start the container if needed
		if startContainer {
			cli.Info("starting container for", db.Name())

			if err := cli.docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}

			cli.SubInfo("container for", db.Name(), "started")
		}

		// remove the filter
		filter.Del("label", DatabaseEngineLabel+"="+db.Engine)
		filter.Del("label", DatabaseVersionLabel+"="+db.Version)
	}

	// TODO(jasonmccallister) get all of the sites, their local path, the php version, and the type of project (nginx or PHP-FPM)
	cli.Info("Checking for existing sites")

	for _, site := range cfg.Sites {
		// add the site filter
		filter.Add("label", HostLabel+"="+site.Hostname)

		// TODO(jasonmccallister) make the php version dynamic based on the site
		image := fmt.Sprintf("docker.io/craftcms/php-fpm:%s-dev", "7.4")

		containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
		if err != nil {
			return fmt.Errorf("error getting a list of containers")
		}

		var containerID string
		var startContainer bool
		switch len(containers) {
		case 1:
			cli.SubInfo("using existing container for", site.Hostname)

			// get the container id
			containerID = containers[0].ID

			// check if the container is running
			if containers[0].State != "running" {
				startContainer = true
			}
		default:
			// TODO(jasonmccallister) make this dynamic
			sourcePath := "~/dev/plugins-dev"
			if site.Hostname == "extendingcaddy.nitro" {
				sourcePath = "~/dev/extendingcaddy"
			}

			// TODO get the complete file path
			if strings.Contains(sourcePath, "~") {
				sourcePath = strings.Replace(sourcePath, "~", home, 1)
			}

			absPath, err := filepath.Abs(sourcePath)
			if err != nil {
				return fmt.Errorf("unable to get the absolute path to the site, %w", err)
			}

			// pull the image
			if _, err := cli.docker.ImagePull(ctx, image, types.ImagePullOptions{All: false}); err != nil {
				return fmt.Errorf("unable to pull the image, %w", err)
			}

			// create the container
			resp, err := cli.docker.ContainerCreate(
				ctx,
				&container.Config{
					Image: image,
					Labels: map[string]string{
						EnvironmentLabel: env,
						HostLabel:        site.Hostname,
					},
				},
				&container.HostConfig{
					Mounts: []mount.Mount{{
						Type: mount.TypeBind,
						// TODO (jasonmccallister) get the source from the site
						Source: absPath,
						//Source: site.Webroot,
						Target: "/app",
					},
					},
				},
				&network.NetworkingConfig{
					EndpointsConfig: map[string]*network.EndpointSettings{
						env: {
							NetworkID: networkID,
						},
					},
				},
				site.Hostname,
			)
			if err != nil {
				return fmt.Errorf("unable to create the container, %w", err)
			}

			containerID = resp.ID

			cli.SubInfo("created container for", site.Hostname)
		}

		// start the container if needed
		if startContainer {
			if err := cli.docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}
		}

		// remove the site filter
		filter.Del("label", HostLabel+"="+site.Hostname)
	}

	// TODO(jasonmccallister) convert the sites into a Caddy json config and send to the API

	cli.Info("All containers are running")

	return nil
}
