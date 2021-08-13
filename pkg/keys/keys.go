package keys

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrNoKeysFound is returned when there are no keys present
var ErrNoKeysFound = fmt.Errorf("Unable to find keys to import")

// Find all the key pairs in a directory
func Find(path string) (map[string]string, error) {
	keys := make(map[string]string)
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || info.Name() == "known_hosts" || strings.Contains(info.Name(), ".pem") || info.Name() == "config" || info.Name() == ".DS_Store" {
			return nil
		}

		sp := strings.Split(info.Name(), ".pub")

		if keys[sp[0]] == "" {
			keys[sp[0]] = sp[0]
		}

		if strings.Contains(info.Name(), ".pub") {
			keys[sp[0]] = info.Name()
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// if there are no keys, return an error
	if len(keys) == 0 {
		return nil, ErrNoKeysFound
	}

	return keys, nil
}
