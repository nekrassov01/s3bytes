package s3bytes

import (
	"bytes"
	"testing"

	"github.com/aws/smithy-go/logging"
	"github.com/charmbracelet/log"
)

func Test_newAppLogger(t *testing.T) {
	tests := []struct {
		name     string
		logLevel log.Level
		msg      string
		want     string
	}{
		{
			name:     "debug",
			logLevel: log.DebugLevel,
			msg:      "This is an debug message",
			want:     "DBG S3BYTES: This is an debug message\n",
		},
		{
			name:     "info",
			logLevel: log.InfoLevel,
			msg:      "This is an info message",
			want:     "INF S3BYTES: This is an info message\n",
		},
		{
			name:     "warn",
			logLevel: log.WarnLevel,
			msg:      "This is a warning message",
			want:     "WRN S3BYTES: This is a warning message\n",
		},
		{
			name:     "error",
			logLevel: log.ErrorLevel,
			msg:      "This is a error message",
			want:     "ERR S3BYTES: This is a error message\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			logger := newAppLogger(w)
			logger.SetLevel(tt.logLevel)
			switch tt.logLevel {
			case log.DebugLevel:
				logger.Debug(tt.msg)
			case log.InfoLevel:
				logger.Info(tt.msg)
			case log.WarnLevel:
				logger.Warn(tt.msg)
			case log.ErrorLevel:
				logger.Error(tt.msg)
			}
			if got := w.String(); got != tt.want {
				t.Errorf("newAppLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newSDKLogger(t *testing.T) {
	type args struct {
		loglevel log.Level
	}
	tests := []struct {
		name   string
		args   args
		class  logging.Classification
		format string
		v      []any
		want   string
	}{
		{
			name:   "warn",
			class:  logging.Warn,
			args:   args{loglevel: log.WarnLevel},
			format: "Warning: %s",
			v:      []any{"something happened"},
			want:   "WRN SDK: Warning: something happened\n",
		},
		{
			name:   "debug",
			class:  logging.Debug,
			args:   args{loglevel: log.DebugLevel},
			format: "Debugging: %s",
			v:      []any{"details"},
			want:   "DBG SDK: Debugging: details\n",
		},
		{
			name:   "default",
			class:  logging.Classification("other"),
			args:   args{loglevel: log.InfoLevel},
			format: "Default: %s",
			v:      []any{"something"},
			want:   "INF SDK: Default: something\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			logger := newSDKLogger(w, tt.args.loglevel)
			logger.Logf(tt.class, tt.format, tt.v...)
			if got := w.String(); got != tt.want {
				t.Errorf("newSDKLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
