package helpers

import "strings"

// RemoveTrailingSlash is used to take a string and remove the
// trailing slash from the provided string
func RemoveTrailingSlash(v string) string {
	if strings.HasSuffix(v, "/") {
		return strings.TrimRight(v, "/")
	}

	return v
}
