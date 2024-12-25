package s3bytes

type Metric interface {
	Label() string
	Value() float64
}

type SizeMetric struct {
	BucketName    string
	Region        string
	StorageType   StorageType
	Bytes         float64
	ReadableBytes string // human readable bytes
}

func (t *SizeMetric) Label() string {
	return t.BucketName
}

func (t *SizeMetric) Value() float64 {
	return t.Bytes
}

type ObjectMetric struct {
	BucketName  string
	Region      string
	StorageType StorageType
	Objects     float64
}

func (t *ObjectMetric) Label() string {
	return t.BucketName
}

func (t *ObjectMetric) Value() float64 {
	return t.Objects
}
