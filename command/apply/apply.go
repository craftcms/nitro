package apply

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
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
	"github.com/craftcms/nitro/sudo"
	"github.com/craftcms/nitro/terminal"
)

var (
	// ErrNoNetwork is used when we cannot find the network
	ErrNoNetwork = fmt.Errorf("Unable to find the network")

	// ErrNoProxyContainer is returned when the proxy container is not found for an environment
	ErrNoProxyContainer = fmt.Errorf("unable to locate the proxy container")

	// NginxImage is the image used for sites, with the PHP version
	NginxImage = "docker.io/craftcms/nginx:%s-dev"

	// DatabaseImage is used for determining the engine and version
	DatabaseImage = "docker.io/library/%s:%s"
)

const exampleText = `  # apply changes from a config
  nitro apply`

func NewCommand(home string, docker client.CommonAPIClient, nitrod protob.NitroClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply changes to an environment",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			ctx := cmd.Context()
			if ctx == nil {
				// when we call commands from other commands (e.g. init)
				// the context could be nil, so we set it to the parent
				// context just in case.
				ctx = cmd.Parent().Context()
			}

			// load the config
			cfg, err := config.Load(home, env)
			if err != nil {
				return err
			}

			// parse flags

			// should we skip pulls?
			skipPulls, err := strconv.ParseBool(cmd.Flag("skip-pulls").Value.String())
			if err != nil {
				skipPulls = false
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)

			output.Info("Checking Network...")

			// check the network
			envNetwork, err := checkNetwork(ctx, docker, env, filter)
			if err != nil {
				return err
			}

			output.Success("network ready")

			output.Info("Checking Proxy...")

			// check the proxy
			proxy, err := checkProxy(ctx, docker, env)
			if err != nil {
				return err
			}

			output.Success("proxy ready")

			output.Info("Checking Databases...")

			// check the databases
			for _, db := range cfg.Databases {
				if err := checkDatabase(ctx, docker, output, filter, db, envNetwork.ID, env, skipPulls); err != nil {
					return err
				}
			}

			// get all of the sites, their local path, the php version, and the type of project (nginx or PHP-FPM)
			output.Info("Checking Sites...")

			// get the envs for the sites
			envs := cfg.AsEnvs()

			for _, site := range cfg.Sites {
				output.Pending("checking", site.Hostname)

				// add the site filter
				filter.Add("label", labels.Host+"="+site.Hostname)

				// look for a container for the site
				containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
				if err != nil {
					return fmt.Errorf("error getting a list of containers")
				}

				// create the site options
				opts := &SiteOptions{
					Site:            site,
					Home:            home,
					Environment:     env,
					EnvironmentVars: envs,
					SkipPulls:       skipPulls,
					Network:         envNetwork,
					Proxy:           &proxy,
				}

				// if there are no containers we need to create one
				switch len(containers) == 0 {
				case true:
					// create the container
					if err := createSiteContainer(ctx, docker, output, opts); err != nil {
						return err
					}

					// remove the site filter
					filter.Del("label", labels.Host+"="+site.Hostname)
				default:
					// there is a running container
					c := containers[0]

					// get the containers details that include environment variables
					details, err := docker.ContainerInspect(ctx, c.ID)
					if err != nil {
						return err
					}

					// make sure container is in sync
					if match.Site(home, site, cfg.PHP, details) == false {
						fmt.Print("- updating... ")
						// stop container
						if err := docker.ContainerStop(ctx, c.ID, nil); err != nil {
							return err
						}

						// remove container
						if err := docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{}); err != nil {
							return err
						}

						// create the site container
						if err := createSiteContainer(ctx, docker, output, opts); err != nil {
							output.Warning()
							return err
						}
					}

					// remove the site filter
					filter.Del("label", labels.Host+"="+site.Hostname)
				}

				output.Done()
			}

			output.Info("Checking Proxy...")

			output.Pending("updating proxy")

			if err := updateProxy(ctx, docker, nitrod, *cfg); err != nil {
				output.Warning()
				return err
			}

			output.Done()

			// update the hosts files
			if os.Getenv("NITRO_EDIT_HOSTS") == "false" || cmd.Flag("skip-hosts").Value.String() == "true" {
				// skip updating the hosts file
				return nil
			}

			// get all possible hostnames
			var hostnames []string
			for _, s := range cfg.Sites {
				hostnames = append(hostnames, s.Hostname)
				hostnames = append(hostnames, s.Aliases...)
			}

			if len(hostnames) > 0 {
				// get the executable
				nitro, err := os.Executable()
				if err != nil {
					return fmt.Errorf("unable to locate the nitro path, %w", err)
				}

				// run the hosts command
				switch runtime.GOOS {
				case "windows":
					return fmt.Errorf("setting hosts file is not yet supported on windows")
				default:
					output.Info("Modifying hosts file (you might be prompted for your password)")

					// add the hosts
					if err := sudo.Run(nitro, "nitro", "hosts", "--hostnames="+strings.Join(hostnames, ",")); err != nil {
						return err
					}
				}
			}

			output.Info(env, "is up and running 😃")

			return nil
		},
	}

	// add flag to skip pulling images
	cmd.Flags().Bool("skip-pull", false, "skip pulling images")
	cmd.Flags().Bool("skip-hosts", false, "skip modifying the hosts file")

	return cmd
}

func createContainer(ctx context.Context, docker client.ContainerAPIClient, config *container.Config, host *container.HostConfig, network *network.NetworkingConfig, name string) (string, error) {
	// create the container
	resp, err := docker.ContainerCreate(ctx, config, host, network, name)
	if err != nil {
		return "", fmt.Errorf("unable to create the container, %w", err)
	}

	// start the container
	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("unable to start the container, %w", err)
	}

	return resp.ID, nil
}

func checkNetwork(ctx context.Context, docker client.NetworkAPIClient, env string, filter filters.Args) (*types.NetworkResource, error) {
	// find networks
	networks, err := docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
	if err != nil {
		return nil, fmt.Errorf("unable to list docker networks\n%w", err)
	}

	// get the network for the environment
	for _, network := range networks {
		if network.Name == env {
			return &network, nil
		}
	}

	// the network is not found
	return nil, ErrNoNetwork
}

func checkProxy(ctx context.Context, docker client.ContainerAPIClient, env string) (types.Container, error) {
	f := filters.NewArgs()
	f.Add("label", labels.Proxy+"="+env)
	// TODO(jasonmccallister) add the type filter as well?

	// check if there is an existing container for the nitro-proxy
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: f, All: true})
	if err != nil {
		return types.Container{}, fmt.Errorf("unable to list the containers\n%w", err)
	}

	// get the container id and determine if the container needs to start
	for _, c := range containers {
		for _, n := range c.Names {
			if n == env || n == "/"+env {
				// check if it is running
				if c.State != "running" {
					if err := docker.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
						return types.Container{}, fmt.Errorf("unable to start the nitro container, %w", err)
					}
				}

				// return the container
				return c, nil
			}
		}
	}

	return types.Container{}, ErrNoProxyContainer
}

func checkDatabase(ctx context.Context, docker client.CommonAPIClient, output terminal.Outputer, filter filters.Args, db config.Database, networkID, env string, skipPull bool) error {
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
	switch len(containers) {
	// the database container exists
	case 1:
		// check if the container is running
		if containers[0].State != "running" {
			output.Pending("starting", hostname)

			// start the container
			if err := docker.ContainerStart(ctx, containers[0].ID, types.ContainerStartOptions{}); err != nil {
				output.Warning()
				return err
			}

			output.Done()
		} else {
			output.Success(hostname, "ready")
		}
	default:
		// database container does not exist, so create the volume and start it
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
		volume, err := docker.VolumeCreate(ctx, volumetypes.VolumesCreateBody{Driver: "local", Name: hostname, Labels: lbls})
		if err != nil {
			return fmt.Errorf("unable to create the volume, %w", err)
		}

		output.Done()

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

		// check for if we should skip pulling an image
		if skipPull {
			output.Pending("pulling", image)

			// pull the image
			rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
			if err != nil {
				output.Warning()
				return fmt.Errorf("unable to pull image %s, %w", image, err)
			}

			// read the output to pull the image
			buf := &bytes.Buffer{}
			if _, err := buf.ReadFrom(rdr); err != nil {
				output.Warning()
				return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
			}

			output.Done()
		}

		output.Pending("creating", hostname)

		// set the port for the database
		port, err := nat.NewPort("tcp", db.Port)
		if err != nil {
			output.Warning()
			return fmt.Errorf("unable to create the port, %w", err)
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
				env: {
					NetworkID: networkID,
				},
			},
		}

		// create the container for the database
		if _, err := createContainer(ctx, docker, containerConfig, hostConfig, networkConfig, hostname); err != nil {
			output.Warning()
			return err
		}

		output.Done()
	}

	// remove the database filters
	filter.Del("label", labels.DatabaseEngine+"="+db.Engine)
	filter.Del("label", labels.DatabaseVersion+"="+db.Version)
	filter.Del("label", labels.Type+"=database")
	return nil
}

// SiteOptions are used to create site containers
type SiteOptions struct {
	Site            config.Site
	Home            string
	Environment     string
	EnvironmentVars []string
	SkipPulls       bool
	Network         *types.NetworkResource
	Proxy           *types.Container
}

func createSiteContainer(ctx context.Context, docker client.CommonAPIClient, output terminal.Outputer, opts *SiteOptions) error {
	image := fmt.Sprintf(NginxImage, opts.Site.PHP)

	// should we skip pulling the image
	if opts.SkipPulls {
		output.Pending("pulling", image)

		// pull the image
		rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
		if err != nil {
			return fmt.Errorf("unable to pull the image, %w", err)
		}

		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(rdr); err != nil {
			return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
		}

		output.Done()
	}

	// get the sites path
	path, err := opts.Site.GetAbsPath(opts.Home)
	if err != nil {
		return err
	}

	// add the site itself to the extra hosts
	extraHosts := []string{fmt.Sprintf("%s:%s", opts.Site.Hostname, "127.0.0.1")}
	for _, s := range opts.Site.Aliases {
		extraHosts = append(extraHosts, fmt.Sprintf("%s:%s", s, "127.0.0.1"))
	}

	// check if xdebug is enabled
	switch opts.Site.Xdebug {
	case false:
		opts.EnvironmentVars = append(opts.EnvironmentVars, "XDEBUG_MODE=off")
	default:
		// opts.Proxy.NetworkSettings.Networks[opts.Environment].IPAddress
		// opts.Network.IPAM.Config[0].Gateway
		opts.EnvironmentVars = append(opts.EnvironmentVars, fmt.Sprintf(`XDEBUG_CONFIG=client_host=%s log=/tmp/xdebug.log start_with_request=yes log_level=10`, opts.Proxy.NetworkSettings.Networks[opts.Environment].IPAddress))
		opts.EnvironmentVars = append(opts.EnvironmentVars, "XDEBUG_SESSION=nitro")
		opts.EnvironmentVars = append(opts.EnvironmentVars, "XDEBUG_MODE=develop,debug")
	}

	// create the container
	resp, err := docker.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			Labels: map[string]string{
				labels.Environment: opts.Environment,
				labels.Host:        opts.Site.Hostname,
			},
			Env: opts.EnvironmentVars,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: path,
					Target: "/app",
				},
			},
			ExtraHosts: extraHosts,
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				opts.Environment: {
					NetworkID: opts.Network.ID,
				},
			},
		},
		opts.Site.Hostname,
	)
	if err != nil {
		return fmt.Errorf("unable to create the container, %w", err)
	}

	// start the container
	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("unable to start the container, %w", err)
	}

	// TODO(jasonmccallister) check for a custom root and copt the template to the container

	return nil
}

func updateProxy(ctx context.Context, docker client.ContainerAPIClient, nitrod protob.NitroClient, cfg config.Config) error {
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

	// if there are no sites, we are done
	if len(sites) == 0 {
		return nil
	}

	// wait for the api to be ready
	wait := true
	for wait {
		_, err := nitrod.Ping(ctx, &protob.PingRequest{})
		if err == nil {
			wait = false
			break
		}
	}

	// configure the proxy with the sites
	if _, err := nitrod.Apply(ctx, &protob.ApplyRequest{Sites: sites}); err != nil {
		return err
	}

	return nil
}
