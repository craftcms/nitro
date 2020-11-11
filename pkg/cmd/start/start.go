package start

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/client"
	"github.com/spf13/cobra"
)

// StartCommand is the command for creating new development environments
var StartCommand = &cobra.Command{
	Use:   "start",
	Short: "Start an environment",
	RunE:  start,
	Example: `  # start containers for the default environment
  nitro start`,
}

func start(cmd *cobra.Command, args []string) error {
	env := cmd.Flag("environment").Value.String()

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	return nitro.Start(cmd.Context(), env, args)
}
