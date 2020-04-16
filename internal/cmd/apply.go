package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/hack"
	"github.com/craftcms/nitro/internal/nitro"
)

var applyCommand = &cobra.Command{
	Use:    "apply",
	Short:  "Apply changes from config",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)
		path, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		c := exec.Command(path, []string{"info", name, "--format=csv"}...)
		output, err := c.Output()
		if err != nil {
			return err
		}

		attachedMounts, err := hack.FindMounts(name, output)
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

		// TODO find sites not created

		fmt.Printf("ok, there are %d mounted directories and %d mounts in the config file. Applying changes now...\n", len(attachedMounts), len(fileMounts))

		// prompt?
		var actions []nitro.Action

		mountActions, err := hack.MountDiffActions(name, attachedMounts, fileMounts)
		if err != nil {
			return err
		}
		actions = append(actions, mountActions...)

		if flagDebug {
			for _, a := range actions {
				fmt.Println(a.Args)
			}

			return nil
		}

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			return err
		}

		fmt.Println("applied changes from", viper.ConfigFileUsed())

		return nil
	},
}
