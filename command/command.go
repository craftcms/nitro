package command

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/nitro/action"
)

func Initialize() *cli.Command {
	return &cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "Initialize a new machine",
		Action: func(c *cli.Context) error {
			return action.Initialize(c)
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "bootstrap",
				Aliases: []string{"b"},
				Usage:   "Bootstrap the machine with installation defaults",
				Value:   true,
			},
		},
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
		// TODO add flags for version and database
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
		Name:    "ssh",
		Aliases: []string{"s", "connect", "login"},
		Usage:   "SSH into a machine as administrator",
		Action: func(c *cli.Context) error {
			return action.SSH(c.String("machine"), e)
		},
	}
}

func AddSite() *cli.Command {
	return &cli.Command{
		Name:    "add-site",
		Aliases: []string{"add", "site"},
		Usage:   "Add a new site to a machine",
		Action: func(context *cli.Context) error {
			return errors.New("not implemented yet")
		},
		Before: func(context *cli.Context) error {
			if context.NArg() != 1 {
				return errors.New("you must pass a domain and php version")
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "php-version",
				Aliases:     []string{"php", "version"},
				Usage:       "Provide the version of PHP",
				DefaultText: "7.4",
			},
		},
	}
}

func Install() *cli.Command {
	return &cli.Command{
		Name:        "install",
		Usage:       "Install common tools such as PHP, web servers, and databases",
		Description: "Install offers common options for installing packages on a machine",
		Action: func(c *cli.Context) error {
			return cli.ShowSubcommandHelp(c)
		},
		Subcommands: []*cli.Command{
			installPHP(),
			installNginx(),
			installMariaDB(),
			installRedis(),
			installPostgres(),
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
