package restart

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/client"
	"github.com/spf13/cobra"
)

// RestartCommand is the command for creating new development environments
var RestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "Restart an environment",
	RunE:  restartMain,
}

func restartMain(cmd *cobra.Command, args []string) error {
	env := cmd.Flag("environment").Value.String()

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	return nitro.Restart(cmd.Context(), env, args)
}
