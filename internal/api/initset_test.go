package api

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
			wantCommands: []string{"sed"},
			wantArgs: []map[string][]string{{
				"sed": []string{"-i", "s|max_execution_time|max_execution_time = 300|g", "/etc/php/7.4/fpm/php.ini"},
			}},
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
			s := &NitrodService{
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
				t.Errorf("expected the commands to be:\n%v\n, got:\n%v", spy.Commands, tt.wantCommands)
			}
			if !reflect.DeepEqual(spy.Args, tt.wantArgs) {
				t.Errorf("expected the args to be:\n%v\ngot:\n%v", spy.Args, tt.wantArgs)
			}
		})
	}
}
