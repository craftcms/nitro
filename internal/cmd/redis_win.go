// +build !linux, !darwin, windows

package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var redisCommand = &cobra.Command{
	Use:   "redis",
	Short: "Enter a redis shell",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		c := exec.Command(mp, "exec", machine, "--", "redis-cli")
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr

		return c.Run()
	},
}
