package hostedit

import (
	"os"
	"testing"
)

func TestGetSection(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name      string
		args      args
		wantStart int
		wantEnd   int
	}{
		{
			name:      "returns correct values when there is an existing section",
			args:      args{f: "testdata/has-section.txt"},
			wantStart: 11,
			wantEnd:   13,
		},
		{
			name:      "returns nil values when there is not an existing section",
			args:      args{f: "testdata/no-section.txt"},
			wantStart: 0,
			wantEnd:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.args.f)
			if err != nil {
				t.Fatal(err)
			}

			got, got1 := find(f)
			if got != tt.wantStart {
				t.Errorf("GetSection() got = %v, want %v", got, tt.wantStart)
			}
			if got1 != tt.wantEnd {
				t.Errorf("GetSection() got1 = %v, want %v", got1, tt.wantEnd)
			}
		})
	}
}
