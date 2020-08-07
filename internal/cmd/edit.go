package cmd

import (
	"errors"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/editor"
)

var editCommand = &cobra.Command{
	Use:   "edit",
	Short: "Edit config",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgFile := viper.ConfigFileUsed()
		if cfgFile == "" {
			return errors.New("unable to find the config file")
		}

		filePath, err := filepath.Abs(cfgFile)
		if err != nil {
			return err
		}

		_, err = editor.CaptureInputFromEditor(filePath, editor.GetPreferredEditorFromEnvironment)

		return err
	},
}
