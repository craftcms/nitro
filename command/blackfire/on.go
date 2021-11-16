package blackfire

import (
	"fmt"
	"os"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const onTest = `  # enable blackfire for an app
  nitro blackfire on`

func onCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "on",
		Short:   "Enables Blackfire for an app.",
		Example: onTest,
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

			// ensure the blackfire credentials are set
			if cfg.Blackfire.ClientID == "" {
				// ask for the server token
				token, err := output.Ask("Enter your Blackfire Client ID", "", ":", nil)
				if err != nil {
					return err
				}

				cfg.Blackfire.ClientID = token

				// save the config file
				if err := cfg.Save(); err != nil {
					return fmt.Errorf("unable to save config, %w", err)
				}
			}

			// ensure the blackfire credentials are set
			if cfg.Blackfire.ClientToken == "" {
				// ask for the server token
				token, err := output.Ask("Enter your Blackfire Client Token", "", ":", nil)
				if err != nil {
					return err
				}

				cfg.Blackfire.ClientToken = token

				// save the config file
				if err := cfg.Save(); err != nil {
					return fmt.Errorf("unable to save config, %w", err)
				}
			}

			var hostname string
			switch flags.AppName == "" {
			case false:
				hostname = flags.AppName
			default:
				// get the current working directory
				wd, err := os.Getwd()
				if err != nil {
					return err
				}

				hostname, err = appaware.Detect(*cfg, wd)
				if err != nil {
					return err
				}
			}

			// find the app by the hostname
			app, err := cfg.FindAppByHostname(hostname)
			if err != nil {
				return err
			}

			// if xdebug is set, we need to disable it to profile the app
			if app.Xdebug {
				// disable xdebug for the app hostname
				if err := cfg.DisableXdebug(app.Hostname); err != nil {
					return err
				}
			}

			// disable blackfire for the app hostname
			if err := cfg.EnableBlackfire(app.Hostname); err != nil {
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
