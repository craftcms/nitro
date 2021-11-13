package alias

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/validate"
)

const exampleText = `  # add alias domains to a site
  nitro alias`

// NewCommand allows users to set aliases or subdomains on an existing site. Useful for multi-site configurations.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "alias",
		Short:   "Adds alias domains.",
		Example: exampleText,
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, false, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// load the configuration
			cfg, err := config.Load(home, false)
			if err != nil {
				return err
			}

			// check for a site arg
			var siteArg string
			if len(args) > 0 {
				siteArg = strings.TrimSpace(args[0])
			}

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			// did they ask for a specific site?
			var site *config.Site
			switch siteArg == "" {
			case true:
				switch len(options) {
				case 1:
					output.Info("adding aliases to", options[0])

					// add the label to get the site
					site, _ = cfg.FindSiteByHostName(options[0])
				default:
					// prompt for the site to alias
					selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
					if err != nil {
						return err
					}

					site, _ = cfg.FindSiteByHostName(options[selected])
				}

				// show aliases if they exist
				if len(site.Aliases) > 0 {
					output.Info("The following aliases are set for", site.Hostname)
					for _, a := range site.Aliases {
						output.Info("  ", a)
					}
				} else {
					output.Info("No existing aliases are set for", site.Hostname)
				}
			default:
				site, err = cfg.FindSiteByHostName(siteArg)
				if err != nil {
					return err
				}
			}

			// prompt the user to add new alias
			v := validate.MultipleHostnameValidator{}
			alias, err := output.Ask("Enter the alias domain for the site (use commas to enter multiple)", "", ":", &v)
			if err != nil {
				return err
			}

			parts, err := v.Parse(alias)
			if err != nil {
				return err
			}

			if len(parts) > 1 {
				output.Info("Adding aliases:")
			} else {
				output.Info("Adding alias:")
			}
			for _, a := range parts {
				output.Info("  ", a)

				// set the alias
				if err := cfg.SetSiteAlias(site.Hostname, a); err != nil {
					return err
				}
			}

			// save the config file
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("unable to save config, %w", err)
			}

			return nil
		},
	}

	return cmd
}
