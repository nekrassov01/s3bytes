package s3bytes

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	_ S3API         = (*S3)(nil)
	_ S3API         = (*mockS3)(nil)
	_ CloudWatchAPI = (*CloudWatch)(nil)
	_ CloudWatchAPI = (*mockCloudWatch)(nil)
)

// S3API is an interface for the s3 client.
type S3API interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
}

// CloudWatchAPI is an interface for the cloudwatch client.
type CloudWatchAPI interface {
	GetMetricData(ctx context.Context, params *cloudwatch.GetMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error)
}

// Client is a wrapper for the s3 and cloudwatch clients.
type Client struct {
	S3API
	CloudWatchAPI
}

// S3 is a wrapper for the s3 client.
type S3 struct {
	*s3.Client
}

// CloudWatch is a wrapper for the cloudwatch client.
type CloudWatch struct {
	*cloudwatch.Client
}

// NewClient creates a new client.
func NewClient(cfg aws.Config) *Client {
	return &Client{
		S3API:         s3.NewFromConfig(cfg),
		CloudWatchAPI: cloudwatch.NewFromConfig(cfg),
	}
}

// mockS3 is a mock for the s3 client.
type mockS3 struct {
	ListBucketsFunc func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
}

// mockCloudWatch is a mock for the cloudwatch client.
type mockCloudWatch struct {
	GetMetricDataFunc func(ctx context.Context, params *cloudwatch.GetMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error)
}

// ListBuckets is a wrapper for the ListBuckets method.
func (m *mockS3) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.ListBucketsFunc(ctx, params, optFns...)
}

// GetMetricData is a wrapper for the GetMetricData method.
func (m *mockCloudWatch) GetMetricData(ctx context.Context, params *cloudwatch.GetMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
	return m.GetMetricDataFunc(ctx, params, optFns...)
}

// newMockClient is a constructor for the mock client.
func newMockClient(s3 S3API, cw CloudWatchAPI) *Client {
	return &Client{
		S3API:         s3,
		CloudWatchAPI: cw,
	}
}
