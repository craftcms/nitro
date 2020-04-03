package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var redisCommand = &cobra.Command{
	Use:   "redis",
	Short: "Enter a redis shell",
	Run: func(cmd *cobra.Command, args []string) {
		name := config.GetString("machine", flagMachineName)

		if err := nitro.Run(
			nitro.NewMultipassRunner("multipass"),
			nitro.Redis(name),
		); err != nil {
			log.Fatal(err)
		}
	},
}
