package certinstall

import "fmt"

var (
	certificatePaths = map[string]string{
		"ubuntu": "/usr/local/share/ca-certificates/",
	}
	certificateTools = map[string]string{
		"ubuntu": "update-ca-certificates",
	}
)

func Linux(system, certFile string) error {
	return nil
}

func findPath(system string) (path string, err error) {
	if path, ok := certificatePaths[system]; ok {
		return path, nil
	}

	return "", fmt.Errorf("unable to find the path for the system %q", system)
}

func findTool(system string) (tool string, err error) {
	if tool, ok := certificateTools[system]; ok {
		return tool, nil
	}

	return "", fmt.Errorf("unable to find the tool for the system %q", system)
}
