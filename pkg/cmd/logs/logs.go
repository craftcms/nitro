package logs

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/client"
	"github.com/spf13/cobra"
)

// LogsCommand is used to retrieve logs from a container
var LogsCommand = &cobra.Command{
	Use:   "logs",
	Short: "View logs",
	RunE:  logsMain,
	Example: `  # list logs for a container
  nitro logs`,
}

func logsMain(cmd *cobra.Command, args []string) error {
	env := cmd.Flag("environment").Value.String()

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	containers, err := nitro.LS(cmd.Context(), env, args)
	if err != nil {
		return err
	}

	for name, id := range containers {
		fmt.Println(name, id)
	}

	return nil
}
