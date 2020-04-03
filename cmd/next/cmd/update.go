package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command"
	"github.com/craftcms/nitro/internal/nitro"
)

func init() {
	rootCmd.AddCommand(updateCommand)
	updateCommand.Flags().StringVar(&flagMachineName, "machine", "", "name of machine")
}

var updateCommand = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "upgrade"},
	Short:   "Update a machine",
	PreRun: func(cmd *cobra.Command, args []string) {
		// set the defaults and load the yaml
		// TODO validate options for php and etc
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := nitro.Run(
			command.NewMultipassRunner("multipass"),
			nitro.Update(flagMachineName),
		); err != nil {
			log.Fatal(err)
		}
	},
}
