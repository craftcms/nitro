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

	"github.com/craftcms/nitro/datetime"
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

		// create a list of containers
		output, err := script.Run(false, scripts.DockerListContainerNames)
		if err != nil {
			return err
		}
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
		dbs := []string{"all-dbs"}
		switch strings.Contains(container, "mysql") {
		case false:
			if output, err := script.Run(false, fmt.Sprintf(`docker exec -i %s psql --username nitro --command "SELECT datname FROM pg_database WHERE datistemplate = false;"`, container)); err == nil {
				sp := strings.Split(output, "\n")
				for i, d := range sp {
					if i == 0 || i == 1 || i == len(sp) || strings.Contains(d, "rows)") {
						continue
					}

					dbs = append(dbs, strings.TrimSpace(d))
				}
			}
		default:
			if output, err := script.Run(false, fmt.Sprintf(`docker exec -i %s mysql -unitro -e "SHOW DATABASES;"`, container)); err != nil {
				for _, db := range strings.Split(output, "\n") {
					// ignore the system defaults
					if db == "Database" || db == "information_schema" || db == "performance_schema" || db == "sys" {
						continue
					}
					dbs = append(dbs, db)
				}
			}
		}

		if len(dbs) == 0 {
			return errors.New("no databases to backup in " + container)
		}

		database, _, err := p.Select("Which database should we backup", dbs, &prompt.SelectOptions{Default: 1})
		if err != nil {
			return err
		}

		// task
		var fullVmBackupPath string
		backupFileName := database + "-" + datetime.Parse(time.Now()) + ".sql"
		switch strings.Contains(container, "mysql") {
		case true:
			fullVmBackupPath = "/home/ubuntu/.nitro/databases/mysql/backups/" + backupFileName

			// if its everything, back them all up
			if database == "all-dbs" {
				if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerBackupAllMysqlDatabases, container, fullVmBackupPath)); err != nil {
					fmt.Println(output)
					return err
				}
			} else {
				// backup a specific database
				if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerBackupIndividualMysqlDatabase, container, database, fullVmBackupPath)); err != nil {
					fmt.Println(output)
					return err
				}
			}
		default:
			fullVmBackupPath = "/home/ubuntu/.nitro/databases/postgres/backups/" + backupFileName

			// if its all the databases
			if database == "all-dbs" {
				if output, err := script.Run(false, fmt.Sprintf(`docker exec -i %s pg_dumpall -U nitro > %s`, container, fullVmBackupPath)); err != nil {
					fmt.Println(output)
					return err
				}
			} else {
				// backup a specific database
				if output, err := script.Run(false, fmt.Sprintf(`docker exec -i %s pg_dump -U nitro %s > %s`, container, database, fullVmBackupPath)); err != nil {
					fmt.Println(output)
					return err
				}
			}
		}

		fmt.Println(fmt.Sprintf("Created backup %q, downloading...", backupFileName))

		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		// make sure the backups folder exists
		backupsFolder := home + "/.nitro/backups/"
		if err := helpers.MkdirIfNotExists(backupsFolder); err != nil {
			return err
		}

		// make sure the machine folder exists
		backupsFolder = backupsFolder + machine
		if err := helpers.MkdirIfNotExists(backupsFolder); err != nil {
			return err
		}

		// create a container name
		backupsFolder = backupsFolder + "/" + container
		if err := helpers.MkdirIfNotExists(backupsFolder); err != nil {
			return err
		}

		// transfer the folder into the host machine
		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{{Type: "transfer", Args: []string{"transfer", machine + ":" + fullVmBackupPath, backupsFolder}}}); err != nil {
			return err
		}

		_, err = script.Run(false, fmt.Sprintf(`rm %s`, fullVmBackupPath))
		if err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("Backup completed and stored in %q", backupsFolder+"/"+backupFileName))
		// end action

		return nil
	},
}
