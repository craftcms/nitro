package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/scripts"
)

var dbBackupCommand = &cobra.Command{
	Use:   "backup",
	Short: "Backup a database",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}
		p := prompt.NewPrompt()
		script := scripts.New(mp, machine)

		output, err := script.Run(false, scripts.DockerListContainerNames)
		if err != nil {
			return err
		}

		// create a list
		containers := strings.Split(output, "\n")
		if len(containers) == 0 {
			return errors.New("there are no databases we can add to")
		}

		// which database
		container, _, err := p.Select("Which database engine", containers, &prompt.SelectOptions{
			Default: 1,
		})
		if err != nil {
			return err
		}

		// get all of the databases from the container
		var dbs []string
		switch strings.Contains(container, "mysql") {
		case false:

		default:
			output, err := script.Run(false, fmt.Sprintf(`docker exec -i %s mysql -unitro -e "SHOW DATABASES;"`, container))
			if err != nil {
				return err
			}

			for _, db := range strings.Split(output, "\n") {
				// ignore the system defaults
				if db == "Database" || db == "information_schema" || db == "performance_schema" || db == "sys" {
					continue
				}
				dbs = append(dbs, db)
			}
		}

		if len(dbs) == 0 {
			return errors.New("no databases to backup for " + container)
		}

		// append the all option
		dbs = append(dbs, "all-databases")

		database, _, err := p.Select("Which database should we backup", dbs, &prompt.SelectOptions{Default: len(dbs)})
		if err != nil {
			return err
		}

		datetime := time.Now().Format("01-01-2020")
		backupFile := container + "-" + database + "-" + datetime + ".sql"
		var fullBackupPath string
		switch strings.Contains(container, "mysql") {
		case true:
			fullBackupPath = "/home/ubuntu/.nitro/databases/mysql/backups/" + backupFile
			if database == "all-databases" {
				output, err := script.Run(false, fmt.Sprintf(`docker exec %s /usr/bin/mysqldump --all-databases -unitro > %s`, container, fullBackupPath))
				if err != nil {
					fmt.Println(output)
					return err
				}
			} else {
				output, err := script.Run(false, fmt.Sprintf(`docker exec %s /usr/bin/mysqldump -unitro %s > %s`, container, database, fullBackupPath))
				if err != nil {
					fmt.Println(output)
					return err
				}
			}
		default:
			fullBackupPath = "/home/ubuntu/.nitro/databases/postgres/backups/" + backupFile
			output, err := script.Run(false, fmt.Sprintf(`echo "missing commands for %s %s"`, container, fullBackupPath))
			if err != nil {
				fmt.Println(output)
				return err
			}
		}

		fmt.Println(fmt.Sprintf("Created backup %q, downloading...", backupFile))

		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		machineFolder := home + "/.nitro/" + machine
		if err := helpers.MkdirIfNotExists(machineFolder); err != nil {
			return err
		}

		backupsFolder := machineFolder + "/backups/"
		if err := helpers.MkdirIfNotExists(backupsFolder); err != nil {
			return err
		}

		// transfer the folder into the host machine
		action := nitro.Action{
			Type: "transfer",
			Args: []string{"transfer", machine + ":" + fullBackupPath, backupsFolder},
		}
		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{action}); err != nil {
			return err
		}

		_, err = script.Run(false, fmt.Sprintf(`rm %s`, fullBackupPath))
		if err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("Backup completed and stored in %q", backupsFolder+backupFile))

		return nil
	},
}
