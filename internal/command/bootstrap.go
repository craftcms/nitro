package command

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal"
	"github.com/craftcms/nitro/internal/validate"
)

// Bootstrap will install the software packages on the machine
func Bootstrap(r internal.Runner) *cli.Command {
	return &cli.Command{
		Name:        "bootstrap",
		Usage:       "Delete machine",
		Description: "Delete a machine when no longer needed, this is recoverable and not permanently deleted.",
		Before:      bootstrapBeforeAction,
		Action: func(c *cli.Context) error {
			return bootstrapAction(c, r)
		},
		After: bootstrapAfterAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "php-version",
				Usage:       "Provide version of PHP",
				Value:       "7.4",
				DefaultText: "7.4",
			},
			&cli.StringFlag{
				Name:        "database",
				Usage:       "Provide version of PHP",
				Value:       "mariadb",
				DefaultText: "mariadb",
			},
		},
	}
}

func bootstrapBeforeAction(c *cli.Context) error {
	if err := validate.PHPVersion(c.String("php-version")); err != nil {
		return err
	}

	if err := validate.Database(c.String("database")); err != nil {
		return err
	}

	return nil
}

func bootstrapAction(c *cli.Context, r internal.Runner) error {
	machine := c.String("machine")
	php := c.String("php-version")
	database := c.String("database")

	return r.Run([]string{"exec", machine, "--", "sudo", "bash", "/opt/nitro/bootstrap.sh", php, database})
}

func bootstrapAfterAction(c *cli.Context) error {
	// TODO make this pull the actual values
	database := c.String("database")

	var port int
	switch database {
	case "postgres":
		port = 5432
	default:
		port = 3306
	}

	fmt.Println("")
	fmt.Println("==== DATABASE ====")
	fmt.Println("Server:", "192.168.x.x")
	fmt.Println("Port:", port)
	fmt.Println("Driver:", database)
	fmt.Println("Username:", "nitro")
	fmt.Println("Password:", "nitro")

	return nil
}
