package envedit

import (
	"io"
	"os"
	"testing"
)

func TestSet(t *testing.T) {
	// setup
	testEnvfile := "./testdata/.env"
	if err := os.Link("./testdata/.env.example", testEnvfile); err != nil {
		t.Fatal(err)
	}

	// cleanup
	if err := os.Remove(testEnvfile); err != nil {
		t.Fatal(err)
	}
}

func TestFileEditer_Set(t *testing.T) {
	type fields struct {
		file *os.File
		rw   io.ReadWriter
	}
	type args struct {
		env string
		val string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "can write to the",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &FileEditer{
				file: tt.fields.file,
			}
		})
	}
}
