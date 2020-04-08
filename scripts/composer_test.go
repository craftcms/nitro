package scripts

import "testing"

func TestInstallComposer(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "returns the script to install composer",
			want: `php -r "readfile('http://getcomposer.org/installer')" | sudo php -- --install-dir=/usr/bin/ --filename=composer`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InstallComposer(); got != tt.want {
				t.Errorf("InstallComposer() = %v, want %v", got, tt.want)
			}
		})
	}
}
