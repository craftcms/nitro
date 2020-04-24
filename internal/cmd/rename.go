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

		// ask to select a site
		i, _ := prompt.Select("Select site to rename", configFile.SitesAsList())

		site := configSites[i]

		// ask for the new hostname
		var hostname string
		hostnamePrompt := promptui.Prompt{
			Label:    fmt.Sprintf("What should the new hostname be? [current: %s]", site.Hostname),
			Validate: validate.Hostname,
		}

		hostnameEntered, err := hostnamePrompt.Run()
		if err != nil {
			return err
		}

		switch hostnameEntered {
		case "":
			hostname = site.Hostname
		default:
			hostname = hostnameEntered
		}

		path, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		php, err := find.PHPVersion(
			exec.Command(path, []string{"exec", machine, "--", "php", "--version"}...),
		)

		actions, err := task.Rename(machine, php, site)

		if site.Hostname == hostname {
			return errors.New("the new and original hostnames match, nothing to do")
		}

		// update the config
		if err := configFile.RenameSite(site, hostname); err != nil {
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

		fmt.Println(fmt.Sprintf("Ok, we renamed the site %s to %s. We are now going to update the hosts file...", site.Hostname, hostname))

		nitro, err := exec.LookPath("nitro")
		if err != nil {
			return err
		}

		// remove the site
		if err := sudo.RunCommand(nitro, machine, "hosts", "remove", site.Hostname); err != nil {
			fmt.Println("Error removing", site.Hostname, "from the hosts file")
			return err
		}

		if err := sudo.RunCommand(nitro, machine, "hosts"); err != nil {
			fmt.Println("Error adding", hostname, "to the hosts file")
			return err
		}

		return nil
	},
}
