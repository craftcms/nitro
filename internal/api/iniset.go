package api

import (
	"context"
	"errors"

	"github.com/craftcms/nitro/validate"
)

func (s *NitrodService) PhpIniSettings(ctx context.Context, request *ChangePhpIniSettingRequest) (*ServiceResponse, error) {
	var setting string
	switch request.GetSetting() {
	case PhpIniSetting_MAX_EXECUTION_TIME:
		if err := validate.MaxExecutionTime(request.GetValue()); err != nil {
			return nil, err
		}
		setting = "max_execution_time"
	case PhpIniSetting_MAX_INPUT_VARS:
		if err := validate.MaxInputVars(request.GetValue()); err != nil {
			return nil, err
		}
		setting = "max_input_vars"
	default:
		e := errors.New("changing this setting is not authorized")
		s.logger.Println(e)
		return nil, e
	}

	if output, err := s.command.Run("sed", []string{"-i", "s|" + setting + "|" + setting + " = " + request.GetValue() + "|g", "/etc/php/" + request.GetVersion() + "/fpm/php.ini"}); err != nil {
		s.logger.Println("error changing ini setting, error:", err)
		s.logger.Println("output:", string(output))
		return nil, err
	}

	if output, err := s.command.Run("service", []string{"php" + request.GetVersion() + "-fpm", "restart"}); err != nil {
		s.logger.Println("error restarting php-fpm, error:", err)
		s.logger.Println("output:", string(output))
		return nil, err
	}

	return &ServiceResponse{Message: "successfully changed the ini setting for " + setting + " to " + request.GetValue()}, nil
}
