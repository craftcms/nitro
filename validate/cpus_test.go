package validate

import "testing"

func TestValidCPUCount_Validate(t *testing.T) {
	type fields struct {
		actual int
	}
	type args struct {
		requested string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "if we don't have the host count return an error",
			fields:  fields{actual: 0},
			args:    args{requested: "6"},
			wantErr: true,
		},
		{
			name:    "setting a CPU count higer than the machine returns error",
			fields:  fields{actual: 4},
			args:    args{requested: "6"},
			wantErr: true,
		},
		{
			name:    "setting a CPU count to the full number of cores available will return an error",
			fields:  fields{actual: 4},
			args:    args{requested: "4"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ValidCPUCount{
				Actual: tt.fields.actual,
			}
			if err := c.Validate(tt.args.requested); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
