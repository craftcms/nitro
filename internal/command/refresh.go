package command

import (
	"fmt"

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
		Usage: "Update scripts on machine",
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

	// make the files in /tmp/opt/
	for _, dir := range []string{"/tmp/opt/nitro/php", "/tmp/opt/nitro/mariadb", "/tmp/opt/nitro/postgres", "/tmp/opt/nitro/nginx"} {
		if err := r.Run([]string{"exec", c.String("machine"), "--", "mkdir", "-p", dir}); err != nil {
			return err
		}
	}
	fmt.Println("created temp directories")

	// grant the ubuntu user and group ownership on /tmp/opt/
	if err := r.Run([]string{"exec", c.String("machine"), "--", "sudo", "chown", "-R", "ubuntu:ubuntu", "/tmp/opt/nitro"}); err != nil {
		return err
	}
	fmt.Println("change ownership to ubuntu:ubuntu")

	// for each of those, override the path with the content
	for _, file := range cfg.WriteFiles {
		// create a tmp file to avoid this issue: https://github.com/canonical/multipass/issues/1434
		if err := r.Run([]string{"exec", c.String("machine"), "--", "touch", "/tmp" + file.Path}); err != nil {
			return err
		}

		if err := r.SetInput(file.Content); err != nil {
			return err
		}

		if err := r.Run([]string{"transfer", "-", c.String("machine") + ":/tmp" + file.Path}); err != nil {
			return err
		}
		fmt.Println("transferred file", file.Path)
	}

	// copy the tmp files
	fmt.Println("moving files...")
	if err := r.Run([]string{"exec", c.String("machine"), "--", "sudo", "cp", "-r", "/tmp/opt/nitro", "/opt/nitro"}); err != nil {
		return err
	}
	fmt.Println("refresh complete!")

	return nil
}
