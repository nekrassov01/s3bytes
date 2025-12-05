package s3bytes

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"golang.org/x/sync/semaphore"
)

func BenchmarkList(b *testing.B) {
	var (
		n       = 10
		buckets = make([]s3types.Bucket, n)
		name    = aws.String("bucket")
		region  = aws.String("us-east-1")
		metrics = make([]cwtypes.MetricDataResult, n)
		id      = aws.String("m0")
		label   = aws.String("bucket")
		values  = []float64{2048}
	)
	man := &Manager{
		client: newMockClient(
			&mockS3{
				ListBucketsFunc: func(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
					for i := 0; i < n; i++ {
						buckets[i] = s3types.Bucket{
							Name:         name,
							BucketRegion: region,
						}
					}
					out := &s3.ListBucketsOutput{
						Buckets: buckets,
					}
					return out, nil
				},
			},
			&mockCloudWatch{
				GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
					for i := 0; i < n; i++ {
						metrics[i] = cwtypes.MetricDataResult{
							Id:     id,
							Label:  label,
							Values: values,
						}
					}
					out := &cloudwatch.GetMetricDataOutput{
						MetricDataResults: metrics,
					}
					return out, nil
				},
			}),
		regions:     []string{"us-east-1"},
		metricName:  MetricNameBucketSizeBytes,
		storageType: StorageTypeStandardStorage,
		sem:         semaphore.NewWeighted(NumWorker),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := man.List(context.Background()); err != nil {
			b.Fatal(err)
		}
	}
}
