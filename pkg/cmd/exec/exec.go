package exec

import (
	"fmt"
	"io/ioutil"

	"github.com/craftcms/nitro/pkg/client"
	"github.com/spf13/cobra"
)

var ExecCommand = &cobra.Command{
	Use:   "exec",
	Short: "Access a shell in a container",
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

	stream, err := nitro.Exec(cmd.Context(), args[0], []string{"ls", "-la"})
	if err != nil {
		return fmt.Errorf("unable to exec command, %w", err)
	}
	defer stream.Close()

	cert, err := ioutil.ReadAll(stream.Reader)
	if err != nil {
		return fmt.Errorf("unable to read response from exec, %w", err)
	}

	// if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, stream.Reader); err != nil {
	// 	return nil, fmt.Errorf("unable to copy the output of the container logs, %w", err)
	// }
	fmt.Println(cert)

	return nil
}
