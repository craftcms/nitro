package archive

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// FromFile takes a file and returns a reader or error.
// It is used for creating a tar/archive to copy items
// into a container.
func FromFile(file *os.File) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	info, err := os.Stat(file.Name())
	if err != nil {
		return nil, err
	}

	header, err := tar.FileInfoHeader(info, file.Name())
	if err != nil {
		return nil, err
	}

	// header.Name = strings.TrimPrefix(strings.Replace(file, path, "", -1), string(filepath.Separator))
	err = tw.WriteHeader(header)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("is directory")
	}

	_, err = io.Copy(tw, file)
	if err != nil {
		return nil, err
	}

	return tar.NewReader(&buf), nil
}

func FromString(filename, content string) (io.Reader, error) {
	// create a temp file
	f, err := ioutil.TempFile(os.TempDir(), "nitro-archive-")
	if err != nil {
		return nil, err
	}

	// write the contents
	if _, err := f.Write([]byte(content)); err != nil {
		return nil, err
	}

	return FromFile(f)
}
