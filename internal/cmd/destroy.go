package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"

	"github.com/pixelandtonic/prompt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/datetime"
	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/scripts"
)

var destroyCommand = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		p := prompt.NewPrompt()
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}
		script := scripts.New(mp, machine)

		// get the sites
		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			return err
		}

		var domains []string
		for _, site := range cfg.Sites {
			domains = append(domains, site.Hostname)
		}

		confirmed, err := p.Confirm("Are you sure you want to permanently destroy the "+machine+" machine", &prompt.InputOptions{
			Default:            "no",
			Validator:          nil,
			AppendQuestionMark: true,
		})
		if err != nil {
			return err
		}

		if !confirmed {
			return nil
		}

		// if we have any containers to backup, do so now
		if flagSkipBackup == false && len(cfg.Databases) != 0 {
			// backup the container
			for _, db := range cfg.Databases {
				container := db.Name()

				// run the script to get all databases
				var dbs []string
				switch db.Engine {
				case "postgres":
					if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerPostgresShowAllDatabases, container)); err == nil {
						sp := strings.Split(output, "\n")
						for i, d := range sp {
							d = strings.TrimSpace(d)
							if i == 0 || i == 1 || i == len(sp) || strings.Contains(d, "rows)") || d == "mysql" {
								continue
							}

							dbs = append(dbs, d)
						}
					}
				default:
					if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerMysqlShowAllDatabases, container)); err == nil {
						for _, db := range strings.Split(output, "\n") {
							// ignore the system defaults
							if db == "Database" || db == "information_schema" || db == "performance_schema" || db == "sys" || strings.Contains(db, "password on the command line") || db == "mysql" {
								continue
							}
							dbs = append(dbs, db)
						}
					}
				}

				if len(dbs) == 0 {
					fmt.Println(fmt.Sprintf("There are no databases in %s %s (port %s) to backup.", db.Engine, db.Version, db.Port))
					continue
				}

				var backupErrorMessage = "There was a problem backing up the databases.\nIf you wish to destroy " + machine + " without backups use --skip-backup."

				// backup each database
				for _, database := range dbs {
					var fullVmBackupPath string
					backupFileName := database + "-" + datetime.Parse(time.Now()) + ".sql"

					switch db.Engine {
					case "postgres":
						// create the backup directory if not found
						if output, err := script.Run(false, fmt.Sprintf(scripts.FmtCreateDirectory, "/home/ubuntu/.nitro/databases/postgres/backups/")); err != nil {
							fmt.Println(output)
							fmt.Println(err)
							fmt.Println(backupErrorMessage)
							return err
						}

						// run the backup
						fullVmBackupPath = "/home/ubuntu/.nitro/databases/postgres/backups/" + backupFileName
						if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerBackupIndividualPostgresDatabase, container, database, fullVmBackupPath)); err != nil {
							fmt.Println(output)
							fmt.Println(err)
							fmt.Println(backupErrorMessage)
							return err
						}
					default:
						// create the backup directory if not found
						if output, err := script.Run(false, fmt.Sprintf(scripts.FmtCreateDirectory, "/home/ubuntu/.nitro/databases/mysql/backups/")); err != nil {
							fmt.Println(output)
							fmt.Println(err)
							fmt.Println(backupErrorMessage)
							return err
						}

						fullVmBackupPath = "/home/ubuntu/.nitro/databases/mysql/backups/" + backupFileName
						if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerBackupAllMysqlDatabases, container, fullVmBackupPath)); err != nil {
							fmt.Println(output)
							fmt.Println(err)
							fmt.Println(backupErrorMessage)
							return err
						}
					}

					// show the status
					fmt.Println(fmt.Sprintf("Created backup %q, downloading...", backupFileName))

					home, err := homedir.Dir()
					if err != nil {
						return err
					}

					// make sure the backups folder exists
					backupsFolder := home + "/.nitro/backups/"
					if err := helpers.MkdirIfNotExists(backupsFolder); err != nil {
						fmt.Println(err)
						fmt.Println(backupErrorMessage)
						return err
					}

					// make sure the machine folder exists
					backupsFolder = backupsFolder + machine
					if err := helpers.MkdirIfNotExists(backupsFolder); err != nil {
						fmt.Println(err)
						fmt.Println(backupErrorMessage)
						return err
					}

					// create a container name
					backupsFolder = backupsFolder + "/" + container
					if err := helpers.MkdirIfNotExists(backupsFolder); err != nil {
						fmt.Println(err)
						fmt.Println(backupErrorMessage)
						return err
					}

					// transfer the folder into the host machine
					if err := nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{{Type: "transfer", Args: []string{"transfer", machine + ":" + fullVmBackupPath, backupsFolder}}}); err != nil {
						fmt.Println(err)
						fmt.Println(backupErrorMessage)
						return err
					}

					fmt.Println(fmt.Sprintf("Backup saved to %q", backupsFolder+"/"+backupFileName))
				}
			}
		}

		if flagDebug {
			fmt.Println("DEBUG: not removing the machine")
			return nil
		}

		destroyAction, err := nitro.Destroy(machine)
		if err != nil {
			return err
		}

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{*destroyAction}); err != nil {
			return err
		}

		if flagClean {
			if err := os.Remove(viper.ConfigFileUsed()); err != nil {
				fmt.Println("Unable to remove the config: ", viper.ConfigFileUsed())
			}
		}

		if len(domains) == 0 {
			fmt.Println("Permanently destroyed", machine)
			return nil
		}

		fmt.Println("Removing sites from your hosts file.")

		return hostsRemoveCommand.RunE(cmd, domains)
	},
}

func init() {
	destroyCommand.Flags().BoolVar(&flagClean, "clean", false, "Remove the config file when destroying the machine")
	destroyCommand.Flags().BoolVar(&flagSkipBackup, "skip-backup", false, "Skip database backups when destroying the machine")
}
