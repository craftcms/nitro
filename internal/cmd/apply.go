package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/hack"
	"github.com/craftcms/nitro/internal/nitro"
)

var applyCommand = &cobra.Command{
	Use:    "apply",
	Short:  "Apply changes from nitro.yaml",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := "nitro-dev"
		if flagMachineName != "" {
			machine = flagMachineName
		}

		path, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		c := exec.Command(path, []string{"info", machine, "--format=csv"}...)
		output, err := c.Output()
		if err != nil {
			return err
		}

		attachedMounts, err := hack.FindMounts(machine, output)
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

		fmt.Printf("There are %d mounted directories and %d mounts in the config file. Applying changes now...\n", len(attachedMounts), len(fileMounts))

		// prompt?
		var actions []nitro.Action

		mountActions, err := hack.MountDiffActions(machine, attachedMounts, fileMounts)
		if err != nil {
			return err
		}
		actions = append(actions, mountActions...)

		// create site actions
		for _, site := range sitesToCreate {
			m := configFile.FindMountBySiteWebroot(site.Webroot)
			mountAction, err := nitro.MountDir(machine, m.AbsSourcePath(), m.Dest)
			if err != nil {
				return err
			}
			actions = append(actions, *mountAction)

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

		if flagDebug {
			for _, a := range actions {
				fmt.Println(a.Args)
			}

			return nil
		}

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			return err
		}

		fmt.Println("Applied changes from", viper.ConfigFileUsed())

		return nil
	},
}
