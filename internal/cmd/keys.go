package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var keysCommand = &cobra.Command{
	Use:   "keys",
	Short: "Import SSH keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		// machine := flagMachineName
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		sshDir := home + "/.ssh/"

		if _, err := os.Stat(sshDir); os.IsNotExist(err) {
			return errors.New("unable to find directory " + sshDir)
		}

		keys := make(map[string]string)
		if err := filepath.Walk(sshDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() || info.Name() == "known_hosts" || strings.Contains(info.Name(), ".pem") {
				return nil
			}

			// TODO create a map for the key (e.g. `is_rsa` or `personal_rsa`)
			// should be is_rsa = id_rsa.pub
			if keys[info.Name()] == "" {
				if strings.Contains(info.Name(), ".pub") {
					return nil
				}
				keys[info.Name()] = info.Name()
			}

			return nil
		}); err != nil {
			panic(err)
		}

		for _, key := range keys {
			fmt.Println(key)
		}

		return nil
	},
}
