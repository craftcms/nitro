package share

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
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// execName is the name of the executable to search for. We make it a variable so we can replace it during tests.
	execName = "ngrok"
)

const exampleText = `  # share a local site with ngrok
  nitro share`

// NewCommand is used to destroy all resources for an environment. It will prompt for
// user verification and defaults to no. Part of the destroy process is to
// perform a backup for all databases in each container database.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "share",
		Short:   "Shares a local site.",
		Example: exampleText,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			cfg, err := config.Load(home)
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}

			var options []string
			for _, s := range cfg.Sites {
				options = append(options, s.Hostname)
			}

			return options, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// find ngrok
			ngrok, err := exec.LookPath(execName)
			if err != nil {
				return err
			}

			// if ngrok is missing, return the error
			if ngrok == "" {
				output.Info("Ngrok is required to share sites, download ngrok from https://ngrok.com")

				return nil
			}

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

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			var individual string
			if len(args) > 0 {
				individual = strings.TrimSpace(args[0])
			}

			var site config.Site
			switch individual == "" {
			case false:
				for k, v := range options {
					if individual == v {
						filter.Add("label", containerlabels.Host+"="+sites[k].Hostname)

						site = sites[k]

						break
					}
				}

			default:
				switch len(sites) {
				case 0:
					// prompt for the site to ssh into
					selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
					if err != nil {
						return err
					}

					// add the label to get the site
					filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)

					site = sites[selected]
				case 1:
					output.Info("connecting to", sites[0].Hostname)

					// add the label to get the site
					filter.Add("label", containerlabels.Host+"="+sites[0].Hostname)
					site = sites[0]
				default:
					// prompt for the site to ssh into
					selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
					if err != nil {
						return err
					}

					// add the label to get the site
					filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)
					site = sites[selected]
				}
			}

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter, All: true})
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

			// parse the flags

			ngrokArgs := []string{"http"}

			// set the main hostname
			ngrokArgs = append(ngrokArgs, "-host-header="+site.Hostname)

			// append the aliases
			for _, a := range site.Aliases {
				ngrokArgs = append(ngrokArgs, "-host-header="+a)
			}

			// set the region
			region, err := cmd.Flags().GetString("region")
			if err != nil {
				region = "us"
			}
			ngrokArgs = append(ngrokArgs, "--region="+region)

			// set the port
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				port = "80"
			}
			ngrokArgs = append(ngrokArgs, port)

			c := exec.Command(ngrok, ngrokArgs...)

			c.Stderr = cmd.ErrOrStderr()
			c.Stdout = cmd.OutOrStdout()

			return c.Run()
		},
	}

	// add flags to the command
	cmd.Flags().String("region", "us", "which ngrok region to use for sharing")
	cmd.Flags().String("port", "80", "which port to use for ngrok")

	return cmd
}
