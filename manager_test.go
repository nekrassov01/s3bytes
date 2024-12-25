package s3bytes

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestNewManager(t *testing.T) {
	type args struct {
		ctx         context.Context
		client      *Client
		region      string
		prefix      string
		expr        string
		metricName  MetricName
		storageType StorageType
	}
	tests := []struct {
		name    string
		args    args
		want    *Manager
		wantErr bool
	}{
		{
			name: "empty client",
			args: args{
				ctx:         context.Background(),
				client:      NewMockClient(&MockS3{}, &MockCW{}),
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
				prefix:      "",
				expr:        "< 100",
				region:      "ap-northeast-1",
			},
			want: &Manager{
				Client:      NewMockClient(&MockS3{}, &MockCW{}),
				Buckets:     []s3types.Bucket{},
				Batches:     [][]cwtypes.MetricDataQuery{},
				Metrics:     []Metric{},
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				MaxQueries:  maxQueries,
				Prefix:      "",
				Region:      "ap-northeast-1",
				ctx:         context.Background(),
			},
			wantErr: false,
		},
		{
			name: "nil client",
			args: args{
				ctx:         context.Background(),
				client:      nil,
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
				prefix:      "",
				expr:        "> 100",
				region:      "ap-northeast-1",
			},
			want: &Manager{
				Client:      nil,
				Buckets:     []s3types.Bucket{},
				Batches:     [][]cwtypes.MetricDataQuery{},
				Metrics:     []Metric{},
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				MaxQueries:  maxQueries,
				Prefix:      "",
				Region:      "ap-northeast-1",
				ctx:         context.Background(),
			},
			wantErr: false,
		},
		{
			name: "invalid expr",
			args: args{
				ctx:         context.Background(),
				client:      nil,
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
				prefix:      "",
				expr:        "abcd",
				region:      "ap-northeast-1",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewManager(tt.args.ctx, tt.args.client, tt.args.region, tt.args.prefix, tt.args.expr, tt.args.metricName, tt.args.storageType)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if !reflect.DeepEqual(got.Client, tt.want.Client) {
					t.Errorf("NewManager() Client = %v, want %v", got.Client, tt.want.Client)
				}
				if len(got.Buckets) != len(tt.want.Buckets) {
					t.Errorf("NewManager() Buckets length = %d, want %d", len(got.Buckets), len(tt.want.Buckets))
				}
				if len(got.Batches) != len(tt.want.Batches) {
					t.Errorf("NewManager() Batches length = %d, want %d", len(got.Batches), len(tt.want.Batches))
				}
				if len(got.Metrics) != len(tt.want.Metrics) {
					t.Errorf("NewManager() Metrics length = %d, want %d", len(got.Metrics), len(tt.want.Metrics))
				}
				if got.MetricName != tt.want.MetricName {
					t.Errorf("NewManager() MetricName = %v, want %v", got.MetricName, tt.want.MetricName)
				}
				if got.StorageType != tt.want.StorageType {
					t.Errorf("NewManager() StorageType = %v, want %v", got.StorageType, tt.want.StorageType)
				}
				if got.MaxQueries != tt.want.MaxQueries {
					t.Errorf("NewManager() MaxQueries = %d, want %d", got.MaxQueries, tt.want.MaxQueries)
				}
				if got.Prefix != tt.want.Prefix {
					t.Errorf("NewManager() Prefix = %v, want %v", got.Prefix, tt.want.Prefix)
				}
				if got.Region != tt.want.Region {
					t.Errorf("NewManager() Prefix = %v, want %v", got.Region, tt.want.Region)
				}
				if got.ctx != tt.want.ctx {
					t.Errorf("NewManager() ctx = %v, want %v", got.ctx, tt.want.ctx)
				}
			}
		})
	}
}

func TestManager_String(t *testing.T) {
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
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal",
			fields: fields{
				Client:      NewMockClient(&MockS3{}, &MockCW{}),
				Buckets:     []s3types.Bucket{{Name: aws.String("bucket0")}},
				Batches:     [][]cwtypes.MetricDataQuery{{}},
				Metrics:     []Metric{},
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				MaxQueries:  maxQueries,
				Prefix:      "",
				Region:      "ap-northeast-1",
				ctx:         context.Background(),
			},
			want: `{
  "Buckets": [
    {
      "BucketRegion": null,
      "CreationDate": null,
      "Name": "bucket0"
    }
  ],
  "Batches": [
    []
  ],
  "Metrics": [],
  "MetricName": "BucketSizeBytes",
  "StorageType": "StandardStorage",
  "MaxQueries": ` + strconv.Itoa(maxQueries) + `,
  "Prefix": "",
  "Region": "ap-northeast-1"
}`,
		},
		{
			name:   "empty",
			fields: fields{},
			want: `{
  "Buckets": null,
  "Batches": null,
  "Metrics": null,
  "MetricName": "BucketSizeBytes",
  "StorageType": "StandardStorage",
  "MaxQueries": 0,
  "Prefix": "",
  "Region": ""
}`,
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
			if got := man.String(); got != tt.want {
				t.Errorf("Manager.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Debug(t *testing.T) {
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
		name   string
		fields fields
	}{
		{
			name: "normal",
			fields: fields{
				Client:      NewMockClient(&MockS3{}, &MockCW{}),
				Buckets:     []s3types.Bucket{{Name: aws.String("bucket1")}},
				Batches:     [][]cwtypes.MetricDataQuery{{}},
				Metrics:     []Metric{},
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				MaxQueries:  maxQueries,
				Prefix:      "",
				Region:      "ap-northeast-1",
				ctx:         context.Background(),
			},
		},
		{
			name:   "empty",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
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
			man.Debug()
		})
	}
}

func TestManager_eval(t *testing.T) {
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
	type args struct {
		expr string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		cases  []struct {
			v       float64
			want    bool
			wantErr bool
		}
		wantErr bool
	}{
		{
			name: "> 100",
			args: args{
				expr: "> 100",
			},
			cases: []struct {
				v       float64
				want    bool
				wantErr bool
			}{
				{v: 150, want: true, wantErr: false},
				{v: 100, want: false, wantErr: false},
				{v: 50, want: false, wantErr: false},
			},
			wantErr: false,
		},
		{
			name: ">= 200",
			args: args{
				expr: ">= 200",
			},
			cases: []struct {
				v       float64
				want    bool
				wantErr bool
			}{
				{v: 250, want: true, wantErr: false},
				{v: 200, want: true, wantErr: false},
				{v: 150, want: false, wantErr: false},
			},
			wantErr: false,
		},
		{
			name: "< 300",
			args: args{
				expr: "< 300",
			},
			cases: []struct {
				v       float64
				want    bool
				wantErr bool
			}{
				{v: 250, want: true, wantErr: false},
				{v: 300, want: false, wantErr: false},
				{v: 350, want: false, wantErr: false},
			},
			wantErr: false,
		},
		{
			name: "<= 400",
			args: args{
				expr: "<= 400",
			},
			cases: []struct {
				v       float64
				want    bool
				wantErr bool
			}{
				{v: 350, want: true, wantErr: false},
				{v: 400, want: true, wantErr: false},
				{v: 450, want: false, wantErr: false},
			},
			wantErr: false,
		},
		{
			name: "== 500",
			args: args{
				expr: "== 500",
			},
			cases: []struct {
				v       float64
				want    bool
				wantErr bool
			}{
				{v: 500, want: true, wantErr: false},
				{v: 450, want: false, wantErr: false},
				{v: 550, want: false, wantErr: false},
			},
			wantErr: false,
		},
		{
			name: "!= 600",
			args: args{
				expr: "!= 600",
			},
			cases: []struct {
				v       float64
				want    bool
				wantErr bool
			}{
				{v: 500, want: true, wantErr: false},
				{v: 700, want: true, wantErr: false},
				{v: 600, want: false, wantErr: false},
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				expr: "",
			},
			cases: []struct {
				v       float64
				want    bool
				wantErr bool
			}{
				{v: 0, want: true, wantErr: false},
			},
			wantErr: false,
		},
		{
			name: "invalid syntax",
			args: args{
				expr: "abcd",
			},
			cases: []struct {
				v       float64
				want    bool
				wantErr bool
			}{},
			wantErr: true,
		},
		{
			name: "invalid operator",
			args: args{
				expr: "=~ 100",
			},
			cases: []struct {
				v       float64
				want    bool
				wantErr bool
			}{},
			wantErr: true,
		},
		{
			name: "non numeric value",
			args: args{
				expr: "> abc",
			},
			cases: []struct {
				v       float64
				want    bool
				wantErr bool
			}{},
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
			got, err := man.eval(tt.args.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, tc := range tt.cases {
				result := got(tc.v)
				if result != tc.want {
					t.Errorf("Manager.eval(): v = %v, got = %v, result = %v", tc.v, tc.want, result)
				}
			}
		})
	}
}
