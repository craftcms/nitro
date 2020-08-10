package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var supportCommand = &cobra.Command{
	Use:   "support",
	Short: "Get support",
	RunE: func(cmd *cobra.Command, args []string) error {
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		output, err := exec.Command(mp, "version").CombinedOutput()
		if err != nil {
			return err
		}

		sp := strings.Split(string(output), "\n")
		mpVersion := strings.Split(sp[0], "  ")

		url := "https://github.com/craftcms/nitro/issues/new?labels=bug&body=" + fmt.Sprintf(`### Description

### Steps to reproduce

1.
2.

### Additional info

- Nitro version: %s
- Multipass version: %s
- Host OS: %s
`, Version, mpVersion[1], runtime.GOOS)

		if err := browser.OpenURL(url); err != nil {
			fmt.Println("Failed to open browser, please use this URL to create a new support ticket:", url)
			return err
		}


		return nil
	},
}
