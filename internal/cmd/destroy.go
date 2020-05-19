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
	"github.com/craftcms/nitro/internal/sudo"
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

		confirmed, err := p.Confirm("Are you sure you want to permanently destroy your "+machine+" machine", &prompt.InputOptions{
			Default:   "no",
			Validator: nil,
		})
		if err != nil {
			return err
		}

		if !confirmed {
			return nil
		}

		// TODO backup all databases in the

		var containers []string
		for _, db := range cfg.Databases {
			containers = append(containers, db.Name())
		}

		// if we have any containers to backup, do so now
		if len(containers) != 0 {
			// backup the container
			for _, container := range containers {
				var fullVmBackupPath string
				backupFileName := "all-dbs-" + datetime.Parse(time.Now()) + ".sql"

				switch strings.Contains(container, "mysql") {
				case false:
					fullVmBackupPath = "/home/ubuntu/.nitro/databases/postgres/backups/" + backupFileName
					if output, err := script.Run(false, fmt.Sprintf(`docker exec -i %s pg_dumpall -U nitro > %s`, container, fullVmBackupPath)); err != nil {
						fmt.Println(output)
						return err
					}
				default:
					fullVmBackupPath = "/home/ubuntu/.nitro/databases/mysql/backups/" + backupFileName
					if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerBackupAllMysqlDatabases, container, fullVmBackupPath)); err != nil {
						fmt.Println(output)
						return err
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

				fmt.Println(fmt.Sprintf("Backup completed and stored in %q", backupsFolder+backupFileName))
			}
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
				fmt.Println("unable to remove the config:", viper.ConfigFileUsed())
			}
		}

		if len(domains) == 0 {
			fmt.Println("Permanently destroyed", machine)
			return nil
		}

		cmds := []string{"hosts", "remove"}
		for _, domain := range domains {
			cmds = append(cmds, domain)
		}

		// prompt to remove hosts file
		nitro, err := exec.LookPath("nitro")
		if err != nil {
			return err
		}

		fmt.Println("Removing sites from your hosts file")

		return sudo.RunCommand(nitro, machine, cmds...)
	},
}

func init() {
	destroyCommand.Flags().BoolVar(&flagClean, "clean", false, "remove the config file when destroying the machine")
}
