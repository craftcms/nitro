package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var (
	flagPermanent bool

	destroyCommand = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy a machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := config.GetString("machine", flagMachineName)

			if flagPermanent {
				fmt.Println("permanently deleting", name)
			} else {
				fmt.Println("soft deleting", name)
			}

			if err := nitro.Run(
				nitro.NewMultipassRunner("multipass"),
				nitro.Destroy(name, flagPermanent),
			); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	destroyCommand.Flags().BoolVarP(&flagPermanent, "permanent", "p", false, "permanently destroy the machine")
}
