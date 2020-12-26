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
	// Input(r io.Reader, validate Validator, msg string) (string, error)
	Warning()
	Done()
}

// Validator is used to pass information into a terminal prompt for
// validating inputs that are strings
type Validator interface {
	// Validate takes an input, illegal chars, and the error message and returns an error
	Validate(input, chars, msg string) error
}

// ValidateString is used to validate terminal input as a string
type ValidateString struct{}

// Validate takes the input and checks if the string contains any chars
func (t *ValidateString) Validate(input, chars, msg string) error {
	if strings.ContainsAny(input, chars) {
		return fmt.Errorf(msg)
	}

	return nil
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

func (t terminal) Warning() {
	fmt.Print("\u2717\n")
}

// func (t terminal) Input(r io.Reader, validate Validator, msg string) (string, error) {
// 	// show the message
// 	fmt.Print(msg)

// 	var input string
// 	w := true
// 	for w {
// 		rdr := bufio.NewReader(r)
// 		char, err := rdr.ReadString('\n')
// 		if err != nil {
// 			return "", err
// 		}

// 		// remove the new line from string
// 		char = strings.TrimSpace(char)

// 		err := validate.Validate(input)
// 	}

// 	return input, nil
// }

func (t terminal) Select(r io.Reader, msg string, opts []string) (int, error) {
	// if the options only have one item, return it
	if len(opts) == 1 {
		return 0, nil
	}

	// show the message
	fmt.Println(msg)
	// show all the options
	for k, v := range opts {
		fmt.Println(fmt.Sprintf("  %d. %s", k+1, v))
	}

	fmt.Print("Enter your selection: ")

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

		// convert the selection to an integer and make sure its a valid option
		s, err := strconv.Atoi(char)
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
