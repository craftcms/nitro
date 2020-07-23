package cmd

import "github.com/spf13/cobra"

var inisetCommand = &cobra.Command{
	Use:       "iniset",
	Short:     "Change php.ini",
	ValidArgs: []string{"max_execution_time"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	inisetCommand.AddCommand(inisetMaxExecutionTimeCommand)
}
