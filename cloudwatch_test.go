package s3bytes

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"golang.org/x/sync/semaphore"
)

func TestManager_getMetrics(t *testing.T) {
	type fields struct {
		Client      *Client
		MetricName  MetricName
		StorageType StorageType
		Prefix      *string
		Regions     []string
		filterFunc  func(float64) bool
		sem         *semaphore.Weighted
		ctx         context.Context
	}
	type args struct {
		buckets []s3types.Bucket
		region  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Metric
		want1   int64
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							return &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{
									{
										Label:  aws.String("bucket0"),
										Values: []float64{1024, 2048},
									},
									{
										Label:  aws.String("bucket1"),
										Values: []float64{0},
									},
								},
								NextToken: nil,
							}, nil
						},
					},
				),
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			args: args{
				buckets: []s3types.Bucket{
					{Name: aws.String("bucket0")},
					{Name: aws.String("bucket1")},
				},
				region: "ap-northeast-1",
			},
			want: []*Metric{
				{
					BucketName:  "bucket0",
					Region:      "ap-northeast-1",
					MetricName:  MetricNameBucketSizeBytes,
					StorageType: StorageTypeStandardStorage,
					Value:       2048,
				},
				{
					BucketName:  "bucket1",
					Region:      "ap-northeast-1",
					MetricName:  MetricNameBucketSizeBytes,
					StorageType: StorageTypeStandardStorage,
					Value:       0,
				},
			},
			want1:   2048,
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							return nil, errors.New("error")
						},
					},
				),
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			args: args{
				buckets: []s3types.Bucket{
					{Name: aws.String("bucket0")},
				},
				region: "ap-northeast-1",
			},
			want:    nil,
			want1:   0,
			wantErr: true,
		},
		{
			name: "filter func returns false",
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							return &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{
									{
										Label:  aws.String("bucket0"),
										Values: []float64{1024, 2048},
									},
								},
								NextToken: nil,
							}, nil
						},
					},
				),
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				filterFunc:  func(float64) bool { return false },
				ctx:         context.Background(),
			},
			args: args{
				buckets: []s3types.Bucket{
					{Name: aws.String("bucket0")},
				},
				region: "ap-northeast-1",
			},
			want:    []*Metric{},
			want1:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				Client:      tt.fields.Client,
				metricName:  tt.fields.MetricName,
				storageType: tt.fields.StorageType,
				prefix:      tt.fields.Prefix,
				regions:     tt.fields.Regions,
				filterFunc:  tt.fields.filterFunc,
				sem:         tt.fields.sem,
				ctx:         tt.fields.ctx,
			}
			got, got1, err := man.getMetrics(tt.args.buckets, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.getMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.getMetrics() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Manager.getMetrics() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestManager_getMetricsFromQueries(t *testing.T) {
	type fields struct {
		Client      *Client
		MetricName  MetricName
		StorageType StorageType
		Prefix      *string
		Regions     []string
		filterFunc  func(float64) bool
		sem         *semaphore.Weighted
		ctx         context.Context
	}
	type args struct {
		queries []cwtypes.MetricDataQuery
		region  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Metric
		want1   int64
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							return &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{
									{
										Label:  aws.String("bucket0"),
										Values: []float64{1024, 2048},
									},
									{
										Label:  aws.String("bucket1"),
										Values: []float64{0},
									},
								},
								NextToken: nil,
							}, nil
						},
					},
				),
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			args: args{
				queries: []cwtypes.MetricDataQuery{
					{
						Id:    aws.String("m0"),
						Label: aws.String("bucket0"),
						MetricStat: &cwtypes.MetricStat{
							Metric: &cwtypes.Metric{
								Namespace:  aws.String("AWS/S3"),
								MetricName: aws.String(MetricNameBucketSizeBytes.String()),
								Dimensions: []cwtypes.Dimension{
									{
										Name:  aws.String("BucketName"),
										Value: aws.String("bucket0"),
									},
									{
										Name:  aws.String("StorageType"),
										Value: aws.String(StorageTypeStandardStorage.String()),
									},
								},
							},
							Period: aws.Int32(86400),
							Stat:   aws.String("Average"),
						},
					},
					{
						Id:    aws.String("m1"),
						Label: aws.String("bucket1"),
						MetricStat: &cwtypes.MetricStat{
							Metric: &cwtypes.Metric{
								Namespace:  aws.String("AWS/S3"),
								MetricName: aws.String(MetricNameBucketSizeBytes.String()),
								Dimensions: []cwtypes.Dimension{
									{
										Name:  aws.String("BucketName"),
										Value: aws.String("bucket1"),
									},
									{
										Name:  aws.String("StorageType"),
										Value: aws.String(StorageTypeStandardStorage.String()),
									},
								},
							},
							Period: aws.Int32(86400),
							Stat:   aws.String("Average"),
						},
					},
				},
				region: "ap-northeast-1",
			},
			want: []*Metric{
				{
					BucketName:  "bucket0",
					Region:      "ap-northeast-1",
					MetricName:  MetricNameBucketSizeBytes,
					StorageType: StorageTypeStandardStorage,
					Value:       2048,
				},
				{
					BucketName:  "bucket1",
					Region:      "ap-northeast-1",
					MetricName:  MetricNameBucketSizeBytes,
					StorageType: StorageTypeStandardStorage,
					Value:       0,
				},
			},
			want1:   2048,
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							return nil, errors.New("error")
						},
					},
				),
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			args: args{
				queries: []cwtypes.MetricDataQuery{
					{
						Id:    aws.String("m0"),
						Label: aws.String("bucket0"),
						MetricStat: &cwtypes.MetricStat{
							Metric: &cwtypes.Metric{
								Namespace:  aws.String("AWS/S3"),
								MetricName: aws.String(MetricNameBucketSizeBytes.String()),
								Dimensions: []cwtypes.Dimension{
									{
										Name:  aws.String("BucketName"),
										Value: aws.String("bucket0"),
									},
									{
										Name:  aws.String("StorageType"),
										Value: aws.String(StorageTypeStandardStorage.String()),
									},
								},
							},
							Period: aws.Int32(86400),
							Stat:   aws.String("Average"),
						},
					},
				},
				region: "ap-northeast-1",
			},
			want:    nil,
			want1:   0,
			wantErr: true,
		},
		{
			name: "filter func returns false",
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							return &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{
									{
										Label:  aws.String("bucket0"),
										Values: []float64{1024, 2048},
									},
								},
								NextToken: nil,
							}, nil
						},
					},
				),
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				filterFunc:  func(float64) bool { return false },
				ctx:         context.Background(),
			},
			args: args{
				queries: []cwtypes.MetricDataQuery{
					{
						Id:    aws.String("m0"),
						Label: aws.String("bucket0"),
						MetricStat: &cwtypes.MetricStat{
							Metric: &cwtypes.Metric{
								Namespace:  aws.String("AWS/S3"),
								MetricName: aws.String(MetricNameBucketSizeBytes.String()),
								Dimensions: []cwtypes.Dimension{
									{
										Name:  aws.String("BucketName"),
										Value: aws.String("bucket0"),
									},
									{
										Name:  aws.String("StorageType"),
										Value: aws.String(StorageTypeStandardStorage.String()),
									},
								},
							},
							Period: aws.Int32(86400),
							Stat:   aws.String("Average"),
						},
					},
				},
				region: "ap-northeast-1",
			},
			want:    []*Metric{},
			want1:   0,
			wantErr: false,
		},
		{
			name: "result values nil",
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							return &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{
									{
										Label:  aws.String("bucket0"),
										Values: nil,
									},
								},
								NextToken: nil,
							}, nil
						},
					},
				),
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			args: args{
				queries: []cwtypes.MetricDataQuery{
					{
						Id:    aws.String("m0"),
						Label: aws.String("bucket0"),
						MetricStat: &cwtypes.MetricStat{
							Metric: &cwtypes.Metric{
								Namespace:  aws.String("AWS/S3"),
								MetricName: aws.String(MetricNameBucketSizeBytes.String()),
								Dimensions: []cwtypes.Dimension{
									{
										Name:  aws.String("BucketName"),
										Value: aws.String("bucket0"),
									},
									{
										Name:  aws.String("StorageType"),
										Value: aws.String(StorageTypeStandardStorage.String()),
									},
								},
							},
							Period: aws.Int32(86400),
							Stat:   aws.String("Average"),
						},
					},
				},
				region: "ap-northeast-1",
			},
			want: []*Metric{
				{
					BucketName:  "bucket0",
					Region:      "ap-northeast-1",
					MetricName:  MetricNameBucketSizeBytes,
					StorageType: StorageTypeStandardStorage,
					Value:       0,
				},
			},
			want1:   0,
			wantErr: false,
		},
		{
			name: "pagination",
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, params *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							if params.NextToken == nil {
								out := &cloudwatch.GetMetricDataOutput{
									MetricDataResults: []cwtypes.MetricDataResult{
										{
											Label:  aws.String("bucket0"),
											Values: []float64{1024},
										},
									},
									NextToken: aws.String("token0"),
								}
								return out, nil
							}
							out := &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{
									{
										Label:  aws.String("bucket1"),
										Values: []float64{2048},
									},
								},
								NextToken: nil,
							}
							return out, nil
						},
					},
				),
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			args: args{
				queries: []cwtypes.MetricDataQuery{
					{
						Id:    aws.String("m0"),
						Label: aws.String("bucket0"),
						MetricStat: &cwtypes.MetricStat{
							Metric: &cwtypes.Metric{
								Namespace:  aws.String("AWS/S3"),
								MetricName: aws.String(MetricNameBucketSizeBytes.String()),
								Dimensions: []cwtypes.Dimension{
									{
										Name:  aws.String("BucketName"),
										Value: aws.String("bucket0"),
									},
									{
										Name:  aws.String("StorageType"),
										Value: aws.String(StorageTypeStandardStorage.String()),
									},
								},
							},
							Period: aws.Int32(86400),
							Stat:   aws.String("Average"),
						},
					},
					{
						Id:    aws.String("m1"),
						Label: aws.String("bucket1"),
						MetricStat: &cwtypes.MetricStat{
							Metric: &cwtypes.Metric{
								Namespace:  aws.String("AWS/S3"),
								MetricName: aws.String(MetricNameBucketSizeBytes.String()),
								Dimensions: []cwtypes.Dimension{
									{
										Name:  aws.String("BucketName"),
										Value: aws.String("bucket1"),
									},
									{
										Name:  aws.String("StorageType"),
										Value: aws.String(StorageTypeStandardStorage.String()),
									},
								},
							},
							Period: aws.Int32(86400),
							Stat:   aws.String("Average"),
						},
					},
				},
				region: "ap-northeast-1",
			},
			want: []*Metric{
				{
					BucketName:  "bucket0",
					Region:      "ap-northeast-1",
					MetricName:  MetricNameBucketSizeBytes,
					StorageType: StorageTypeStandardStorage,
					Value:       1024,
				},
				{
					BucketName:  "bucket1",
					Region:      "ap-northeast-1",
					MetricName:  MetricNameBucketSizeBytes,
					StorageType: StorageTypeStandardStorage,
					Value:       2048,
				},
			},
			want1:   3072,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				Client:      tt.fields.Client,
				metricName:  tt.fields.MetricName,
				storageType: tt.fields.StorageType,
				prefix:      tt.fields.Prefix,
				regions:     tt.fields.Regions,
				filterFunc:  tt.fields.filterFunc,
				sem:         tt.fields.sem,
				ctx:         tt.fields.ctx,
			}
			got, got1, err := man.getMetricsFromQueries(tt.args.queries, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.getMetricsFromQueries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.getMetricsFromQueries() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Manager.getMetricsFromQueries() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
