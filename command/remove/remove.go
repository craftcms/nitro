package remove

import (
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # remove a site from the config
  nitro remove`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Removes a site.",
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

			return options, cobra.ShellCompDirectiveDefault
		},
		Aliases: []string{"rm"},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, false, output)
		},
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

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			var siteArg string
			if len(args) > 0 {
				siteArg = strings.TrimSpace(args[0])
			}

			var site *config.Site
			switch siteArg == "" {
			case true:
				switch len(sites) {
				case 0:
					selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
					if err != nil {
						return err
					}

					site = &sites[selected]
				case 1:
					site = &sites[0]
				default:
					selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
					if err != nil {
						return err
					}

					site = &sites[selected]
				}
			default:
				site, err = cfg.FindSiteByHostName(siteArg)
				if err != nil {
					return err
				}
			}

			output.Info("Removing", site.Hostname)

			// remove the site
			if err := cfg.RemoveSite(site); err != nil {
				return err
			}

			// save the config
			if err := cfg.Save(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
