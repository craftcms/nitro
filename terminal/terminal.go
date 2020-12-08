package terminal

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Outputer is an interface that captures the output to a terminal.
// It is used to make our output consistent in the command line.
type Outputer interface {
	Info(s ...string)
	Success(s ...string)
	Pending(s ...string)
	Select(r io.Reader, msg string, opts []string) (int, error)
	Done()
}

type terminal struct{}

// New returns an Outputer interface
func New() Outputer {
	return terminal{}
}

func (t terminal) Info(s ...string) {
	fmt.Printf("%s\n", strings.Join(s, " "))
}

func (t terminal) Success(s ...string) {
	fmt.Printf("  \u2713 %s\n", strings.Join(s, " "))
}

func (t terminal) Pending(s ...string) {
	fmt.Printf("  … %s ", strings.Join(s, " "))
}

func (t terminal) Done() {
	fmt.Print("\u2713\n")
}

func (t terminal) Select(r io.Reader, msg string, opts []string) (int, error) {
	// show all the options
	for k, v := range opts {
		fmt.Println(fmt.Sprintf("  %d. %s", k+1, v))
	}

	// show the message
	fmt.Print(msg)

	// create for loop until the input is valid
	var selection int
	wait := true
	for wait {
		rdr := bufio.NewReader(r)
		char, err := rdr.ReadString('\n')
		if err != nil {
			return 0, err
		}

		// remove the new line from string
		char = strings.TrimSpace(char)

		// convert the selection to an integer
		s, err := strconv.Atoi(char)
		// make sure its a valid option
		if err != nil || len(opts) < s {
			wait = true
			fmt.Println("Please choose a valid option 🙄...")

			for k, v := range opts {
				fmt.Println(fmt.Sprintf("  %d. %s", k+1, v))
			}

			fmt.Print(msg)
		} else {
			// take away one from the selection
			selection = s - 1
			wait = false
		}
	}

	return selection, nil
}
