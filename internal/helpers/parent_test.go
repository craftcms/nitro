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
			name: "can get the parent directory name",
			args: args{
				path: "./testdata/good-example/web",
			},
			want:    "good-example",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParentPathName(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParentPathName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParentPathName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
