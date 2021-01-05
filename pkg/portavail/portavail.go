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
		return fmt.Errorf("it appears port %s, is already in use", port)
	}

	// check the close error
	if err := lis.Close(); err != nil {
		w := fmt.Errorf("unable to close the listener, %w", err)

		fmt.Println(w.Error())

		return w
	}

	return nil
}
