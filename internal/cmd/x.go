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

var xCommand = &cobra.Command{
	Use:    "x",
	Short:  "Examine machine and config file",
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

		if len(attachedMounts) == 0 {
			fmt.Println("ok, there are no already mounted directories")
		}
		if len(attachedMounts) == 1 {
			fmt.Println("ok, we found 1 already mounted directory")
		}
		if len(attachedMounts) > 1 {
			fmt.Printf("ok, we found %d already mounted directories\n", len(attachedMounts))
		}

		// prompt?
		var actions []nitro.Action

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

		mountActions, err := hack.MountDiffActions(name, attachedMounts, fileMounts)
		if err != nil {
			return err
		}
		actions = append(actions, mountActions...)

		for _, a := range actions {
			fmt.Println(a.Args)
		}

		return nil
	},
}
