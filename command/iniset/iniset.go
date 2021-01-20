package iniset

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrUnknownSetting is used when an unknown service is requested
	ErrUnknownSetting = fmt.Errorf("unknown setting requested")
)

const exampleText = `  # change PHP settings for a site
  nitro iniset`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "iniset",
		Short:   "Change PHP setting",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// load the configuration
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro)

			// get all of the sites
			var sites, found []string
			for _, s := range cfg.Sites {
				p, _ := s.GetAbsPath(home)

				// check if the path matches a sites path, then we are in a known site
				if strings.Contains(wd, p) {
					found = append(found, s.Hostname)
				}

				// add the site to the list in case we cannot find the directory
				sites = append(sites, s.Hostname)
			}

			// if there are found sites we want to show or connect to the first one, otherwise prompt for
			// which site to connect to.
			switch len(found) {
			case 0:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", sites)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[selected])
			case 1:
				// add the label to get the site
				filter.Add("label", labels.Host+"="+found[0])
			default:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", found)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", labels.Host+"="+found[selected])
			}

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find an matching site")
			}

			// start the container if its not running
			if containers[0].State != "running" {
				if err := docker.ContainerStart(ctx, containers[0].ID, types.ContainerStartOptions{}); err != nil {
					return err
				}
			}

			hostname := strings.TrimLeft(containers[0].Names[0], "/")

			settings := []string{
				"display_errors",
				"max_execution_time",
				"max_input_vars",
				"max_input_time",
				"max_file_upload",
				"memory_limit",
				"opcache_enable",
				"opcache_revalidate_freq",
				"post_max_size",
				"upload_max_file_size",
			}

			// which setting to change
			selected, err := output.Select(cmd.InOrStdin(), "Which PHP setting would you like to change? ", settings)
			if err != nil {
				return err
			}

			// get the specific setting to change
			setting := settings[selected]

			// find the site by the selected hostname
			site, err := cfg.FindSiteByHostName(hostname)
			if err != nil {
				return err
			}

			// prompt the user for the setting to change
			switch setting {
			case "memory_limit":
				// TODO(jasonmccallister) create and add a validator
				answer, err := output.Ask("What should the new memeory limit be set to", "512M", ":", nil)
				if err != nil {
					return err
				}

				// change the value because its validated
				site.ChangePHPMemoryLimit(answer)
			default:
				return ErrUnknownSetting
			}

			// save the config file
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("unable to save config, %w", err)
			}

			// ask if the apply command should run
			apply, err := output.Confirm("Apply changes now", true, "?")
			if err != nil {
				return err
			}

			// if apply is false return nil
			if !apply {
				return nil
			}

			// run the apply command
			for _, c := range cmd.Parent().Commands() {
				// set the apply command
				if c.Use == "apply" {
					if err := c.RunE(c, args); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	return cmd
}
