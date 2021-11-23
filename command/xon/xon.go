package xon

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
)

const exampleText = `  # example command
  nitro xon`

// NewCommand returns the command that is used to enable xdebug for a specific app. It will first check
// if the current working directory or prompt the user for an app.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "xon",
		Short:   "Enables Xdebug for an app.",
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

			output.Info("Disabling xdebug for", app.Hostname)

			// php 7.0 does not support xdebug
			if app.PHPVersion == "7.0" {
				return fmt.Errorf("xdebug with PHP 7.0 is not supported")
			}

			// if blackfire is set, we need to disable it to profile the app
			if app.Blackfire {
				// disable blackfire for the app hostname
				if err := cfg.DisableBlackfire(app.Hostname); err != nil {
					return err
				}
			}

			// enable xdebug for the app hostname
			if err := cfg.EnableXdebug(app.Hostname); err != nil {
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
