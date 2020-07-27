package nitrod

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/validate"
)

func (s *Service) PhpFpmService(ctx context.Context, request *PhpFpmServiceRequest) (*ServiceResponse, error) {
	// validate the request
	if err := validate.PHPVersion(request.GetVersion()); err != nil {
		s.logger.Println(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	var action string
	var message string
	switch request.GetAction() {
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
	if output, err := s.command.Run("service", []string{"php" + request.GetVersion() + "-fpm", action}); err != nil {
		s.logger.Println(err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	msg := "successfully " + message + " php-fpm " + request.GetVersion()

	s.logger.Println(msg)

	return &ServiceResponse{Message: msg}, nil
}
