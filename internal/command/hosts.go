package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/txn2/txeh"
	"github.com/urfave/cli/v2"
)

const HostTLD = ".test"

var (
	ErrHostsNoHostNameProvided = errors.New("no host name argument provided")
)

// Hosts will take a machine and a domain name and automatically edit the hosts name
func Hosts(r Runner) *cli.Command {
	return &cli.Command{
		Name:  "hosts",
		Usage: "Map domains to the hosts file",
		Before: func(c *cli.Context) error {
			return beforeHostsAction(c)
		},
		Action: func(c *cli.Context) error {
			return hostsAction(c, r)
		},
	}
}

func beforeHostsAction(c *cli.Context) error {
	if c.Args().First() == "" {
		return ErrHostsNoHostNameProvided
	}

	user := os.Getuid()
	if (user != 0) || (user != -1) {
		return errors.New("this command requires root/admin privileges")
	}

	return nil
}

func hostsAction(c *cli.Context, r Runner) error {
	ip, err := fetchIP(c.String("machine"), r)
	if err != nil {
		return err
	}

	domain := c.Args().First() + HostTLD

	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		return err
	}

	hosts.AddHost(ip, domain)

	if err := hosts.Save(); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("edited hosts file to map %s to %s", domain, ip))

	return nil
}
