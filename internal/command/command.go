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
		Name:    "bootstrap",
		Aliases: []string{"b"},
		Usage:   "Bootstrap the machine with installation defaults",
		Value:   true,
	}
	phpVersionFlag = &cli.StringFlag{
		Name:        "php-version",
		Aliases:     []string{"php", "version"},
		Usage:       "Provide version of PHP",
		DefaultText: "7.4",
	}
)

func Initialize() *cli.Command {
	return &cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "Initialize a new machine",
		Action: func(c *cli.Context) error {
			return action.Initialize(c)
		},
		Flags: []cli.Flag{bootstrapFlag},
	}
}

func Bootstrap(executor action.CommandLineExecutor) *cli.Command {
	return &cli.Command{
		Name:    "bootstrap",
		Aliases: []string{"b", "boot"},
		Usage:   "Bootstrap the installation of a new machine",
		Action: func(context *cli.Context) error {
			return action.Bootstrap(context, executor)
		},
		Flags: []cli.Flag{phpVersionFlag},
		Before: func(context *cli.Context) error {
			// if bootstrap, add-site
			fmt.Println("add validation")
			return nil
		},
	}
}

func Update() *cli.Command {
	return &cli.Command{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "Update a machine with the latest software",
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
		Name:    "add-host",
		Aliases: []string{"add", "host"},
		Usage:   "Add a new virtual host to a machine",
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

func Logs() *cli.Command {
	return &cli.Command{
		Name:        "logs",
		Aliases:     []string{"log", "l"},
		Description: "Show a machines logs",
		Action: func(context *cli.Context) error {
			return cli.ShowSubcommandHelp(context)
		},
		Subcommands: []*cli.Command{
			{
				Name:        "nginx",
				Description: "Show logs from nginx",
				Action: func(context *cli.Context) error {
					// TODO tail multiple files at once with
					// tail -f /var/log/syslog -f /var/log/auth.log
					return errors.New("not implemented")
				},
			},
		},
	}
}
func Stop() *cli.Command {
	return &cli.Command{
		Name:        "stop",
		Aliases:     []string{"shutdown"},
		Usage:       "Stop a machine",
		Description: "Stop a machine when not in use",
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
