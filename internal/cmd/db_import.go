package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/normalize"
	"github.com/craftcms/nitro/internal/scripts"
)

var dbImportCommand = &cobra.Command{
	Use:   "import my-backup.sql",
	Short: "Import database",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		// get the filename
		filename, fileAbsPath, err := normalize.Path(args[0], home)
		if err != nil {
			return err
		}

		// make sure the file exists
		if !helpers.FileExists(fileAbsPath) {
			return errors.New(fmt.Sprintf("Unable to locate the file %q", fileAbsPath))
		}

		// which database engine?
		var databases []config.Database
		if err := viper.UnmarshalKey("databases", &databases); err != nil {
			return err
		}
		var dbs []string
		for _, db := range databases {
			dbs = append(dbs, db.Name())
		}

		if len(dbs) == 0 {
			return errors.New("there are no databases that we can import the file into")
		}

		p := prompt.NewPrompt()

		containerName, _, err := p.Select("Which database engine to import the backup", dbs, &prompt.SelectOptions{Default: 1})

		databaseName, err := p.Ask("What is the database name to create for the import", &prompt.InputOptions{Default: "", Validator: nil})
		if err != nil {
			return err
		}

		var actions []nitro.Action

		// syntax is strange, see this issue: https://github.com/canonical/multipass/issues/1165#issuecomment-548763143
		fileFullPath := "/home/ubuntu/.nitro/databases/imports/" + filename
		transferAction := nitro.Action{
			Type:       "transfer",
			UseSyscall: false,
			Args:       []string{"transfer", fileAbsPath, machine + ":" + fileFullPath},
		}
		actions = append(actions, transferAction)

		fmt.Printf("Uploading %q into %q (large files may take a while)...\n", filename, machine)

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			return err
		}

		// run the import scripts

		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		script := scripts.New(mp, machine)

		if strings.Contains(containerName, "mysql") {
			_, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerMysqlCreateDatabaseIfNotExists, containerName, databaseName))
			if err != nil {
				return err
			}
			fmt.Println("Created database", databaseName)

			_, err = script.Run(false, fmt.Sprintf(scripts.FmtDockerMysqlImportDatabase, fileFullPath, containerName, databaseName))
			if err != nil {
				return err
			}

			fmt.Println("Successfully imported the database backup into", containerName)

			return nil
		}

		fmt.Println("Successfully imported the database backup into", containerName)

		return nil
	},
}
