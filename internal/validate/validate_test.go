package validate

import "testing"

func TestPHPVersionFlag(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "7.4",
			args: args{
				"7.4",
			},
			wantErr: false,
		},
		{
			name: "7.3",
			args: args{
				"7.3",
			},
			wantErr: false,
		},
		{
			name: "7.2",
			args: args{
				"7.2",
			},
			wantErr: false,
		},
		{
			name: "junk",
			args: args{
				"junk",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PHPVersion(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("PHPVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseFlag(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "mariadb",
			args:    args{v: "mariadb"},
			wantErr: false,
		},
		{
			name:    "maria",
			args:    args{v: "maria"},
			wantErr: true,
		},
		{
			name:    "postgres",
			args:    args{v: "postgres"},
			wantErr: false,
		},
		{
			name:    "postgresql",
			args:    args{v: "postgresql"},
			wantErr: true,
		},
		{
			name:    "pgsql",
			args:    args{v: "pgsql"},
			wantErr: true,
		},
		{
			name:    "cockroach",
			args:    args{v: "cockroach"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Database(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Database() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
