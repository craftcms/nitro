package xon

import (
	"fmt"
	"os"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const exampleText = `  # example command
  nitro xon`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "xon",
		Short:   "Enable xdebug for a site",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// load the config
			cfg, err := config.Load(home, env)
			if err != nil {
				return err
			}

			// create a filter for the enviroment
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)

			// check if we are a known site
			var sites []string
			for _, site := range cfg.Sites {
				// get the path
				path, _ := site.GetAbsPath(home)

				// add the site as an option
				sites = append(sites, site.Hostname)

				// are we in a current project/site?
				if strings.Contains(wd, path) {
					filter.Add("label", labels.Host+"="+site.Hostname)
				}
			}

			// find all of the containers, there should only be one if we are in a known directory
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			switch len(containers) {
			case 1:
				// get the containers details
				details, err := docker.ContainerInspect(cmd.Context(), containers[0].ID)
				if err != nil {
					return err
				}

				// find the environment variable for xdebug
				for _, e := range details.Config.Env {
					env := strings.Split(e, "=")

					fmt.Println(env[0])
				}
			default:
				selected, err := output.Select(cmd.InOrStdin(), "Select a site to enable xdebug:", sites)
				if err != nil {
					return err
				}

				filter.Add("label", labels.Host+"="+sites[selected])

				// get the containers details
				details, err := docker.ContainerInspect(cmd.Context(), containers[0].ID)
				if err != nil {
					return err
				}

				// find the environment variable for xdebug
				for _, e := range details.Config.Env {
					env := strings.Split(e, "=")

					fmt.Println(env[0])
				}
			}

			// otherwise show the list of sites

			// get the "selected" site and get the container from the list

			// get the containers details docker.ContainerJSON

			// get the envs

			// find the xdebug environment variable

			// stop the container

			// remove the container

			// create the new container

			return nil
		},
	}

	// set flags for the command
	cmd.Flags().String("example", "example", "an example flag")

	return cmd
}
