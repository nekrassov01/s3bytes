package s3bytes

import (
	"cmp"
	"slices"
)

// SortMetrics sorts the metrics by value and bucket name.
func SortMetrics(data *MetricData) {
	slices.SortFunc(data.Metrics, func(a, b *Metric) int {
		if n := cmp.Compare(b.Value, a.Value); n != 0 {
			return n
		}
		return cmp.Compare(a.BucketName, b.BucketName)
	})
}
