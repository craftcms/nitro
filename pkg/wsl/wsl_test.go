package wsl

import (
	"os"
	"testing"
)

func TestIsWSL(t *testing.T) {
	tests := []struct {
		name  string
		env   string
		empty bool
		want  bool
	}{
		{
			name:  "returns true if the environment variables are empty",
			env:   "WSL_DISTRO_NAME",
			empty: true,
			want:  true,
		},
		{
			name: "returns false if the environment variables are not set",
			want: false,
		},
		{
			name: "returns true if the environment variable is set",
			env:  "WSLENV",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.env) > 0 {
				switch tt.empty {
				case true:
					os.Setenv(tt.env, "")
				default:
					os.Setenv(tt.env, "somerandomthing")
				}

				defer os.Unsetenv(tt.env)
			}

			if got := IsWSL(); got != tt.want {
				t.Errorf("IsWSL() = %v, want %v", got, tt.want)
			}
		})
	}
}
