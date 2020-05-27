package cmd

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/craftcms/nitro/config"
)

var removeCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove site",
	RunE: func(cmd *cobra.Command, args []string) error {
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		sites := configFile.GetSites()

		if len(sites) == 0 {
			return errors.New("there are no sites to remove")
		}

		p := prompt.NewPrompt()

		var site config.Site
		_, i, err := p.Select("Select a site to remove", configFile.SitesAsList(), &prompt.SelectOptions{
			Default: 1,
		})
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

		apply, err := p.Confirm("Apply changes from config now", &prompt.InputOptions{
			Default:   "yes",
			Validator: nil,
			AppendQuestionMark: true,
		})
		if err != nil {
			return err
		}

		if apply {
			fmt.Println("Applying changes from the config file...")
			return applyCommand.RunE(cmd, args)
		}

		return nil
	},
}
