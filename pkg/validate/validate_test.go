package validate

import (
	"testing"
)

func TestIsMegabytes(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "must end in M",
			args:    args{v: "256"},
			wantErr: true,
		},
		{
			name:    "valid values pass",
			args:    args{v: "256M"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := isMegabytes(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("isMegabytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHostnameValidator_Validate(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		v       *HostnameValidator
		args    args
		wantErr bool
	}{
		{
			name: "valid hostnames do not return an err",
			args: args{
				input: "validhostname.tld",
			},
			wantErr: false,
		},
		{
			name: "special chars returns an err",
			args: args{
				input: "somehostname!",
			},
			wantErr: true,
		},
		{
			name: "spaces returns an err",
			args: args{
				input: "12    ",
			},
			wantErr: true,
		},
		{
			name: "less than 3 chars returns an err",
			args: args{
				input: "12",
			},
			wantErr: true,
		},
		{
			name: "comma separated list returns an err",
			args: args{
				input: "host1.tld,host.tld",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &HostnameValidator{}
			if err := v.Validate(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("HostnameValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPHPVersionValidator_Validate(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		v       *PHPVersionValidator
		args    args
		wantErr bool
	}{
		{
			name: "8.0 version is valid",
			args: args{
				input: "8.0",
			},
			wantErr: false,
		},
		{
			name: "7.4 version is valid",
			args: args{
				input: "7.4",
			},
			wantErr: false,
		},
		{
			name: "7.3 version is valid",
			args: args{
				input: "7.3",
			},
			wantErr: false,
		},
		{
			name: "7.2 version is valid",
			args: args{
				input: "7.2",
			},
			wantErr: false,
		},
		{
			name: "7.1 version is valid",
			args: args{
				input: "7.1",
			},
			wantErr: false,
		},
		{
			name: "7.0 version is valid",
			args: args{
				input: "7.0",
			},
			wantErr: false,
		},
		{
			name: "invalid versions returns an error",
			args: args{
				input: "nope",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &PHPVersionValidator{}
			if err := v.Validate(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("PHPVersionValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsBoolean_Validate(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		v       *IsBoolean
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &IsBoolean{}
			if err := v.Validate(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("IsBoolean.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
