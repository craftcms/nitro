package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/runas"
)

var (
	xCommand = &cobra.Command{
		Use:   "x",
		Short: "working on windows hosts",
		RunE: func(cmd *cobra.Command, args []string) error {

			args = append(args, "--machine=nitro-dev", "hosts")
			if err := runas.Elevated("nitro-dev", args); err != nil {
				return err
			}

			return hostsCommand.RunE(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(xCommand)
}
