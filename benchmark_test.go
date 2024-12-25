package s3bytes

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func generateBuckets(n int) []s3types.Bucket {
	buckets := make([]s3types.Bucket, n)
	for i := 0; i < n; i++ {
		buckets[i] = s3types.Bucket{
			Name:         aws.String("bucket" + strconv.Itoa(i)),
			BucketRegion: aws.String("ap-northeast-1"),
		}
	}
	return buckets
}

func generateQueries(n, maxQueriesPerBatch int) [][]cwtypes.MetricDataQuery {
	batches := [][]cwtypes.MetricDataQuery{}
	current := make([]cwtypes.MetricDataQuery, 0, maxQueriesPerBatch)
	for i := 0; i < n; i++ {
		query := cwtypes.MetricDataQuery{
			Id:    aws.String(fmt.Sprintf("m%d", i)),
			Label: aws.String(fmt.Sprintf("bucket%d", i)),
			MetricStat: &cwtypes.MetricStat{
				Metric: &cwtypes.Metric{
					Namespace:  aws.String("AWS/S3"),
					MetricName: aws.String("BucketSizeBytes"),
					Dimensions: []cwtypes.Dimension{
						{
							Name:  aws.String("BucketName"),
							Value: aws.String(fmt.Sprintf("bucket%d", i)),
						},
						{
							Name:  aws.String("StorageType"),
							Value: aws.String("StandardStorage"),
						},
					},
				},
				Period: aws.Int32(86400),
				Stat:   aws.String("Average"),
			},
		}
		current = append(current, query)
		if len(current) == maxQueriesPerBatch {
			batches = append(batches, current)
			current = make([]cwtypes.MetricDataQuery, 0, maxQueriesPerBatch)
		}
	}
	if len(current) > 0 {
		batches = append(batches, current)
	}
	return batches
}

//	func BenchmarkSetQueries(b *testing.B) {
//		maxQueriesPerBatch := 100
//		man := &Manager{
//			MetricName:  MetricNameBucketSizeBytes,
//			StorageType: StorageTypeStandardStorage,
//			Buckets:     generateBuckets(300),
//			MaxQueries:  maxQueriesPerBatch,
//		}
//		b.ResetTimer()
//		for i := 0; i < b.N; i++ {
//			if err := man.SetQueries(); err != nil {
//				b.Fatal(err)
//			}
//		}
//	}

func BenchmarkSetData(b *testing.B) {
	maxQueriesPerBatch := 100
	man := &Manager{
		Client: NewMockClient(nil, &MockCW{
			GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
				out := &cloudwatch.GetMetricDataOutput{
					MetricDataResults: []cwtypes.MetricDataResult{
						{
							Label:  aws.String("bucket0"),
							Values: []float64{1024, 2048},
						},
						{
							Label:  aws.String("bucket1"),
							Values: []float64{0},
						},
						{
							Label:  aws.String("bucket2"),
							Values: []float64{4096, 2048},
						},
						{
							Label:  aws.String("bucket3"),
							Values: []float64{100},
						},
					},
					NextToken: nil,
				}
				return out, nil
			},
		}),
		MetricName:  MetricNameBucketSizeBytes,
		StorageType: StorageTypeStandardStorage,
		Buckets:     generateBuckets(300),
		Batches:     generateQueries(300, maxQueriesPerBatch),
		MaxQueries:  maxQueriesPerBatch,
		filterFunc:  func(float64) bool { return true },
		ctx:         context.Background(),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := man.SetData(); err != nil {
			b.Fatal(err)
		}
	}
}
