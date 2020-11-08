package portavail

import (
	"fmt"
	"net"
)

// Check takes ports and will check for use against the localhost:port. If any port provided
// is in use, it will return an error.
func Check(ports ...string) error {
	if len(ports) == 0 {
		return fmt.Errorf("expected a list of ports to check, nothing was provided")
	}

	for _, port := range ports {
		lis, err := net.Listen("tcp", "localhost:"+port)
		if err != nil {
			return fmt.Errorf("It appears port %q, is already in use", port)
		}

		if err := lis.Close(); err != nil {
			return fmt.Errorf("unable to close the listener after checking the ports, %w", err)
		}
	}

	return nil
}
