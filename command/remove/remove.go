package remove

import (
	"os"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # remove an app from the config
  nitro remove`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Removes an app.",
		Example: exampleText,
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

			// get the app
			appName := flags.AppName
			if appName == "" {
				// get the current working directory
				wd, err := os.Getwd()
				if err != nil {
					return err
				}

				appName, err = appaware.Detect(*cfg, wd)
				if err != nil {
					return err
				}
			}

			app, err := cfg.FindAppByHostname(appName)
			if err != nil {
				return err
			}

			output.Info("Removing", appName)

			// remove the app
			if err := cfg.RemoveApp(app); err != nil {
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
