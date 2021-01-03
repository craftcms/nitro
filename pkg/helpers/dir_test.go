package helpers

import "testing"

func TestMkdirIfNotExists(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "returns an error if a directory exists",
			args:    args{dir: "testdata/existing-dir"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MkdirIfNotExists(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("MkdirIfNotExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_dirExists(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "existing directory returns true",
			args: args{dir: "testdata/existing-dir"},
			want: true,
		},
		{
			name: "non-existing directory returns false",
			args: args{dir: "testdata/non-existing-dir"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DirExists(tt.args.dir); got != tt.want {
				t.Errorf("DirExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
