package nitrod

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/validate"
)

// SystemService is used to start, stop, and restart
// system services on the virtual machine such as
// nginx, php-fpm, docker, and etc.
type SystemService struct {
	command Runner
	logger  *log.Logger
}

// Nginx is used to manage the nginx service.
func (s *SystemService) Nginx(ctx context.Context, req *NginxServiceRequest) (*ServiceResponse, error) {
	message, action := s.messageAndAction(req.GetAction())

	// perform the action on the nginx service
	if output, err := s.command.Run("service", []string{"nginx", action}); err != nil {
		s.logger.Println(err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	msg := "successfully " + message + " nginx"

	s.logger.Println(msg)

	return &ServiceResponse{Message: msg}, nil
}

// PhpFpm is used to manage the php<version>-fpm service.
func (s *SystemService) PhpFpm(ctx context.Context, req *PhpFpmServiceRequest) (*ServiceResponse, error) {
	// validate the request
	if err := validate.PHPVersion(req.GetVersion()); err != nil {
		s.logger.Println(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	message, action := s.messageAndAction(req.GetAction())

	// perform the action on the nginx service
	if output, err := s.command.Run("service", []string{"php" + req.GetVersion() + "-fpm", action}); err != nil {
		s.logger.Println(err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	msg := "successfully " + message + " php-fpm " + req.GetVersion()

	s.logger.Println(msg)

	return &ServiceResponse{Message: msg}, nil
}

func (s *SystemService) messageAndAction(action ServiceAction) (string, string) {
	switch action {
	case ServiceAction_START:
		return "started", "start"
	case ServiceAction_STOP:
		return "stopped", "stop"
	default:
		return "restarted", "restart"
	}
}

// NewSystemService will create a new
// service with the default command
// runner and logging to stdout
func NewSystemService() *SystemService {
	return &SystemService{
		command: &ServiceRunner{},
		logger:  log.New(os.Stdout, "nitrod ", 0),
	}
}
