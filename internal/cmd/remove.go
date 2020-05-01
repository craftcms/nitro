package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/pixelandtonic/go-input"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/prompt"
)

var removeCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove site from a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		sites := configFile.GetSites()

		if len(sites) == 0 {
			return errors.New("there are no sites to remove")
		}

		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		var site config.Site
		_, i, err := prompt.Select(ui, "Select a site to remove:", sites[0].Hostname, configFile.SitesAsList())
		if err != nil {
			return err
		}
		site = sites[i]

		// find the mount
		mount := configFile.FindMountBySiteWebroot(site.Webroot)
		if mount == nil {
			return errors.New("unable to find an associated mount")
		}

		// remove site
		if err := configFile.RemoveSite(site.Hostname); err != nil {
			return err
		}

		// remove the mount
		if err := configFile.RemoveMountBySiteWebroot(site.Webroot); err != nil {
			return err
		}

		// START HACK
		// use viper to ensure consistency when saving
		c, err := yaml.Marshal(configFile)
		if err := viper.ReadConfig(bytes.NewBuffer(c)); err != nil {
			return err
		}

		if !flagDebug {
			if err := viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
				return err
			}
		}

		// unmarshal the messy config into a config
		var messyConfig config.Config
		if err := viper.Unmarshal(&messyConfig); err != nil {
			return err
		}

		if !flagDebug {
			// save that config in the right order
			if err := messyConfig.Save(viper.ConfigFileUsed()); err != nil {
				return err
			}
		}
		// END HACK

		applyChanges, err := prompt.Verify(ui, "Apply changes from config now?", "y")
		if err != nil {
			return err
		}

		if applyChanges {
			fmt.Println("Applying changes from the config file...")
			return applyCommand.RunE(cmd, args)
		}

		return nil
	},
}
