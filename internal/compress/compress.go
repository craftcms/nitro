package compress

import (
	"compress/gzip"
	"io/ioutil"
	"os"
	"time"
)

// ToTempFile takes a file path and a pattern to create a temporary
// file in the OS temp directory. The file contents will then be
// read and compressed into the temp file and return the temp file.
func ToTempFile(file, pattern string) (*os.File, error)  {
	t, err := ioutil.TempFile(os.TempDir(), pattern)
	if err != nil {
		return nil, err
	}

	z := gzip.NewWriter(t)
	z.Name = file
	z.ModTime = time.Now()

	c, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	_, err = z.Write(c)
	if err != nil {
		return nil, err
	}

	if err := z.Close(); err != nil {
		return nil, err
	}

	return t, nil
}