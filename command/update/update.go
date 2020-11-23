package update

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// New returns the update command for updating images on the
// local machine
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update Docker images",
		Example: `  # update docker images
  nitro update`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), strings.Join(args, " "))
			return nil
		},
	}

	// set the flags

	return cmd
}
