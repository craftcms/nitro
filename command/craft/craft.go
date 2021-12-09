package craft

import (
	"fmt"
	"os"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/appcontainer"
	"github.com/craftcms/nitro/pkg/dockerexec"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # run craft console command in the current app container
  nitro craft update all

  # view craft console help command
  nitro craft

  # enter the craft shell for a specific app
  nitro --app myapp.nitro craft shell`

// NewCommand returns the craft command which allows users to pass craft specific commands to a sites
// container. Its context aware and will prompt the user for the site if its not in a directory.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "craft",
		Short:   "Runs a Craft console command.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

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

			output.Info("connecting to", appName)

			// add the label to get the site
			filter.Add("label", containerlabels.Host+"="+appName)

			// find the containers but limited to the app label
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find a matching app")
			}

			// start the container if it's not running
			if containers[0].State != "running" {
				for _, command := range cmd.Root().Commands() {
					if command.Use == "start" {
						if err := command.RunE(cmd, []string{}); err != nil {
							return err
						}
					}
				}
			}

			app, err := cfg.FindAppByHostname(appName)
			if err != nil {
				return err
			}
			path := appcontainer.ContainerPath(*app)

			cmds := []string{"php", fmt.Sprintf("%s/%s", path, "craft")}

			switch len(args) == 0 {
			case true:
				// no args were provided, use the help command
				cmds = append(cmds, "help")
			default:
				// append the provided args to the command
				cmds = append(cmds, args...)
			}

			return dockerexec.Exec(cmd.InOrStdin(), cmd.OutOrStdout(), "nitro", appName, cmds...)
		},
	}

	return cmd
}
