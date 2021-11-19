package alias

import (
	"fmt"
	"os"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/flags"
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
			// load the configuration
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			var appName string
			switch flags.AppName == "" {
			case true:
				wd, err := os.Getwd()
				if err != nil {
					return err
				}

				appName, err = appaware.Detect(*cfg, wd)
				if err != nil {
					return err
				}
			default:
				appName = flags.AppName
			}

			// find the app by the hostname
			app, err := cfg.FindAppByHostname(appName)
			if err != nil {
				return err
			}

			// prompt the user to add new alias
			v := validate.MultipleHostnameValidator{}
			alias, err := output.Ask("Enter the alias domain for the app (use commas to enter multiple)", "", ":", &v)
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

				// set the alias for the app
				if err := cfg.SetAppAliases(app.Hostname, a); err != nil {
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
