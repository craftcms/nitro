package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/internal/client"
	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/nitrod"
	"github.com/craftcms/nitro/internal/normalize"
)

var dbImportCommand = &cobra.Command{
	Use:   "import my-backup.sql",
	Short: "Import database",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"sql"}, cobra.ShellCompDirectiveFilterFileExt
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		c, err := client.NewClient(ip, "50051")
		if err != nil {
			return err
		}

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
			return errors.New(fmt.Sprintf("Unable to locate the file %q.", fileAbsPath))
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

		// if there is only one
		var containerName string
		switch len(dbs) {
		case 1:
			containerName = dbs[0]
		default:
			containerName, _, err = p.Select("Select database engine", dbs, &prompt.SelectOptions{
				Default: 1,
			})
			if err != nil {
				return err
			}
		}

		databaseName, err := p.Ask("Enter the database name to create for the import", &prompt.InputOptions{Validator: nil})
		if err != nil {
			return err
		}

		f, err := os.Open(filename)
		if err != nil {
			return err
		}

		// create the stream
		stream, err := c.ImportDatabase(cmd.Context())
		if err != nil {
			return err
		}

		fmt.Printf("Uploading %q into %q (large files may take a while)...\n", filename, machine)

		rdr := bufio.NewReader(f)
		buf := make([]byte, 1024)

		for {
			n, err := rdr.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			req := &nitrod.ImportDatabaseRequest{
				Container: containerName,
				Database:  databaseName,
				Data:    buf[:n],
			}

			err = stream.Send(req)
			if err == io.EOF {
				fmt.Println("eof on client, breaking")
				break
			}
			if err != nil {
				return err
			}
		}

		res, err := stream.CloseAndRecv()
		if err != nil {
			return err
		}

		fmt.Println(res.Message)

		return nil
	},
}
