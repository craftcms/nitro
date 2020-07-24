package api

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/validate"
)

func (s *NitrodService) PhpIniSettings(ctx context.Context, request *ChangePhpIniSettingRequest) (*ServiceResponse, error) {
	var setting string
	value := request.GetValue()
	version := request.GetVersion()

	// validate the php version
	if err := validate.PHPVersion(version); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// validate the setting and value
	switch request.GetSetting() {
	case PhpIniSetting_MAX_EXECUTION_TIME:
		if err := validate.MaxExecutionTime(value); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		setting = "max_execution_time"
	case PhpIniSetting_UPLOAD_MAX_FILESIZE:
		if err := validate.IsMegabytes(value); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		setting = "upload_max_filesize"
	case PhpIniSetting_MAX_INPUT_TIME:
		if _, err := strconv.Atoi(value); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		setting = "max_input_time"
	case PhpIniSetting_MAX_INPUT_VARS:
		if err := validate.MaxInputVars(value); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		setting = "max_input_vars"
	case PhpIniSetting_MAX_FILE_UPLOADS:
		if err := validate.PhpMaxFileUploads(value); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		setting = "max_file_uploads"
	case PhpIniSetting_MEMORY_LIMIT:
		if err := validate.IsMegabytes(value); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		setting = "memory_limit"
	default:
		msg := "changing " + PhpIniSetting_name[int32(request.GetSetting())] + " is not authorized"
		s.logger.Println(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}

	// change setting using sed
	if output, err := s.command.Run("sed", []string{"-i", "s|" + setting + "|" + setting + " = " + value + "|g", "/etc/php/" + version + "/fpm/php.ini"}); err != nil {
		s.logger.Println("error changing ini setting, error:", err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	// restart php-fpm using service
	if output, err := s.command.Run("service", []string{"php" + version + "-fpm", "restart"}); err != nil {
		s.logger.Println("error restarting php-fpm, error:", err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	return &ServiceResponse{Message: "successfully changed the ini setting for " + setting + " to " + value}, nil
}
