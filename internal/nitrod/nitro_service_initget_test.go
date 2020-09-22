package nitrod

import (
	"context"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func TestNitrodService_GetPhpIniSetting(t *testing.T) {
	type fields struct {
		command Runner
		logger  *log.Logger
	}
	type args struct {
		ctx     context.Context
		request *GetPhpIniSettingRequest
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
			name: "can get the php ini setting for max_file_uploads",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &GetPhpIniSettingRequest{Version: "7.4", Setting: "max_file_uploads"},
			},
			want:         &ServiceResponse{Message: "The setting \"max_file_uploads\" is currently set to 512M"},
			wantErr:      false,
			wantCommands: []string{"bash"},
			wantArgs: []map[string][]string{
				{
					"bash": {"-c", "php-fpm7.4 -i | grep 'max_file_uploads'"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &spyChainRunner{
				Output: "memory_limit => 512M => 512M",
			}
			s := &NitroService{
				command: spy,
				logger:  tt.fields.logger,
			}
			got, err := s.GetPhpIniSetting(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPhpIniSetting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPhpIniSetting() got = %v, want %v", got, tt.want)
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
