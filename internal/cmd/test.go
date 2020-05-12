package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/nitro"
)

var testCommand = &cobra.Command{
	Use:   "test",
	Short: "Testing",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		var actions []nitro.Action

		complexAction := nitro.Action{
			Type:       "exec",
			UseSyscall: false,
			Args:       []string{"exec", machine, "--", `bash`, `-c`, `if test -f 'test'; then echo 'exists'; fi`},
		}

		c := strings.Join(complexAction.Args, " ")

		fmt.Println(c)

		actions = append(actions, complexAction)

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCommand)
}
