package ls

import (
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # view information about your nitro environment
  nitro ls`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "Show nitro info",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			// cfg, err := config.Load(home)
			// if err != nil {
			// 	return err
			// }

			// add filters to show only the environment and database containers
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// get a list of all the databases
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{All: true, Filters: filter})
			if err != nil {
				return err
			}

			// sort containers by the name
			sort.SliceStable(containers, func(i, j int) bool {
				return containers[i].Names[0] < containers[j].Names[0]
			})

			// print the table headers
			tbl := table.New("Hostname", "Path", "State")
			tbl.WithWriter(cmd.OutOrStdout())

			// generate a list of engines for the prompt
			for _, c := range containers {
				// start the container if not running
				if c.State != "running" {
					for _, command := range cmd.Root().Commands() {
						if command.Use == "start" {
							if err := command.RunE(cmd, []string{}); err != nil {
								return err
							}
						}
					}
				}

				if c.Labels[containerlabels.Host] != "" {
					tbl.AddRow(strings.TrimLeft(c.Names[0], "/"), c.Mounts[0].Mode, c.State)
				}
			}

			tbl.Print()

			return nil
		},
	}

	return cmd
}
