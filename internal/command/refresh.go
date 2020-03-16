package command

import (
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type config struct {
	WriteFiles []struct {
		Path    string `yaml:"path"`
		Content string `yaml:"content"`
	} `yaml:"write_files"`
}

func Refresh(r Runner) *cli.Command {
	return &cli.Command{
		Name:  "refresh",
		Usage: "",
		Action: func(c *cli.Context) error {
			return refreshAction(c, r)
		},
	}
}

func refreshAction(c *cli.Context, r Runner) error {
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

		if err := r.SetInput(file.Content); err != nil {
			return err
		}

		args := []string{"transfer", "-", c.String("machine") + ":/home/ubuntu" + file.Path}
		if err := r.Run(args); err != nil {
			return err
		}
	}

	return nil
}
