package s3bytes

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/charmbracelet/log"
)

const (
	Version string = "0.0.1"
)

const (
	appName        string = "s3bytes"
	canonicalName  string = "S3BYTES"
	fallbackRegion string = "us-east-1"
	newLine        string = "\n"
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

func LoadAWSConfig(ctx context.Context, w io.Writer, profile string, loglevel log.Level) (aws.Config, error) {
	logger := newSDKLogger(w, loglevel)
	mode := aws.LogRequest | aws.LogResponse | aws.LogRetries | aws.LogSigning | aws.LogDeprecatedUsage
	opts := []func(*config.LoadOptions) error{
		config.WithLogger(logger),
		config.WithClientLogMode(mode),
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
