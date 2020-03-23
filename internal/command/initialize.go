package command

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"
)

// Initialize it the main entry point when calling the init command to start a new machine
func Initialize(r Runner) *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize new machine",
		Action: func(c *cli.Context) error {
			return initializeAction(c, r)
		},
		After: func(c *cli.Context) error {
			return initializeAfterAction(c, r)
		},
		Flags: []cli.Flag{
			bootstrapFlag,
			phpVersionFlag,
			databaseFlag,
			cpusFlag,
			memoryFlag,
			diskFlag,
		},
	}
}

func initializeAction(c *cli.Context, r Runner) error {
	machine := c.String("machine")
	cpus := strconv.Itoa(c.Int("cpus"))
	disk := c.String("disk")
	mem := c.String("memory")

	// pass the cloud init as stdin
	if err := r.SetInput(cloudInit); err != nil {
		return err
	}

	return r.Run([]string{"launch", "--name", machine, "--cpus", cpus, "--disk", disk, "--mem", mem, "--cloud-init", "-"})
}

func initializeAfterAction(c *cli.Context, r Runner) error {
	// if we are bootstrapping, call the command
	if c.Bool("bootstrap") {
		machine := c.String("machine")
		php := c.String("php-version")
		database := c.String("database")

		if err := r.Run([]string{"exec", machine, "--", "sudo", "bash", "/opt/nitro/bootstrap.sh", php, database}); err != nil {
			return err
		}

		// print the system information
		ip, err := fetchIP(machine, r)
		if err != nil {
			return err
		}

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
		fmt.Println("")
		fmt.Println("To get started, you can add a new site, which will mount a directory into the virtual machine and configure nginx.\n Run the command:")
		fmt.Println("nitro --machine",c.String("machine"), "site --path /path/to/website project.test")

		return nil
	}

	return nil
}
