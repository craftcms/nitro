package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/mitchellh/go-homedir"
	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/keys"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/scripts"
)

var keysCommand = &cobra.Command{
	Use:   "keys",
	Short: "Add keys to machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		p := prompt.NewPrompt()
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}
		script := scripts.New(mp, machine)
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		path := home + "/.ssh/"

		if _, err := os.Stat(path); os.IsNotExist(err) {
			return errors.New("unable to find directory " + path)
		}

		// find all of the keys
		keys, err := keys.Find(path)
		if err != nil {
			return err
		}

		// create the options to present
		var opts []string
		for k, v := range keys {
			if v == "" {
				continue
			}
			opts = append(opts, k)
		}

		// ask the user which key to select
		selected, _, err := p.Select(fmt.Sprintf("Select the key to add to %q", machine), opts, &prompt.SelectOptions{
			Default: 1,
		})
		if err != nil {
			return err
		}

		// create the actions to transfer the keys
		var actions []nitro.Action
		for k, v := range keys {
			if k == selected {
				// transfer the selected key to /home/ubuntu/.ssh/<file>
				actions = append(actions, nitro.Action{
					Type: "transfer",
					Args: []string{"transfer", path + k, machine + ":/home/ubuntu/.ssh/" + k},
				})

				actions = append(actions, nitro.Action{
					Type: "transfer",
					Args: []string{"transfer", path + v, machine + ":/home/ubuntu/.ssh/" + v},
				})
			}
		}

		if flagDebug {
			for _, action := range actions {
				fmt.Println(action.Args)
			}

			return nil
		}

		// run the actions
		if err := nitro.Run(nitro.NewMultipassRunner(mp), actions); err != nil {
			return err
		}

		// create the script to add to known hosts
		if output, err := script.Run(false, fmt.Sprintf(`cat /home/ubuntu/.ssh/%s >> /home/ubuntu/.ssh/authorized_keys`, keys[selected])); err != nil {
			fmt.Println(output)
			return err
		}

		fmt.Println(fmt.Sprintf("Transferred the key %q into %q.", selected, machine))

		return nil
	},
}
