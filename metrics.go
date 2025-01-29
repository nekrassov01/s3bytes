package s3bytes

import (
	"strconv"
)

var header = []string{
	"BucketName",
	"Region",
	"MetricName",
	"StorageType",
	"Value",
}

type MetricData struct {
	Header  []string
	Metrics []*Metric
	Total   int64
}

type Metric struct {
	BucketName  string
	Region      string
	MetricName  MetricName
	StorageType StorageType
	Value       float64
}

func (t *Metric) toInput() []any {
	return []any{
		t.BucketName,
		t.Region,
		t.MetricName,
		t.StorageType,
		t.Value,
	}
}

func (t *Metric) toTSV() []string {
	return []string{
		t.BucketName,
		t.Region,
		t.MetricName.String(),
		t.StorageType.String(),
		strconv.FormatFloat(t.Value, 'f', 0, 64),
	}
}
