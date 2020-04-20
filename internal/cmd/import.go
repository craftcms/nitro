package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var importCommand = &cobra.Command{
	Use:   "import",
	Short: "Import database into machine",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		// verify the file exists
		if !fileExists(args[0]) {
			return errors.New(fmt.Sprintf("unable to located the file %q to import", args[0]))
		}

		// which database engine?
		var databases []config.Database
		if err := viper.UnmarshalKey("databases", &databases); err != nil {
			return err
		}
		var dbs []string
		for _, db := range databases {
			dbs = append(dbs, fmt.Sprintf("%s_%s_%s", db.Engine, db.Version, db.Port))
		}
		databaseContainerName := promptui.Select{
			Label: "Select database",
			Items: dbs,
		}

		_, containerName, err := databaseContainerName.Run()
		if err != nil {
			return err
		}

		var actions []nitro.Action

		// syntax is strange, see this issue: https://github.com/canonical/multipass/issues/1165#issuecomment-548763143
		transferAction := nitro.Action{
			Type:       "transfer",
			UseSyscall: false,
			Args:       []string{"transfer", args[0], name + ":" + args[0]},
		}
		actions = append(actions, transferAction)

		engine := "mysql"
		if strings.Contains(containerName, "postgres") {
			engine = "postgres"
		}

		importArgs := []string{"exec", name, "--", "bash", "/opt/nitro/scripts/docker-exec-import.sh", containerName, "nitro", args[0], engine}
		dockerExecAction := nitro.Action{
			Type:       "exec",
			UseSyscall: false,
			Args:       importArgs,
		}
		actions = append(actions, dockerExecAction)

		fmt.Printf("Importing %q into %q (large files may take a while)...\n", args[0], containerName)

		return nitro.Run(nitro.NewMultipassRunner("multipass"), actions)
	},
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dirExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsExist(err) {
		return false
	}
	return info.IsDir()
}
