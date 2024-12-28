package s3bytes

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Manager is a manager struct for the s3bytes package.
type Manager struct {
	*Client

	Buckets     []s3types.Bucket
	Batches     [][]cwtypes.MetricDataQuery
	Metrics     []Metric
	MetricName  MetricName
	StorageType StorageType
	MaxQueries  int
	Prefix      string             // filter prefix for bucket names
	Region      string             // current region state in process
	filterFunc  func(float64) bool // filter function for metrics
	ctx         context.Context
}

// NewManager creates a new manager.
func NewManager(ctx context.Context, client *Client, region, prefix, expr string, metricName MetricName, storageType StorageType) (*Manager, error) {
	man := &Manager{
		Client:      client,
		Buckets:     make([]s3types.Bucket, 0, maxQueries*2),
		Batches:     make([][]cwtypes.MetricDataQuery, 0, 2),
		Metrics:     make([]Metric, 0, maxQueries*2),
		MetricName:  metricName,
		StorageType: storageType,
		MaxQueries:  maxQueries,
		Prefix:      prefix,
		Region:      region,
		filterFunc:  func(float64) bool { return true },
		ctx:         ctx,
	}
	fn, err := man.eval(expr)
	if err != nil {
		return nil, err
	}
	man.filterFunc = fn
	return man, nil
}

// String returns a string representation of the manager.
func (man *Manager) String() string {
	b, _ := json.MarshalIndent(man, "", "  ")
	return string(b)
}

// Debug prints a debug message.
func (man *Manager) Debug() {
	logger.Debug(man.Region + "\n" + man.String() + "\n")
}

func (man *Manager) eval(expr string) (func(float64) bool, error) {
	if expr == "" {
		return func(float64) bool { return true }, nil
	}
	tokens := strings.SplitN(expr, " ", 2)
	if len(tokens) < 2 {
		return nil, fmt.Errorf("invalid syntax: %q", expr)
	}
	operator := tokens[0]
	v, err := strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		return nil, err
	}
	operators := map[string]func(float64, float64) bool{
		">": func(a, b float64) bool {
			return a > b
		},
		">=": func(a, b float64) bool {
			return a >= b
		},
		"<": func(a, b float64) bool {
			return a < b
		},
		"<=": func(a, b float64) bool {
			return a <= b
		},
		"==": func(a, b float64) bool {
			return a == b
		},
		"!=": func(a, b float64) bool {
			return a != b
		},
	}
	compare, ok := operators[operator]
	if !ok {
		return nil, fmt.Errorf("invalid operator: %q", operator)
	}
	fn := func(f float64) bool {
		return compare(f, v)
	}
	return fn, nil
}
