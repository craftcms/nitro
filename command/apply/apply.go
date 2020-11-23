package apply

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
)

var (
	// ErrNoNetwork is used when we cannot find the network
	ErrNoNetwork = fmt.Errorf("Unable to find the network")
)

const exampleText = `  # apply changes from a config to the environment
  nitro apply`

// New takes a docker client and the terminal output to run the apply actions
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply changes to an environment",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			ctx := cmd.Context()

			_, cfg, err := config.Load()
			if err != nil {
				return err
			}

			// create a filter for the network
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)

			output.Info(fmt.Sprintf("Checking %s Network...", env))

			// find networks
			networks, err := docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
			if err != nil {
				return fmt.Errorf("unable to list docker networks\n%w", err)
			}

			// get the network for the environment
			var networkID string
			for _, n := range networks {
				if n.Name == env {
					networkID = n.ID
				}
			}

			// if the network is not found
			if networkID == "" {
				return ErrNoNetwork
			}

			output.Success("using", networkID)

			// check the databases
			output.Info("Checking Databases...")
			for _, db := range cfg.Databases {
				// add filters to check for the container
				filter.Add("label", labels.DatabaseEngine+"="+db.Engine)
				filter.Add("label", labels.DatabaseVersion+"="+db.Version)

				// get the containers for databases
				containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
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
					// set the container id
					containerID = containers[0].ID

					// check if the container is running
					if containers[0].State != "running" {
						startContainer = true
					} else {
						output.Success(hostname, "ready")
					}

					// TODO(jasonmccallister) check is the mounts expects match whats there
				default:
					output.Pending("creating volume", hostname)

					// create the labels
					lbls := map[string]string{
						labels.Environment:     env,
						labels.DatabaseEngine:  db.Engine,
						labels.DatabaseVersion: db.Version,
					}

					// create the volume
					volResp, err := docker.VolumeCreate(ctx, volumetypes.VolumesCreateBody{
						Driver: "local",
						Name:   hostname,
						Labels: lbls,
					})
					if err != nil {
						return fmt.Errorf("unable to create the volume, %w", err)
					}

					output.Done()

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

					output.Pending("pulling", image)

					// pull the image
					rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
					if err != nil {
						return fmt.Errorf("unable to pull image %s, %w", image, err)
					}

					output.Done()

					buf := &bytes.Buffer{}
					if _, err := buf.ReadFrom(rdr); err != nil {
						return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
					}

					// set the port for the database
					port, err := nat.NewPort("tcp", db.Port)
					if err != nil {
						return fmt.Errorf("unable to create the port, %w", err)
					}

					// create the container
					output.Pending("creating", hostname)

					conResp, err := docker.ContainerCreate(
						ctx,
						&container.Config{
							Image:  image,
							Labels: lbls,
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

					output.Done()
				}

				// start the container if needed
				if startContainer {
					output.Pending("starting", hostname)

					if err := docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
						return fmt.Errorf("unable to start the container, %w", err)
					}

					output.Done()
				}

				// remove the filters
				filter.Del("label", labels.DatabaseEngine+"="+db.Engine)
				filter.Del("label", labels.DatabaseVersion+"="+db.Version)
			}

			// get all of the sites, their local path, the php version, and the type of project (nginx or PHP-FPM)
			output.Info("Checking Sites...")

			// get the envs for the sites
			envs := cfg.AsEnvs()
			for _, e := range envs {
				output.Info(e)
			}

			for _, site := range cfg.Sites {
				// add the site filter
				filter.Add("label", labels.Host+"="+site.Hostname)

				containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
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
						output.Pending(site.Hostname, "out of sync, applying")

						path, err := site.GetAbsPath()
						if err != nil {
							return err
						}

						// stop container
						if err := docker.ContainerStop(ctx, c.ID, nil); err != nil {
							return err
						}

						// remove container
						if err := docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{}); err != nil {
							return err
						}

						// pull the image
						output.Pending("pulling", image)

						rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
						if err != nil {
							return fmt.Errorf("unable to pull image, %w", err)
						}

						buf := &bytes.Buffer{}
						if _, err := buf.ReadFrom(rdr); err != nil {
							return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
						}

						output.Done()

						// create new container, will have a new container id
						// create the container
						resp, err := docker.ContainerCreate(
							ctx,
							&container.Config{
								Image: image,
								Labels: map[string]string{
									labels.Environment: env,
									labels.Host:        site.Hostname,
								},
								Env: envs,
							},
							&container.HostConfig{
								Mounts: []mount.Mount{
									{
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

						output.Done()

						break
					}

					output.Success(site.Hostname, "ready")

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
					output.Pending("pulling", image)

					rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
					if err != nil {
						return fmt.Errorf("unable to pull the image, %w", err)
					}

					buf := &bytes.Buffer{}
					if _, err := buf.ReadFrom(rdr); err != nil {
						return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
					}

					output.Done()

					output.Pending("creating", site.Hostname)

					// create the container
					resp, err := docker.ContainerCreate(
						ctx,
						&container.Config{
							Image: image,
							Labels: map[string]string{
								labels.Environment: env,
								labels.Host:        site.Hostname,
							},
							Env: envs,
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

					output.Done()
				}

				// start the container if needed
				if startContainer {
					output.Pending("starting", site.Hostname)

					if err := docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
						return fmt.Errorf("unable to start the container, %w", err)
					}

					output.Done()
				}

				// remove the site filter
				filter.Del("label", labels.Host+"="+site.Hostname)
			}

			// TODO(jasonmccallister) convert the sites into a Caddy json config and send to the API

			output.Info("Everything for", env, "is up and running 😃")

			return nil
		},
	}

	return cmd
}
