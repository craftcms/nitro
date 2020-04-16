package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/hack"
)

var (
	Version = "0.0.0"

	versionCommand = &cobra.Command{
		Use:   "version",
		Short: "View Nitro version",
		RunE: func(cmd *cobra.Command, args []string) error {
			// show current version
			fmt.Printf("nitro %s\n", Version)
			fmt.Println("")

			latest, err := hack.GetLatestVersion(http.DefaultClient, "https://api.github.com/repos/pixelandtonic/nitro/releases")
			if err != nil {
				return err
			}

			fmt.Println("The latest version of nitro is", latest)
			fmt.Println("Visit https://github.com/craftcms/nitro for more details or \nrun `nitro self-update` to perform an upgrade")

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCommand)
}
