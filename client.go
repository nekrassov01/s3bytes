package s3bytes

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	_ IS3 = (*S3)(nil)
	_ IS3 = (*MockS3)(nil)
	_ ICW = (*CW)(nil)
	_ ICW = (*MockCW)(nil)
)

type (
	// IS3 is an interface for the s3 client.
	IS3 interface {
		ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	}

	// ICW is an interface for the cloudwatch client.
	ICW interface {
		GetMetricData(ctx context.Context, params *cloudwatch.GetMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error)
	}
)

type (
	// S3 is a wrapper for the s3 client.
	S3 struct {
		*s3.Client
	}

	// CW is a wrapper for the cloudwatch client.
	CW struct {
		*cloudwatch.Client
	}
)

type (
	// MockS3 is a mock for the s3 client.
	MockS3 struct {
		ListBucketsFunc func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	}

	// MockCW is a mock for the cloudwatch client.
	MockCW struct {
		GetMetricDataFunc func(ctx context.Context, params *cloudwatch.GetMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error)
	}
)

// ListBuckets is a wrapper for the ListBuckets method.
func (m *MockS3) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.ListBucketsFunc(ctx, params, optFns...)
}

// GetMetricData is a wrapper for the GetMetricData method.
func (m *MockCW) GetMetricData(ctx context.Context, params *cloudwatch.GetMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
	return m.GetMetricDataFunc(ctx, params, optFns...)
}

// NewS3 is a constructor for the S3 client.
type Client struct {
	s3 IS3
	cw ICW
}

// NewClient is a constructor for the Client.
func NewClient(cfg aws.Config) *Client {
	return &Client{
		s3: s3.NewFromConfig(cfg),
		cw: cloudwatch.NewFromConfig(cfg),
	}
}

// NewMockClient is a constructor for the mock client.
func NewMockClient(s3 IS3, cw ICW) *Client {
	return &Client{
		s3: s3,
		cw: cw,
	}
}
