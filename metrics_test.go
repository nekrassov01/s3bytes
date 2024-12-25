package s3bytes

import "testing"

func TestSizeMetric_Label(t *testing.T) {
	type fields struct {
		BucketName    string
		Region        string
		StorageType   StorageType
		Bytes         float64
		ReadableBytes string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal",
			fields: fields{
				BucketName: "bucket0",
			},
			want: "bucket0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &SizeMetric{
				BucketName:    tt.fields.BucketName,
				Region:        tt.fields.Region,
				StorageType:   tt.fields.StorageType,
				Bytes:         tt.fields.Bytes,
				ReadableBytes: tt.fields.ReadableBytes,
			}
			if got := tr.Label(); got != tt.want {
				t.Errorf("SizeMetric.Label() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSizeMetric_Value(t *testing.T) {
	type fields struct {
		BucketName    string
		Region        string
		StorageType   StorageType
		Bytes         float64
		ReadableBytes string
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "normal",
			fields: fields{
				Bytes: 1024,
			},
			want: 1024,
		},
		{
			name: "zero",
			fields: fields{
				Bytes: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &SizeMetric{
				BucketName:    tt.fields.BucketName,
				Region:        tt.fields.Region,
				StorageType:   tt.fields.StorageType,
				Bytes:         tt.fields.Bytes,
				ReadableBytes: tt.fields.ReadableBytes,
			}
			if got := tr.Value(); got != tt.want {
				t.Errorf("SizeMetric.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectMetric_Label(t *testing.T) {
	type fields struct {
		BucketName  string
		Region      string
		StorageType StorageType
		Objects     float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal",
			fields: fields{
				BucketName: "bucket0",
			},
			want: "bucket0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &ObjectMetric{
				BucketName:  tt.fields.BucketName,
				Region:      tt.fields.Region,
				StorageType: tt.fields.StorageType,
				Objects:     tt.fields.Objects,
			}
			if got := tr.Label(); got != tt.want {
				t.Errorf("ObjectMetric.Label() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectMetric_Value(t *testing.T) {
	type fields struct {
		BucketName  string
		Region      string
		StorageType StorageType
		Objects     float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "normal",
			fields: fields{
				Objects: 100,
			},
			want: 100,
		},
		{
			name: "zero",
			fields: fields{
				Objects: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &ObjectMetric{
				BucketName:  tt.fields.BucketName,
				Region:      tt.fields.Region,
				StorageType: tt.fields.StorageType,
				Objects:     tt.fields.Objects,
			}
			if got := tr.Value(); got != tt.want {
				t.Errorf("ObjectMetric.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
