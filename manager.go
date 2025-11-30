package s3bytes

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"golang.org/x/sync/semaphore"
)

// Manager is a manager struct for the s3bytes package.
type Manager struct {
	client      *Client `json:"-"`
	metricName  MetricName
	storageType StorageType
	prefix      *string
	regions     []string
	filterFunc  func(float64) bool
	sem         *semaphore.Weighted
}

// NewManager creates a new manager.
func NewManager(client *Client) *Manager {
	return &Manager{
		client:     client,
		regions:    DefaultRegions,
		filterFunc: func(float64) bool { return true },
		sem:        semaphore.NewWeighted(NumWorker),
	}
}

// SetRegion sets the specified regions.
func (man *Manager) SetRegion(regions []string) error {
	if len(regions) == 0 {
		return nil
	}
	for _, region := range regions {
		if _, ok := allowedRegions[region]; !ok {
			return fmt.Errorf("unsupported region: %s", region)
		}
	}
	man.regions = regions
	return nil
}

// SetPrefix sets the prefix.
func (man *Manager) SetPrefix(prefix string) error {
	if prefix == "" {
		return nil
	}
	if !bucketPrefixPattern.MatchString(prefix) {
		return fmt.Errorf("invalid prefix: %q", prefix)
	}
	man.prefix = aws.String(prefix)
	return nil
}

// SetFilter sets the filter expressions.
func (man *Manager) SetFilter(expr string) error {
	if expr == "" {
		return nil
	}
	tokens := strings.SplitN(expr, " ", 2)
	if len(tokens) < 2 {
		return fmt.Errorf("invalid syntax: %q", expr)
	}
	operator := tokens[0]
	v, err := strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		return err
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
		return fmt.Errorf("invalid operator: %q", operator)
	}
	man.filterFunc = func(f float64) bool {
		return compare(f, v)
	}
	return nil
}

// SetMetric sets the metric name and storage type.
func (man *Manager) SetMetric(metricName MetricName, storageType StorageType) error {
	if metricName == MetricNameBucketSizeBytes && storageType == StorageTypeAllStorageTypes {
		return errors.New("BucketSizeBytes metric does not support AllStorageTypes")
	}
	if metricName == MetricNameNumberOfObjects && storageType != StorageTypeAllStorageTypes {
		return errors.New("NumberOfObjects metric only supports AllStorageTypes")
	}
	man.metricName = metricName
	man.storageType = storageType
	return nil
}

// String returns a string representation of the manager.
func (man *Manager) String() string {
	s := struct {
		MetricName  string   `json:"metricName"`
		StorageType string   `json:"storageType"`
		Prefix      *string  `json:"prefix"`
		Regions     []string `json:"regions"`
	}{
		MetricName:  man.metricName.String(),
		StorageType: man.storageType.String(),
		Prefix:      man.prefix,
		Regions:     man.regions,
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}
