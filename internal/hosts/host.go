package hosts

import (
	"fmt"
)

// AdderSaver is an interface for adding and saving hosts
type AdderSaver interface {
	AddHosts(address string, domains []string)
	Save() error
}

// Add takes an address and a list of hosts, or domains, to
// add to a hosts file. It uses an AdderSaver, which uses
// the functionality of the library used to edit hosts.
func Add(a AdderSaver, address string, domains []string) error {
	if domains == nil {
		fmt.Println("No sites to add to hosts, skipping...")
		return nil
	}

	a.AddHosts(address, domains)

	return a.Save()
}
