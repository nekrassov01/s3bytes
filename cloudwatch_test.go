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
)

func TestManager_SetQueries(t *testing.T) {
	type fields struct {
		Client      *Client
		Buckets     []s3types.Bucket
		Batches     [][]cwtypes.MetricDataQuery
		Metrics     []Metric
		MetricName  MetricName
		StorageType StorageType
		MaxQueries  int
		Prefix      string
		Region      string
		filterFunc  func(float64) bool
		ctx         context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    [][]cwtypes.MetricDataQuery
		wantErr bool
	}{
		{
			name: "single metric",
			fields: fields{
				Buckets: []s3types.Bucket{
					{
						Name:         aws.String("bucket0"),
						BucketRegion: aws.String("ap-northeast-1"),
					},
				},
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				Batches:     [][]cwtypes.MetricDataQuery{},
				MaxQueries:  maxQueries,
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
			wantErr: false,
		},
		{
			name: "multiple metrics",
			fields: fields{
				Buckets: []s3types.Bucket{
					{
						Name:         aws.String("bucket0"),
						BucketRegion: aws.String("ap-northeast-1"),
					},
					{
						Name:         aws.String("bucket1"),
						BucketRegion: aws.String("ap-northeast-2"),
					},
				},
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				Batches:     [][]cwtypes.MetricDataQuery{},
				MaxQueries:  maxQueries,
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
			wantErr: false,
		},
		{
			name: "max queries exceeded",
			fields: fields{
				Buckets: []s3types.Bucket{
					{
						Name:         aws.String("bucket0"),
						BucketRegion: aws.String("ap-northeast-1"),
					},
					{
						Name:         aws.String("bucket1"),
						BucketRegion: aws.String("ap-northeast-2"),
					},
				},
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				Batches:     [][]cwtypes.MetricDataQuery{},
				MaxQueries:  1,
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
				{
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
			wantErr: false,
		},
		{
			name: "unsupported metric name and storage type combinations 1",
			fields: fields{
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeAllStorageTypes,
				MaxQueries:  maxQueries,
				Batches:     [][]cwtypes.MetricDataQuery{},
			},
			want:    [][]cwtypes.MetricDataQuery{},
			wantErr: true,
		},
		{
			name: "unsupported metric name and storage type combinations 2",
			fields: fields{
				MetricName:  MetricNameNumberOfObjects,
				StorageType: StorageTypeStandardStorage,
				MaxQueries:  maxQueries,
				Batches:     [][]cwtypes.MetricDataQuery{},
			},
			want:    [][]cwtypes.MetricDataQuery{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				Client:      tt.fields.Client,
				Buckets:     tt.fields.Buckets,
				Batches:     tt.fields.Batches,
				Metrics:     tt.fields.Metrics,
				MetricName:  tt.fields.MetricName,
				StorageType: tt.fields.StorageType,
				MaxQueries:  tt.fields.MaxQueries,
				Prefix:      tt.fields.Prefix,
				Region:      tt.fields.Region,
				filterFunc:  tt.fields.filterFunc,
				ctx:         tt.fields.ctx,
			}
			if err := man.SetQueries(); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetQueries() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(man.Batches, tt.want) {
				t.Errorf("Manager.SetQueries() = %v, want %v", man.Batches, tt.want)
			}
		})
	}
}

func TestManager_SetData(t *testing.T) {
	type fields struct {
		Client      *Client
		Buckets     []s3types.Bucket
		Batches     [][]cwtypes.MetricDataQuery
		Metrics     []Metric
		MetricName  MetricName
		StorageType StorageType
		MaxQueries  int
		Prefix      string
		Region      string
		filterFunc  func(float64) bool
		ctx         context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Metric
		wantErr bool
	}{
		{
			name: "bytes",
			fields: fields{
				Client: NewMockClient(
					nil,
					&MockCW{
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
				Batches: [][]cwtypes.MetricDataQuery{
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
				Metrics:     []Metric{},
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				MaxQueries:  maxQueries,
				Region:      "ap-northeast-1",
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			want: []Metric{
				&SizeMetric{
					BucketName:    "bucket0",
					Region:        "ap-northeast-1",
					StorageType:   StorageTypeStandardStorage,
					Bytes:         2048,
					ReadableBytes: "2.0 KiB",
				},
				&SizeMetric{
					BucketName:    "bucket1",
					Region:        "ap-northeast-1",
					StorageType:   StorageTypeStandardStorage,
					Bytes:         0,
					ReadableBytes: "0 B",
				},
			},
			wantErr: false,
		},
		{
			name: "objects",
			fields: fields{
				Client: NewMockClient(
					nil,
					&MockCW{
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
				Batches: [][]cwtypes.MetricDataQuery{
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
				Metrics:     []Metric{},
				MetricName:  MetricNameNumberOfObjects,
				StorageType: StorageTypeAllStorageTypes,
				MaxQueries:  maxQueries,
				Region:      "ap-northeast-1",
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			want: []Metric{
				&ObjectMetric{
					BucketName:  "bucket0",
					Region:      "ap-northeast-1",
					StorageType: StorageTypeAllStorageTypes,
					Objects:     200,
				},
				&ObjectMetric{
					BucketName:  "bucket1",
					Region:      "ap-northeast-1",
					StorageType: StorageTypeAllStorageTypes,
					Objects:     0,
				},
			},
			wantErr: false,
		},
		{
			name: "filter func returns false",
			fields: fields{
				Client: NewMockClient(
					nil,
					&MockCW{
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
				Batches: [][]cwtypes.MetricDataQuery{
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
				Metrics:     []Metric{},
				MetricName:  MetricNameNumberOfObjects,
				StorageType: StorageTypeAllStorageTypes,
				MaxQueries:  maxQueries,
				Region:      "ap-northeast-1",
				filterFunc:  func(float64) bool { return false },
				ctx:         context.Background(),
			},
			want:    []Metric{},
			wantErr: false,
		},
		{
			name: "no value in result",
			fields: fields{
				Client: NewMockClient(
					nil,
					&MockCW{
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
				Batches: [][]cwtypes.MetricDataQuery{
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
				Metrics:     []Metric{},
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				MaxQueries:  maxQueries,
				Region:      "ap-northeast-1",
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			want: []Metric{
				&SizeMetric{
					BucketName:    "bucket0",
					Region:        "ap-northeast-1",
					StorageType:   StorageTypeStandardStorage,
					Bytes:         2048,
					ReadableBytes: "2.0 KiB",
				},
				&SizeMetric{
					BucketName:    "bucket1",
					Region:        "ap-northeast-1",
					StorageType:   StorageTypeStandardStorage,
					Bytes:         0,
					ReadableBytes: "0 B",
				},
			},
			wantErr: false,
		},
		{
			name: "pagination",
			fields: fields{
				Client: NewMockClient(
					nil,
					&MockCW{
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
				Batches: [][]cwtypes.MetricDataQuery{
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
				Metrics:     []Metric{},
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				MaxQueries:  maxQueries,
				Region:      "ap-northeast-1",
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			want: []Metric{
				&SizeMetric{
					BucketName:    "bucket0",
					Region:        "ap-northeast-1",
					StorageType:   StorageTypeStandardStorage,
					Bytes:         1024,
					ReadableBytes: "1.0 KiB",
				},
				&SizeMetric{
					BucketName:    "bucket1",
					Region:        "ap-northeast-1",
					StorageType:   StorageTypeStandardStorage,
					Bytes:         0,
					ReadableBytes: "0 B",
				},
				&SizeMetric{
					BucketName:    "bucket2",
					Region:        "ap-northeast-1",
					StorageType:   StorageTypeStandardStorage,
					Bytes:         2048,
					ReadableBytes: "2.0 KiB",
				},
				&SizeMetric{
					BucketName:    "bucket3",
					Region:        "ap-northeast-1",
					StorageType:   StorageTypeStandardStorage,
					Bytes:         4096,
					ReadableBytes: "4.0 KiB",
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				Client: NewMockClient(
					nil,
					&MockCW{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							return nil, errors.New("error")
						},
					},
				),
				Batches: [][]cwtypes.MetricDataQuery{
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
				Metrics:     []Metric{},
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				MaxQueries:  maxQueries,
				Region:      "ap-northeast-1",
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			want:    []Metric{},
			wantErr: true,
		},
		{
			name: "no metric data",
			fields: fields{
				Client: NewMockClient(
					nil,
					&MockCW{
						GetMetricDataFunc: func(_ context.Context, _ *cloudwatch.GetMetricDataInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
							out := &cloudwatch.GetMetricDataOutput{
								MetricDataResults: []cwtypes.MetricDataResult{},
								NextToken:         nil,
							}
							return out, nil
						},
					},
				),
				Batches:     [][]cwtypes.MetricDataQuery{},
				Metrics:     []Metric{},
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				MaxQueries:  maxQueries,
				Region:      "ap-northeast-1",
				filterFunc:  func(float64) bool { return true },
				ctx:         context.Background(),
			},
			want:    []Metric{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				Client:      tt.fields.Client,
				Buckets:     tt.fields.Buckets,
				Batches:     tt.fields.Batches,
				Metrics:     tt.fields.Metrics,
				MetricName:  tt.fields.MetricName,
				StorageType: tt.fields.StorageType,
				MaxQueries:  tt.fields.MaxQueries,
				Prefix:      tt.fields.Prefix,
				Region:      tt.fields.Region,
				filterFunc:  tt.fields.filterFunc,
				ctx:         tt.fields.ctx,
			}
			if err := man.SetData(); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetData() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(man.Metrics, tt.want) {
				t.Errorf("Manager.Metrics() = %v, want %v", man.Metrics, tt.want)
			}
		})
	}
}
