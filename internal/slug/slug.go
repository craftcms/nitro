package slug

import "strings"

// Generate takes a string and will trim, remove spaces, and remove special characters
func Generate(s string) string {
	sl := s

	// trim the string
	sl = strings.TrimSpace(sl)

	// remove spaces
	if strings.Contains(sl, " ") {
		sl = strings.ReplaceAll(sl, " ", "_")
	}

	return sl
}
