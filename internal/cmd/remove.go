package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/prompt"
)

var removeCommand = &cobra.Command{
	Use:   "remove",
	Short: "Manage your nitro sites",
	RunE: func(cmd *cobra.Command, args []string) error {
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}
		name := configFile.Name

		sites := configFile.GetSites()
		i, _ := prompt.Select("Select site to remove", configFile.SitesAsList())

		site := sites[i]
		if !prompt.Verify(fmt.Sprintf("this will remove %s from %s, do you want to apply the changes? [y]", site.Hostname, name)) {
			fmt.Println("ok, you can apply new nitro.yaml changes later by running `nitro apply`.")

			return nil
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

		if !flagDebug {
			if err := configFile.Save(viper.ConfigFileUsed()); err != nil {
				return err
			}
		}

		// remove the mount
		if err := configFile.RemoveMountBySiteWebroot(site.Webroot); err != nil {
			return err
		}

		if !flagDebug {
			if err := configFile.Save(viper.ConfigFileUsed()); err != nil {
				return err
			}
		}

		actions, err := removeActions(name, *mount, site)
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
			fmt.Println("failed to remove the site:", err)
			return err
		}

		fmt.Println("ok, we removed the site from your nitro.yaml and made the changes")

		return nil
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
