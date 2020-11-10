package stop

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/client"
	"github.com/spf13/cobra"
)

// StopCommand is the command for creating new development environments
var StopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop environment",
	RunE:  stopMain,
	Example: `  # stop containers for the default environment
  nitro stop`,
}

func stopMain(cmd *cobra.Command, args []string) error {
	env := cmd.Flag("environment").Value.String()

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	return nitro.Stop(cmd.Context(), env, args)
}
