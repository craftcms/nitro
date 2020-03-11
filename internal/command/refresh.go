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

	for _, file := range cfg.WriteFiles {
		// for each of those, override the path with the content
		// TODO this works for now, but need to determine how to pass the content without executing the file.Content
		e := fmt.Sprintf(`echo "content" | sudo tee %s`, file.Path)
		args := []string{"exec", c.String("machine"), "--", "bash", "-c", e}
		// fmt.Println(args)
		if err := r.Run(args); err != nil {
			return err
		}
	}

	return nil
}
