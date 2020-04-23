package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/diff"
	"github.com/craftcms/nitro/internal/find"
	"github.com/craftcms/nitro/internal/nitro"
)

var applyCommand = &cobra.Command{
	Use:    "apply",
	Short:  "Apply changes from config",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		path, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		c := exec.Command(path, []string{"info", machine, "--format=csv"}...)
		output, err := c.Output()
		if err != nil {
			return err
		}

		attachedMounts, err := find.Mounts(machine, output)
		if err != nil {
			return err
		}

		// load the config file
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		// get abs path for file sources
		var fileMounts []config.Mount
		for _, m := range configFile.Mounts {
			fileMounts = append(fileMounts, config.Mount{Source: m.AbsSourcePath(), Dest: m.Dest})
		}

		// find sites not created
		var sitesToCreate []config.Site
		for _, site := range configFile.Sites {
			c := exec.Command(path, "exec", machine, "--", "sudo", "bash", "/opt/nitro/scripts/site-exists.sh", site.Hostname)
			output, err := c.Output()
			if err != nil {
				return err
			}
			if !strings.Contains(string(output), "exists") {
				sitesToCreate = append(sitesToCreate, site)
			}
		}

		// check for new dbs
		dbsToCreate, err := find.ContainersToCreate(machine, configFile)
		if err != nil {
			return err
		}

		// prompt?
		var actions []nitro.Action

		mountActions, err := diff.MountActions(machine, attachedMounts, fileMounts)
		if err != nil {
			return err
		}
		actions = append(actions, mountActions...)

		// create site actions
		for _, site := range sitesToCreate {
			// TODO abstract this logic into a func that takes mountActions and sites to return the mount action
			for _, ma := range mountActions {
				// break the string
				mnt := strings.Split(ma.Args[2], ":")

				// if the webroot is not of the mounts, then we should create an action
				if !strings.Contains(mnt[1], site.Webroot) {
					m := configFile.FindMountBySiteWebroot(site.Webroot)
					mountAction, err := nitro.MountDir(machine, m.AbsSourcePath(), m.Dest)
					if err != nil {
						return err
					}
					actions = append(actions, *mountAction)
				}
			}

			copyTemplateAction, err := nitro.CopyNginxTemplate(machine, site.Hostname)
			if err != nil {
				return err
			}
			actions = append(actions, *copyTemplateAction)

			// copy the nginx template
			changeNginxVariablesAction, err := nitro.ChangeTemplateVariables(machine, site.Webroot, site.Hostname, configFile.PHP, site.Aliases)
			if err != nil {
				return err
			}
			actions = append(actions, *changeNginxVariablesAction...)

			createSymlinkAction, err := nitro.CreateSiteSymllink(machine, site.Hostname)
			if err != nil {
				return err
			}
			actions = append(actions, *createSymlinkAction)
		}

		if len(sitesToCreate) > 0 {
			restartNginxAction, err := nitro.NginxReload(machine)
			if err != nil {
				return err
			}
			actions = append(actions, *restartNginxAction)
		}

		// create database actions
		for _, database := range dbsToCreate {
			volumeAction, err := nitro.CreateDatabaseVolume(machine, database.Engine, database.Version, database.Port)
			if err != nil {
				return err
			}
			actions = append(actions, *volumeAction)

			createDatabaseAction, err := nitro.CreateDatabaseContainer(machine, database.Engine, database.Version, database.Port)
			if err != nil {
				return err
			}
			actions = append(actions, *createDatabaseAction)
		}

		if flagDebug {
			for _, a := range actions {
				fmt.Println(a.Args)
			}

			return nil
		}

		fmt.Printf("There are %d mounted directories and %d mounts in the config file. Applying changes now...\n", len(attachedMounts), len(fileMounts))
		fmt.Printf("There are %d sites to create and %d sites in the config file. Applying changes now...\n", len(sitesToCreate), len(configFile.Sites))

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			return err
		}

		fmt.Println("Applied changes from", viper.ConfigFileUsed())

		return nil
	},
}
