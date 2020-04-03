package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var (
	servicesCommand = &cobra.Command{
		Use:   "services",
		Short: "Start, stop, or restart services on machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := config.GetString("machine", flagMachineName)

			if err := nitro.Run(
				nitro.NewMultipassRunner("multipass"),
				nitro.Empty(name),
			); err != nil {
				return err
			}

			return nil
		},
	}
)
