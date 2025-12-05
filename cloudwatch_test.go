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
	"github.com/nekrassov01/filter"
	"golang.org/x/sync/semaphore"
)

func TestManager_getMetrics(t *testing.T) {
	type fields struct {
		client      *Client
		metricName  MetricName
		storageType StorageType
		prefix      *string
		regions     []string
		filterExpr  filterExpr
		filterRaw   string
		sem         *semaphore.Weighted
	}
	type args struct {
		ctx     context.Context
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
				client: newMockClient(
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
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
			},
			args: args{
				ctx: context.Background(),
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
				client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							return nil, errors.New("error")
						},
					},
				),
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
			},
			args: args{
				ctx: context.Background(),
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
			name: "filter returns false",
			fields: fields{
				client: newMockClient(
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
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
				filterExpr:  func() filterExpr { expr, _ := filter.Parse(`bytes == 0`); return expr }(),
				filterRaw:   "bytes == 0",
			},
			args: args{
				ctx: context.Background(),
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
				client:      tt.fields.client,
				metricName:  tt.fields.metricName,
				storageType: tt.fields.storageType,
				prefix:      tt.fields.prefix,
				regions:     tt.fields.regions,
				filterExpr:  tt.fields.filterExpr,
				filterRaw:   tt.fields.filterRaw,
				sem:         tt.fields.sem,
			}
			got, got1, err := man.getMetrics(tt.args.ctx, tt.args.buckets, tt.args.region)
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
		client      *Client
		metricName  MetricName
		storageType StorageType
		prefix      *string
		regions     []string
		filterExpr  filterExpr
		filterRaw   string
		sem         *semaphore.Weighted
	}
	type args struct {
		ctx     context.Context
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
				client: newMockClient(
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
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
			},
			args: args{
				ctx: context.Background(),
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
				client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							return nil, errors.New("error")
						},
					},
				),
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
			},
			args: args{
				ctx: context.Background(),
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
			name: "filter returns false",
			fields: fields{
				client: newMockClient(
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
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
				filterExpr:  func() filterExpr { expr, _ := filter.Parse(`bytes == 0`); return expr }(),
				filterRaw:   "bytes == 0",
			},
			args: args{
				ctx: context.Background(),
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
				client: newMockClient(
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
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
				filterExpr:  nil,
				filterRaw:   "",
			},
			args: args{
				ctx: context.Background(),
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
				client: newMockClient(
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
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
			},
			args: args{
				ctx: context.Background(),
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
				client:      tt.fields.client,
				metricName:  tt.fields.metricName,
				storageType: tt.fields.storageType,
				prefix:      tt.fields.prefix,
				regions:     tt.fields.regions,
				filterExpr:  tt.fields.filterExpr,
				filterRaw:   tt.fields.filterRaw,
				sem:         tt.fields.sem,
			}
			got, got1, err := man.getMetricsFromQueries(tt.args.ctx, tt.args.queries, tt.args.region)
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
