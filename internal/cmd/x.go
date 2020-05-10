package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	xCommand = &cobra.Command{
		Use:   "x",
		Short: "working on windows hosts",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(runtime.GOOS)

			return nil

			// TODO follow the same "sudo" process but use the runas syntax
			// from https://www.windows-commandline.com/windows-runas-command-prompt/

			// runas /user:"Jason McCallister":administrator nitro.exe

			//args = append(args, "--machine=nitro-dev", "hosts")
			//if err := runas.Elevated("nitro-dev", args); err != nil {
			//	return err
			//}
			//
			//return hostsCommand.RunE(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(xCommand)
}
