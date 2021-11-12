//go:build darwin || !linux
// +build darwin !linux

package ssh

import (
	"fmt"
	"os"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/dockerexec"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/proxycontainer"
	"github.com/craftcms/nitro/pkg/terminal"
)

// NewCommand returns the ssh command to get a shell in a container. The command is context aware and if
// it is not in a known project directory, it will provide a list of known sites to the user.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssh",
		Short:   "Opens a shell in a container.",
		Example: exampleText,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// is the docker api alive?
			if _, err := docker.Ping(cmd.Context()); err != nil {
				return fmt.Errorf("couldn’t connect to Docker; please make sure Docker is running")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// check if the root user should be used
			containerUser := "nitro"
			if RootUser || ProxyContainer {
				containerUser = "root"
			}

			// show a notice about changes
			if containerUser == "root" {
				output.Info("using root… system changes are ephemeral…")
			}

			shell := "bash"
			if ProxyContainer {
				shell = "sh"
			}

			var containerID string
			switch ProxyContainer {
			case true:
				// file by the container name
				filter.Add("name", proxycontainer.ProxyName)

				// find the containers but limited to the site label
				containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
				if err != nil {
					return err
				}

				if len(containers) == 0 {
					return fmt.Errorf("no containers found")
				}

				// start the container if not running
				if containers[0].State != "running" {
					if err := docker.ContainerStart(cmd.Context(), containers[0].ID, types.ContainerStartOptions{}); err != nil {
						return err
					}
				}

				containerID = containers[0].ID
			default:
				// is there a global flag for the app?
				if flags.AppName != "" {
					// find the app by the hostname
					app, err :=cfg.FindAppByHostname(flags.AppName)
					if err != nil {
						return err
					}

					// add the label to get the site
					filter.Add("label", containerlabels.Host+"="+app.Hostname)

					// find the containers but limited to the app label
					containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
					if err != nil {
						return err
					}

					// are there any containers??
					if len(containers) == 0 {
						return fmt.Errorf("unable to find an matching site")
					}

					// start the container if not running
					if containers[0].State != "running" {
						if err := docker.ContainerStart(cmd.Context(), containers[0].ID, types.ContainerStartOptions{}); err != nil {
							return err
						}
					}

					containerID = containers[0].ID

					break
				}

				appName, err := appaware.Detect(*cfg, wd)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+appName)

				// find the containers but limited to the app label
				containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
				if err != nil {
					return err
				}

				// are there any containers??
				if len(containers) == 0 {
					return fmt.Errorf("unable to find an matching site")
				}

				// start the container if not running
				if containers[0].State != "running" {
					if err := docker.ContainerStart(cmd.Context(), containers[0].ID, types.ContainerStartOptions{}); err != nil {
						return err
					}
				}

				containerID = containers[0].ID
			}

			return dockerexec.Connect(cmd.InOrStdin(), cmd.OutOrStdout(), containerUser, containerID, shell)
		},
	}

	cmd.Flags().BoolVar(&RootUser, "root", false, "connect as root user")
	cmd.Flags().BoolVar(&ProxyContainer, "proxy", false, "connect to proxy container")

	return cmd
}
