package command

import (
	"bufio"
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
			// initialize that machine
			if err := action.Initialize(c); err != nil {
				return err
			}

			// if we are bootstrapping, call that command
			if c.Bool("bootstrap") {
				return c.App.RunContext(c.Context, []string{c.App.Name, "--machine", c.String("machine"), "bootstrap"})
			}

			return nil
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
			if err := action.Bootstrap(context, executor); err != nil {
				return err
			}

			return nil
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
			if err := action.Update(c); err != nil {
				return err
			}
			return nil
		},
	}
}

func SSH(e action.CommandLineExecutor) *cli.Command {
	return &cli.Command{
		Name:    "ssh",
		Aliases: []string{"s", "connect", "login"},
		Usage:   "SSH into a machine as administrator",
		Action: func(c *cli.Context) error {
			if err := action.SSH(c.String("machine"), e); err != nil {
				return err
			}
			return nil
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
			if err := action.Stop(c); err != nil {
				return err
			}

			return nil
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

			if err := action.Delete(c); err != nil {
				return err
			}

			return nil
		},
	}
}
