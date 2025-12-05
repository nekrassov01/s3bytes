package s3bytes

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/nekrassov01/filter"
	"golang.org/x/sync/semaphore"
)

type (
	filterExpr   = *filter.Expr
	filterTarget = filter.Target
)

// Manager is a manager struct for the s3bytes package.
type Manager struct {
	client      *Client `json:"-"`
	metricName  MetricName
	storageType StorageType
	prefix      *string
	regions     []string
	filterExpr  filterExpr
	filterRaw   string
	sem         *semaphore.Weighted
}

// NewManager creates a new manager.
func NewManager(client *Client) *Manager {
	return &Manager{
		client:  client,
		regions: DefaultRegions,
		sem:     semaphore.NewWeighted(NumWorker),
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
func (man *Manager) SetFilter(raw string) error {
	if raw == "" {
		return nil
	}
	expr, err := filter.Parse(raw)
	if err != nil {
		return fmt.Errorf("failed to parse filter: %w", err)
	}
	man.filterExpr = expr
	man.filterRaw = raw
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
