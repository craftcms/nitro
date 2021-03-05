package volumename

import (
	"os"
	"strings"
)

// FromPath takes a file path and returns a name that is friendly to use for docker volumes
func FromPath(path string) string {
	// make it lower case
	path = strings.ToLower(path)

	// replace path separators with underscores
	path = strings.Replace(path, string(os.PathSeparator), "_", -1)

	// replace spaces with underscores
	path = strings.Replace(path, " ", "_", -1)

	// remove : to prevent errors on windows
	path = strings.Replace(path, ":", "_", -1)

	// remove the first underscore
	return strings.TrimLeft(path, "_")
}
