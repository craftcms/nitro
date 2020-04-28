package cmd

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/find"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/prompt"
	"github.com/craftcms/nitro/internal/sudo"
	"github.com/craftcms/nitro/internal/task"
	"github.com/craftcms/nitro/validate"
)

var renameCommand = &cobra.Command{
	Use:   "rename",
	Short: "Rename a site",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		configSites := configFile.GetSites()

		if len(configSites) == 0 {
			return errors.New("there are no sites to rename")
		}

		// ask to select a existingSite
		i, _ := prompt.Select("Select existingSite to rename", configFile.SitesAsList())

		existingSite := configSites[i]

		// ask for the new newHostname
		var newHostname string
		hostnamePrompt := promptui.Prompt{
			Label:    fmt.Sprintf("What should the new newHostname be? [current: %s]", existingSite.Hostname),
			Validate: validate.Hostname,
		}

		hostnameEntered, err := hostnamePrompt.Run()
		if err != nil {
			return err
		}

		switch hostnameEntered {
		case "":
			newHostname = existingSite.Hostname
		default:
			newHostname = hostnameEntered
		}

		if existingSite.Hostname == newHostname {
			return errors.New("the new and original hostnames match, nothing to do")
		}

		path, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		php, err := find.PHPVersion(
			exec.Command(path, []string{"exec", machine, "--", "php", "--version"}...),
		)

		renamedSite := config.Site{Hostname: newHostname, Webroot: existingSite.Webroot}
		mount := configFile.FindMountBySiteWebroot(existingSite.Webroot)

		actions, err := task.Rename(machine, php, existingSite, renamedSite, mount)

		// update the config
		if err := configFile.RenameSite(existingSite, newHostname); err != nil {
			return err
		}

		// save the file
		if !flagDebug {
			if err := configFile.Save(viper.ConfigFileUsed()); err != nil {
				return err
			}
		}

		if flagDebug {
			for _, action := range actions {
				fmt.Println(action.Args)
			}

			return nil
		}

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("Ok, we renamed the existingSite %s to %s. We are now going to update the hosts file...", existingSite.Hostname, newHostname))

		nitro, err := exec.LookPath("nitro")
		if err != nil {
			return err
		}

		// remove the existingSite
		if err := sudo.RunCommand(nitro, machine, "hosts", "remove", existingSite.Hostname); err != nil {
			fmt.Println("Error removing", existingSite.Hostname, "from the hosts file")
			return err
		}

		if err := sudo.RunCommand(nitro, machine, "hosts"); err != nil {
			fmt.Println("Error adding", newHostname, "to the hosts file")
			return err
		}

		return nil
	},
}
