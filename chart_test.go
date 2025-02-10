package s3bytes

import (
	"reflect"
	"testing"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func Test_getPieItems(t *testing.T) {
	type args struct {
		data *MetricData
	}
	type want struct {
		title string
		items []opts.PieData
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "size",
			args: args{
				data: &MetricData{
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
				},
			},
			want: want{
				title: "Bucket Size Bytes",
				items: []opts.PieData{
					{
						Name:  "bucket0",
						Value: float64(1024),
					},
					{
						Name:  "bucket1",
						Value: float64(4096),
					},
				},
			},
		},
		{
			name: "objects",
			args: args{
				data: &MetricData{
					Header: header,
					Metrics: []*Metric{
						{
							BucketName:  "bucket0",
							Region:      "ap-northeast-1",
							MetricName:  MetricNameNumberOfObjects,
							StorageType: StorageTypeAllStorageTypes,
							Value:       10,
						},
						{
							BucketName:  "bucket1",
							Region:      "ap-northeast-2",
							MetricName:  MetricNameNumberOfObjects,
							StorageType: StorageTypeAllStorageTypes,
							Value:       20,
						},
					},
				},
			},
			want: want{
				title: "Number Of Objects",
				items: []opts.PieData{
					{
						Name:  "bucket0",
						Value: float64(10),
					},
					{
						Name:  "bucket1",
						Value: float64(20),
					},
				},
			},
		},
		{
			name: "include zero",
			args: args{
				data: &MetricData{
					Header: header,
					Metrics: []*Metric{
						{
							BucketName:  "bucket0",
							Region:      "ap-northeast-1",
							MetricName:  MetricNameBucketSizeBytes,
							StorageType: StorageTypeStandardStorage,
							Value:       0,
						},
						{
							BucketName:  "bucket1",
							Region:      "ap-northeast-2",
							MetricName:  MetricNameBucketSizeBytes,
							StorageType: StorageTypeGlacierStorage,
							Value:       4096,
						},
					},
				},
			},
			want: want{
				title: "Bucket Size Bytes",
				items: []opts.PieData{
					{
						Name:  "bucket1",
						Value: float64(4096),
					},
				},
			},
		},
		{
			name: "others",
			args: args{
				data: &MetricData{
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
						{
							BucketName:  "bucket2",
							Region:      "us-east-1",
							MetricName:  MetricNameBucketSizeBytes,
							StorageType: StorageTypeGlacierStorage,
							Value:       256,
						},
						{
							BucketName:  "bucket3",
							Region:      "us-east-1",
							MetricName:  MetricNameBucketSizeBytes,
							StorageType: StorageTypeGlacierStorage,
							Value:       512,
						},
					},
				},
			},
			want: want{
				title: "Bucket Size Bytes",
				items: []opts.PieData{
					{
						Name:  "bucket0",
						Value: float64(1024),
					},
					{
						Name:  "bucket1",
						Value: float64(4096),
					},
					{
						Name:  "others",
						Value: float64(768),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, items := getPieItems(tt.args.data)
			if title != tt.want.title {
				t.Errorf("getPieItems() title = %v, want %v", title, tt.want)
			}
			if !reflect.DeepEqual(items, tt.want.items) {
				t.Errorf("getPieItems() items = %v, want %v", items, tt.want.items)
			}
		})
	}
}

func Test_render(t *testing.T) {
	type args struct {
		pie *charts.Pie
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "render",
			args: args{
				pie: newPie("Bucket Size Bytes", []opts.PieData{
					{
						Name:  "bucket1",
						Value: float64(4096),
					},
					{
						Name:  "bucket0",
						Value: float64(1024),
					},
					{
						Name:  "others",
						Value: float64(768),
					},
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := render(tt.args.pie); (err != nil) != tt.wantErr {
				t.Errorf("render() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
