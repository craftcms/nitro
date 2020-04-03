package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagMachineName string

	rootCmd = &cobra.Command{
		Use:   "nitro",
		Short: "Local Craft CMS on tap",
		Long:  `TODO add the long description`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			if len(args) == 0 {
				cmd.Help()
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&flagMachineName, "machine", "", "name of machine")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
