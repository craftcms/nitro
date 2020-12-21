package hostedit

import (
	"testing"
)

func TestUpdate(t *testing.T) {
	type args struct {
		file  string
		addr  string
		hosts []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "can append to the file when no section is found",
			args: args{file: "testdata/no-section.txt", addr: "127.0.0.1", hosts: []string{"one", "two", "three"}},
			want: `##
# Host Database
#
# localhost is used to configure the loopback interface
# when the system is booting.  Do not change this entry.
##
127.0.0.1        localhost
255.255.255.255  broadcasthost
::1              localhost

127.0.0.1        kubernetes.docker.internal
# Added by Docker Desktop
# To allow the same kube context to work on the host and the container:
127.0.0.1        kubernetes.docker.internal
# End of section

# <nitro>
127.0.0.1	one two three
# </nitro>
`,
		},
		{
			name: "can update the right section",
			args: args{file: "testdata/has-section.txt", addr: "127.0.0.1", hosts: []string{"one", "two", "three"}},
			want: `##
# Host Database
#
# localhost is used to configure the loopback interface
# when the system is booting.  Do not change this entry.
##
127.0.0.1        localhost
255.255.255.255  broadcasthost
::1              localhost

# <nitro>
127.0.0.1	one two three
# </nitro>

127.0.0.1        kubernetes.docker.internal
# Added by Docker Desktop
# To allow the same kube context to work on the host and the container:
127.0.0.1        kubernetes.docker.internal
# End of section
`,
		},
		{
			name:    "no file returns an error",
			args:    args{file: "testdata/empty"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Update(tt.args.file, tt.args.addr, tt.args.hosts...)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUpdated(t *testing.T) {
	type args struct {
		file  string
		addr  string
		hosts []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "returns true when nothing needs to occur",
			args: args{
				file:  "testdata/up-to-date.txt",
				addr:  "127.0.0.1",
				hosts: []string{"one", "two", "three"},
			},
			want: true,
		},
		{
			name: "returns false when file is not updated",
			args: args{
				file:  "testdata/up-to-date.txt",
				addr:  "127.0.0.1",
				hosts: []string{"one", "two", "three", "four"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsUpdated(tt.args.file, tt.args.addr, tt.args.hosts...)

			if (err != nil) != tt.wantErr {
				t.Errorf("IsUpdated() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("IsUpdated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexes(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name   string
		args   args
		start  int
		middle int
		end    int
	}{
		{
			name: "can find the start, middle, and end",
			args: args{content: []byte(`# <nitro>
127.0.0.1 host
# </nitro>`)},
			start:  0,
			middle: 1,
			end:    2,
		},
		{
			name: "can find the start, middle, and end in random",
			args: args{content: []byte(`something host


# <nitro>
127.0.0.1 host
# </nitro>`)},
			start:  3,
			middle: 4,
			end:    5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := indexes(tt.args.content)
			if got != tt.start {
				t.Errorf("indexes() got = %v, want %v", got, tt.start)
			}
			if got1 != tt.middle {
				t.Errorf("indexes() got1 = %v, want %v", got1, tt.middle)
			}
			if got2 != tt.end {
				t.Errorf("indexes() got2 = %v, want %v", got2, tt.end)
			}
		})
	}
}
