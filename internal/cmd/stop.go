package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop a machine",
	Run: func(cmd *cobra.Command, args []string) {
		name := config.GetString("machine", flagMachineName)

		if err := nitro.Run(
			nitro.NewMultipassRunner("multipass"),
			nitro.Stop(name),
		); err != nil {
			log.Fatal(err)
		}
	},
}
