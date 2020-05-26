package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/nitro"
)

var redisCommand = &cobra.Command{
	Use:   "redis",
	Short: "Enter Redis",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		redisAction, err := nitro.Redis(machine)
		if err != nil {
			return err
		}

		return nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{*redisAction})
	},
}
