package npm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/client"
)

var NPMCommand = &cobra.Command{
	Use:   "npm",
	Short: "Run npm actions",
	RunE:  npmMain,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveFilterDirs
	},
	Example: `  # run node install in a current directory
  nitro npm

  # updating a node project outside of the current directory
  nitro npm ./project-dir --version 14 --update`,
}

func npmMain(cmd *cobra.Command, args []string) error {
	var path string
	switch len(args) {
	case 0:
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to get the current directory, %w", err)
		}

		path, err = filepath.Abs(wd)
		if err != nil {
			return fmt.Errorf("unable to find the absolute path, %w", err)
		}
	default:
		var err error
		path, err = filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("unable to find the absolute path, %w", err)
		}
	}

	// determine the default action
	action := "install"
	if cmd.Flag("update").Value.String() == "true" {
		action = "update"
	}

	// get the full file path
	nodeFile := "package.json"
	var nodePath string
	switch action {
	case "install":
		nodePath = fmt.Sprintf("%s%c%s", path, os.PathSeparator, "package.json")
	default:
		nodeFile = "package-lock.json"
		nodePath = fmt.Sprintf("%s%c%s", path, os.PathSeparator, "package-lock.json")
	}

	// make sure the file exists
	fmt.Println("Checking for", nodeFile, "file in:")
	fmt.Println("  ==>", nodePath)
	_, err := os.Stat(nodePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("unable to locate a node file at %s", path)
	}

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	return nitro.Node(cmd.Context(), path, cmd.Flag("version").Value.String(), action)
}

func init() {
	flags := NPMCommand.Flags()

	// set the flags for this command
	flags.BoolP("update", "u", false, "run node update instead of install")
	flags.StringP("version", "v", "14", "which node version to use")
}
