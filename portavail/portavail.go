package portavail

import (
	"fmt"
	"net"
)

// Check takes ports and will check for use against the localhost:port. If any port provided
// is in use, it will return an error.
func Check(port string) error {
	lis, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		return fmt.Errorf("It appears port %s, is already in use", port)
	}

	return lis.Close()
}
