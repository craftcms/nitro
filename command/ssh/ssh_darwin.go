// +build darwin, !linux

package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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

			var site string
			if len(args) > 0 {
				site = strings.TrimSpace(args[0])
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

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

				// start the container if its not running
				if containers[0].State != "running" {
					for _, command := range cmd.Root().Commands() {
						if command.Use == "start" {
							if err := command.RunE(cmd, []string{}); err != nil {
								return err
							}
						}
					}
				}

				containerID = containers[0].ID
			default:
				// get a context aware list of sites
				sites := cfg.ListOfSitesByDirectory(home, wd)

				// create the options for the sites
				var options []string
				for _, s := range sites {
					options = append(options, s.Hostname)
				}

				// did they ask for a specific site?
				switch site != "" {
				case true:
					for k, v := range options {
						if site == v {
							// add the label to get the site
							filter.Add("label", containerlabels.Host+"="+sites[k].Hostname)
							break
						}
					}
				default:
					// if there are found sites we want to show or connect to the first one, otherwise prompt for which site to connect to.
					switch len(sites) {
					case 0:
						// prompt for the site to ssh into
						selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
						if err != nil {
							return err
						}

						// add the label to get the site
						filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)
					case 1:
						output.Info("connecting to", sites[0].Hostname)

						// add the label to get the site
						filter.Add("label", containerlabels.Host+"="+sites[0].Hostname)
					default:
						// prompt for the site to ssh into
						selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
						if err != nil {
							return err
						}

						// add the label to get the site
						filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)
					}
				}

				// find the containers but limited to the site label
				containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
				if err != nil {
					return err
				}

				// are there any containers??
				if len(containers) == 0 {
					return fmt.Errorf("unable to find an matching site")
				}

				// start the container if its not running
				if containers[0].State != "running" {
					for _, command := range cmd.Root().Commands() {
						if command.Use == "start" {
							if err := command.RunE(cmd, []string{}); err != nil {
								return err
							}
						}
					}
				}

				containerID = containers[0].ID
			}

			// find the docker executable
			cli, err := exec.LookPath("docker")
			if err != nil {
				return err
			}

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

			c := exec.Command(cli, "exec", "-u", containerUser, "-it", containerID, shell)

			c.Stdin = cmd.InOrStdin()
			c.Stderr = cmd.ErrOrStderr()
			c.Stdout = cmd.OutOrStdout()

			return c.Run()
		},
	}

	cmd.Flags().BoolVar(&RootUser, "root", false, "connect as root user")
	cmd.Flags().BoolVar(&ProxyContainer, "proxy", false, "connect to proxy container")

	return cmd
}
