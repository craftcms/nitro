package nitrod

import (
	"log"
	"testing"
)

func TestNitroService_ImportDatabase(t *testing.T) {
	type fields struct {
		command Runner
		logger  *log.Logger
	}
	type args struct {
		stream NitroService_ImportDatabaseServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &NitroService{
				command: tt.fields.command,
				logger:  tt.fields.logger,
			}
			if err := s.ImportDatabase(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("ImportDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
