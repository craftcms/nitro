package cmd

import (
	"errors"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/internal/editor"
)

var xCommand = &cobra.Command{
	Use:    "x",
	Hidden: true,
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
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(xCommand)
}
