package s3bytes

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (man *Manager) SetBuckets() error {
	var (
		prefix = aws.String(man.Prefix)
		region = aws.String(man.Region)
		token  = (*string)(nil)
		opt    = func(o *s3.Options) { o.Region = man.Region }
	)
	for {
		in := &s3.ListBucketsInput{
			BucketRegion:      region,
			ContinuationToken: token,
		}
		if man.Prefix != "" {
			in.Prefix = prefix
		}
		out, err := man.s3.ListBuckets(man.ctx, in, opt)
		if err != nil {
			return err
		}
		man.Buckets = append(man.Buckets, out.Buckets...)
		token = out.ContinuationToken
		if token == nil {
			break
		}
	}
	return nil
}
