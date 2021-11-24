package dockererrors

import "strings"

// IsPortError takes an error from the Docker API and checks the content
// to determine if the error is a port allocation issue.
func IsPortError(err error) bool {
	if err == nil {
		return false
	}

	// check the content of the error message
	if strings.Contains(err.Error(), "port is already allocated") {
		return true
	}

	return false
}
