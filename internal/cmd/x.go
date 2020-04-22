package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/find"
)

var xCommand = &cobra.Command{
	Use:    "x",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := "nitro-dev"
		if flagMachineName != "" {
			machine = flagMachineName
		}

		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			return err
		}

		dbs, err := find.ContainersToCreate(machine, cfg)
		if err != nil {
			return err
		}

		for _, db := range dbs {
			fmt.Println(db.Engine)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(xCommand)
}
