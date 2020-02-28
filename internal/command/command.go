package command

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/nitro/internal/action"
	"github.com/pixelandtonic/nitro/internal/validate"
)

var (
	bootstrapFlag = &cli.BoolFlag{
		Name:        "bootstrap",
		Usage:       "Bootstrap the machine with defaults",
		Value:       false,
		DefaultText: "false",
	}
	phpVersionFlag = &cli.StringFlag{
		Name:        "php-version",
		Usage:       "Provide version of PHP",
		DefaultText: "7.4",
	}
	cpuFlag = &cli.Int64Flag{
		Name:        "cpus",
		Usage:       "The number of CPUs to assign the machine",
		Required:    false,
		Value:       1,
		DefaultText: "1",
	}
	memoryFlag = &cli.StringFlag{
		Name:        "memory",
		Usage:       "The amount of memory to assign the machine",
		Required:    false,
		Value:       "1G",
		DefaultText: "1G",
	}
	diskFlag = &cli.StringFlag{
		Name:        "disk",
		Usage:       "The amount of disk to assign the machine",
		Required:    false,
		Value:       "5G",
		DefaultText: "5G",
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
		Flags: []cli.Flag{
			bootstrapFlag,
			cpuFlag,
			memoryFlag,
			diskFlag,
		},
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
		Before: func(c *cli.Context) error {
			if c.Args().First() == "" {
				// TODO validate the domain name with validate.Domain(d)
				return errors.New("you must pass a domain name")
			}

			if err := validate.PHPVersion(c.String("php-version")); err != nil {
				return err
			}

			if err := validate.Path(c.String("path")); err != nil {
				return err
			}

			return nil
		},
		Flags: []cli.Flag{phpVersionFlag, &cli.StringFlag{
			Name:     "path",
			Usage:    "The path to the directory to mount",
			Required: true,
		}},
	}
}

func Logs(e action.CommandLineExecutor) *cli.Command {
	return &cli.Command{
		Name:        "logs",
		Description: "Show a machines logs",
		Action: func(c *cli.Context) error {
			return cli.ShowSubcommandHelp(c)
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

func IP() *cli.Command {
	return &cli.Command{
		Name:        "ip",
		Usage:       "Show the IP address of a machine",
		Description: "Show a machines IP address",
		Action: func(c *cli.Context) error {
			return action.IP(c)
		},
	}
}

func Database(e action.CommandLineExecutor) *cli.Command {
	return &cli.Command{
		Name:        "database",
		Description: "Perform actions related to the database",
		Action: func(c *cli.Context) error {
			return cli.ShowSubcommandHelp(c)
		},
		Subcommands: []*cli.Command{
			databasePassword(e),
		},
	}
}

func databasePassword(e action.CommandLineExecutor) *cli.Command {
	return &cli.Command{
		Name:        "password",
		Description: "Show the user password for the database",
		Action: func(c *cli.Context) error {
			return action.DatabasePassword(c, e)
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
