package start

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/proxycontainer"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/contextor"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrNoContainers is returned when no containers are running for an environment
	ErrNoContainers = fmt.Errorf("there are no running containers")
)

const exampleText = `  # start all containers
  nitro start

  # start a specific site
  nitro start my-site.nitro`

var site *config.Site

// NewCommand returns the command used to start containers.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Starts containers.",
		Example: exampleText,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			cfg, err := config.Load(home)
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}

			var options []string
			for _, s := range cfg.Sites {
				options = append(options, s.Hostname)
			}

			return options, cobra.ShellCompDirectiveDefault
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// is the docker api alive?
			if _, err := docker.Ping(cmd.Context()); err != nil {
				return fmt.Errorf("Couldn’t connect to Docker; please make sure Docker is running.")
			}

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// is there a site as the first arg?
			if len(args) > 0 {
				site, err = cfg.FindSiteByHostName(args[0])
				if err != nil {
					return err
				}

				return nil
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := contextor.New(cmd.Context())

			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			output.Info("Starting Nitro…")

			// inspect the container
			info, err := docker.ContainerInspect(ctx, proxycontainer.ProxyName)
			if err != nil {
				return err
			}

			// if the container is already running
			switch info.State.Status == "running" {
			case false:
				output.Pending("starting", proxycontainer.ProxyName)

				// start the proxy container
				if err := docker.ContainerStart(ctx, proxycontainer.ProxyName, types.ContainerStartOptions{}); err != nil {
					return fmt.Errorf("unable to start the proxy container")
				}

				output.Done()
			default:
				output.Success(proxycontainer.ProxyName)
			}

			// start the databases
			for _, db := range cfg.Databases {
				name, err := db.GetHostname()
				if err != nil {
					return err
				}

				// inspect the container
				info, err := docker.ContainerInspect(ctx, name)
				if err != nil {
					return err
				}

				// if the container is already running
				switch info.State.Status == "running" {
				case false:
					output.Pending("starting", name)

					// start the proxy container
					if err := docker.ContainerStart(ctx, name, types.ContainerStartOptions{}); err != nil {
						return fmt.Errorf("unable to start the database container: err %w", err)
					}

					output.Done()
				default:
					output.Success(name)
				}
			}

			// start each environment container
			for _, s := range cfg.Sites {
				// if the user wants a single site only, skip all the other sites
				if (site != nil) && s.Hostname != site.Hostname {
					continue
				}

				// make sure the site is not disabled
				if s.Disabled {
					output.Info(fmt.Sprintf("%s is disabled, skipping", s.Hostname))
					continue
				}

				// inspect the container
				info, err := docker.ContainerInspect(ctx, s.Hostname)
				if err != nil {
					return err
				}

				// if the container is already running
				switch info.State.Status == "running" {
				case false:
					output.Pending("starting", s.Hostname)

					// start the proxy container
					if err := docker.ContainerStart(ctx, s.Hostname, types.ContainerStartOptions{}); err != nil {
						return fmt.Errorf("unable to start the database container: err %w", err)
					}

					output.Done()
				default:
					output.Success(s.Hostname)
				}
			}

			output.Info("Nitro started 👍")

			return nil
		},
	}

	return cmd
}
