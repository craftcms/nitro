package command

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

// Bootstrap will install the requirements
func Bootstrap() *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Usage:       "Delete machine",
		Description: "Delete a machine when no longer needed, this is recoverable and not permanently deleted.",
		Before:      bootstrapBeforeAction,
		Action:      bootstrapAction,
	}
}

func bootstrapBeforeAction(c *cli.Context) error {
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

	return nil
}

func bootstrapAction(c *cli.Context) error {
	machine := c.String("machine")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	cmd := exec.Command(multipass, "delete", machine)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
