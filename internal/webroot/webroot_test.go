package webroot

import "testing"

func TestFindWebRoot(t *testing.T) {
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
			name: "can find the web dir",
			args: args{
				path: "./testdata/good-example",
			},
			want: "web",
			wantErr: false,
		},
		{
			name: "can find the public dir",
			args: args{
				path: "./testdata/public-example",
			},
			want: "public",
			wantErr: false,
		},
		{
			name: "can find the public_html dir",
			args: args{
				path: "./testdata/public_html-example",
			},
			want: "public_html",
			wantErr: false,
		},
		{
			name: "can find the www dir",
			args: args{
				path: "./testdata/www-example",
			},
			want: "www",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Find(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}
