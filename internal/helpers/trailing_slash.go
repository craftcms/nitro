package helpers

import "strings"

func RemoveTrailingSlash(v string) string {
	if strings.HasSuffix(v, "/") {
		return strings.TrimRight(v, "/")
	}

	return v
}
