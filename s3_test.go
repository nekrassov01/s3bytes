package s3bytes

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestManager_SetBuckets(t *testing.T) {
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
		ctx         context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    []s3types.Bucket
		wantErr bool
	}{
		{
			name: "single bucket",
			fields: fields{
				Client: NewMockClient(
					&MockS3{
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
					nil,
				),
				Buckets: []s3types.Bucket{},
				Prefix:  "",
				Region:  fallbackRegion,
				ctx:     context.Background(),
			},
			want: []s3types.Bucket{
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
				Client: NewMockClient(
					&MockS3{
						ListBucketsFunc: func(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							out := &s3.ListBucketsOutput{
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
								ContinuationToken: nil,
							}
							return out, nil
						},
					},
					nil,
				),
				Buckets: []s3types.Bucket{},
				Prefix:  "",
				Region:  fallbackRegion,
				ctx:     context.Background(),
			},
			want: []s3types.Bucket{
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
			name: "prefix",
			fields: fields{
				Client: NewMockClient(
					&MockS3{
						ListBucketsFunc: func(_ context.Context, params *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							params.Prefix = aws.String("bucket")
							out := &s3.ListBucketsOutput{
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
								ContinuationToken: nil,
							}
							return out, nil
						},
					},
					nil,
				),
				Buckets: []s3types.Bucket{},
				Prefix:  "bucket",
				Region:  fallbackRegion,
				ctx:     context.Background(),
			},
			want: []s3types.Bucket{
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
			name: "pagination",
			fields: fields{
				Client: NewMockClient(
					&MockS3{
						ListBucketsFunc: func(_ context.Context, params *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							if params.ContinuationToken == nil {
								out := &s3.ListBucketsOutput{
									Buckets: []s3types.Bucket{
										{
											Name:         aws.String("bucket0"),
											BucketRegion: aws.String("ap-northeast-1"),
										},
									},
									ContinuationToken: aws.String("token0"),
								}
								return out, nil
							}
							out := &s3.ListBucketsOutput{
								Buckets: []s3types.Bucket{
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
				Buckets: []s3types.Bucket{},
				Prefix:  "",
				Region:  fallbackRegion,
				ctx:     context.Background(),
			},
			want: []s3types.Bucket{
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
				Client: NewMockClient(
					&MockS3{
						ListBucketsFunc: func(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							return nil, errors.New("failed to list buckets")
						},
					},
					nil,
				),
				Buckets: []s3types.Bucket{},
			},
			want:    []s3types.Bucket{},
			wantErr: true,
		},
		{
			name: "no buckets",
			fields: fields{
				Client: NewMockClient(
					&MockS3{
						ListBucketsFunc: func(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
							out := &s3.ListBucketsOutput{
								Buckets:           []s3types.Bucket{},
								ContinuationToken: nil,
							}
							return out, nil
						},
					},
					nil,
				),
				Buckets: []s3types.Bucket{},
			},
			want:    []s3types.Bucket{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{
				Client:  tt.fields.Client,
				Buckets: tt.fields.Buckets,
				Prefix:  tt.fields.Prefix,
				Region:  tt.fields.Region,
				ctx:     tt.fields.ctx,
			}
			if err := man.SetBuckets(); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetBuckets() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(man.Buckets, tt.want) {
				t.Errorf("Manager.SetBuckets() = %v, want %v", man.Buckets, tt.want)
			}
		})
	}
}
