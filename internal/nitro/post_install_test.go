package nitro

import (
	"reflect"
	"testing"
)

func TestConfigurePHPMemoryLimit(t *testing.T) {
	type args struct {
		name  string
		php   string
		limit string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "run the sed command for a specific version",
			args: args{
				name:  "somename",
				php:   "7.4",
				limit: "256M",
			},
			want: &Action{
				Type:       "exec",
				Output:     "Configuring PHP 7.4 memory limit to 256M",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "sudo", "sed", "-i", "s|memory_limit = 128M|memory_limit = 256M|g", "/etc/php/7.4/fpm/php.ini"},
			},
			wantErr: false,
		},
		{
			name: "wrong php returns error",
			args: args{
				name:  "somename",
				php:   "7.9",
				limit: "256M",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad name returns error",
			args: args{
				name:  "",
				php:   "7.9",
				limit: "256M",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConfigurePHPMemoryLimit(tt.args.name, tt.args.php, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigurePHPMemoryLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigurePHPMemoryLimit() got = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}

func TestConfigurePHPExecutionTimeLimit(t *testing.T) {
	type args struct {
		name  string
		php   string
		limit string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "run the sed command for a specific version",
			args: args{
				name:  "somename",
				php:   "7.4",
				limit: "240",
			},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "sudo", "sed", "-i", "s|max_execution_time = 30|max_execution_time = 240|g", "/etc/php/7.4/fpm/php.ini"},
			},
			wantErr: false,
		},
		{
			name: "wrong php returns error",
			args: args{
				name:  "somename",
				php:   "7.9",
				limit: "240",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad name returns error",
			args: args{
				name:  "",
				php:   "7.9",
				limit: "240",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConfigurePHPExecutionTimeLimit(tt.args.name, tt.args.php, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigurePHPExecutionTimeLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigurePHPExecutionTimeLimit() got = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}

func TestConfigureXdebug(t *testing.T) {
	type args struct {
		name string
		php  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "returns the printf command to set the config",
			args: args{
				name: "somename",
				php:  "7.4",
			},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "sudo", "cp", "/opt/nitro/php-xdebug.ini", "/etc/php/7.4/mods-available/xdebug.ini"},
			},
			wantErr: false,
		},
		{
			name: "bad name returns error",
			args: args{
				name: "",
				php:  "7.4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad php version returns error",
			args: args{
				name: "somename",
				php:  "7.9",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConfigureXdebug(tt.args.name, tt.args.php)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigureXdebug() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigureXdebug() got = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}
