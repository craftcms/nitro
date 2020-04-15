package helpers

import "testing"

func TestParentPathName(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "can get the path directory name",
			args: args{
				path: "./testdata/good-example",
			},
			want:    "good-example",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PathName(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("PathName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PathName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
