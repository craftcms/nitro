package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/scripts"
)

var (
	refreshCommand = &cobra.Command{
		Use:   "refresh",
		Short: "Refresh machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			machine := flagMachineName
			mp, err := exec.LookPath("multipass")
			if err != nil {
				return err
			}

			script := scripts.New(mp, machine)

			fmt.Println("Downloading the latest refresh script")

			_, err = script.Run(false, `wget https://raw.githubusercontent.com/craftcms/nitro/master/scripts/refresh.sh -O /tmp/refresh.sh`)
			if err != nil {
				return err
			}

			// run the script
			output, err := script.Run(true, `bash /tmp/refresh.sh `+Version)
			if err != nil {
				return err
			}

			fmt.Println(output)

			fmt.Println("Refreshed the templates and configs for the machine to", Version)

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(refreshCommand)
}
