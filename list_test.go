package s3bytes

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/semaphore"
)

func TestManager_List(t *testing.T) {
	type fields struct {
		Client      *Client
		metricName  MetricName
		storageType StorageType
		prefix      *string
		regions     []string
		filterFunc  func(float64) bool
		sem         *semaphore.Weighted
		ctx         context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    *MetricData
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Client: newMockClient(
					&mockS3{
						ListBucketsFunc: func(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							out := &s3.ListBucketsOutput{
								Buckets: []s3types.Bucket{
									{
										Name:         aws.String("bucket0"),
										BucketRegion: aws.String("ap-northeast-1"),
									},
								},
							}
							return out, nil
						},
					},
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							return &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{
									{
										Id:     aws.String("m0"),
										Label:  aws.String("bucket0"),
										Values: []float64{2048},
									},
								},
							}, nil
						},
					},
				),
				regions:     []string{"ap-northeast-1"},
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
				filterFunc:  func(float64) bool { return true },
				sem:         semaphore.NewWeighted(NumWorker),
				ctx:         context.Background(),
			},
			want: &MetricData{
				Header: header,
				Metrics: []*Metric{
					{
						BucketName:  "bucket0",
						Region:      "ap-northeast-1",
						MetricName:  MetricNameBucketSizeBytes,
						StorageType: StorageTypeStandardStorage,
						Value:       2048,
					},
				},
				Total: 2048,
			},
			wantErr: false,
		},
		{
			name: "metric error",
			fields: fields{
				Client: newMockClient(
					&mockS3{
						ListBucketsFunc: func(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							out := &s3.ListBucketsOutput{
								Buckets: []s3types.Bucket{
									{
										Name:         aws.String("bucket0"),
										BucketRegion: aws.String("ap-northeast-1"),
									},
								},
								ContinuationToken: nil,
							}
							return out, nil
						},
					},
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							return nil, errors.New("error")
						},
					},
				),
				regions:     []string{"ap-northeast-1"},
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
				filterFunc:  func(float64) bool { return true },
				sem:         semaphore.NewWeighted(NumWorker),
				ctx:         context.Background(),
			},
			wantErr: true,
		},
		{
			name: "bucket error",
			fields: fields{
				Client: newMockClient(
					&mockS3{
						ListBucketsFunc: func(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							return nil, errors.New("error")
						},
					},
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							out := &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{
									{
										Id:     aws.String("m0"),
										Label:  aws.String("bucket"),
										Values: []float64{2048},
									},
								},
							}
							return out, nil
						},
					},
				),
				regions:     []string{"ap-northeast-1"},
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
				filterFunc:  func(float64) bool { return true },
				sem:         semaphore.NewWeighted(NumWorker),
				ctx:         context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				Client:      tt.fields.Client,
				metricName:  tt.fields.metricName,
				storageType: tt.fields.storageType,
				prefix:      tt.fields.prefix,
				regions:     tt.fields.regions,
				filterFunc:  tt.fields.filterFunc,
				sem:         tt.fields.sem,
				ctx:         tt.fields.ctx,
			}
			got, err := man.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//	if !reflect.DeepEqual(got, tt.want) {
			//		t.Errorf("Manager.List() = %v, want %v", got, tt.want)
			//	}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Manager.List() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
