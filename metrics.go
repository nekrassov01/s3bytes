package s3bytes

import (
	"fmt"
	"strconv"
)

var header = []string{
	"BucketName",
	"Region",
	"MetricName",
	"StorageType",
	"Value",
}

var _ filterTarget = (*Metric)(nil)

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

func (t *Metric) GetField(key string) (any, error) {
	switch key {
	case "bytes", "Bytes", "value", "Value":
		return t.Value, nil
	default:
		return 0, fmt.Errorf("field not found: %q", key)
	}
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
