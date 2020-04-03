package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var updateCommand = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "upgrade"},
	Short:   "Update a machine",
	Run: func(cmd *cobra.Command, args []string) {
		name := config.GetString("machine", flagMachineName)

		if err := nitro.Run(
			command.NewMultipassRunner("multipass"),
			nitro.Update(name),
		); err != nil {
			log.Fatal(err)
		}
	},
}
