package s3bytes

import (
	"encoding/json"
	"fmt"
)

// OutputType represents the output type of the renderer.
type OutputType int

const (
	// OutputTypeNone is the output type that means none.
	OutputTypeNone OutputType = iota

	// OutputTypeJSON is the output type that means JSON format.
	OutputTypeJSON

	// OutputTypePrettyJSON is the output type that means pretty JSON format.
	OutputTypePrettyJSON

	// OutputTypeText is the output type that means text format.
	OutputTypeText

	// OutputTypeCompressedText is the output type that means compressed text format.
	OutputTypeCompressedText

	// OutputTypeMarkdown is the output type that means markdown format.
	OutputTypeMarkdown

	// OutputTypeBacklog is the output type that means backlog format.
	OutputTypeBacklog

	// OutputTypeTSV is the output type that means TSV format.
	OutputTypeTSV

	// OutputTypeChart is the output type that means pie chart.
	OutputTypeChart
)

// String returns the string representation of the output type.
func (t OutputType) String() string {
	switch t {
	case OutputTypeNone:
		return "none"
	case OutputTypeJSON:
		return "json"
	case OutputTypePrettyJSON:
		return "prettyjson"
	case OutputTypeText:
		return "text"
	case OutputTypeCompressedText:
		return "compressedtext"
	case OutputTypeMarkdown:
		return "markdown"
	case OutputTypeBacklog:
		return "backlog"
	case OutputTypeTSV:
		return "tsv"
	case OutputTypeChart:
		return "chart"
	default:
		return ""
	}
}

// MarshalJSON returns the JSON representation of the output type.
func (t OutputType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// ParseOutputType parses the output type from the string representation.
func ParseOutputType(s string) (OutputType, error) {
	switch s {
	case OutputTypeJSON.String():
		return OutputTypeJSON, nil
	case OutputTypePrettyJSON.String():
		return OutputTypePrettyJSON, nil
	case OutputTypeText.String():
		return OutputTypeText, nil
	case OutputTypeCompressedText.String():
		return OutputTypeCompressedText, nil
	case OutputTypeMarkdown.String():
		return OutputTypeMarkdown, nil
	case OutputTypeBacklog.String():
		return OutputTypeBacklog, nil
	case OutputTypeTSV.String():
		return OutputTypeTSV, nil
	case OutputTypeChart.String():
		return OutputTypeChart, nil
	default:
		return OutputTypeNone, fmt.Errorf("unsupported output type: %q", s)
	}
}

// MetricName represents the metric name.
type MetricName int

const (
	// MetricNameNone is the metric name that means none.
	MetricNameNone MetricName = iota

	//MetricNameBucketSizeBytes is the metric name that means bucket size in bytes.
	MetricNameBucketSizeBytes

	// MetricNameNumberOfObjects is the metric name that means number of objects.
	MetricNameNumberOfObjects
)

// String returns the string representation of the metric name.
func (t MetricName) String() string {
	switch t {
	case MetricNameNone:
		return "none"
	case MetricNameBucketSizeBytes:
		return "BucketSizeBytes"
	case MetricNameNumberOfObjects:
		return "NumberOfObjects"
	default:
		return ""
	}
}

// MarshalJSON returns the JSON representation of the metric name.
func (t MetricName) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// ParseMetricName parses the metric name from the string representation.
func ParseMetricName(s string) (MetricName, error) {
	switch s {
	case MetricNameBucketSizeBytes.String():
		return MetricNameBucketSizeBytes, nil
	case MetricNameNumberOfObjects.String():
		return MetricNameNumberOfObjects, nil
	default:
		return MetricNameNone, fmt.Errorf("unsupported metrics name: %q", s)
	}
}

// StorageType represents the storage type.
// See: https://docs.aws.amazon.com/AmazonS3/latest/userguide/metrics-dimensions.html#s3-cloudwatch-metrics
type StorageType int

const (

	// StorageTypeNone is the storage type that means none.
	StorageTypeNone StorageType = iota

	// StorageTypeStandardStorage is the storage type that means standard storage.
	StorageTypeStandardStorage

	// StorageTypeIntelligentTieringFAStorage is the storage type that means intelligent tiering FA storage.
	StorageTypeIntelligentTieringFAStorage

	// StorageTypeIntelligentTieringIAStorage is the storage type that means intelligent tiering IA storage.
	StorageTypeIntelligentTieringIAStorage

	// StorageTypeIntelligentTieringAAStorage is the storage type that means intelligent tiering AA storage.
	StorageTypeIntelligentTieringAAStorage

	// StorageTypeIntelligentTieringAIAStorage is the storage type that means intelligent tiering AIA storage.
	StorageTypeIntelligentTieringAIAStorage

	// StorageTypeIntelligentTieringDAAStorage is the storage type that means intelligent tiering DAA storage.
	StorageTypeIntelligentTieringDAAStorage

	// StorageTypeStandardIAStorage is the storage type that means standard IA storage.
	StorageTypeStandardIAStorage

	// StorageTypeStandardIASizeOverhead is the storage type that means standard IA size overhead.
	StorageTypeStandardIASizeOverhead

	// StorageTypeStandardIAObjectOverhead is the storage type that means standard IA object overhead.
	StorageTypeStandardIAObjectOverhead

	// StorageTypeOneZoneIAStorage is the storage type that means one zone IA storage.
	StorageTypeOneZoneIAStorage

	// StorageTypeOneZoneIASizeOverhead is the storage type that means one zone IA size overhead.
	StorageTypeOneZoneIASizeOverhead

	// StorageTypeOneZoneIAObjectOverhead is the storage type that means one zone IA object overhead.
	StorageTypeOneZoneIAObjectOverhead

	// StorageTypeReducedRedundancyStorage is the storage type that means reduced redundancy storage.
	StorageTypeReducedRedundancyStorage

	// StorageTypeGlacierIRSizeOverhead is the storage type that means glacier IR size overhead.
	StorageTypeGlacierIRSizeOverhead

	// StorageTypeGlacierInstantRetrievalStorage is the storage type that means glacier instant retrieval storage.
	StorageTypeGlacierInstantRetrievalStorage

	// StorageTypeGlacierStorage is the storage type that means glacier storage.
	StorageTypeGlacierStorage

	// StorageTypeGlacierStagingStorage is the storage type that means glacier staging storage.
	StorageTypeGlacierStagingStorage

	// StorageTypeGlacierObjectOverhead is the storage type that means glacier object overhead.
	StorageTypeGlacierObjectOverhead

	// StorageTypeGlacierS3ObjectOverhead is the storage type that means glacier S3 object overhead.
	StorageTypeGlacierS3ObjectOverhead

	// StorageTypeDeepArchiveStorage is the storage type that means deep archive storage.
	StorageTypeDeepArchiveStorage

	// StorageTypeDeepArchiveObjectOverhead is the storage type that means deep archive object overhead.
	StorageTypeDeepArchiveObjectOverhead

	// StorageTypeDeepArchiveS3ObjectOverhead is the storage type that means deep archive S3 object overhead.
	StorageTypeDeepArchiveS3ObjectOverhead

	// StorageTypeDeepArchiveStagingStorage is the storage type that means deep archive staging storage.
	StorageTypeDeepArchiveStagingStorage

	// StorageTypeAllStorageTypes is the storage type that means all storage types.
	StorageTypeAllStorageTypes
)

// String returns the string representation of the storage type.
func (t StorageType) String() string {
	switch t {
	case StorageTypeNone:
		return "none"

	// S3 Standard:
	case StorageTypeStandardStorage:
		return "StandardStorage"

	// S3 Intelligent-Tiering:
	case StorageTypeIntelligentTieringFAStorage:
		return "IntelligentTieringFAStorage"
	case StorageTypeIntelligentTieringIAStorage:
		return "IntelligentTieringIAStorage"
	case StorageTypeIntelligentTieringAAStorage:
		return "IntelligentTieringAAStorage"
	case StorageTypeIntelligentTieringAIAStorage:
		return "IntelligentTieringAIAStorage"
	case StorageTypeIntelligentTieringDAAStorage:
		return "IntelligentTieringDAAStorage"

	// S3 Standard-Infrequent Access:
	case StorageTypeStandardIAStorage:
		return "StandardIAStorage"
	case StorageTypeStandardIASizeOverhead:
		return "StandardIASizeOverhead"
	case StorageTypeStandardIAObjectOverhead:
		return "StandardIAObjectOverhead"

	// S3 One Zone-Infrequent Access:
	case StorageTypeOneZoneIAStorage:
		return "OneZoneIAStorage"
	case StorageTypeOneZoneIASizeOverhead:
		return "OneZoneIASizeOverhead"

	// S3 Reduced Redundancy Storage:
	case StorageTypeReducedRedundancyStorage:
		return "ReducedRedundancyStorage"

	// S3 Glacier Instant Retrieval:
	case StorageTypeGlacierIRSizeOverhead:
		return "GlacierIRSizeOverhead"
	case StorageTypeGlacierInstantRetrievalStorage:
		return "GlacierInstantRetrievalStorage"

	// S3 Glacier Flexible Retrieval:
	case StorageTypeGlacierStorage:
		return "GlacierStorage"
	case StorageTypeGlacierStagingStorage:
		return "GlacierStagingStorage"
	case StorageTypeGlacierObjectOverhead:
		return "GlacierObjectOverhead"
	case StorageTypeGlacierS3ObjectOverhead:
		return "GlacierS3ObjectOverhead"

	// S3 Glacier Deep Archive:
	case StorageTypeDeepArchiveStorage:
		return "DeepArchiveStorage"
	case StorageTypeDeepArchiveObjectOverhead:
		return "DeepArchiveObjectOverhead"
	case StorageTypeDeepArchiveS3ObjectOverhead:
		return "DeepArchiveS3ObjectOverhead"
	case StorageTypeDeepArchiveStagingStorage:
		return "DeepArchiveStagingStorage"

	// S3 Express One Zone:
	// case StorageTypeExpressOneZoneStorage:
	// 	return "ExpressOneZoneStorage"

	// Fixed value for metric of NumberOfObjects:
	case StorageTypeAllStorageTypes:
		return "AllStorageTypes"

	default:
		return ""
	}
}

// MarshalJSON returns the JSON representation of the storage type.
func (t StorageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// ParseStorageType parses the storage type from the string representation.
func ParseStorageType(s string) (StorageType, error) {
	switch s {
	// S3 Standard.String():
	case StorageTypeStandardStorage.String():
		return StorageTypeStandardStorage, nil

	// S3 Intelligent-Tiering.String():
	case StorageTypeIntelligentTieringFAStorage.String():
		return StorageTypeIntelligentTieringFAStorage, nil
	case StorageTypeIntelligentTieringIAStorage.String():
		return StorageTypeIntelligentTieringIAStorage, nil
	case StorageTypeIntelligentTieringAAStorage.String():
		return StorageTypeIntelligentTieringAAStorage, nil
	case StorageTypeIntelligentTieringAIAStorage.String():
		return StorageTypeIntelligentTieringAIAStorage, nil
	case StorageTypeIntelligentTieringDAAStorage.String():
		return StorageTypeIntelligentTieringDAAStorage, nil

	// S3 Standard-Infrequent Access.String():
	case StorageTypeStandardIAStorage.String():
		return StorageTypeStandardIAStorage, nil
	case StorageTypeStandardIASizeOverhead.String():
		return StorageTypeStandardIASizeOverhead, nil
	case StorageTypeStandardIAObjectOverhead.String():
		return StorageTypeStandardIAObjectOverhead, nil

	// S3 One Zone-Infrequent Access.String():
	case StorageTypeOneZoneIAStorage.String():
		return StorageTypeOneZoneIAStorage, nil
	case StorageTypeOneZoneIASizeOverhead.String():
		return StorageTypeOneZoneIASizeOverhead, nil

	// S3 Reduced Redundancy Storage.String():
	case StorageTypeReducedRedundancyStorage.String():
		return StorageTypeReducedRedundancyStorage, nil

	// S3 Glacier Instant Retrieval.String():
	case StorageTypeGlacierIRSizeOverhead.String():
		return StorageTypeGlacierIRSizeOverhead, nil
	case StorageTypeGlacierInstantRetrievalStorage.String():
		return StorageTypeGlacierInstantRetrievalStorage, nil

	// S3 Glacier Flexible Retrieval.String():
	case StorageTypeGlacierStorage.String():
		return StorageTypeGlacierStorage, nil
	case StorageTypeGlacierStagingStorage.String():
		return StorageTypeGlacierStagingStorage, nil
	case StorageTypeGlacierObjectOverhead.String():
		return StorageTypeGlacierObjectOverhead, nil
	case StorageTypeGlacierS3ObjectOverhead.String():
		return StorageTypeGlacierS3ObjectOverhead, nil

	// S3 Glacier Deep Archive.String():
	case StorageTypeDeepArchiveStorage.String():
		return StorageTypeDeepArchiveStorage, nil
	case StorageTypeDeepArchiveObjectOverhead.String():
		return StorageTypeDeepArchiveObjectOverhead, nil
	case StorageTypeDeepArchiveS3ObjectOverhead.String():
		return StorageTypeDeepArchiveS3ObjectOverhead, nil
	case StorageTypeDeepArchiveStagingStorage.String():
		return StorageTypeDeepArchiveStagingStorage, nil

	// S3 Express One Zone.String():
	// case StorageTypeExpressOneZoneStorage.String():
	// 	return StorageTypeExpressOneZoneStorage, nil

	// Fixed value for metric of NumberOfObjects.String():
	case StorageTypeAllStorageTypes.String():
		return StorageTypeAllStorageTypes, nil

	default:
		return StorageTypeNone, fmt.Errorf("unsupported storage type: %q", s)
	}
}
