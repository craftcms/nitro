package slug

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile("[^a-z0-9]+")

// Generate is used to create a slugged version of a string.
func Generate(s string) string {
	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "_"), "-")
}
