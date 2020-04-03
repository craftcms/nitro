package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command"
	"github.com/craftcms/nitro/internal/nitro"
)

func init() {
	siteCommand.AddCommand(siteRemoveCommand)
	siteRemoveCommand.Flags().StringVar(&flagMachineName, "machine", "", "name of machine")
}

var siteRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove a site from machine",
	Run: func(cmd *cobra.Command, args []string) {
		if err := nitro.Run(
			command.NewMultipassRunner("multipass"),
			nitro.Empty(flagMachineName),
		); err != nil {
			log.Fatal(err)
		}
	},
}
