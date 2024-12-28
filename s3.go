package s3bytes

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// SetBuckets sets the buckets.
func (man *Manager) SetBuckets() error {
	in := &s3.ListBucketsInput{
		BucketRegion: aws.String(man.Region),
	}
	if man.Prefix != "" {
		in.Prefix = aws.String(man.Prefix)
	}
	opt := func(o *s3.Options) {
		o.Region = man.Region
	}
	out, err := man.s3.ListBuckets(man.ctx, in, opt)
	if err != nil {
		return err
	}
	man.Buckets = append(man.Buckets, out.Buckets...)
	return nil
}
