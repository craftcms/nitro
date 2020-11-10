package destroy

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/client"
	"github.com/spf13/cobra"
)

// DestroyCommand is the command for creating new development environments
var DestroyCommand = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy an environment",
	RunE:  destroy,
}

func destroy(cmd *cobra.Command, args []string) error {
	env := cmd.Flag("environment").Value.String()

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	return nitro.Destroy(cmd.Context(), env, args)
}
