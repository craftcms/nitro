package apply

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/apply/internal/match"
	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/protob"
	"github.com/craftcms/nitro/terminal"
)

var (
	// ErrNoNetwork is used when we cannot find the network
	ErrNoNetwork = fmt.Errorf("Unable to find the network")

	// NginxImage is the image used for sites, with the PHP version
	NginxImage = "docker.io/craftcms/nginx:%s"

	// DatabaseImage is used for determining the engine and version
	DatabaseImage = "docker.io/library/%s:%s"
)

const exampleText = `  # apply changes from a config
  nitro apply`

// New takes a docker client and the terminal output to run the apply actions
func New(docker client.CommonAPIClient, nitrod protob.NitroClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply changes to an environment",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			ctx := cmd.Context()

			cfg, err := config.Load(env)
			if err != nil {
				return err
			}

			// create a filter for the network
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)

			output.Info("Checking Network...")

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
					output.Success("network ready")
				}
			}

			// if the network is not found
			if networkID == "" {
				return ErrNoNetwork
			}

			output.Info("Checking Proxy...")

			proxyFilter := filters.NewArgs()
			proxyFilter.Add("label", labels.Proxy+"="+env)

			// check if there is an existing container for the nitro-proxy
			var proxyContainerID string
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: proxyFilter, All: true})
			if err != nil {
				return fmt.Errorf("unable to list the containers\n%w", err)
			}

			var proxyRunning bool
			for _, c := range containers {
				for _, n := range c.Names {
					if n == env || n == "/"+env {
						proxyContainerID = c.ID

						// check if it is running
						if c.State == "running" {
							output.Success("proxy ready")
							proxyRunning = true
						}
					}
				}
			}

			if !proxyRunning {
				output.Pending("starting proxy")
				if err := docker.ContainerStart(ctx, proxyContainerID, types.ContainerStartOptions{}); err != nil {
					return fmt.Errorf("unable to start the nitro container, %w", err)
				}
				output.Done()
			}

			// check the databases
			output.Info("Checking Databases...")
			for _, db := range cfg.Databases {
				// add filters to check for the container
				filter.Add("label", labels.DatabaseEngine+"="+db.Engine)
				filter.Add("label", labels.DatabaseVersion+"="+db.Version)
				filter.Add("label", labels.Type+"=database")

				// get the containers for the database
				containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
				if err != nil {
					return fmt.Errorf("error getting a list of containers")
				}

				// set the hostname for the database container
				hostname, err := db.GetHostname()
				if err != nil {
					return err
				}

				// if there is not a container for the database, create a volume, container, and start the container
				var containerID string
				var startContainer bool
				switch len(containers) {
				// there database container exists
				case 1:
					// set the container id
					containerID = containers[0].ID

					// check if the container is running
					if containers[0].State != "running" {
						startContainer = true
						output.Pending("starting", hostname)
					} else {
						output.Success(hostname, "ready")
					}
					// database container does not exist, so create the volume and start it
				default:
					output.Pending("creating volume", hostname)

					// create the database labels
					lbls := map[string]string{
						labels.Environment:     env,
						labels.DatabaseEngine:  db.Engine,
						labels.DatabaseVersion: db.Version,
						labels.Type:            "database",
					}

					// if the database is mysql or mariadb, mark them as
					// mysql compatible (used for importing backups)
					if db.Engine == "mariadb" || db.Engine == "mysql" {
						lbls[labels.DatabaseCompatability] = "mysql"
					}

					// if the database is postgres, mark it as compatible
					// with postgres. This is not needed but a place holder
					// if cockroachdb is ever supported by craft.
					if db.Engine == "postgres" {
						lbls[labels.DatabaseCompatability] = "postgres"
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
					image := fmt.Sprintf(DatabaseImage, db.Engine, db.Version)

					target := "/var/lib/mysql"
					var envs []string
					if strings.Contains(image, "postgres") {
						target = "/var/lib/postgresql/data"
						envs = []string{"POSTGRES_USER=nitro", "POSTGRES_DB=nitro", "POSTGRES_PASSWORD=nitro"}
					} else {
						envs = []string{"MYSQL_ROOT_PASSWORD=nitro", "MYSQL_DATABASE=nitro", "MYSQL_USER=nitro", "MYSQL_PASSWORD=nitro"}
					}

					// TODO(jasonmccallister) check for skip apply
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

					output.Pending("creating", hostname)

					// tell docker to create the database container
					resp, err := docker.ContainerCreate(
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

					// set the database container id to start
					containerID = resp.ID
					startContainer = true
				}

				// start the database container if needed
				if startContainer {
					if err := docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
						return fmt.Errorf("unable to start the container, %w", err)
					}

					output.Done()
				}

				// remove the database filters
				filter.Del("label", labels.DatabaseEngine+"="+db.Engine)
				filter.Del("label", labels.DatabaseVersion+"="+db.Version)
				filter.Del("label", labels.Type+"=database")
			}

			// get all of the sites, their local path, the php version, and the type of project (nginx or PHP-FPM)
			if err := checkSites(ctx, docker, output, filter, env, networkID, cfg); err != nil {
				return err
			}

			// convert the sites into the gRPC API Apply request
			sites := make(map[string]*protob.Site)
			for _, s := range cfg.Sites {
				hosts := []string{s.Hostname}

				// if there are aliases lets append them to the hosts
				if len(s.Aliases) > 0 {
					hosts = append(hosts, s.Aliases...)
				}

				// create the site
				sites[s.Hostname] = &protob.Site{
					Hostname: s.Hostname,
					Aliases:  strings.Join(hosts, ","),
					Port:     8080,
				}
			}

			// ping the nitrod API until its ready...
			output.Pending("waiting for api")
			ping := &protob.PingRequest{}
			waiting := true
			for waiting {
				_, err := nitrod.Ping(ctx, ping)
				if err == nil {
					waiting = false
				}
			}

			output.Done()

			// configure the proxy with the sites
			output.Info("Configuring Proxy...")
			if _, err = nitrod.Apply(ctx, &protob.ApplyRequest{Sites: sites}); err != nil {
				return err
			}

			output.Success("proxy ready")

			output.Info(env, "is up and running 😃")

			return nil
		},
	}

	// add flag to skip pulling images
	cmd.Flags().BoolP("skip-pull", "s", false, "skip pulling images")

	return cmd
}

func checkSites(
	ctx context.Context,
	docker client.CommonAPIClient,
	output terminal.Outputer,
	filter filters.Args,
	env, networkID string,
	cfg *config.Config,
) error {
	output.Info("Checking Sites...")

	// get the envs for the sites
	envs := cfg.AsEnvs()
	for _, site := range cfg.Sites {
		// add the site filter
		filter.Add("label", labels.Host+"="+site.Hostname)

		// look for a container for the site
		containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
		if err != nil {
			return fmt.Errorf("error getting a list of containers")
		}

		var containerID string
		var startContainer bool
		switch len(containers) {
		case 1:
			// there is a running container
			c := containers[0]
			image := fmt.Sprintf(NginxImage, site.PHP)
			path, err := site.GetAbsPath()
			if err != nil {
				return err
			}

			expected, err := site.GetAbsMountPaths()
			if err != nil {
				return err
			}
			// hard code the path to the first site mount
			expected[path] = "/app"

			// make sure the images and mounts match, if they don't stop, remove, and create the container
			// with the new image
			if c.Image != image || match.Mounts(c.Mounts, expected) == false {
				output.Pending(site.Hostname, "out of sync")

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

				output.Done()

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

				// add the path mount
				mounts := []mount.Mount{}
				mounts = append(mounts, mount.Mount{
					Type:   mount.TypeBind,
					Source: path,
					Target: "/app",
				})

				// get additional site mounts
				siteMounts, err := site.GetAbsMountPaths()
				if err != nil {
					return err
				}

				for k, v := range siteMounts {
					mounts = append(mounts, mount.Mount{
						Type:   mount.TypeBind,
						Source: k,
						Target: v,
					})
				}

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
						Mounts: mounts,
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

				break
			}

			// get the container id
			containerID = c.ID

			// check if the container is running
			if containers[0].State != "running" {
				startContainer = true
				output.Pending("starting", site.Hostname)
			} else {
				output.Success(site.Hostname, "ready")
			}
		default:
			// create a brand new container since there is not an existing one
			image := fmt.Sprintf(NginxImage, site.PHP)

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

			// add the path mount
			mounts := []mount.Mount{}
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: path,
				Target: "/app",
			})

			// get additional site mounts
			siteMounts, err := site.GetAbsMountPaths()
			if err != nil {
				return err
			}

			for k, v := range siteMounts {
				mounts = append(mounts, mount.Mount{
					Type:   mount.TypeBind,
					Source: k,
					Target: v,
				})
			}

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
					Mounts: mounts,
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
		}

		// start the container if needed
		if startContainer {
			if err := docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}

			output.Done()
		}

		// remove the site filter
		filter.Del("label", labels.Host+"="+site.Hostname)
	}

	return nil
}
