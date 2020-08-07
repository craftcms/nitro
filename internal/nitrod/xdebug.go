package nitrod

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/internal/validate"
)

func (s *NitroService) DisableXdebug(ctx context.Context, req *DisableXdebugRequest) (*ServiceResponse, error) {
	// validate the php version
	if err := validate.PHPVersion(req.GetVersion()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// disable xdebug
	if output, err := s.command.Run("phpdismod", []string{"-v", "7.4", "xdebug"}); err != nil {
		s.logger.Println("error disabling xdebug, error:", err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	// restart php-fpm using service
	if output, err := s.command.Run("service", []string{"php" + req.GetVersion() + "-fpm", "restart"}); err != nil {
		s.logger.Println("error restarting php-fpm, error:", err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	return &ServiceResponse{Message: "disabled xdebug for PHP " + req.GetVersion()}, nil
}

func (s *NitroService) EnableXdebug(ctx context.Context, req *EnableXdebugRequest) (*ServiceResponse, error) {
	// validate the php version
	if err := validate.PHPVersion(req.GetVersion()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// enable xdebug
	if output, err := s.command.Run("phpenmod", []string{"-v", "7.4", "xdebug"}); err != nil {
		s.logger.Println("error enabling xdebug, error:", err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	// restart php-fpm using service
	if output, err := s.command.Run("service", []string{"php" + req.GetVersion() + "-fpm", "restart"}); err != nil {
		s.logger.Println("error restarting php-fpm, error:", err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	return &ServiceResponse{Message: "enabled xdebug for PHP " + req.GetVersion()}, nil
}
