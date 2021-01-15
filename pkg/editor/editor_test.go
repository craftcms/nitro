package editor

import (
	"os"
	"reflect"
	"testing"
)

func TestGetPreferredEditorFromEnvironment(t *testing.T) {
	type args struct {
		goos string
		env  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "linux returns default",
			args: args{goos: "linux"},
			want: "vim",
		},
		{
			name: "linux returns editor from env",
			args: args{goos: "linux", env: "nano"},
			want: "nano",
		},
		{
			name: "windows returns notepad",
			args: args{goos: "windows"},
			want: "notepad.exe",
		},
	}
	for _, tt := range tests {
		// set the env if defined
		if tt.args.env != "" {
			os.Setenv("EDITOR", tt.args.env)
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := GetPreferredEditorFromEnvironment(tt.args.goos); got != tt.want {
				t.Errorf("GetPreferredEditorFromEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resolveEditorArguments(t *testing.T) {
	type args struct {
		executable string
		filename   string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "vscode appends args",
			args: args{executable: "Visual Studio Code.app", filename: "some-file"},
			want: []string{"--wait", "some-file"},
		},
		{
			name: "returns args",
			args: args{executable: "someexec", filename: "some-file"},
			want: []string{"some-file"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveEditorArguments(tt.args.executable, tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveEditorArguments() = %v, want %v", got, tt.want)
			}
		})
	}
}
