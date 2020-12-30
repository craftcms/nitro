package edit

import (
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/editor"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # edit the config file
  nitro edit`

// NewCommand returns the command to edit a config file with the users default editor as defined by the
// $EDITOR variable.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit",
		Short:   "Edit the nitro config",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			_, err = editor.CaptureInputFromEditor(cfg.GetFile(), editor.GetPreferredEditorFromEnvironment)

			return err
		},
	}

	return cmd
}
