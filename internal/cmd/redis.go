package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var redisCommand = &cobra.Command{
	Use:   "redis",
	Short: "Enter a redis shell",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("machine", flagMachineName)

		return nitro.Run(nitro.NewMultipassRunner("multipass"), nitro.Redis(name))
	},
}
