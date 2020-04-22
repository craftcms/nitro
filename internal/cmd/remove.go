package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/prompt"
	"github.com/craftcms/nitro/internal/sudo"
)

var removeCommand = &cobra.Command{
	Use:   "remove",
	Short: "Manage your nitro sites",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := "nitro-dev"
		if flagMachineName != "" {
			machine = flagMachineName
		}

		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		sites := configFile.GetSites()

		if len(sites) == 0 {
			return errors.New("there are no sites to remove")
		}

		i, _ := prompt.Select("Select site to remove", configFile.SitesAsList())

		site := sites[i]

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
		if err := viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
			return err
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

		actions, err := removeActions(machine, *mount, site)
		if err != nil {
			return err
		}

		// save the config
		if flagDebug {
			for _, a := range actions {
				fmt.Println(a.Args)
			}
			return nil
		}

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			fmt.Println("Failed to remove the site:", err)
			return err
		}

		fmt.Println("Removed the site from your config and applied the changes.")

		// prompt to remove hosts file
		nitro, err := exec.LookPath("nitro")
		if err != nil {
			return err
		}

		fmt.Println("Removing site from your hosts file")

		return sudo.RunCommand(nitro, machine, "hosts", "remove", site.Hostname)
	},
}

func removeActions(name string, mount config.Mount, site config.Site) ([]nitro.Action, error) {
	var actions []nitro.Action

	// unmount
	unmountAction, err := nitro.UnmountDir(name, mount.Dest)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *unmountAction)

	// remove nginx symlink
	removeSymlinkAction, err := nitro.RemoveSymlink(name, site.Hostname)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *removeSymlinkAction)

	restartNginxAction, err := nitro.NginxReload(name)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *restartNginxAction)

	return actions, nil
}
