package php

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/dockerexec"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # run php version for the current app
  nitro php

  # get an interactive shell for an app
  nitro --app myapp.nitro php -- -a`

// NewCommand returns the php command which allows users to pass php specific commands to a sites
// container. Its context aware and will prompt the user for the site if its not in a directory.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "php",
		Short:   "Runs PHP commands.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			appName := flags.AppName
			if appName == "" {
				// get the current working directory
				wd, err := os.Getwd()
				if err != nil {
					return err
				}

				appName, err = appaware.Detect(*cfg, wd)
				if err != nil {
					return err
				}
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			output.Info("connecting to", appName)

			// add the label to get the site
			filter.Add("label", containerlabels.Host+"="+appName)

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

			// create the command for running the php command
			cmds := []string{"php"}

			// append the provided args to the command
			if len(args) == 0 {
				cmds = append(cmds, "-v")
			} else {
				cmds = append(cmds, args...)
			}

			return dockerexec.Exec(cmd.InOrStdin(), cmd.OutOrStdout(), "nitro", appName, cmds...)
		},
	}

	return cmd
}
