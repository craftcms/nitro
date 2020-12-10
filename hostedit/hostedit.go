package hostedit

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

const (
	StartText = "# <nitro>"
	EndText   = "# </nitro>"
)

func Update(file, addr string, hosts ...string) (string, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	// split the file into multiple lines
	lines := strings.Split(string(f), "\n")

	// the index represents where the content (addr and hosts) should be placed
	// which is in between the start and end text comment
	var index int
	for l, t := range lines {
		// look for the beginning text
		if strings.Contains(t, StartText) {
			// the next line is the empty line
			index = l + 1
		}

		// look for the end text
		if strings.Contains(t, EndText) {
			// we want the previous line
			index = l - 1
		}
	}

	switch index {
	// if there is not a comment section, we need to create one
	case 0:
		if lines[len(lines)-1] != "" {
			// add a blank line
			lines = append(lines, "")
		}
		lines = append(lines, StartText)
		lines = append(lines, fmt.Sprintf("%s\t%s", addr, strings.Join(hosts, " ")))
		lines = append(lines, EndText+"\n")
	default:
		// replace the line between the start and end text with the contents of the address and hosts
		lines[index] = fmt.Sprintf("%s\t%s", addr, strings.Join(hosts, " "))
	}

	return strings.Join(lines, "\n"), nil
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
