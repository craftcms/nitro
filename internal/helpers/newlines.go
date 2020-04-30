package helpers

import "strings"

func normalizeNewlines(v string) string {
	return strings.Replace(v, "\r", "", -1)
}
