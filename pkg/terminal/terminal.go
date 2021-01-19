package terminal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Outputer is an interface that captures the output to a terminal.
// It is used to make our output consistent in the command line.
type Outputer interface {
	Ask(message, fallback, sep string, validator Validator) (string, error)
	Info(s ...string)
	Success(s ...string)
	Pending(s ...string)
	Select(r io.Reader, msg string, opts []string) (int, error)
	Warning()
	Done()
}

type Asker interface {
	Ask(message, fallback, sep string, validator Validator) (string, error)
}

// Validator is used to pass information into a terminal prompt for
// validating inputs that are strings
type Validator interface {
	Validate(input string) error
}

type terminal struct{}

// New returns an Outputer interface
func New() *terminal {
	return &terminal{}
}

func (t *terminal) Ask(message, fallback, sep string, validator Validator) (string, error) {
	t.printMessage(message, fallback, sep)

	// create a new scanner
	s := bufio.NewScanner(os.Stdin)

	// split on lines so we can listen for carriage returns
	s.Split(bufio.ScanLines)

	var a string
	for s.Scan() {
		txt := s.Text()

		// return the fallback if the text is blank
		if txt == "" && fallback != "" {
			a = fallback
			break
		}

		// validate the input
		if validator != nil {
			if err := validator.Validate(txt); err != nil {
				// there is an error, so display the error and show the message again
				t.printValidatorError(err)

				t.printMessage(message, fallback, sep)

				continue
			}
		}

		// if there is text, set and break because it passed validation
		if txt != "" {
			a = txt
			break
		}

		// last ditch effort here, check both the text and fallback
		if txt == "" && fallback == "" {
			t.printMessage(message, fallback, sep)

			continue
		}
	}

	// handle error from the scanner
	if err := s.Err(); err != nil {
		return "", fmt.Errorf("unable to handle input, %w", err)
	}

	return a, nil
}

func (t *terminal) printMessage(message, fallback, sep string) {
	if fallback == "" {
		fmt.Fprintf(os.Stdout, "%s%s ", message, sep)
		return
	}

	fmt.Fprintf(os.Stdout, "%s [%s]%s ", message, fallback, sep)
}

func (t *terminal) printValidatorError(err error) {
	fmt.Fprintf(os.Stdin, " \u2717 %s\n", err.Error())
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

func (t terminal) Select(r io.Reader, msg string, opts []string) (int, error) {
	// if the options only have one item, return it
	if len(opts) == 1 {
		return 0, nil
	}

	// show the message
	fmt.Println(msg)

	// show all the options
	for k, v := range opts {
		fmt.Printf("  %d. %s\n", k+1, v)
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
			fmt.Println("Please choose a valid option…")

			for k, v := range opts {
				fmt.Printf("  %d. %s\n", k+1, v)
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
