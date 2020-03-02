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
