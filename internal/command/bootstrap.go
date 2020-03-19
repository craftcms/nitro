package command

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/validate"
)

// Bootstrap will install the software packages on the machine
func Bootstrap(r Runner) *cli.Command {
	return &cli.Command{
		Name:        "bootstrap",
		Usage:       "Bootstrap machine",
		Description: "Bootstrap machine",
		Before:      bootstrapBeforeAction,
		Action: func(c *cli.Context) error {
			return bootstrapAction(c, r)
		},
		After: func(c *cli.Context) error {
			return bootstrapAfterAction(c, r)
		},
		Flags: []cli.Flag{
			phpVersionFlag,
			databaseFlag,
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

func bootstrapAction(c *cli.Context, r Runner) error {
	machine := c.String("machine")
	php := c.String("php-version")
	database := c.String("database")

	return r.Run([]string{"exec", machine, "--", "sudo", "bash", "/opt/nitro/bootstrap.sh", php, database})
}

func bootstrapAfterAction(c *cli.Context, r Runner) error {
	ip, err := fetchIP(c.String("machine"), r)
	if err != nil {
		return err
	}

	database := c.String("database")

	var port int
	var driver string
	switch database {
	case "postgres":
		driver = "pgsql"
		port = 5432
	default:
		driver = "mysql"
		port = 3306
	}

	fmt.Println("")
	fmt.Println("==== SERVER INFO ====")
	fmt.Println("server:", "http://"+ip)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("==== DATABASE INFO ====")
	fmt.Println("server:", ip)
	fmt.Println("port:", port)
	fmt.Println("driver:", driver)
	fmt.Println("database:", "craftcms")
	fmt.Println("username:", "nitro")
	fmt.Println("password:", "nitro")
	fmt.Println("")
	fmt.Println("For additional information on nitro, visit https://docs.craftcms.com/v3/nitro")

	return nil
}
