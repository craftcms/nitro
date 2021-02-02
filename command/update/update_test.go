package update

import "testing"

func Test_versionFromName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "can get the php version from the name",
			args: args{
				name: "nginx:7.3-dev",
			},
			want: "7.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := versionFromName(tt.args.name); got != tt.want {
				t.Errorf("versionFromName() = %v, want %v", got, tt.want)
			}
		})
	}
}
