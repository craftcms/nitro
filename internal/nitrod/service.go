package nitrod

import (
	"log"
	"os"
)

// Version is the nitro version, not using it in the
// api yet.
var Version string

// NitroService is the struct that runs the gRPC API
type NitroService struct {
	command Runner
	logger  *log.Logger
}

// NewNitroService will create a new service
// with the default command and logger
func NewNitroService() *NitroService {
	return &NitroService{
		command: &ServiceRunner{},
		logger:  log.New(os.Stdout, "nitrod ", 0),
	}
}
