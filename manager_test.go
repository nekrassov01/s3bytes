package s3bytes

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/sync/semaphore"
)

func TestNewManager(t *testing.T) {
	type args struct {
		ctx    context.Context
		client *Client
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		{
			name: "empty client",
			args: args{
				ctx:    context.Background(),
				client: newMockClient(&mockS3{}, &mockCloudWatch{}),
			},
			want: &Manager{
				Client:      newMockClient(&mockS3{}, &mockCloudWatch{}),
				metricName:  MetricNameNone,
				storageType: StorageTypeNone,
				prefix:      nil,
				regions:     DefaultRegions,
				ctx:         context.Background(),
			},
		},
		{
			name: "nil client",
			args: args{
				ctx:    context.Background(),
				client: nil,
			},
			want: &Manager{
				Client:      nil,
				metricName:  MetricNameNone,
				storageType: StorageTypeNone,
				prefix:      nil,
				regions:     DefaultRegions,
				ctx:         context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewManager(tt.args.ctx, tt.args.client)
			opts := cmpopts.IgnoreUnexported(*got)
			if diff := cmp.Diff(got, tt.want, opts); diff != "" {
				t.Errorf("NewManager() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestManager_SetRegions(t *testing.T) {
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
		regions []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				regions: nil,
			},
			wantErr: false,
		},
		{
			name: "valid",
			args: args{
				regions: []string{"ap-northeast-1", "ap-northeast-2"},
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				regions: []string{"ap-northeast-1", "invalid"},
			},
			wantErr: true,
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
			if err := man.SetRegion(tt.args.regions); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetRegion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_SetPrefix(t *testing.T) {
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
		prefix string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				prefix: "",
			},
			wantErr: false,
		},
		{
			name: "valid",
			args: args{
				prefix: "test",
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				prefix: "test/",
			},
			wantErr: true,
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
			if err := man.SetPrefix(tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetPrefix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_SetFilter(t *testing.T) {
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
				metricName:  tt.fields.MetricName,
				storageType: tt.fields.StorageType,
				prefix:      tt.fields.Prefix,
				regions:     tt.fields.Regions,
				filterFunc:  tt.fields.filterFunc,
				sem:         tt.fields.sem,
				ctx:         tt.fields.ctx,
			}
			if err := man.SetFilter(tt.args.expr); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_SetMetric(t *testing.T) {
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
		metricName  MetricName
		storageType StorageType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeStandardStorage,
			},
			wantErr: false,
		},
		{
			name: "invalid combination",
			args: args{
				metricName:  MetricNameBucketSizeBytes,
				storageType: StorageTypeAllStorageTypes,
			},
			wantErr: true,
		},
		{
			name: "invalid combination",
			args: args{
				metricName:  MetricNameNumberOfObjects,
				storageType: StorageTypeStandardStorage,
			},
			wantErr: true,
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
			if err := man.SetMetric(tt.args.metricName, tt.args.storageType); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_String(t *testing.T) {
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
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal",
			fields: fields{
				Client:      newMockClient(&mockS3{}, &mockCloudWatch{}),
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				Prefix:      nil,
				ctx:         context.Background(),
			},
			want: `{
  "metricName": "BucketSizeBytes",
  "storageType": "StandardStorage",
  "prefix": null,
  "regions": null
}`,
		},
		{
			name: "prefixed",
			fields: fields{
				Client:      newMockClient(&mockS3{}, &mockCloudWatch{}),
				MetricName:  MetricNameBucketSizeBytes,
				StorageType: StorageTypeStandardStorage,
				Prefix:      aws.String("test"),
				ctx:         context.Background(),
			},
			want: `{
  "metricName": "BucketSizeBytes",
  "storageType": "StandardStorage",
  "prefix": "test",
  "regions": null
}`,
		},
		{
			name:   "empty",
			fields: fields{},
			want: `{
  "metricName": "none",
  "storageType": "none",
  "prefix": null,
  "regions": null
}`,
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
			if diff := cmp.Diff(man.String(), tt.want); diff != "" {
				t.Errorf("Manager.String() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
