package s3bytes

// Metric is an interface for the metrics.
type Metric interface {
	Label() string
	Value() float64
}

// SizeMetric is a struct for the size metric.
type SizeMetric struct {
	BucketName    string
	Region        string
	StorageType   StorageType
	Bytes         float64
	ReadableBytes string // human readable bytes
}

// Label returns the label of the size metric.
func (t *SizeMetric) Label() string {
	return t.BucketName
}

// Value returns the value of the size metric.
func (t *SizeMetric) Value() float64 {
	return t.Bytes
}

// ObjectMetric is a struct for the object metric.
type ObjectMetric struct {
	BucketName  string
	Region      string
	StorageType StorageType
	Objects     float64
}

// Label returns the label of the object metric.
func (t *ObjectMetric) Label() string {
	return t.BucketName
}

// Value returns the value of the object metric.
func (t *ObjectMetric) Value() float64 {
	return t.Objects
}
