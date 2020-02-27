package command

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/nitro/internal/action"
)

var (
	bootstrapFlag = &cli.BoolFlag{
		Name:  "bootstrap",
		Usage: "Bootstrap the machine with defaults",
		Value: true,
	}
	phpVersionFlag = &cli.StringFlag{
		Name:        "php-version",
		Usage:       "Provide version of PHP",
		DefaultText: "7.4",
	}
)

func Initialize() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize a new machine",
		Action: func(c *cli.Context) error {
			return action.Initialize(c)
		},
		After: func(c *cli.Context) error {
			// if we are bootstrapping, call the command
			if c.Bool("bootstrap") {
				return c.App.RunContext(c.Context, []string{c.App.Name, "--machine", c.String("machine"), "bootstrap"})
			}

			return nil
		},
		Flags: []cli.Flag{bootstrapFlag},
	}
}

func Bootstrap(executor action.CommandLineExecutor) *cli.Command {
	return &cli.Command{
		Name:  "bootstrap",
		Usage: "Bootstrap the installation of a new machine",
		Action: func(context *cli.Context) error {
			return action.Bootstrap(context, executor)
		},
		Flags: []cli.Flag{phpVersionFlag},
	}
}

func Update() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update a machine with the latest software",
		Action: func(c *cli.Context) error {
			return action.Update(c)
		},
	}
}

func SSH(e action.CommandLineExecutor) *cli.Command {
	return &cli.Command{
		Name:  "ssh",
		Usage: "SSH into a machine as administrator",
		Action: func(c *cli.Context) error {
			return action.SSH(c.String("machine"), e)
		},
	}
}

func AddHost(e action.CommandLineExecutor) *cli.Command {
	return &cli.Command{
		Name:  "add-host",
		Usage: "Add a new virtual host to a machine",
		Action: func(context *cli.Context) error {
			return action.AddHost(context, e)
		},
		Before: func(context *cli.Context) error {
			if context.NArg() != 1 {
				return errors.New("you must pass a domain name")
			}

			// TODO validate the domain name with validate.Domain(d)

			return nil
		},
		Flags: []cli.Flag{phpVersionFlag},
	}
}

func Logs(e action.CommandLineExecutor) *cli.Command {
	return &cli.Command{
		Name:        "logs",
		Description: "Show a machines logs",
		Action: func(context *cli.Context) error {
			return cli.ShowSubcommandHelp(context)
		},
		Subcommands: []*cli.Command{
			{
				Name:        "nginx",
				Description: "Show logs from nginx",
				Action: func(c *cli.Context) error {
					return action.LogsNginx(c, e)
				},
			},
		},
	}
}
func Stop() *cli.Command {
	return &cli.Command{
		Name:        "stop",
		Usage:       "Stop a machine",
		Description: "Stop a machine when not in use (this does not delete the machine)",
		Action: func(c *cli.Context) error {
			return action.Stop(c)
		},
	}
}

func Delete() *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Usage:       "Delete a machine",
		Description: "Delete a machine when no longer needed",
		Action: func(c *cli.Context) error {
			rdr := bufio.NewReader(os.Stdin)
			fmt.Print("WARNING: Are you sure you wish to perform this task (y/N)? ")

			answer, err := rdr.ReadString(' ')
			if err != nil {
				return err
			}

			if strings.Contains(answer, "n") {
				fmt.Println("Skipping the deletion of", c.String("machine"))
				return nil
			}

			return action.Delete(c)
		},
	}
}
