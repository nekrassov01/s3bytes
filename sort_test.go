package s3bytes

import "testing"

func TestApp_sort(t *testing.T) {
	data := &MetricData{
		Metrics: []*Metric{
			{
				BucketName:  "bucket-c",
				Region:      "ap-northeast-1",
				StorageType: StorageTypeStandardStorage,
				Value:       150,
			},
			{
				BucketName:  "bucket-a",
				Region:      "ap-northeast-1",
				StorageType: StorageTypeStandardStorage,
				Value:       200,
			},
			{
				BucketName:  "bucket-b",
				Region:      "ap-northeast-1",
				StorageType: StorageTypeStandardStorage,
				Value:       200,
			},
			{
				BucketName:  "bucket-d",
				Region:      "ap-northeast-1",
				StorageType: StorageTypeStandardStorage,
				Value:       100,
			},
		},
	}
	sorted := &MetricData{
		Metrics: []*Metric{
			{
				BucketName:  "bucket-a",
				Region:      "ap-northeast-1",
				StorageType: StorageTypeStandardStorage,
				Value:       200,
			},
			{
				BucketName:  "bucket-b",
				Region:      "ap-northeast-1",
				StorageType: StorageTypeStandardStorage,
				Value:       200,
			},
			{
				BucketName:  "bucket-c",
				Region:      "ap-northeast-1",
				StorageType: StorageTypeStandardStorage,
				Value:       150,
			},
			{
				BucketName:  "bucket-d",
				Region:      "ap-northeast-1",
				StorageType: StorageTypeStandardStorage,
				Value:       100,
			},
		},
	}
	SortMetrics(data)
	for i, metric := range data.Metrics {
		var (
			got  = metric
			want = sorted.Metrics[i]
		)
		if got.Value != want.Value {
			t.Errorf("Metric[%d] Bytes = %v, want %v", i, got.Value, want.Value)
		}
		if got.BucketName != want.BucketName {
			t.Errorf("Metric[%d] BucketName = %v, want %v", i, got.BucketName, want.BucketName)
		}
	}
}
