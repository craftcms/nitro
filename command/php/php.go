package php

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # run php version
  nitro php -v

  # view php info
  nitro php -i`

// NewCommand returns the php command which allows users to pass php specific commands to a sites
// container. Its context aware and will prompt the user for the site if its not in a directory.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "php",
		Short:              "Run PHP commands",
		Example:            exampleText,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
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
			filter.Add("label", labels.Nitro)

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			var site config.Site
			switch len(sites) {
			case 0:
				// prompt for the site
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// set the site we selected
				site = sites[selected]

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[selected].Hostname)
			case 1:
				output.Info("connecting to", sites[0].Hostname)

				// set the site we selected
				site = sites[0]

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[0].Hostname)
			default:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// set the site we selected
				site = sites[selected]

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[selected].Hostname)
			}

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
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

			// create the command for running the craft console
			cmds := []string{"exec", "-it", containers[0].ID}

			// get the container path
			path := site.GetContainerPath()
			if path != "" {
				cmds = append(cmds, fmt.Sprintf("%s/%s", path, "php"))
			} else {
				cmds = append(cmds, "php")
			}

			// append the provided args to the command
			if len(args) == 0 {
				cmds = append(cmds, "-v")
			} else {
				cmds = append(cmds, args...)
			}

			// find the docker executable
			cli, err := exec.LookPath("docker")
			if err != nil {
				return err
			}

			// create the command
			c := exec.Command(cli, cmds...)

			c.Stdin = cmd.InOrStdin()
			c.Stderr = cmd.ErrOrStderr()
			c.Stdout = cmd.OutOrStdout()

			if err := c.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
