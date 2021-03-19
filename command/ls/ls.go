package ls

import (
	"context"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # view information about your nitro environment
  nitro ls`

var (
	flagCustom, flagDatabases, flagProxy, flagSites, flagVerbose bool
)

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "Show Nitro info",
		Example: exampleText,
		// Deprecated: true,
		// Aliases: []string{"context"},
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
			tbl := table.New("Hostname", "Type", "State")
			tbl.WithWriter(cmd.OutOrStdout()).WithPadding(10)

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

				hostname := strings.TrimLeft(c.Names[0], "/")

				containerType := "site"
				if c.Labels[containerlabels.DatabaseEngine] != "" {
					containerType = "database"
				}

				if c.Labels[containerlabels.NitroContainer] != "" {
					containerType = "custom"
				}

				if c.Labels[containerlabels.Proxy] != "" {
					containerType = "proxy"
				}

				tbl.AddRow(hostname, containerType, c.State)
			}

			tbl.Print()

			return nil
		},
	}

	cmd.Flags().BoolVar(&flagCustom, "custom", false, "Show only custom containers")
	cmd.Flags().BoolVar(&flagDatabases, "database", false, "Show only database containers")
	cmd.Flags().BoolVar(&flagProxy, "proxy", false, "Show only the proxy container")
	cmd.Flags().BoolVar(&flagSites, "site", false, "Show only site containers")
	cmd.Flags().BoolVar(&flagVerbose, "verbose", false, "Show extended information")

	return cmd
}

func sitesTable(ctx context.Context, home string, cfg *config.Config, docker client.CommonAPIClient) table.Table {
	return table.New()
}
