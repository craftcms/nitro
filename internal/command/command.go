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
		Name:  "logs",
		Usage: "Show logs for a machine",
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

func MariaDB() *cli.Command {
	return &cli.Command{
		Name:        "mariadb",
		Usage:       "Enter a root shell for mariadb",
		Description: "Allows a user to quickly access a mariadb shell as the root user",
		Category:    "databases",
		Action: func(c *cli.Context) error {
			return errors.New("not implemented")
		},
		OnUsageError: nil,
		Subcommands:  nil,
	}
}
