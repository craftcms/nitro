package blackfire

import (
	"os"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const offText = `  # disable blackfire for a site
  nitro blackfire off`

func offCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "off",
		Short:   "Disables Blackfire for a site.",
		Example: offText,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.VerifyInit(cmd, args, home, output)
		},
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

			var site config.Site
			switch len(sites) {
			case 0:
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				site = sites[selected]
			case 1:
				output.Info("Disabling Blackfire for", sites[0].Hostname)

				site = sites[0]
			default:
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				site = sites[selected]
			}

			// if xdebug is set, we need to disable it to profile the site
			if site.Xdebug {
				// disable blackfire for the sites hostname
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
