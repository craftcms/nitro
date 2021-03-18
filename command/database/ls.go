package database

import (
	"sort"
	"strings"

	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var lsExampleText = `  # show the database containers
  nitro db ls`

func lsCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "Show database containers",
		Example: lsExampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// add filters to show only the environment and database containers
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)
			filter.Add("label", containerlabels.Type+"=database")

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
			tbl := table.New("Hostname", "External Port", "Internal Port", "Username/Password", "State")
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

				// if there is more than one port, grab the first one
				var external, internal uint16
				switch len(c.Ports) {
				case 1:
					external = c.Ports[0].PublicPort
					internal = c.Ports[0].PrivatePort
				default:
					for _, p := range c.Ports {
						if p.PublicPort != 0 {
							external = p.PublicPort
							internal = p.PrivatePort
						}
					}
				}

				tbl.AddRow(strings.TrimLeft(c.Names[0], "/"), external, internal, "nitro/nitro", c.State)
			}

			tbl.Print()

			return nil
		},
	}

	return cmd
}
