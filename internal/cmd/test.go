package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/scripts"
)

var testCommand = &cobra.Command{
	Use:   "test",
	Short: "Testing",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		mp, err := exec.LookPath("multipass")
		if err != nil {
			fmt.Println("error with executable")
			return err
		}

		output, err := scripts.Run(mp, []string{"exec", machine, "--", `bash`, `-c`, `if test -f 'test'; then echo 'exists'; fi`})
		if err != nil {
			return err
		}

		fmt.Println(output)

		return nil

		//var actions []nitro.Action

		//complexAction := nitro.Action{
		//	Type:       "exec",
		//	UseSyscall: false,
		//	Args:       []string{"exec", machine, "--", `bash`, `-c`, `if test -f 'test'; then echo 'exists'; fi`},
		//}

		//actions = append(actions, complexAction)

		//if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
		//	return err
		//}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCommand)
}
