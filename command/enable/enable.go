package enable

import (
	"os"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # enable the app in the current directory
  nitro enable

  # enable a specific app using the global flag
  nitro --app myapp.nitro enable`

// NewCommand returns the command to enable an app from automatically starting.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "enable",
		Short:   "Enables an app.",
		Example: exampleText,
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
			name := flags.AppName
			if name == "" {
				// get the current working directory
				wd, err := os.Getwd()
				if err != nil {
					return err
				}

				name, err = appaware.Detect(*cfg, wd)
				if err != nil {
					return err
				}
			}

			// enable the app
			if err := cfg.EnableApp(name); err != nil {
				return err
			}

			// save the config
			if err := cfg.Save(); err != nil {
				return err
			}

			output.Info("Enabled", name)

			return nil
		},
	}

	return cmd
}
