package envedit

import (
	"fmt"
	"io/ioutil"
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
func Edit(file string, updates map[string]string) (string, error) {
	// make sure the file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return "", ErrNoEnvFile
	}

	// read the file
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	// split the file into multiple lines
	lines := strings.Split(string(f), "\n")
	for line, txt := range lines {
		// split using =
		sp := strings.Split(txt, "=")

		// check if this is a thing we should modify
		if _, ok := updates[sp[0]]; ok {
			// replace the line
			lines[line] = strings.Join([]string{sp[0], updates[sp[0]]}, "=")
		}
	}

	return strings.Join(lines, "\n"), nil
}

// EnvExists takes an existing env file and key and checks if the env var has already been defined. If it has been defined
// it will return true otherwise it will return false. If the file does not exist, it will return false.
func EnvExists(file, key string) bool {
	// make sure the file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	// read the file
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return false
	}

	// split the file into multiple lines
	lines := strings.Split(string(f), "\n")
	for _, txt := range lines {
		// split using =
		sp := strings.Split(txt, "=")

		if sp[0] == key && sp[1] != "" {
			return true
		}
	}

	return false
}
