package s3bytes

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"golang.org/x/sync/semaphore"
)

func TestManager_getBuckets(t *testing.T) {
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
		region string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []types.Bucket
		wantErr bool
	}{
		{
			name: "single bucket",
			fields: fields{
				Client: newMockClient(
					&mockS3{
						ListBucketsFunc: func(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							out := &s3.ListBucketsOutput{
								Buckets: []types.Bucket{
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
					nil,
				),
				Prefix: nil,
				ctx:    context.Background(),
			},
			want: []types.Bucket{
				{
					Name:         aws.String("bucket0"),
					BucketRegion: aws.String("ap-northeast-1"),
				},
			},
			wantErr: false,
		},
		{
			name: "multiple bucket",
			fields: fields{
				Client: newMockClient(
					&mockS3{
						ListBucketsFunc: func(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							out := &s3.ListBucketsOutput{
								Buckets: []types.Bucket{
									{
										Name:         aws.String("bucket0"),
										BucketRegion: aws.String("ap-northeast-1"),
									},
									{
										Name:         aws.String("bucket1"),
										BucketRegion: aws.String("ap-northeast-2"),
									},
								},
								ContinuationToken: nil,
							}
							return out, nil
						},
					},
					nil,
				),
				Prefix: nil,
				ctx:    context.Background(),
			},
			want: []types.Bucket{
				{
					Name:         aws.String("bucket0"),
					BucketRegion: aws.String("ap-northeast-1"),
				},
				{
					Name:         aws.String("bucket1"),
					BucketRegion: aws.String("ap-northeast-2"),
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				Client: newMockClient(
					&mockS3{
						ListBucketsFunc: func(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							return nil, errors.New("failed to list buckets")
						},
					},
					nil,
				),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no buckets",
			fields: fields{
				Client: newMockClient(
					&mockS3{
						ListBucketsFunc: func(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							out := &s3.ListBucketsOutput{
								Buckets:           []types.Bucket{},
								ContinuationToken: nil,
							}
							return out, nil
						},
					},
					nil,
				),
			},
			want:    []types.Bucket{},
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
			got, err := man.getBuckets(tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.getBuckets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.getBuckets() = %v, want %v", got, tt.want)
			}
		})
	}
}
