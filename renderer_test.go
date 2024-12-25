package s3bytes

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/dustin/go-humanize"
)

var sampleSizeMetrics = []Metric{
	&SizeMetric{
		BucketName:    "bucket0",
		Region:        "ap-northeast-1",
		StorageType:   StorageTypeStandardStorage,
		Bytes:         1024,
		ReadableBytes: humanize.IBytes(1024),
	},
	&SizeMetric{
		BucketName:    "bucket1",
		Region:        "ap-northeast-2",
		StorageType:   StorageTypeGlacierStorage,
		Bytes:         4096,
		ReadableBytes: humanize.IBytes(4096),
	},
}

var sampleObjectMetrics = []Metric{
	&ObjectMetric{
		BucketName:  "bucket0",
		Region:      "ap-northeast-1",
		StorageType: StorageTypeAllStorageTypes,
		Objects:     20,
	},
	&ObjectMetric{
		BucketName:  "bucket1",
		Region:      "ap-northeast-2",
		StorageType: StorageTypeAllStorageTypes,
		Objects:     0,
	},
}

func TestNewRenderer(t *testing.T) {
	type args struct {
		metrics    []Metric
		metricName MetricName
		outputType OutputType
	}
	tests := []struct {
		name string
		args args
		want *Renderer
	}{
		{
			name: "normal",
			args: args{
				metrics:    sampleSizeMetrics,
				metricName: MetricNameBucketSizeBytes,
				outputType: OutputTypeJSON,
			},
			want: &Renderer{
				Metrics:    sampleSizeMetrics,
				MetricName: MetricNameBucketSizeBytes,
				OutputType: OutputTypeJSON,
				w:          &bytes.Buffer{},
			},
		},
		{
			name: "empty",
			args: args{
				metrics:    []Metric{},
				metricName: MetricNameNumberOfObjects,
				outputType: OutputTypeText,
			},
			want: &Renderer{
				Metrics:    []Metric{},
				MetricName: MetricNameNumberOfObjects,
				OutputType: OutputTypeText,
				w:          &bytes.Buffer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got := NewRenderer(w, tt.args.metrics, tt.args.metricName, tt.args.outputType)
			if !reflect.DeepEqual(got.Metrics, tt.want.Metrics) {
				t.Errorf("NewRenderer() Metrics = %v, want %v", got.Metrics, tt.want.Metrics)
			}
			if got.MetricName != tt.want.MetricName {
				t.Errorf("NewRenderer() MetricName = %v, want %v", got.MetricName, tt.want.MetricName)
			}
			if got.OutputType != tt.want.OutputType {
				t.Errorf("NewRenderer() OutputType = %v, want %v", got.OutputType, tt.want.OutputType)
			}
		})
	}
}

func TestRenderer_String(t *testing.T) {
	type fields struct {
		Metrics    []Metric
		MetricName MetricName
		OutputType OutputType
		w          io.Writer
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal",
			fields: fields{
				Metrics:    sampleSizeMetrics,
				MetricName: MetricNameBucketSizeBytes,
				OutputType: OutputTypeJSON,
				w:          &bytes.Buffer{},
			},
			want: `{
  "Metrics": [
    {
      "BucketName": "bucket0",
      "Region": "ap-northeast-1",
      "StorageType": "StandardStorage",
      "Bytes": 1024,
      "ReadableBytes": "1.0 KiB"
    },
    {
      "BucketName": "bucket1",
      "Region": "ap-northeast-2",
      "StorageType": "GlacierStorage",
      "Bytes": 4096,
      "ReadableBytes": "4.0 KiB"
    }
  ],
  "MetricName": "BucketSizeBytes",
  "OutputType": "json"
}`,
		},
		{
			name: "empty",
			fields: fields{
				Metrics:    []Metric{},
				MetricName: MetricNameNumberOfObjects,
				OutputType: OutputTypeText,
				w:          &bytes.Buffer{},
			},
			want: `{
  "Metrics": [],
  "MetricName": "NumberOfObjects",
  "OutputType": "text"
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ren := &Renderer{
				Metrics:    tt.fields.Metrics,
				MetricName: tt.fields.MetricName,
				OutputType: tt.fields.OutputType,
				w:          tt.fields.w,
			}
			if got := ren.String(); got != tt.want {
				t.Errorf("Renderer.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderer_Render(t *testing.T) {
	type fields struct {
		Metrics    []Metric
		MetricName MetricName
		OutputType OutputType
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "json",
			fields: fields{
				Metrics:    sampleSizeMetrics,
				MetricName: MetricNameBucketSizeBytes,
				OutputType: OutputTypeJSON,
			},
			want: `[
  {
    "BucketName": "bucket0",
    "Region": "ap-northeast-1",
    "StorageType": "StandardStorage",
    "Bytes": 1024,
    "ReadableBytes": "1.0 KiB"
  },
  {
    "BucketName": "bucket1",
    "Region": "ap-northeast-2",
    "StorageType": "GlacierStorage",
    "Bytes": 4096,
    "ReadableBytes": "4.0 KiB"
  }
]
`,
			wantErr: false,
		},
		{
			name: "text for size metric",
			fields: fields{
				Metrics:    sampleSizeMetrics,
				MetricName: MetricNameBucketSizeBytes,
				OutputType: OutputTypeText,
			},
			want: `+------------+----------------+-----------------+-------+---------------+
| BucketName | Region         | StorageType     | Bytes | ReadableBytes |
+------------+----------------+-----------------+-------+---------------+
| bucket0    | ap-northeast-1 | StandardStorage |  1024 | 1.0 KiB       |
+------------+----------------+-----------------+-------+---------------+
| bucket1    | ap-northeast-2 | GlacierStorage  |  4096 | 4.0 KiB       |
+------------+----------------+-----------------+-------+---------------+
`,
			wantErr: false,
		},
		{
			name: "text for object metric",
			fields: fields{
				Metrics:    sampleObjectMetrics,
				MetricName: MetricNameNumberOfObjects,
				OutputType: OutputTypeText,
			},
			want: `+------------+----------------+-----------------+---------+
| BucketName | Region         | StorageType     | Objects |
+------------+----------------+-----------------+---------+
| bucket0    | ap-northeast-1 | AllStorageTypes |      20 |
+------------+----------------+-----------------+---------+
| bucket1    | ap-northeast-2 | AllStorageTypes |       0 |
+------------+----------------+-----------------+---------+
`,
			wantErr: false,
		},
		{
			name: "markdown",
			fields: fields{
				Metrics:    sampleSizeMetrics,
				MetricName: MetricNameBucketSizeBytes,
				OutputType: OutputTypeMarkdown,
			},
			want: `| BucketName | Region         | StorageType     | Bytes | ReadableBytes |
|------------|----------------|-----------------|-------|---------------|
| bucket0    | ap-northeast-1 | StandardStorage |  1024 | 1.0 KiB       |
| bucket1    | ap-northeast-2 | GlacierStorage  |  4096 | 4.0 KiB       |
`,
			wantErr: false,
		},
		{
			name: "backlog",
			fields: fields{
				Metrics:    sampleSizeMetrics,
				MetricName: MetricNameBucketSizeBytes,
				OutputType: OutputTypeBacklog,
			},
			want: `| BucketName | Region         | StorageType     | Bytes | ReadableBytes |h
| bucket0    | ap-northeast-1 | StandardStorage |  1024 | 1.0 KiB       |
| bucket1    | ap-northeast-2 | GlacierStorage  |  4096 | 4.0 KiB       |
`,
			wantErr: false,
		},
		{
			name: "tsv for size metric",
			fields: fields{
				Metrics:    sampleSizeMetrics,
				MetricName: MetricNameBucketSizeBytes,
				OutputType: OutputTypeTSV,
			},
			want: `BucketName	Region	StorageType	Bytes	ReadableBytes
bucket0	ap-northeast-1	StandardStorage	1024	1.0 KiB
bucket1	ap-northeast-2	GlacierStorage	4096	4.0 KiB
`,
			wantErr: false,
		},
		{
			name: "tsv for object metric",
			fields: fields{
				Metrics:    sampleObjectMetrics,
				MetricName: MetricNameNumberOfObjects,
				OutputType: OutputTypeTSV,
			},
			want: `BucketName	Region	StorageType	Objects
bucket0	ap-northeast-1	AllStorageTypes	20
bucket1	ap-northeast-2	AllStorageTypes	0
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			ren := &Renderer{
				Metrics:    tt.fields.Metrics,
				MetricName: tt.fields.MetricName,
				OutputType: tt.fields.OutputType,
				w:          w,
			}
			if err := ren.Render(); (err != nil) != tt.wantErr {
				t.Errorf("Renderer.Render() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := w.String(); got != tt.want {
				t.Errorf("Renderer.Render() = %v, want %v", got, tt.want)
			}
		})
	}
}
