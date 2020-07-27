package nitrod

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/validate"
)

type SystemService struct {
	command Runner
	logger  *log.Logger
}

func (s *SystemService) Nginx(ctx context.Context, req *NginxServiceRequest) (*ServiceResponse, error) {
	var action string
	var message string
	switch req.GetAction() {
	case ServiceAction_START:
		message = "started"
		action = "start"
	case ServiceAction_STOP:
		message = "stopped"
		action = "stop"
	default:
		message = "restarted"
		action = "restart"
	}

	// perform the action on the php-fpm service
	if output, err := s.command.Run("service", []string{"nginx", action}); err != nil {
		s.logger.Println(err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	return &ServiceResponse{Message: "successfully " + message + " nginx"}, nil
}

func (s *SystemService) PhpFpm(ctx context.Context, req *PhpFpmServiceRequest) (*ServiceResponse, error) {
	// validate the request
	if err := validate.PHPVersion(req.GetVersion()); err != nil {
		s.logger.Println(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	var action string
	var message string
	switch req.GetAction() {
	case ServiceAction_START:
		message = "started"
		action = "start"
	case ServiceAction_STOP:
		message = "stopped"
		action = "stop"
	default:
		message = "restarted"
		action = "restart"
	}

	// perform the action on the php-fpm service
	if output, err := s.command.Run("service", []string{"php" + req.GetVersion() + "-fpm", action}); err != nil {
		s.logger.Println(err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	msg := "successfully " + message + " php-fpm " + req.GetVersion()

	s.logger.Println(msg)

	return &ServiceResponse{Message: msg}, nil
}

func NewSystemService() *SystemService {
	return &SystemService{
		command: &ServiceRunner{},
		logger:  log.New(os.Stdout, "nitrod ", 0),
	}
}
