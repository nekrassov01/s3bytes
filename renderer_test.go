package s3bytes

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var testSizeMetricData = &MetricData{
	Header: header,
	Metrics: []*Metric{
		{
			BucketName:  "bucket0",
			Region:      "ap-northeast-1",
			MetricName:  MetricNameBucketSizeBytes,
			StorageType: StorageTypeStandardStorage,
			Value:       1024,
		},
		{
			BucketName:  "bucket1",
			Region:      "ap-northeast-2",
			MetricName:  MetricNameBucketSizeBytes,
			StorageType: StorageTypeGlacierStorage,
			Value:       4096,
		},
	},
}

var testObjectMetricData = &MetricData{
	Header: header,
	Metrics: []*Metric{
		{
			BucketName:  "bucket0",
			Region:      "ap-northeast-1",
			MetricName:  MetricNameNumberOfObjects,
			StorageType: StorageTypeAllStorageTypes,
			Value:       20,
		},
		{
			BucketName:  "bucket1",
			Region:      "ap-northeast-2",
			MetricName:  MetricNameNumberOfObjects,
			StorageType: StorageTypeAllStorageTypes,
			Value:       0,
		},
	},
}

func TestNewRenderer(t *testing.T) {
	type args struct {
		data       *MetricData
		outputType OutputType
	}
	tests := []struct {
		name  string
		args  args
		want  *Renderer
		wantW string
	}{
		{
			name: "normal",
			args: args{
				data:       testSizeMetricData,
				outputType: OutputTypeJSON,
			},
			want: &Renderer{
				Data:       testSizeMetricData,
				OutputType: OutputTypeJSON,
				w:          &bytes.Buffer{},
			},
		},
		{
			name: "empty",
			args: args{
				data:       nil,
				outputType: OutputTypeText,
			},
			want: &Renderer{
				Data:       nil,
				OutputType: OutputTypeText,
				w:          &bytes.Buffer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if got := NewRenderer(w, tt.args.data, tt.args.outputType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRenderer() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("NewRenderer() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestRenderer_String(t *testing.T) {
	type fields struct {
		Data       *MetricData
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
				Data:       testSizeMetricData,
				OutputType: OutputTypeJSON,
				w:          &bytes.Buffer{},
			},
			want: `{
  "Data": {
    "Header": [
      "BucketName",
      "Region",
      "MetricName",
      "StorageType",
      "Value"
    ],
    "Metrics": [
      {
        "BucketName": "bucket0",
        "Region": "ap-northeast-1",
        "MetricName": "BucketSizeBytes",
        "StorageType": "StandardStorage",
        "Value": 1024
      },
      {
        "BucketName": "bucket1",
        "Region": "ap-northeast-2",
        "MetricName": "BucketSizeBytes",
        "StorageType": "GlacierStorage",
        "Value": 4096
      }
    ],
    "Total": 0
  },
  "OutputType": "json"
}`,
		},
		{
			name: "empty",
			fields: fields{
				Data:       nil,
				OutputType: OutputTypeText,
				w:          &bytes.Buffer{},
			},
			want: `{
  "Data": null,
  "OutputType": "text"
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ren := &Renderer{
				Data:       tt.fields.Data,
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
		Data       *MetricData
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
				Data:       testSizeMetricData,
				OutputType: OutputTypeJSON,
			},
			want: `[
  {
    "BucketName": "bucket0",
    "Region": "ap-northeast-1",
    "MetricName": "BucketSizeBytes",
    "StorageType": "StandardStorage",
    "Value": 1024
  },
  {
    "BucketName": "bucket1",
    "Region": "ap-northeast-2",
    "MetricName": "BucketSizeBytes",
    "StorageType": "GlacierStorage",
    "Value": 4096
  }
]
`,
			wantErr: false,
		},
		{
			name: "text for size metric",
			fields: fields{
				Data:       testSizeMetricData,
				OutputType: OutputTypeText,
			},
			want: `+------------+----------------+-----------------+-----------------+-------+
| BucketName | Region         | MetricName      | StorageType     | Value |
+------------+----------------+-----------------+-----------------+-------+
| bucket0    | ap-northeast-1 | BucketSizeBytes | StandardStorage |  1024 |
+------------+----------------+-----------------+-----------------+-------+
| bucket1    | ap-northeast-2 | BucketSizeBytes | GlacierStorage  |  4096 |
+------------+----------------+-----------------+-----------------+-------+
`,
			wantErr: false,
		},
		{
			name: "text for object metric",
			fields: fields{
				Data:       testObjectMetricData,
				OutputType: OutputTypeText,
			},
			want: `+------------+----------------+-----------------+-----------------+-------+
| BucketName | Region         | MetricName      | StorageType     | Value |
+------------+----------------+-----------------+-----------------+-------+
| bucket0    | ap-northeast-1 | NumberOfObjects | AllStorageTypes |    20 |
+------------+----------------+-----------------+-----------------+-------+
| bucket1    | ap-northeast-2 | NumberOfObjects | AllStorageTypes |     0 |
+------------+----------------+-----------------+-----------------+-------+
`,
			wantErr: false,
		},
		{
			name: "compressed text",
			fields: fields{
				Data:       testSizeMetricData,
				OutputType: OutputTypeCompressedText,
			},
			want: `+------------+----------------+-----------------+-----------------+-------+
| BucketName | Region         | MetricName      | StorageType     | Value |
+------------+----------------+-----------------+-----------------+-------+
| bucket0    | ap-northeast-1 | BucketSizeBytes | StandardStorage |  1024 |
| bucket1    | ap-northeast-2 | BucketSizeBytes | GlacierStorage  |  4096 |
+------------+----------------+-----------------+-----------------+-------+
`,
			wantErr: false,
		},
		{
			name: "markdown",
			fields: fields{
				Data:       testSizeMetricData,
				OutputType: OutputTypeMarkdown,
			},
			want: `| BucketName | Region         | MetricName      | StorageType     | Value |
|------------|----------------|-----------------|-----------------|-------|
| bucket0    | ap-northeast-1 | BucketSizeBytes | StandardStorage |  1024 |
| bucket1    | ap-northeast-2 | BucketSizeBytes | GlacierStorage  |  4096 |
`,
			wantErr: false,
		},
		{
			name: "backlog",
			fields: fields{
				Data:       testSizeMetricData,
				OutputType: OutputTypeBacklog,
			},
			want: `| BucketName | Region         | MetricName      | StorageType     | Value |h
| bucket0    | ap-northeast-1 | BucketSizeBytes | StandardStorage |  1024 |
| bucket1    | ap-northeast-2 | BucketSizeBytes | GlacierStorage  |  4096 |
`,
			wantErr: false,
		},
		{
			name: "tsv for size metric",
			fields: fields{
				Data:       testSizeMetricData,
				OutputType: OutputTypeTSV,
			},
			want: `BucketName	Region	MetricName	StorageType	Value
bucket0	ap-northeast-1	BucketSizeBytes	StandardStorage	1024
bucket1	ap-northeast-2	BucketSizeBytes	GlacierStorage	4096
`,
			wantErr: false,
		},
		{
			name: "tsv for object metric",
			fields: fields{
				Data:       testObjectMetricData,
				OutputType: OutputTypeTSV,
			},
			want: `BucketName	Region	MetricName	StorageType	Value
bucket0	ap-northeast-1	NumberOfObjects	AllStorageTypes	20
bucket1	ap-northeast-2	NumberOfObjects	AllStorageTypes	0
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			ren := &Renderer{
				Data:       tt.fields.Data,
				OutputType: tt.fields.OutputType,
				w:          w,
			}
			if err := ren.Render(); (err != nil) != tt.wantErr {
				t.Errorf("Renderer.Render() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, w.String()); diff != "" {
				t.Errorf("Renderer.Render() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
