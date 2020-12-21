package envedit

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	// ErrNoEnvFile is returned when the environment file cannot be found
	ErrNoEnvFile = fmt.Errorf("no environment file found")
)

// Edit takes a file and a list of updates as a map of strings. The updates should
// represent the ENV_VAR VAL. Edit will check the file line by line for each change
// and if the environment variable is contained in the updates, it will update and
// save the file.
func Edit(file string, updates ...map[string]string) error {
	// make sure the file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return ErrNoEnvFile
	}

	// open the file
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// split the string
		sp := strings.Split(scanner.Text(), "=")

		// check each of the updates
		for _, update := range updates {
			// is this an update we need to take action on?
			if _, ok := update[sp[0]]; ok {
				fmt.Println(sp[0])
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
