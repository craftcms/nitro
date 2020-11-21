package client

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/go-connections/nat"
)

// Apply is used to create a
func (cli *Client) Apply(ctx context.Context, env string, cfg *config.Config) error {
	// get the network for the environment
	var networkID string

	// create a filter for the network
	filter := filters.NewArgs()
	filter.Add("label", EnvironmentLabel+"="+env)

	cli.Info(fmt.Sprintf("Checking %s Network...", env))

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

	cli.InfoSuccess("using", networkID)

	cli.Info("Checking Databases...")

	for _, db := range cfg.Databases {
		// add filters to check for the container
		filter.Add("label", DatabaseEngineLabel+"="+db.Engine)
		filter.Add("label", DatabaseVersionLabel+"="+db.Version)

		// get the containers for databases
		containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
		if err != nil {
			return fmt.Errorf("error getting a list of containers")
		}

		// set the hostname
		hostname, err := db.GetHostname()
		if err != nil {
			return err
		}

		// if there are no containers, create a volume, container, and start the container
		var containerID string
		var startContainer bool
		switch len(containers) {
		case 1:
			cli.InfoSuccess(hostname, "ready")

			// set the container id
			containerID = containers[0].ID

			// check if the container is running
			if containers[0].State != "running" {
				startContainer = true
			}
		default:
			cli.InfoPending("creating volume", hostname)

			// create the labels
			labels := map[string]string{
				EnvironmentLabel:     env,
				DatabaseEngineLabel:  db.Engine,
				DatabaseVersionLabel: db.Version,
			}

			// create the volume
			volResp, err := cli.docker.VolumeCreate(ctx, volumetypes.VolumesCreateBody{
				Driver: "local",
				Name:   hostname,
				Labels: labels,
			})
			if err != nil {
				return fmt.Errorf("unable to create the volume, %w", err)
			}

			cli.InfoDone()

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
			cli.InfoPending("pulling", image)

			rdr, err := cli.docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
			if err != nil {
				return fmt.Errorf("unable to pull image %s, %w", image, err)
			}

			cli.InfoDone()

			buf := &bytes.Buffer{}
			if _, err := buf.ReadFrom(rdr); err != nil {
				return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
			}

			port, err := nat.NewPort("tcp", db.Port)
			if err != nil {
				return fmt.Errorf("unable to create the port, %w", err)
			}

			// create the container
			cli.InfoPending("creating", hostname)

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
				hostname,
			)
			if err != nil {
				return fmt.Errorf("unable to create the container, %w", err)
			}

			// set the container id to start
			containerID = conResp.ID
			startContainer = true

			cli.InfoDone()
		}

		// start the container if needed
		if startContainer {
			cli.InfoPending("starting", hostname)

			if err := cli.docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}

			cli.InfoDone()
		}

		// remove the filter
		filter.Del("label", DatabaseEngineLabel+"="+db.Engine)
		filter.Del("label", DatabaseVersionLabel+"="+db.Version)
	}

	// get all of the sites, their local path, the php version, and the type of project (nginx or PHP-FPM)
	cli.Info("Checking Sites...")

	for _, site := range cfg.Sites {
		// add the site filter
		filter.Add("label", HostLabel+"="+site.Hostname)

		containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
		if err != nil {
			return fmt.Errorf("error getting a list of containers")
		}

		var containerID string
		var startContainer bool
		switch len(containers) {
		case 1:
			c := containers[0]
			image := fmt.Sprintf("docker.io/craftcms/nginx:%s", site.PHP)

			// make sure the images match, if they don't stop, remove, and create the container
			// with the new image
			if c.Image != image {
				cli.InfoPending(site.Hostname, "out of sync, applying")

				path, err := site.GetAbsPath()
				if err != nil {
					return err
				}

				// stop container
				if err := cli.docker.ContainerStop(ctx, c.ID, nil); err != nil {
					return err
				}

				// remove container
				if err := cli.docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{}); err != nil {
					return err
				}

				// pull the image
				// pull the image
				cli.InfoPending("pulling", image)

				rdr, err := cli.docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
				if err != nil {
					return fmt.Errorf("unable to pull image, %w", err)
				}

				buf := &bytes.Buffer{}
				if _, err := buf.ReadFrom(rdr); err != nil {
					return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
				}

				cli.InfoDone()

				// create new container, will have a new container id
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
							Type:   mount.TypeBind,
							Source: path,
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
				startContainer = true

				cli.InfoDone()

				break
			}

			cli.InfoSuccess(site.Hostname, "ready")

			// get the container id
			containerID = c.ID

			// check if the container is running
			if containers[0].State != "running" {
				startContainer = true
			}
		default:
			image := fmt.Sprintf("docker.io/craftcms/nginx:%s", site.PHP)

			path, err := site.GetAbsPath()
			if err != nil {
				return err
			}

			// pull the image
			cli.InfoPending("pulling", image)

			rdr, err := cli.docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
			if err != nil {
				return fmt.Errorf("unable to pull the image, %w", err)
			}

			buf := &bytes.Buffer{}
			if _, err := buf.ReadFrom(rdr); err != nil {
				return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
			}

			cli.InfoDone()

			cli.InfoPending("creating", site.Hostname)

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
						Type:   mount.TypeBind,
						Source: path,
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
			startContainer = true

			cli.InfoDone()
		}

		// start the container if needed
		if startContainer {
			cli.InfoPending("starting", site.Hostname)

			if err := cli.docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}

			cli.InfoDone()
		}

		// remove the site filter
		filter.Del("label", HostLabel+"="+site.Hostname)
	}

	// TODO(jasonmccallister) convert the sites into a Caddy json config and send to the API

	cli.Info("Everything for", env, "is up and running 😃")

	return nil
}
