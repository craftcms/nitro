package xon

import (
	"fmt"
	"os"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/terminal"
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

			// check each of the sites for a match
			var currentSite string
			var sites []string
			for _, site := range cfg.Sites {
				// get the path
				path, _ := site.GetAbsPath(home)

				// see if the sites path matches the current directory
				if strings.Contains(wd, path) {
					currentSite = site.Hostname
					break
				}

				// add the site as an option
				sites = append(sites, site.Hostname)
			}

			if currentSite == "" {
				// show all of the sites to the user
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", sites)
				if err != nil {
					return err
				}

				currentSite = sites[selected]
			}

			return fmt.Errorf("use config.EnableXdebug(site) and call apply")
		},
	}

	return cmd
}
