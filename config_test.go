package s3bytes

import (
	"bytes"
	"context"
	"testing"

	"github.com/aws/smithy-go/logging"
)

func TestLoadAWSConfig(t *testing.T) {
	type args struct {
		ctx     context.Context
		profile string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				profile: "",
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				profile: "invalid-profile",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			_, err := LoadAWSConfig(tt.args.ctx, tt.args.profile, logging.NewStandardLogger(w), logMode)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAWSConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.want {
				t.Errorf("LoadAWSConfig() = %v, want %v", gotW, tt.want)
			}
		})
	}
}
