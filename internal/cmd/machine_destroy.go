package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var destroyCommand = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		if flagPermanent {
			fmt.Println("Permanently destroying", name)
		} else {
			fmt.Println("Gently destroying", name)
		}

		destroyAction, err := nitro.Destroy(name, flagPermanent)
		if err != nil {
			return err
		}

		return nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{*destroyAction})
	},
}

func init() {
	destroyCommand.Flags().BoolVarP(&flagPermanent, "permanent", "p", false, "Permanently destroy the machine")
}
