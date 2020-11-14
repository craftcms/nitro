package destroy

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/client"
)

// DestroyCommand is the command for creating new development environments
var DestroyCommand = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy an environment",
	RunE:  destroy,
	Example: `  # remove all resources (networks, containers, and volumes) for an environment
  nitro destroy

  # destroy resources for a specific environment
  nitro destroy --environment my-testing-environment`,
}

func destroy(cmd *cobra.Command, args []string) error {
	env := cmd.Flag("environment").Value.String()

	fmt.Print("Are you sure? this will remove all containers, volumes, and networks (Y/n) ")

	// prompt the user for confirmation
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return fmt.Errorf("unable to provide a prompt, %w", err)
	}

	var confirm bool
	resp := strings.TrimSpace(response)
	for _, answer := range []string{"y", "Y", "yes", "Yes", "YES"} {
		if resp == answer {
			confirm = true
		}
	}

	if !confirm {
		fmt.Println("Skipping destroy action")
		fmt.Println("  ==> all resources related to", env, "will remain")

		return nil
	}

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	return nitro.Destroy(cmd.Context(), env, args)
}
