package cmd

import (
	"github.com/craftcms/nitro/internal/runas"
	"github.com/spf13/cobra"
)

var (
	xCommand = &cobra.Command{
		Use:   "x",
		Short: "working on windows hosts",
		RunE: func(cmd *cobra.Command, args []string) error {

			args = append(args, "--machine=nitro-dev", "hosts")
			if err := runas.Elevated(args); err != nil {
				return err
			}

			return hostsCommand.RunE(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(xCommand)
}
