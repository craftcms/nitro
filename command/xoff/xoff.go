package xoff

import (
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # example command
  nitro xoff`

// NewCommand returns the command that is used to disable xdebug for a specific site. It will first check
// if the current working directory or prompt the user for a site.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "xoff",
		Short:   "Disable xdebug for a site",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// check each of the sites for a match
			var site string
			var sites []string
			for _, s := range cfg.Sites {
				// get the path
				path, _ := s.GetAbsPath(home)

				// see if the sites path matches the current directory
				if strings.Contains(wd, path) {
					site = s.Hostname
					break
				}

				// add the site as an option
				sites = append(sites, s.Hostname)
			}

			// if its not the current site
			if site == "" {
				// show all of the sites to the user
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", sites)
				if err != nil {
					return err
				}

				site = sites[selected]
			}

			// diable xdebug for the sites hostname
			if err := cfg.DisableXdebug(site); err != nil {
				return err
			}

			// save the config
			if err := cfg.Save(); err != nil {
				return err
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
