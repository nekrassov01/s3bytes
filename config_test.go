package s3bytes

import (
	"bytes"
	"context"
	"testing"

	"github.com/charmbracelet/log"
)

func TestLoadAWSConfig(t *testing.T) {
	type args struct {
		ctx      context.Context
		profile  string
		loglevel log.Level
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
				profile:  "",
				loglevel: log.DebugLevel,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				profile:  "invalid-profile",
				loglevel: log.DebugLevel,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			_, err := LoadAWSConfig(tt.args.ctx, w, tt.args.profile, tt.args.loglevel)
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
