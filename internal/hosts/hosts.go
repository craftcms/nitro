package hosts

import (
	"fmt"
)

// AdderSaver is an interface for adding and saving hosts
type AdderSaver interface {
	AddHosts(address string, domains []string)
	Save() error
}

type RemoverSaver interface {
	RemoveHosts(domains []string)
	Save() error
}

// Add takes an address and a list of hosts, or domains, to
// add to a hosts file. It uses an AdderSaver, which uses
// the functionality of the library used to edit hosts.
func Add(a AdderSaver, address string, domains []string) error {
	if domains == nil {
		fmt.Println("No sites to add to hosts file, skipping...")
		return nil
	}

	a.AddHosts(address, domains)

	return a.Save()
}

// Remove takes a list of domains that should be removed
// from the hosts file. It uses a RemoverSaver, which
// uses the functionality of the lib to edit hosts.
func Remove(r RemoverSaver, domains []string) error {
	if domains == nil {
		fmt.Println("No sites to remove from hosts file, skipping...")
		return nil
	}

	r.RemoveHosts(domains)

	return r.Save()
}
