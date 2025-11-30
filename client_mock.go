package s3bytes

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	_ S3API         = (*mockS3)(nil)
	_ CloudWatchAPI = (*mockCloudWatch)(nil)
)

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
