package blackfire

import (
	"fmt"
	"os"
	"strings"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const offExampleText = `  # disable blackfire for a site
  nitro blackfire off`

func offCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "off",
		Short:   "Disables Blackfire for a site.",
		Example: offExampleText,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			cfg, err := config.Load(home, false)
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}

			var options []string
			for _, s := range cfg.Sites {
				options = append(options, s.Hostname)
			}

			return options, cobra.ShellCompDirectiveDefault
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.VerifyInit(cmd, args, home, output)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, false, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			cfg, err := config.Load(home, false)
			if err != nil {
				return err
			}

			// ensure the blackfire credentials are set
			if cfg.Blackfire.ServerID == "" {
				// ask for the server id
				id, err := output.Ask("Enter your Blackfire Server ID", "", ":", nil)
				if err != nil {
					return err
				}

				cfg.Blackfire.ServerID = id

				// save the config file
				if err := cfg.Save(); err != nil {
					return fmt.Errorf("unable to save config, %w", err)
				}
			}

			// ensure the blackfire credentials are set
			if cfg.Blackfire.ServerToken == "" {
				// ask for the server token
				token, err := output.Ask("Enter your Blackfire Server Token", "", ":", nil)
				if err != nil {
					return err
				}

				cfg.Blackfire.ServerToken = token

				// save the config file
				if err := cfg.Save(); err != nil {
					return fmt.Errorf("unable to save config, %w", err)
				}
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

			// did they ask for a specific site?
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
					output.Info("Disabling Blackfire for", sites[0].Hostname)

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

			// if xdebug is set, we need to disable it to profile the site
			if site.Xdebug {
				// disable xdebug for the sites hostname
				if err := cfg.DisableXdebug(site.Hostname); err != nil {
					return err
				}
			}

			// disable blackfire for the sites hostname
			if err := cfg.DisableBlackfire(site.Hostname); err != nil {
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
