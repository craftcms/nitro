package hostedit

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// load the file

const (
	StartText = "# <nitro>"
	EndText   = "# </nitro>"
)

func Update(file string) (string, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	start, end := find(bytes.NewBuffer(f))

	fmt.Println(start)
	fmt.Println(end)

	return "", nil
}

func find(r io.Reader) (int, int) {
	scanner := bufio.NewScanner(r)

	// get the
	var start, end, index int
	for scanner.Scan() {
		index = index + 1
		if strings.Contains(scanner.Text(), StartText) {
			start = index
		}

		if strings.Contains(scanner.Text(), EndText) {
			end = index
			// once we have the end, we can stop here
			break
		}
	}

	return start, end
}
