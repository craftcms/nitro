package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/h2non/filetype"
	"github.com/mitchellh/go-homedir"
	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/internal/client"
	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/internal/database"
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
		return []string{"sql", "gz"}, cobra.ShellCompDirectiveFilterFileExt
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}
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

		// create the request
		req := &nitrod.ImportDatabaseRequest{}

		// check if the file is compressed
		if err := isCompressed(filename, req); err != nil {
			fmt.Println("Error checking if the file is compressed, error:", err.Error())
			return nil
		}

		// create a new prompt
		p := prompt.NewPrompt()

		// open the file
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		reader := bufio.NewReader(file)

		// try to determine the database engine
		detected, err := database.DetermineEngine(reader)
		if err != nil {
			fmt.Println("Unable to determine the database engine from the file", filename)
		}

		// get the databases as a list, limiting to the detected engine
		engines := configFile.DatabaseEnginesAsList(detected)
		if len(engines) == 0 {
			fmt.Println("Unable to get a list of the database engines")
			return nil
		}

		// if there is only on database engine
		var container string
		switch len(engines) {
		case 1:
			container = engines[0]
		default:
			container, _, err = p.Select("Select database engine:", engines, &prompt.SelectOptions{
				Default: 1,
			})
			if err != nil {
				return err
			}
		}
		req.Container = container

		// if the detect engine is mysql
		showCreatePrompt := true
		if detected == "mysql" {
			// check if there is a create database statement
			showCreatePrompt, err = database.HasCreateStatement(reader)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Printf("The file %q will create a database during import...\n", filename)
		}

		// prompt for the database name to create if needed
		if showCreatePrompt {
			db, err := p.Ask("Enter the database name to create for the import:", &prompt.InputOptions{Validator: nil})
			if err != nil {
				return err
			}
			req.Database = db
		}

		// create the stream
		stream, err := c.ImportDatabase(cmd.Context())
		if err != nil {
			fmt.Println("Error creating a connection to the nitro server, error:", err.Error())
			return nil
		}

		fmt.Printf("Uploading %q into %q (large files may take a while)...\n", filename, machine)

		start := time.Now()
		buffer := make([]byte, 1024*20)

		for {
			n, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			req.Data = buffer[:n]
			err = stream.Send(req)
			if err == io.EOF {
				return stream.CloseSend()
			}
			if err != nil {
				return err
			}
		}

		res, err := stream.CloseAndRecv()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		fmt.Println(res.Message+".", fmt.Sprintf("Import took %f seconds...", time.Since(start).Seconds()))

		return nil
	},
}

func isCompressed(file string, req *nitrod.ImportDatabaseRequest) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	kind, _ := filetype.Match(b)

	if filetype.IsArchive(b) {
		req.Compressed = true
		req.Extension = kind.Extension
	}

	return nil
}
