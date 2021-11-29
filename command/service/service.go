package service

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrUnknownService is used when an unknown service is requested
	ErrUnknownService = fmt.Errorf("unknown service requested")
)

const exampleText = `  # enable blackfire
  nitro service enable blackfire`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service",
		Aliases: []string{"svc"},
		Short:   "Manage services.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		enableCommand(home, docker, output),
		disableCommand(home, docker, output),
		lsCommand(home, docker, output),
	)

	return cmd
}