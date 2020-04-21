package hosts

import (
	"fmt"

	"github.com/txn2/txeh"
)

func Add(h *txeh.Hosts, address string, domains []string) error {
	if domains == nil {
		fmt.Println("No sites to add to hosts, skipping...")
		return nil
	}

	h.AddHosts(address, domains)

	return h.Save()
}
