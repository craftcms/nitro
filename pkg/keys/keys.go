package keys

import (
	"os"
	"path/filepath"
	"strings"
)

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

	return keys, nil
}
