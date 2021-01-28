package remove

import (
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # remove a site from the config
  nitro remove`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a site",
		Example: exampleText,
		PostRunE: func(cmd *cobra.Command, args []string) error {
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
					return c.RunE(c, args)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// get all of the sites
			var sites []string
			for _, s := range cfg.Sites {
				// add the site to the list
				sites = append(sites, s.Hostname)
			}

			// prompt for the site to remove
			selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", sites)
			if err != nil {
				return err
			}

			site, err := cfg.FindSiteByHostName(sites[selected])
			if err != nil {
				return err
			}

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
