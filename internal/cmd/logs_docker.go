package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var logsDockerCommand = &cobra.Command{
	Use:   "docker",
	Short: "Show docker logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented yet")
	},
}
