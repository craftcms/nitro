package exec

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/client"
	"github.com/spf13/cobra"
)

var ExecCommand = &cobra.Command{
	Use:   "exec",
	Short: "Access container environment",
	RunE:  execCommand,
	Args:  cobra.MinimumNArgs(1),
	Example: `  # get access to a container
  nitro exec example.nitro`,
}

func execCommand(cmd *cobra.Command, args []string) error {
	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	content, err := nitro.Exec(cmd.Context(), args[0], []string{"ls", "-la"})
	if err != nil {
		return fmt.Errorf("unable to exec command, %w", err)
	}

	fmt.Println(string(content))

	return nil
}
