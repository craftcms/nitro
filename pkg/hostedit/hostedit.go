package hostedit

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	startText = "# <nitro>"
	endText   = "# </nitro>"
)

var ErrNotNitroEntries = fmt.Errorf("there are no nitro entries to remove from the hosts file")

// Update takes a file, reads the content and updates or appends
// the addr and hosts for the sites.
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
		if strings.Contains(t, startText) {
			// the next line is the empty line
			index = l + 1
		}

		// look for the end text
		if strings.Contains(t, endText) {
			// we want the previous line
			index = l - 1
		}
	}

	switch index {
	// if there is not a comment section, we need to create one
	case 0:
		lines = append(lines, startText)
		lines = append(lines, fmt.Sprintf("%s\t%s", addr, strings.Join(hosts, " ")))
		lines = append(lines, endText+"\n")
	default:
		// replace the line between the start and end text with the contents of the address and hosts
		lines[index] = fmt.Sprintf("%s\t%s", addr, strings.Join(hosts, " "))
	}

	return strings.Join(lines, "\n"), nil
}

// IsUpdated is used to check if an update will make any changes
// to the hosts file and return true if there is nothing to change
func IsUpdated(file, addr string, hosts ...string) (bool, error) {
	// open the original file
	orig, err := ioutil.ReadFile(file)
	if err != nil {
		// we could not open the file so just assume its good
		return false, err
	}

	// perform the update
	updated, err := Update(file, addr, hosts...)
	if err != nil {
		return false, err
	}

	// compare the two to see if they are updated
	return string(orig) == updated, nil
}

// Remove is responsible for removing all of the hosts entries
// for the nitro config from the hosts file.
func Remove(file string) (string, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	// split the lines
	lines := strings.Split(string(f), "\n")

	// get the indexes to remove (start, middle and end)
	start, middle, end := indexes(f)

	// if there are no entries, return a specific error
	if start == 0 && middle == 0 && end == 0 {
		return "", ErrNotNitroEntries
	}

	// create a new hosts file in memory
	new := []string{}
	for i, v := range lines {
		// if this is one of the indexes, remove
		if i == start || i == middle || i == end {
			continue
		}

		// append the entry to the new file
		new = append(new, v)
	}

	return strings.Join(new, "\n"), nil
}

func indexes(content []byte) (int, int, int) {
	// split the file into multiple lines
	lines := strings.Split(string(content), "\n")

	// the index represents where the content (addr and hosts) should be placed
	// which is in between the start and end text comment
	var start, middle, end int
	for l, t := range lines {
		// look for the beginning text
		if strings.Contains(t, startText) {
			start = l
			// the next line is the empty line
			middle = l + 1
		}

		// look for the end text
		if strings.Contains(t, endText) {
			// we want the previous line
			end = l
			middle = l - 1
		}
	}

	return start, middle, end
}
