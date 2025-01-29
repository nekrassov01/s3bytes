package s3bytes

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// getBuckets returns the buckets in the specified region.
func (man *Manager) getBuckets(region string) ([]types.Bucket, error) {
	in := &s3.ListBucketsInput{
		BucketRegion: aws.String(region),
	}
	if man.prefix != nil {
		in.Prefix = man.prefix
	}
	opt := func(o *s3.Options) {
		o.Region = region
	}
	out, err := man.ListBuckets(man.ctx, in, opt)
	if err != nil {
		return nil, err
	}
	return out.Buckets, nil
}
