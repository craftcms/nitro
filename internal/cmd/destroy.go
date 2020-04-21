package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var destroyCommand = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)
		if len(args) > 0 {
			name = args[0]
		}

		destroyAction, err := nitro.Destroy(name)
		if err != nil {
			return err
		}

		return nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{*destroyAction})
	},
}
