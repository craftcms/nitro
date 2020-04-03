package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command"
	"github.com/craftcms/nitro/internal/nitro"
)

var redisCommand = &cobra.Command{
	Use:   "redis",
	Short: "Enter a redis shell",
	PreRun: func(cmd *cobra.Command, args []string) {
		// set the defaults and load the yaml
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := nitro.Run(
			command.NewMultipassRunner("multipass"),
			nitro.Redis(flagMachineName),
		); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(redisCommand)
}
