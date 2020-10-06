package cmd

import "github.com/spf13/cobra"

var inisetCommand = &cobra.Command{
	Use:       "iniset",
	Short:     "Change PHP settings",
	ValidArgs: []string{"display_errors", "max_execution_time", "max_input_vars", "max_input_time", "upload_max_filesize", "max_file_uploads", "memory_limit"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	inisetCommand.AddCommand(
		inisetMaxExecutionTimeCommand,
		inisetMaxInputVarsCommand,
		inisetUploadMaxFilesizeCommand,
		inisetMaxInputTimeCommand,
		inisetMaxFileUploadsCommand,
		inisetMemoryLimitCommand,
		inisetDisplayErrorsCommand,
	)
}
