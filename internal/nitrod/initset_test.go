package nitrod

import (
	"context"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func TestNitrodService_PhpIniSettings(t *testing.T) {
	type fields struct {
		command Runner
		logger  *log.Logger
	}
	type args struct {
		ctx     context.Context
		request *ChangePhpIniSettingRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         *ServiceResponse
		wantErr      bool
		wantCommands []string
		wantArgs     []map[string][]string
	}{
		{
			name: "can modify the php ini setting for memory_limit",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &ChangePhpIniSettingRequest{Version: "7.4", Setting: PhpIniSetting_MEMORY_LIMIT, Value: "512M"},
			},
			want:         &ServiceResponse{Message: "successfully changed the ini setting for memory_limit to 512M"},
			wantErr:      false,
			wantCommands: []string{"sed", "service"},
			wantArgs: []map[string][]string{
				{
					"sed": {"-i", "s|memory_limit|memory_limit = 512M|g", "/etc/php/7.4/fpm/php.ini"},
				},
				{
					"service": {"php7.4-fpm", "restart"},
				},
			},
		},
		{
			name: "can modify the php ini setting for max_file_uploads",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &ChangePhpIniSettingRequest{Version: "7.4", Setting: PhpIniSetting_MAX_FILE_UPLOADS, Value: "400"},
			},
			want:         &ServiceResponse{Message: "successfully changed the ini setting for max_file_uploads to 400"},
			wantErr:      false,
			wantCommands: []string{"sed", "service"},
			wantArgs: []map[string][]string{
				{
					"sed": {"-i", "s|max_file_uploads|max_file_uploads = 400|g", "/etc/php/7.4/fpm/php.ini"},
				},
				{
					"service": {"php7.4-fpm", "restart"},
				},
			},
		},
		{
			name: "can modify the php ini setting for max_input_time",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &ChangePhpIniSettingRequest{Version: "7.4", Setting: PhpIniSetting_MAX_INPUT_TIME, Value: "4000"},
			},
			want:         &ServiceResponse{Message: "successfully changed the ini setting for max_input_time to 4000"},
			wantErr:      false,
			wantCommands: []string{"sed", "service"},
			wantArgs: []map[string][]string{
				{
					"sed": {"-i", "s|max_input_time|max_input_time = 4000|g", "/etc/php/7.4/fpm/php.ini"},
				},
				{
					"service": {"php7.4-fpm", "restart"},
				},
			},
		},
		{
			name: "can modify the php ini setting for upload_max_filesize",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &ChangePhpIniSettingRequest{Version: "7.4", Setting: PhpIniSetting_UPLOAD_MAX_FILESIZE, Value: "10M"},
			},
			want:         &ServiceResponse{Message: "successfully changed the ini setting for upload_max_filesize to 10M"},
			wantErr:      false,
			wantCommands: []string{"sed", "service"},
			wantArgs: []map[string][]string{
				{
					"sed": {"-i", "s|upload_max_filesize|upload_max_filesize = 10M|g", "/etc/php/7.4/fpm/php.ini"},
				},
				{
					"service": {"php7.4-fpm", "restart"},
				},
			},
		},
		{
			name: "can modify the php ini setting for upload_max_filesize",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &ChangePhpIniSettingRequest{Version: "7.4", Setting: PhpIniSetting_UPLOAD_MAX_FILESIZE, Value: "10M"},
			},
			want:         &ServiceResponse{Message: "successfully changed the ini setting for upload_max_filesize to 10M"},
			wantErr:      false,
			wantCommands: []string{"sed", "service"},
			wantArgs: []map[string][]string{
				{
					"sed": {"-i", "s|upload_max_filesize|upload_max_filesize = 10M|g", "/etc/php/7.4/fpm/php.ini"},
				},
				{
					"service": {"php7.4-fpm", "restart"},
				},
			},
		},
		{
			name: "can modify the php ini setting for max_input_vars",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &ChangePhpIniSettingRequest{Version: "7.4", Setting: PhpIniSetting_MAX_INPUT_VARS, Value: "1000"},
			},
			want:         &ServiceResponse{Message: "successfully changed the ini setting for max_input_vars to 1000"},
			wantErr:      false,
			wantCommands: []string{"sed", "service"},
			wantArgs: []map[string][]string{
				{
					"sed": {"-i", "s|max_input_vars|max_input_vars = 1000|g", "/etc/php/7.4/fpm/php.ini"},
				},
				{
					"service": {"php7.4-fpm", "restart"},
				},
			},
		},
		{
			name: "setting max_input_vars to a non-integer returns an error",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &ChangePhpIniSettingRequest{Version: "7.4", Setting: PhpIniSetting_MAX_INPUT_VARS, Value: "300b"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "setting max_input_vars must be less than 10000",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &ChangePhpIniSettingRequest{Version: "7.3", Setting: PhpIniSetting_MAX_INPUT_VARS, Value: "10000"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "can modify the php ini setting for max_execution_time",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &ChangePhpIniSettingRequest{Version: "7.4", Setting: PhpIniSetting_MAX_EXECUTION_TIME, Value: "300"},
			},
			want:         &ServiceResponse{Message: "successfully changed the ini setting for max_execution_time to 300"},
			wantErr:      false,
			wantCommands: []string{"sed", "service"},
			wantArgs: []map[string][]string{
				{
					"sed": {"-i", "s|max_execution_time|max_execution_time = 300|g", "/etc/php/7.4/fpm/php.ini"},
				},
				{
					"service": {"php7.4-fpm", "restart"},
				},
			},
		},
		{
			name: "setting max_execution_time to a non-integer returns an error",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &ChangePhpIniSettingRequest{Version: "7.4", Setting: PhpIniSetting_MAX_EXECUTION_TIME, Value: "300b"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &spyChainRunner{}
			s := &Service{
				command: spy,
				logger:  tt.fields.logger,
			}
			got, err := s.PhpIniSettings(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("PhpIniSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PhpIniSettings() got = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(spy.Commands, tt.wantCommands) {
				t.Errorf("expected the commands to be:\n%v\n, got:\n%v", tt.wantCommands, spy.Commands)
			}

			if !reflect.DeepEqual(spy.Args, tt.wantArgs) {
				t.Errorf("expected the args to be:\n%v\ngot:\n%v", tt.wantArgs, spy.Args)
			}
		})
	}
}
