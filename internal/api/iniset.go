package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/validate"
)

func (s *NitrodService) PhpIniSettings(ctx context.Context, request *ChangePhpIniSettingRequest) (*ServiceResponse, error) {
	var setting string

	switch request.GetSetting() {
	case PhpIniSetting_MAX_EXECUTION_TIME:
		if err := validate.MaxExecutionTime(request.GetValue()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		setting = "max_execution_time"
	case PhpIniSetting_MAX_INPUT_VARS:
		if err := validate.MaxInputVars(request.GetValue()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		setting = "max_input_vars"
	default:
		msg := "changing " + PhpIniSetting_name[int32(request.GetSetting())] + " setting is not authorized"
		s.logger.Println(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}

	// change setting using sed
	if output, err := s.command.Run("sed", []string{"-i", "s|" + setting + "|" + setting + " = " + request.GetValue() + "|g", "/etc/php/" + request.GetVersion() + "/fpm/php.ini"}); err != nil {
		s.logger.Println("error changing ini setting, error:", err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	// restart php-fpm using service
	if output, err := s.command.Run("service", []string{"php" + request.GetVersion() + "-fpm", "restart"}); err != nil {
		s.logger.Println("error restarting php-fpm, error:", err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	return &ServiceResponse{Message: "successfully changed the ini setting for " + setting + " to " + request.GetValue()}, nil
}
