package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/terminal"
)

func lsCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"l"},
		Short:   "Lists details for Nitro’s services.",
		Example: `  # view information about your nitro environment
  nitro service ls`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// add filters to show only the environment and database containers
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// get a list of all the databases
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
			if err != nil {
				return err
			}

			// sort containers by the name
			sort.SliceStable(containers, func(i, j int) bool {
				return containers[i].Names[0] < containers[j].Names[0]
			})

			// define the table headers
			tbl := table.New("Hostname", "Type", "External Ports", "Internal Ports", "Status").WithWriter(cmd.OutOrStdout()).WithPadding(2)

			for _, c := range containers {
				status := "running"
				if c.State == "exited" {
					status = "stopped"
				}

				// show only the service containers
				if c.Labels[containerlabels.Type] != "dynamodb" && c.Labels[containerlabels.Type] != "mailhog" && c.Labels[containerlabels.Type] != "redis" && c.Labels[containerlabels.Type] != "blackfire" {
					continue
				}

				// get the ports
				var intPorts, extPorts []string

				// get ports
				for _, p := range c.Ports {
					// get the external ports and assign if not 0
					e := p.PublicPort
					if e != 0 {
						extPorts = append(extPorts, fmt.Sprintf("%d", e))
					}

					// get the internal ports and assign if not 0
					pr := p.PrivatePort
					if e != 0 {
						intPorts = append(intPorts, fmt.Sprintf("%d", pr))
					}
				}

				internalPorts := strings.Join(intPorts, ",")
				externalPorts := strings.Join(extPorts, ",")

				tbl.AddRow(strings.TrimLeft(c.Names[0], "/"), containerlabels.Identify(c), externalPorts, internalPorts, status)
			}

			tbl.Print()

			fmt.Println("\nNote: You can enable or disable services using `nitro service enable <service>`.")
			fmt.Println("")

			return nil
		},
	}

	return cmd
}
