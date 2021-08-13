package keys

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/keys"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrExample is used when we want to share an error
	ErrExample = fmt.Errorf("not yet implemented")
)

const exampleText = `  # keys command
  nitro keys`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Short:   "Adds SSH keys to a site container.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := filepath.Join(home, ".ssh")

			if _, err := os.Stat(path); os.IsNotExist(err) {
				return errors.New("unable to find directory " + path)
			}

			// find all of the keys
			keys, err := keys.Find(path)
			if err != nil {
				return err
			}

			// if there are no keys
			if len(keys) == 0 {
				fmt.Println("Unable to find keys to add")
				return nil
			}

			var opts []string
			for k, _ := range keys {
				opts = append(opts, k)
			}

			_, err = output.Select(os.Stdin, "Which key should we import?", opts)
			if err != nil {
				return err
			}

			return ErrExample
		},
	}

	return cmd
}
