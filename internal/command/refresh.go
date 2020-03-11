package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"

	"github.com/craftcms/nitro/internal"
)

type config struct {
	WriteFiles []struct {
		Path    string `yaml:"path"`
		Content string `yaml:"content"`
	} `yaml:"write_files"`
}

func Refresh(r internal.Runner) *cli.Command {
	return &cli.Command{
		Name:  "refresh",
		Usage: "",
		Action: func(c *cli.Context) error {
			return refreshAction(c, r)
		},
	}
}

func refreshAction(c *cli.Context, r internal.Runner) error {
	cfg := config{}
	// parse the cloudInit variable
	if err := yaml.Unmarshal([]byte(cloudInit), &cfg); err != nil {
		return err
	}

	// for each of those, override the path with the content
	for _, file := range cfg.WriteFiles {
		// skip the script that does the updates
		if file.Path == "/opt/nitro/refresh.sh" {
			continue
		}

		fmt.Println("SKIPPED: updating file", file.Path)
		//args := []string{"exec", c.String("machine"), "--", "sudo", "bash", "/opt/nitro/refresh.sh", file.Content, file.Path}
		//if err := r.Run(args); err != nil {
		//	return err
		//}
	}

	return nil
}
