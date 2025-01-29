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

func TestManager_buildQueries(t *testing.T) {
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
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [][]cwtypes.MetricDataQuery
	}{
		{
			name: "single metric",
			args: args{
				[]s3types.Bucket{
					{
						Name:         aws.String("bucket0"),
						BucketRegion: aws.String("ap-northeast-1"),
					},
				},
			},
			fields: fields{
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
			},
			want: [][]cwtypes.MetricDataQuery{
				{
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
			},
		},
		{
			name: "multiple metrics",
			args: args{
				[]s3types.Bucket{
					{
						Name:         aws.String("bucket0"),
						BucketRegion: aws.String("ap-northeast-1"),
					},
					{
						Name:         aws.String("bucket1"),
						BucketRegion: aws.String("ap-northeast-2"),
					},
				},
			},
			fields: fields{
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
			},
			want: [][]cwtypes.MetricDataQuery{
				{
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
			},
		},
		{
			name: "max queries exceeded",
			args: args{
				[]s3types.Bucket{
					{
						Name:         aws.String("bucket0"),
						BucketRegion: aws.String("ap-northeast-1"),
					},
					{
						Name:         aws.String("bucket1"),
						BucketRegion: aws.String("ap-northeast-2"),
					},
					{
						Name:         aws.String("bucket2"),
						BucketRegion: aws.String("us-east-1"),
					},
				},
			},
			fields: fields{
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
			},
			want: [][]cwtypes.MetricDataQuery{
				{
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
				{
					{
						Id:    aws.String("m2"),
						Label: aws.String("bucket2"),
						MetricStat: &cwtypes.MetricStat{
							Metric: &cwtypes.Metric{
								Namespace:  aws.String("AWS/S3"),
								MetricName: aws.String(MetricNameBucketSizeBytes.String()),
								Dimensions: []cwtypes.Dimension{
									{
										Name:  aws.String("BucketName"),
										Value: aws.String("bucket2"),
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := MaxQueries
			MaxQueries = 2
			defer func() { MaxQueries = original }()
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
			got := man.buildQueries(tt.args.buckets)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.buildQueries() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		batches [][]cwtypes.MetricDataQuery
		region  string
	}
	type want struct {
		metrics []*Metric
		total   int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "bytes",
			args: args{
				batches: [][]cwtypes.MetricDataQuery{
					{
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
				},
				region: "ap-northeast-1",
			},
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
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
			want: want{
				metrics: []*Metric{
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
				total: 2048,
			},
			wantErr: false,
		},
		{
			name: "objects",
			args: args{
				batches: [][]cwtypes.MetricDataQuery{
					{
						{
							Id:    aws.String("m0"),
							Label: aws.String("bucket0"),
							MetricStat: &cwtypes.MetricStat{
								Metric: &cwtypes.Metric{
									Namespace:  aws.String("AWS/S3"),
									MetricName: aws.String(MetricNameNumberOfObjects.String()),
									Dimensions: []cwtypes.Dimension{
										{
											Name:  aws.String("BucketName"),
											Value: aws.String("bucket0"),
										},
										{
											Name:  aws.String("StorageType"),
											Value: aws.String(StorageTypeAllStorageTypes.String()),
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
									MetricName: aws.String(MetricNameNumberOfObjects.String()),
									Dimensions: []cwtypes.Dimension{
										{
											Name:  aws.String("BucketName"),
											Value: aws.String("bucket1"),
										},
										{
											Name:  aws.String("StorageType"),
											Value: aws.String(StorageTypeAllStorageTypes.String()),
										},
									},
								},
								Period: aws.Int32(86400),
								Stat:   aws.String("Average"),
							},
						},
					},
				},
				region: "ap-northeast-1",
			},
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							out := &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{
									{
										Label:  aws.String("bucket0"),
										Values: []float64{200, 100},
									},
									{
										Label:  aws.String("bucket1"),
										Values: []float64{0},
									},
								},
								NextToken: nil,
							}
							return out, nil
						},
					},
				),
				MetricName:  MetricNameNumberOfObjects,
				StorageType: StorageTypeAllStorageTypes,
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			want: want{
				metrics: []*Metric{
					{
						BucketName:  "bucket0",
						Region:      "ap-northeast-1",
						MetricName:  MetricNameNumberOfObjects,
						StorageType: StorageTypeAllStorageTypes,
						Value:       200,
					},
					{
						BucketName:  "bucket1",
						Region:      "ap-northeast-1",
						MetricName:  MetricNameNumberOfObjects,
						StorageType: StorageTypeAllStorageTypes,
						Value:       0,
					},
				},
				total: 200,
			},
			wantErr: false,
		},
		{
			name: "filter func returns false",
			args: args{
				batches: [][]cwtypes.MetricDataQuery{
					{
						{
							Id:    aws.String("m0"),
							Label: aws.String("bucket0"),
							MetricStat: &cwtypes.MetricStat{
								Metric: &cwtypes.Metric{
									Namespace:  aws.String("AWS/S3"),
									MetricName: aws.String(MetricNameNumberOfObjects.String()),
									Dimensions: []cwtypes.Dimension{
										{
											Name:  aws.String("BucketName"),
											Value: aws.String("bucket0"),
										},
										{
											Name:  aws.String("StorageType"),
											Value: aws.String(StorageTypeAllStorageTypes.String()),
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
									MetricName: aws.String(MetricNameNumberOfObjects.String()),
									Dimensions: []cwtypes.Dimension{
										{
											Name:  aws.String("BucketName"),
											Value: aws.String("bucket1"),
										},
										{
											Name:  aws.String("StorageType"),
											Value: aws.String(StorageTypeAllStorageTypes.String()),
										},
									},
								},
								Period: aws.Int32(86400),
								Stat:   aws.String("Average"),
							},
						},
					},
				},
				region: "ap-northeast-1",
			},
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							out := &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{
									{
										Label:  aws.String("bucket0"),
										Values: []float64{200, 100},
									},
									{
										Label:  aws.String("bucket1"),
										Values: []float64{0},
									},
								},
								NextToken: nil,
							}
							return out, nil
						},
					},
				),
				MetricName:  MetricNameNumberOfObjects,
				StorageType: StorageTypeAllStorageTypes,
				filterFunc:  func(float64) bool { return false },
				ctx:         context.Background(),
			},
			want: want{
				metrics: []*Metric{},
				total:   0,
			},
			wantErr: false,
		},
		{
			name: "no value in result",
			args: args{
				batches: [][]cwtypes.MetricDataQuery{
					{
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
				},
				region: "ap-northeast-1",
			},
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							out := &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{
									{
										Label:  aws.String("bucket0"),
										Values: []float64{1024, 2048},
									},
									{
										Label:  aws.String("bucket1"),
										Values: nil,
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
			want: want{
				metrics: []*Metric{
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
				total: 2048,
			},
			wantErr: false,
		},
		{
			name: "pagination",
			args: args{
				batches: [][]cwtypes.MetricDataQuery{
					{
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
				},
				region: "ap-northeast-1",
			},
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
										{
											Label:  aws.String("bucket1"),
											Values: []float64{0},
										},
									},
									NextToken: aws.String("token0"),
								}
								return out, nil
							}
							out := &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{
									{
										Label:  aws.String("bucket2"),
										Values: []float64{2048},
									},
									{
										Label:  aws.String("bucket3"),
										Values: []float64{4096},
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
			want: want{
				metrics: []*Metric{
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
						Value:       0,
					},
					{
						BucketName:  "bucket2",
						Region:      "ap-northeast-1",
						MetricName:  MetricNameBucketSizeBytes,
						StorageType: StorageTypeStandardStorage,
						Value:       2048,
					},
					{
						BucketName:  "bucket3",
						Region:      "ap-northeast-1",
						MetricName:  MetricNameBucketSizeBytes,
						StorageType: StorageTypeStandardStorage,
						Value:       4096,
					},
				},
				total: 7168,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				batches: [][]cwtypes.MetricDataQuery{
					{
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
				},
				region: "ap-northeast-1",
			},
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
			want: want{
				metrics: nil,
				total:   0,
			},
			wantErr: true,
		},
		{
			name: "no metric data",
			args: args{
				batches: [][]cwtypes.MetricDataQuery{},
				region:  "ap-northeast-1",
			},
			fields: fields{
				Client: newMockClient(
					nil,
					&mockCloudWatch{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							out := &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{},
								NextToken:         nil,
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
			want: want{
				metrics: []*Metric{},
				total:   0,
			},
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
			got, total, err := man.getMetrics(tt.args.batches, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.getMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.metrics) {
				t.Errorf("Manager.getMetrics() metrics = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(total, tt.want.total) {
				t.Errorf("Manager.getMetrics() total = %v, want %v", got, tt.want)
			}
		})
	}
}
