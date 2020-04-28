package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"
	"gopkg.in/yaml.v3"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/task"
)

var removeCommand = &cobra.Command{
	Use:   "remove",
	Short: "Manage your nitro sites",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

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

		selectedSite, err := ui.Select("Select a site to remove:", configFile.SitesAsList(), &input.Options{
			Required: true,
		})
		if err != nil {
			return err
		}

		var site config.Site
		for i, s := range sites {
			if s.Hostname == selectedSite {
				site = sites[i]
			}
		}
		if site.Hostname == "" {
			return errors.New("error selecting a site")
		}

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

		actions, err := task.Remove(machine, *mount, site)
		if err != nil {
			return err
		}

		applyChanges := false
		answer, err := ui.Ask("Apply changes from config now?", &input.Options{
			Default:  "y",
			Required: true,
			Loop:     true,
		})
		if err != nil {
			return err
		}

		if strings.ContainsAny(answer, "y") {
			applyChanges = true
		}

		// save the config
		if flagDebug {
			if applyChanges {
				fmt.Println("Ok, applying changes from the config file...")
			}
			for _, a := range actions {
				fmt.Println(a.Args)
			}
			return nil
		}

		if applyChanges {
			fmt.Println("Ok, applying changes from the config file...")
			return applyCommand.RunE(cmd, args)
		}

		return nil
	},
}
