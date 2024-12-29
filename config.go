package s3bytes

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go/logging"
)

const (
	// Version is the current version of s3bytes.
	Version string = "0.0.4"
)

const (
	appName        string            = "s3bytes"
	canonicalName  string            = "S3BYTES"
	fallbackRegion string            = "us-east-1"
	logMode        aws.ClientLogMode = aws.LogRequest | aws.LogResponse | aws.LogRetries | aws.LogSigning | aws.LogDeprecatedUsage
)

var (
	// https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_MetricDataQuery.html
	maxQueries = 500

	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html#concepts-regions
	regions = []string{
		"ap-south-1",
		"eu-north-1",
		"eu-west-3",
		"eu-west-2",
		"eu-west-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"ap-northeast-1",
		"ca-central-1",
		"sa-east-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"eu-central-1",
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
	}
)

// LoadAWSConfig loads the AWS configuration.
func LoadAWSConfig(ctx context.Context, profile string, logger logging.Logger, logMode aws.ClientLogMode) (aws.Config, error) {
	opts := make([]func(*config.LoadOptions) error, 0, 3)
	if logger != nil {
		opts = append(opts, config.WithLogger(logger), config.WithClientLogMode(logMode))
	}
	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, err
	}
	if cfg.Region == "" {
		cfg.Region = fallbackRegion
	}
	return cfg, nil
}
