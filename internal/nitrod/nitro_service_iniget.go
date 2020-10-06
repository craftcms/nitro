package nitrod

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/internal/validate"
)

func (s *NitroService) GetPhpIniSetting(ctx context.Context, req *GetPhpIniSettingRequest) (*ServiceResponse, error) {
	// validate php version
	if err := validate.PHPVersion(req.GetVersion()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// get the setting from the php ini_get function
	output, err := s.command.Run("bash", []string{"-c", fmt.Sprintf("php-fpm%s -i | grep '%s'", req.GetVersion(), req.GetSetting())})
	if err != nil {
		s.logger.Println("error getting ini setting:", err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	// get the output in a normal format
	var settingValue string
	sp := strings.Split(string(output), " ")
	if len(sp) == 5 {
		settingValue = sp[len(sp)-1]
	}

	if settingValue == "" {
		return &ServiceResponse{Message: fmt.Sprintf("Unable to find the PHP setting %q", req.GetSetting())}, nil
	}

	return &ServiceResponse{Message: fmt.Sprintf("The setting %q is currently set to %v", req.GetSetting(), settingValue)}, nil
}
