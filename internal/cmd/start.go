package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var startCommand = &cobra.Command{
	Use:     "start",
	Aliases: []string{"up"},
	Short:   "Start a machine",
	Run: func(cmd *cobra.Command, args []string) {
		name := config.GetString("machine", flagMachineName)

		if err := nitro.Run(
			nitro.NewMultipassRunner("multipass"),
			nitro.Start(name),
		); err != nil {
			log.Fatal(err)
		}
	},
}
