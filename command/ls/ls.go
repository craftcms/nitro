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

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # view information about your nitro environment
  nitro ls`

var (
	flagCustom, flagDatabases, flagProxy, flagServices, flagSites, flagVerbose bool
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
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

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

			switch flagVerbose {
			case false:
				// print the table headers
				tbl := table.New("Hostname", "Type", "Status").WithWriter(cmd.OutOrStdout()).WithPadding(10)

				// generate a list of engines for the prompt
				for _, c := range containers {
					tbl.AddRow(strings.TrimLeft(c.Names[0], "/"), getContainerType(c), c.Status)
				}

				tbl.Print()
			default:
				customTbl := table.New("Hostname", "External Port", "Internal Port", "Username/Password", "Status").WithWriter(cmd.OutOrStdout())
				databaseTbl := table.New("Hostname", "External Port", "Internal Port", "Username/Password", "Status").WithWriter(cmd.OutOrStdout())
				sitesTbl := table.New("Hostname", "PHP", "Xdebug", "Path", "Webroot", "Status").WithWriter(cmd.OutOrStdout())

				var showCustom, showDatabases, showSites bool
				for _, c := range containers {
					// get the container type
					containerType := getContainerType(c)

					// get information based on the container type
					switch containerType {
					case "custom":
						showCustom = true
						customTbl.AddRow(strings.TrimLeft(c.Names[0], "/"), c.Status)
					case "database":
						showDatabases = true
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

						databaseTbl.AddRow(strings.TrimLeft(c.Names[0], "/"), external, internal, "nitro/nitro", c.Status)
					default:
						// find the site by the hostname
						site, err := cfg.FindSiteByHostName(strings.TrimLeft(c.Names[0], "/"))
						if err != nil {
							break
						}

						sitesTbl.AddRow(fmt.Sprintf("https://%s", site.Hostname), site.Version, site.Xdebug, site.Path, site.Webroot, c.Status)
						showSites = true
					}
				}

				// output the tables
				fmt.Println("")

				if showCustom {
					customTbl.Print()
					fmt.Println("")
				}

				if showDatabases {
					databaseTbl.Print()
					fmt.Println("")
				}

				if showSites {
					sitesTbl.Print()
					fmt.Println("")
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&flagCustom, "custom", false, "Show only custom containers")
	cmd.Flags().BoolVar(&flagDatabases, "database", false, "Show only database containers")
	cmd.Flags().BoolVar(&flagProxy, "proxy", false, "Show only the proxy container")
	cmd.Flags().BoolVar(&flagSites, "site", false, "Show only site containers")
	cmd.Flags().BoolVar(&flagServices, "service", false, "Show only service containers")
	cmd.Flags().BoolVar(&flagVerbose, "verbose", false, "Show extended information")

	return cmd
}

func getContainerType(c types.Container) string {
	if c.Labels[containerlabels.DatabaseEngine] != "" {
		return "database"
	}

	if c.Labels[containerlabels.NitroContainer] != "" {
		return "custom"
	}

	if c.Labels[containerlabels.Proxy] != "" {
		return "proxy"
	}

	return "site"
}
