package xon

import (
	"os"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const exampleText = `  # example command
  nitro xon`

// NewCommand returns the command that is used to enable xdebug for a sepcific site. It will first check
// if the current working directory or prompt the user for a site.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "xon",
		Short:   "Enable xdebug for a site",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()

			// load the config
			cfg, err := config.Load(home, env)
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

			// enable xdebug for the sites hostname
			if err := cfg.EnableXdebug(site); err != nil {
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
