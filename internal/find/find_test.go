package find

import (
	"reflect"
	"testing"

	"github.com/craftcms/nitro/config"
)

func TestMounts(t *testing.T) {
	type args struct {
		name string
		b    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []config.Mount
		wantErr bool
	}{
		{
			name: "returns mounts from the machine",
			args: args{name: "nitro-dev", b: []byte("Name,State,Ipv4,Ipv6,Release,Image hash,Image release,Load,Disk usage,Disk total,Memory usage,Memory total,Mounts\nnitro-dev,Running,192.168.64.21,,Ubuntu 18.04.4 LTS,3b2e3aaaebf2bc364da70fbc7e9619a7c0bb847932496a1903cd4913cf9b1a26,18.04 LTS,0.01 0.21 0.22,2544762880,41442029568,393465856,4136783872,/Users/jasonmccallister/go/src/github.com/craftcms/nitro/production-site => /nitro/sites/production-site;/Users/jasonmccallister/go/src/github.com/craftcms/nitro/demo-site => /nitro/sites/demo-site;\n")},
			want: []config.Mount{
				{
					Source: "/Users/jasonmccallister/go/src/github.com/craftcms/nitro/production-site",
					Dest:   "/nitro/sites/production-site",
				},
				{
					Source: "/Users/jasonmccallister/go/src/github.com/craftcms/nitro/demo-site",
					Dest:   "/nitro/sites/demo-site",
				},
			},
			wantErr: false,
		},
		{
			name:    "returns no mounts from the machine if they are not present",
			args:    args{name: "nitro-dev", b: []byte("Name,State,Ipv4,Ipv6,Release,Image hash,Image release,Load,Disk usage,Disk total,Memory usage,Memory total,Mounts\n")},
			want:    nil,
			wantErr: false,
		},
		{
			name: "sources with a space in the path still return a mount",
			args: args{name: "nitro-dev", b: []byte("Name,State,Ipv4,Ipv6,Release,Image hash,Image release,Load,Disk usage,Disk total,Memory usage,Memory total,Mounts\nnitro-dev,Running,192.168.64.21,,Ubuntu 18.04.4 LTS,3b2e3aaaebf2bc364da70fbc7e9619a7c0bb847932496a1903cd4913cf9b1a26,18.04 LTS,0.01 0.21 0.22,2544762880,41442029568,393465856,4136783872,/Users/jason mccallister/go/src/github.com/craftcms/nitro/production-site => /nitro/sites/production-site;\n")},
			want: []config.Mount{
				{
					Source: "/Users/jason mccallister/go/src/github.com/craftcms/nitro/production-site",
					Dest:   "/nitro/sites/production-site",
				},
			},
			wantErr: false,
		},
		{
			name: "sources with a space at the end return a mount without the space",
			args: args{name: "nitro-dev", b: []byte("Name,State,Ipv4,Ipv6,Release,Image hash,Image release,Load,Disk usage,Disk total,Memory usage,Memory total,Mounts\nnitro-dev,Running,192.168.64.21,,Ubuntu 18.04.4 LTS,3b2e3aaaebf2bc364da70fbc7e9619a7c0bb847932496a1903cd4913cf9b1a26,18.04 LTS,0.01 0.21 0.22,2544762880,41442029568,393465856,4136783872,/Users/jasonmccallister/go/src/github.com/craftcms/nitro/production-site  => /nitro/sites/production-site ;\n")},
			want: []config.Mount{
				{
					Source: "/Users/jasonmccallister/go/src/github.com/craftcms/nitro/production-site",
					Dest:   "/nitro/sites/production-site",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mounts(tt.args.name, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				if len(tt.want) != len(got) {
					t.Errorf("expected the length to match, got %d but want %d", len(got), len(tt.want))
				}

				for i, mount := range tt.want {
					if !reflect.DeepEqual(mount, got[i]) {
						t.Errorf("Mounts() got = \n%v, \nwant \n%v", got[i], tt.want[i])
					}
				}
			}
		})
	}
}
