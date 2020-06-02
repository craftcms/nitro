package validate

import (
	"testing"
)

func TestPHPVersionFlag(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "7.4",
			args: args{
				"7.4",
			},
			wantErr: false,
		},
		{
			name: "7.3",
			args: args{
				"7.3",
			},
			wantErr: false,
		},
		{
			name: "7.2",
			args: args{
				"7.2",
			},
			wantErr: false,
		},
		{
			name: "junk",
			args: args{
				"junk",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PHPVersion(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("PHPVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDomain(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "is valid",
			args:    args{v: "example.test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Hostname(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Hostname() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "integers only return error",
			args:    args{v: "2"},
			wantErr: true,
		},
		{
			name:    "lower case G return error",
			args:    args{v: "2g"},
			wantErr: true,
		},
		{
			name:    "proper values do not return error",
			args:    args{v: "2G"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Memory(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Memory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPath(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid path returns nil",
			args:    args{p: "/tmp"},
			wantErr: false,
		},
		{
			name:    "invalid path returns error",
			args:    args{p: "does-not-exist"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Path(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("Path() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
