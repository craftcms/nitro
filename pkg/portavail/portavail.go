package portavail

import (
	"fmt"
	"net"
	"strconv"
)

// Check takes ports and will check for use against the localhost:port. If any port provided
// is in use, it will return an error.
func Check(host, port string) error {
	hostname := "localhost"
	if host != "" {
		hostname = host
	}

	// create a new listener
	lis, err := net.Listen("tcp", hostname+":"+port)
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

// FindNext takes a host and port and will find the next available port
func FindNext(host, port string) (string, error) {
	// convert the port to an integer
	p, err := strconv.Atoi(port)
	if err != nil {
		return "", err
	}

	// check if the port is available until it is :)
	for {
		if err := Check(host, strconv.Itoa(p)); err != nil {
			p = p + 1
			continue
		}

		break
	}

	return strconv.Itoa(p), nil
}
