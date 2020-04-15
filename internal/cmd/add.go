package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/helpers"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Add sites and mounts to machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		var wd string
		if len(args) > 0 {
			wd = args[0]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			wd = cwd
		}

		parent, err := helpers.ParentPathName(wd)
		if err != nil {
			return err
		}

		hostnamePrompt := promptui.Prompt{
			Label: fmt.Sprintf("what should the hostname be? [%s]", parent),
		}

		hostname, err := hostnamePrompt.Run()
		if err != nil {
			return err
		}
		if hostname == "" {
			hostname = parent
		}

		// mountPath := "/nitro/sites/" + parent

		foundDir, err := helpers.FindWebRoot(wd)
		if err != nil {
			return err
		}
		webRootPrompt := promptui.Prompt{
			Label: fmt.Sprintf("where is the webroot? [%s]", foundDir),
		}

		webroot, err := webRootPrompt.Run()
		if err != nil {
			return err
		}
		if webroot == "" {
			webroot = foundDir
		}

		// TODO create the mount action
		fmt.Printf("%s has been added to nitro.yaml", hostname)
		applyPrompt := promptui.Prompt{
			Label: "apply nitro.yaml changes now? [y]",
		}

		apply, err := applyPrompt.Run()
		if err != nil {
			return err
		}

		if apply != "y" {
			fmt.Println("ok, you can apply new nitro.yaml changes later by running `nitro apply`.")
		}

		return nil
	},
}
