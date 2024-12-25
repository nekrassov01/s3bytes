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
	IS3 interface {
		ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	}

	ICW interface {
		GetMetricData(ctx context.Context, params *cloudwatch.GetMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error)
	}
)

type (
	S3 struct {
		*s3.Client
	}

	CW struct {
		*cloudwatch.Client
	}
)

type (
	MockS3 struct {
		ListBucketsFunc func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	}

	MockCW struct {
		GetMetricDataFunc func(ctx context.Context, params *cloudwatch.GetMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error)
	}
)

func (m *MockS3) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.ListBucketsFunc(ctx, params, optFns...)
}

func (m *MockCW) GetMetricData(ctx context.Context, params *cloudwatch.GetMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
	return m.GetMetricDataFunc(ctx, params, optFns...)
}

type Client struct {
	s3 IS3
	cw ICW
}

func NewClient(cfg aws.Config) *Client {
	return &Client{
		s3: s3.NewFromConfig(cfg),
		cw: cloudwatch.NewFromConfig(cfg),
	}
}

func NewMockClient(s3 IS3, cw ICW) *Client {
	return &Client{
		s3: s3,
		cw: cw,
	}
}
