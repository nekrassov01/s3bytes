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
		client      *Client
		metricName  MetricName
		storageType StorageType
		prefix      *string
		regions     []string
		sem         *semaphore.Weighted
	}
	type args struct {
		ctx    context.Context
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
				client: newMockClient(
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
				prefix: nil,
			},
			args: args{
				ctx:    context.Background(),
				region: "ap-northeast-1",
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
				client: newMockClient(
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
				prefix: nil,
			},
			args: args{
				ctx:    context.Background(),
				region: "ap-northeast-1",
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
				client: newMockClient(
					&mockS3{
						ListBucketsFunc: func(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							return nil, errors.New("failed to list buckets")
						},
					},
					nil,
				),
			},
			args: args{
				ctx:    context.Background(),
				region: "ap-northeast-1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no buckets",
			fields: fields{
				client: newMockClient(
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
			args: args{
				ctx:    context.Background(),
				region: "ap-northeast-1",
			},
			want:    []types.Bucket{},
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
				sem:         tt.fields.sem,
			}
			got, err := man.getBuckets(tt.args.ctx, tt.args.region)
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
