// +build !linux, !darwin, windows

package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var sshCommand = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		c := exec.Command(mp, "shell", machine)
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr

		return c.Run()
	},
}
