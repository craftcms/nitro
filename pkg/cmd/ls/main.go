package ls

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/client"
	"github.com/spf13/cobra"
)

// LSCommand is the command for creating new development environments
var LSCommand = &cobra.Command{
	Use:   "ls",
	Short: "List environment containers",
	RunE:  lsMain,
}

func lsMain(cmd *cobra.Command, args []string) error {
	env := cmd.Flag("environment").Value.String()

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	return nitro.LS(cmd.Context(), env, args)
}
