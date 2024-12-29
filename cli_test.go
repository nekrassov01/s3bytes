package s3bytes

import (
	"context"
	"io"
	"testing"
)

func Test_cli(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "completion bash",
			args:    []string{appName, "-c", bash.String()},
			wantErr: false,
		},
		{
			name:    "completion zsh",
			args:    []string{appName, "-c", zsh.String()},
			wantErr: false,
		},
		{
			name:    "completion pwsh",
			args:    []string{appName, "-c", pwsh.String()},
			wantErr: false,
		},
		{
			name:    "completion unknown",
			args:    []string{appName, "-c", "fish"},
			wantErr: true,
		},
		{
			name:    "unknown profile",
			args:    []string{appName, "-p", "unknown"},
			wantErr: true,
		},
		// In CI/CD, attempting to ListBuckets fails due to access denial, so exclude from testing
		// {
		// 	name:    "unknown log level",
		// 	args:    []string{appName, "-l", "unknown"},
		// 	wantErr: false, // if the log level is invalid, default to info
		// },
		{
			name:    "unknown region",
			args:    []string{appName, "-r", "unknown"},
			wantErr: true,
		},
		{
			name:    "unknown metric name",
			args:    []string{appName, "-m", "unknown"},
			wantErr: true,
		},
		{
			name:    "unknown storage type",
			args:    []string{appName, "-s", "unknown"},
			wantErr: true,
		},
		{
			name:    "unknown output type",
			args:    []string{appName, "-o", "unknown"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newApp(io.Discard, io.Discard).RunContext(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestApp_sort(t *testing.T) {
	app := &app{}
	metrics := []Metric{
		&SizeMetric{
			BucketName:    "bucket-c",
			Region:        "ap-northeast-1",
			StorageType:   StorageTypeStandardStorage,
			Bytes:         150,
			ReadableBytes: "150B",
		},
		&SizeMetric{
			BucketName:    "bucket-a",
			Region:        "ap-northeast-1",
			StorageType:   StorageTypeStandardStorage,
			Bytes:         200,
			ReadableBytes: "200B",
		},
		&SizeMetric{
			BucketName:    "bucket-b",
			Region:        "ap-northeast-1",
			StorageType:   StorageTypeStandardStorage,
			Bytes:         200,
			ReadableBytes: "200B",
		},
		&SizeMetric{
			BucketName:    "bucket-d",
			Region:        "ap-northeast-1",
			StorageType:   StorageTypeStandardStorage,
			Bytes:         100,
			ReadableBytes: "100B",
		},
	}
	sorted := []Metric{
		&SizeMetric{
			BucketName:    "bucket-a",
			Region:        "ap-northeast-1",
			StorageType:   StorageTypeStandardStorage,
			Bytes:         200,
			ReadableBytes: "200 B",
		},
		&SizeMetric{
			BucketName:    "bucket-b",
			Region:        "ap-northeast-1",
			StorageType:   StorageTypeStandardStorage,
			Bytes:         200,
			ReadableBytes: "200 B",
		},
		&SizeMetric{
			BucketName:    "bucket-c",
			Region:        "ap-northeast-1",
			StorageType:   StorageTypeStandardStorage,
			Bytes:         150,
			ReadableBytes: "150 B",
		},
		&SizeMetric{
			BucketName:    "bucket-d",
			Region:        "ap-northeast-1",
			StorageType:   StorageTypeStandardStorage,
			Bytes:         100,
			ReadableBytes: "100 B",
		},
	}
	app.sort(metrics)
	for i, metric := range metrics {
		var (
			got  = metric.(*SizeMetric)
			want = sorted[i].(*SizeMetric)
		)
		if got.Bytes != want.Bytes {
			t.Errorf("Metric[%d] Bytes = %v, want %v", i, got.Bytes, want.Bytes)
		}
		if got.BucketName != want.BucketName {
			t.Errorf("Metric[%d] BucketName = %v, want %v", i, got.BucketName, want.BucketName)
		}
	}
}
