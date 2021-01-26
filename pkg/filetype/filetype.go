package filetype

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Determine takes a file path and will determine
// if the file is plain, zip, or a tar type of
// file. If the path is not found it will return
// an error.
func Determine(file string) (string, error) {
	// stat the file to make sure it exists
	stat, err := os.Stat(file)
	if err != nil {
		return "", err
	}

	// make sure its not a directory
	if stat.IsDir() {
		return "", fmt.Errorf("file provided is a directory")
	}

	// read the contents
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	// detect the type
	kind := http.DetectContentType(data)

	switch kind {
	case "text/plain; charset=utf-8":
		return "text", nil
	case "application/zip":
		return "zip", nil
	case "application/x-gzip":
		return "tar", nil
	}

	return "", fmt.Errorf("unknown file type: %s", kind)
}
