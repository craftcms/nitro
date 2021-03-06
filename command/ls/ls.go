package ls

import (
	"fmt"
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
  nitro ls

  # show only databases
  nitro ls --databases

  # show only sites
  nitro ls --sites`

var (
	flagCustom, flagDatabases, flagProxy, flagServices, flagSites bool
)

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "Lists details for Nitro’s containers.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {

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

			// define the table headers
			tbl := table.New("Hostname", "Type", "Internal Ports", "External Ports", "Status").WithWriter(cmd.OutOrStdout()).WithPadding(2)

			for _, c := range containers {
				status := "running"
				if c.State == "exited" {
					status = "stopped"
				}

				// if we only want databases
				if cmd.Flag("databases").Value.String() == "true" {
					if c.Labels[containerlabels.Type] != "database" {
						continue
					}
				}

				// show sites
				if cmd.Flag("sites").Value.String() == "true" {
					if c.Labels[containerlabels.Host] == "" {
						continue
					}
				}

				if cmd.Flag("services").Value.String() == "true" {
					if c.Labels[containerlabels.Type] != "dynamodb" && c.Labels[containerlabels.Type] != "mailhog" && c.Labels[containerlabels.Type] != "redis" {
						continue
					}
				}

				if cmd.Flag("custom").Value.String() == "true" {
					if c.Labels[containerlabels.Type] != "custom" {
						continue
					}
				}

				if cmd.Flag("proxy").Value.String() == "true" {
					if c.Labels[containerlabels.Type] != "proxy" {
						continue
					}
				}

				// get the ports
				var intPorts, extPorts []string

				// get ports for the non-site containers
				switch c.Labels[containerlabels.Host] == "" {
				case false:
					intPorts = append(intPorts, "8080", "3000", "3001")
					extPorts = append(extPorts, "(uses proxy ports)")
				default:
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
				}

				// sort the ports
				sort.Slice(intPorts, func(i, j int) bool {
					return intPorts[i] < intPorts[j]
				})

				sort.Slice(extPorts, func(i, j int) bool {
					return extPorts[i] < extPorts[j]
				})

				internalPorts := strings.Join(intPorts, ",")
				externalPorts := strings.Join(extPorts, ",")

				tbl.AddRow(strings.TrimLeft(c.Names[0], "/"), containerlabels.Identify(c), internalPorts, externalPorts, status)
			}

			tbl.Print()

			return nil
		},
	}

	cmd.Flags().BoolVarP(&flagDatabases, "databases", "d", false, "show only databases")
	cmd.Flags().BoolVarP(&flagSites, "sites", "s", false, "show only sites")
	cmd.Flags().BoolVarP(&flagServices, "services", "v", false, "show only services")
	cmd.Flags().BoolVarP(&flagCustom, "custom", "c", false, "show only custom containers")
	cmd.Flags().BoolVarP(&flagProxy, "proxy", "p", false, "show only proxy container")

	return cmd
}
