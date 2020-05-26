package cmd

import (
	"errors"
	"fmt"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/validate"
)

var renameCommand = &cobra.Command{
	Use:   "rename",
	Short: "Rename site",
	RunE: func(cmd *cobra.Command, args []string) error {
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		sites := configFile.GetSites()

		if len(sites) == 0 {
			return errors.New("there are no sites to rename")
		}

		p := prompt.NewPrompt()

		// ask to select a site
		var site config.Site
		_, i, err := p.Select("Which site do you want to rename", configFile.SitesAsList(), &prompt.SelectOptions{
			Default: 1,
		})
		if err != nil {
			return err
		}
		site = sites[i]

		// ask for the new newHostname
		var newHostname string
		newHostname, err = p.Ask("What should the new hostname be", &prompt.InputOptions{
			Validator: validate.Hostname,
		})
		if err != nil {
			return err
		}
		if site.Hostname == newHostname {
			return errors.New("the new and original hostnames match, nothing to do")
		}

		// update the config
		if err := configFile.RenameSite(site, newHostname); err != nil {
			return err
		}

		// save the file
		if !flagDebug {
			if err := configFile.Save(viper.ConfigFileUsed()); err != nil {
				return err
			}
		}

		apply, err := p.Confirm("Apply changes from config", &prompt.InputOptions{Default: "yes"})
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
