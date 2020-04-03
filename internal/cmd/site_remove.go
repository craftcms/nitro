package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var siteRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove a site from machine",
	Run: func(cmd *cobra.Command, args []string) {
		name := config.GetString("machine", flagMachineName)

		if err := nitro.Run(
			command.NewMultipassRunner("multipass"),
			nitro.Empty(name),
		); err != nil {
			log.Fatal(err)
		}
	},
}
