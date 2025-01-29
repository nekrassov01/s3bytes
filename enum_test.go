package s3bytes

import (
	_ "embed"
	"reflect"
	"testing"
)

func TestOutputType_String(t *testing.T) {
	tests := []struct {
		name string
		tr   OutputType
		want string
	}{
		{
			name: "json",
			tr:   OutputTypeJSON,
			want: "json",
		},
		{
			name: "text",
			tr:   OutputTypeText,
			want: "text",
		},
		{
			name: "markdown",
			tr:   OutputTypeMarkdown,
			want: "markdown",
		},
		{
			name: "backlog",
			tr:   OutputTypeBacklog,
			want: "backlog",
		},
		{
			name: "tsv",
			tr:   OutputTypeTSV,
			want: "tsv",
		},
		{
			name: "none",
			tr:   OutputTypeNone,
			want: "none",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("OutputType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		tr      OutputType
		want    []byte
		wantErr bool
	}{
		{
			name: "json",
			tr:   OutputTypeJSON,
			want: []byte(`"json"`),
		},
		{
			name: "text",
			tr:   OutputTypeText,
			want: []byte(`"text"`),
		},
		{
			name: "markdown",
			tr:   OutputTypeMarkdown,
			want: []byte(`"markdown"`),
		},
		{
			name: "backlog",
			tr:   OutputTypeBacklog,
			want: []byte(`"backlog"`),
		},
		{
			name: "tsv",
			tr:   OutputTypeTSV,
			want: []byte(`"tsv"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("OutputType.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OutputType.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseOutputType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    OutputType
		wantErr bool
	}{
		{
			name: "json",
			args: args{
				s: "json",
			},
			want:    OutputTypeJSON,
			wantErr: false,
		},
		{
			name: "text",
			args: args{
				s: "text",
			},
			want:    OutputTypeText,
			wantErr: false,
		},
		{
			name: "compressed",
			args: args{
				s: "compressed",
			},
			want:    OutputTypeCompressedText,
			wantErr: false,
		},
		{
			name: "markdown",
			args: args{
				s: "markdown",
			},
			want:    OutputTypeMarkdown,
			wantErr: false,
		},
		{
			name: "backlog",
			args: args{
				s: "backlog",
			},
			want:    OutputTypeBacklog,
			wantErr: false,
		},
		{
			name: "tsv",
			args: args{
				s: "tsv",
			},
			want:    OutputTypeTSV,
			wantErr: false,
		},
		{
			name: "unsupported",
			args: args{
				s: "unsupported",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOutputType(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOutputType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseOutputType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricName_String(t *testing.T) {
	tests := []struct {
		name string
		tr   MetricName
		want string
	}{
		{
			name: "bucket size bytes",
			tr:   MetricNameBucketSizeBytes,
			want: "BucketSizeBytes",
		},
		{
			name: "number of objects",
			tr:   MetricNameNumberOfObjects,
			want: "NumberOfObjects",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("MetricName.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricName_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		tr      MetricName
		want    []byte
		wantErr bool
	}{
		{
			name: "bucket size bytes",
			tr:   MetricNameBucketSizeBytes,
			want: []byte(`"BucketSizeBytes"`),
		},
		{
			name: "number of objects",
			tr:   MetricNameNumberOfObjects,
			want: []byte(`"NumberOfObjects"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricName.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricName.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMetricName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    MetricName
		wantErr bool
	}{
		{
			name: "bucket size bytes",
			args: args{
				s: "BucketSizeBytes",
			},
			want:    MetricNameBucketSizeBytes,
			wantErr: false,
		},
		{
			name: "number of objects",
			args: args{
				s: "NumberOfObjects",
			},
			want:    MetricNameNumberOfObjects,
			wantErr: false,
		},
		{
			name: "unsupported",
			args: args{
				s: "unsupported",
			},
			want:    MetricNameNone,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMetricName(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMetricName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseMetricName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorageType_String(t *testing.T) {
	tests := []struct {
		name string
		tr   StorageType
		want string
	}{
		{
			name: "StandardStorage",
			tr:   StorageTypeStandardStorage,
			want: "StandardStorage",
		},
		{
			name: "IntelligentTieringFAStorage",
			tr:   StorageTypeIntelligentTieringFAStorage,
			want: "IntelligentTieringFAStorage",
		},
		{
			name: "IntelligentTieringIAStorage",
			tr:   StorageTypeIntelligentTieringIAStorage,
			want: "IntelligentTieringIAStorage",
		},
		{
			name: "IntelligentTieringAAStorage",
			tr:   StorageTypeIntelligentTieringAAStorage,
			want: "IntelligentTieringAAStorage",
		},
		{
			name: "IntelligentTieringAIAStorage",
			tr:   StorageTypeIntelligentTieringAIAStorage,
			want: "IntelligentTieringAIAStorage",
		},
		{
			name: "IntelligentTieringDAAStorage",
			tr:   StorageTypeIntelligentTieringDAAStorage,
			want: "IntelligentTieringDAAStorage",
		},
		{
			name: "StandardIAStorage",
			tr:   StorageTypeStandardIAStorage,
			want: "StandardIAStorage",
		},
		{
			name: "StandardIASizeOverhead",
			tr:   StorageTypeStandardIASizeOverhead,
			want: "StandardIASizeOverhead",
		},
		{
			name: "StandardIAObjectOverhead",
			tr:   StorageTypeStandardIAObjectOverhead,
			want: "StandardIAObjectOverhead",
		},
		{
			name: "OneZoneIAStorage",
			tr:   StorageTypeOneZoneIAStorage,
			want: "OneZoneIAStorage",
		},
		{
			name: "OneZoneIASizeOverhead",
			tr:   StorageTypeOneZoneIASizeOverhead,
			want: "OneZoneIASizeOverhead",
		},
		{
			name: "ReducedRedundancyStorage",
			tr:   StorageTypeReducedRedundancyStorage,
			want: "ReducedRedundancyStorage",
		},
		{
			name: "GlacierIRSizeOverhead",
			tr:   StorageTypeGlacierIRSizeOverhead,
			want: "GlacierIRSizeOverhead",
		},
		{
			name: "GlacierInstantRetrievalStorage",
			tr:   StorageTypeGlacierInstantRetrievalStorage,
			want: "GlacierInstantRetrievalStorage",
		},
		{
			name: "GlacierStorage",
			tr:   StorageTypeGlacierStorage,
			want: "GlacierStorage",
		},
		{
			name: "GlacierStagingStorage",
			tr:   StorageTypeGlacierStagingStorage,
			want: "GlacierStagingStorage",
		},
		{
			name: "GlacierObjectOverhead",
			tr:   StorageTypeGlacierObjectOverhead,
			want: "GlacierObjectOverhead",
		},
		{
			name: "GlacierS3ObjectOverhead",
			tr:   StorageTypeGlacierS3ObjectOverhead,
			want: "GlacierS3ObjectOverhead",
		},
		{
			name: "DeepArchiveStorage",
			tr:   StorageTypeDeepArchiveStorage,
			want: "DeepArchiveStorage",
		},
		{
			name: "DeepArchiveObjectOverhead",
			tr:   StorageTypeDeepArchiveObjectOverhead,
			want: "DeepArchiveObjectOverhead",
		},
		{
			name: "DeepArchiveS3ObjectOverhead",
			tr:   StorageTypeDeepArchiveS3ObjectOverhead,
			want: "DeepArchiveS3ObjectOverhead",
		},
		{
			name: "DeepArchiveStagingStorage",
			tr:   StorageTypeDeepArchiveStagingStorage,
			want: "DeepArchiveStagingStorage",
		},
		{
			name: "AllStorageTypes",
			tr:   StorageTypeAllStorageTypes,
			want: "AllStorageTypes",
		},
		{
			name: "default",
			tr:   StorageType(999),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("StorageType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorageType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		tr      StorageType
		want    []byte
		wantErr bool
	}{
		{
			name: "StandardStorage",
			tr:   StorageTypeStandardStorage,
			want: []byte(`"StandardStorage"`),
		},
		{
			name: "IntelligentTieringFAStorage",
			tr:   StorageTypeIntelligentTieringFAStorage,
			want: []byte(`"IntelligentTieringFAStorage"`),
		},
		{
			name: "IntelligentTieringIAStorage",
			tr:   StorageTypeIntelligentTieringIAStorage,
			want: []byte(`"IntelligentTieringIAStorage"`),
		},
		{
			name: "IntelligentTieringAAStorage",
			tr:   StorageTypeIntelligentTieringAAStorage,
			want: []byte(`"IntelligentTieringAAStorage"`),
		},
		{
			name: "IntelligentTieringAIAStorage",
			tr:   StorageTypeIntelligentTieringAIAStorage,
			want: []byte(`"IntelligentTieringAIAStorage"`),
		},
		{
			name: "IntelligentTieringDAAStorage",
			tr:   StorageTypeIntelligentTieringDAAStorage,
			want: []byte(`"IntelligentTieringDAAStorage"`),
		},
		{
			name: "StandardIAStorage",
			tr:   StorageTypeStandardIAStorage,
			want: []byte(`"StandardIAStorage"`),
		},
		{
			name: "StandardIASizeOverhead",
			tr:   StorageTypeStandardIASizeOverhead,
			want: []byte(`"StandardIASizeOverhead"`),
		},
		{
			name: "StandardIAObjectOverhead",
			tr:   StorageTypeStandardIAObjectOverhead,
			want: []byte(`"StandardIAObjectOverhead"`),
		},
		{
			name: "OneZoneIAStorage",
			tr:   StorageTypeOneZoneIAStorage,
			want: []byte(`"OneZoneIAStorage"`),
		},
		{
			name: "OneZoneIASizeOverhead",
			tr:   StorageTypeOneZoneIASizeOverhead,
			want: []byte(`"OneZoneIASizeOverhead"`),
		},
		{
			name: "ReducedRedundancyStorage",
			tr:   StorageTypeReducedRedundancyStorage,
			want: []byte(`"ReducedRedundancyStorage"`),
		},
		{
			name: "GlacierIRSizeOverhead",
			tr:   StorageTypeGlacierIRSizeOverhead,
			want: []byte(`"GlacierIRSizeOverhead"`),
		},
		{
			name: "GlacierInstantRetrievalStorage",
			tr:   StorageTypeGlacierInstantRetrievalStorage,
			want: []byte(`"GlacierInstantRetrievalStorage"`),
		},
		{
			name: "GlacierStorage",
			tr:   StorageTypeGlacierStorage,
			want: []byte(`"GlacierStorage"`),
		},
		{
			name: "GlacierStagingStorage",
			tr:   StorageTypeGlacierStagingStorage,
			want: []byte(`"GlacierStagingStorage"`),
		},
		{
			name: "GlacierObjectOverhead",
			tr:   StorageTypeGlacierObjectOverhead,
			want: []byte(`"GlacierObjectOverhead"`),
		},
		{
			name: "GlacierS3ObjectOverhead",
			tr:   StorageTypeGlacierS3ObjectOverhead,
			want: []byte(`"GlacierS3ObjectOverhead"`),
		},
		{
			name: "DeepArchiveStorage",
			tr:   StorageTypeDeepArchiveStorage,
			want: []byte(`"DeepArchiveStorage"`),
		},
		{
			name: "DeepArchiveObjectOverhead",
			tr:   StorageTypeDeepArchiveObjectOverhead,
			want: []byte(`"DeepArchiveObjectOverhead"`),
		},
		{
			name: "DeepArchiveS3ObjectOverhead",
			tr:   StorageTypeDeepArchiveS3ObjectOverhead,
			want: []byte(`"DeepArchiveS3ObjectOverhead"`),
		},
		{
			name: "DeepArchiveStagingStorage",
			tr:   StorageTypeDeepArchiveStagingStorage,
			want: []byte(`"DeepArchiveStagingStorage"`),
		},
		{
			name: "AllStorageTypes",
			tr:   StorageTypeAllStorageTypes,
			want: []byte(`"AllStorageTypes"`),
		},
		{
			name: "default",
			tr:   StorageType(999),
			want: []byte(`""`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("StorageType.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StorageType.MarshalJSON() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}

func TestParseStorageType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    StorageType
		wantErr bool
	}{
		{
			name: "StandardStorage",
			args: args{
				s: "StandardStorage",
			},
			want:    StorageTypeStandardStorage,
			wantErr: false,
		},
		{
			name: "IntelligentTieringFAStorage",
			args: args{
				s: "IntelligentTieringFAStorage",
			},
			want:    StorageTypeIntelligentTieringFAStorage,
			wantErr: false,
		},
		{
			name: "IntelligentTieringIAStorage",
			args: args{
				s: "IntelligentTieringIAStorage",
			},
			want:    StorageTypeIntelligentTieringIAStorage,
			wantErr: false,
		},
		{
			name: "IntelligentTieringAAStorage",
			args: args{
				s: "IntelligentTieringAAStorage",
			},
			want:    StorageTypeIntelligentTieringAAStorage,
			wantErr: false,
		},
		{
			name: "IntelligentTieringAIAStorage",
			args: args{
				s: "IntelligentTieringAIAStorage",
			},
			want:    StorageTypeIntelligentTieringAIAStorage,
			wantErr: false,
		},
		{
			name: "IntelligentTieringDAAStorage",
			args: args{
				s: "IntelligentTieringDAAStorage",
			},
			want:    StorageTypeIntelligentTieringDAAStorage,
			wantErr: false,
		},
		{
			name: "StandardIAStorage",
			args: args{
				s: "StandardIAStorage",
			},
			want:    StorageTypeStandardIAStorage,
			wantErr: false,
		},
		{
			name: "StandardIASizeOverhead",
			args: args{
				s: "StandardIASizeOverhead",
			},
			want:    StorageTypeStandardIASizeOverhead,
			wantErr: false,
		},
		{
			name: "StandardIAObjectOverhead",
			args: args{
				s: "StandardIAObjectOverhead",
			},
			want:    StorageTypeStandardIAObjectOverhead,
			wantErr: false,
		},
		{
			name: "OneZoneIAStorage",
			args: args{
				s: "OneZoneIAStorage",
			},
			want:    StorageTypeOneZoneIAStorage,
			wantErr: false,
		},
		{
			name: "OneZoneIASizeOverhead",
			args: args{
				s: "OneZoneIASizeOverhead",
			},
			want:    StorageTypeOneZoneIASizeOverhead,
			wantErr: false,
		},
		{
			name: "ReducedRedundancyStorage",
			args: args{
				s: "ReducedRedundancyStorage",
			},
			want:    StorageTypeReducedRedundancyStorage,
			wantErr: false,
		},
		{
			name: "GlacierIRSizeOverhead",
			args: args{
				s: "GlacierIRSizeOverhead",
			},
			want:    StorageTypeGlacierIRSizeOverhead,
			wantErr: false,
		},
		{
			name: "GlacierInstantRetrievalStorage",
			args: args{
				s: "GlacierInstantRetrievalStorage",
			},
			want:    StorageTypeGlacierInstantRetrievalStorage,
			wantErr: false,
		},
		{
			name: "GlacierStorage",
			args: args{
				s: "GlacierStorage",
			},
			want:    StorageTypeGlacierStorage,
			wantErr: false,
		},
		{
			name: "GlacierStagingStorage",
			args: args{
				s: "GlacierStagingStorage",
			},
			want:    StorageTypeGlacierStagingStorage,
			wantErr: false,
		},
		{
			name: "GlacierObjectOverhead",
			args: args{
				s: "GlacierObjectOverhead",
			},
			want:    StorageTypeGlacierObjectOverhead,
			wantErr: false,
		},
		{
			name: "GlacierS3ObjectOverhead",
			args: args{
				s: "GlacierS3ObjectOverhead",
			},
			want:    StorageTypeGlacierS3ObjectOverhead,
			wantErr: false,
		},
		{
			name: "DeepArchiveStorage",
			args: args{
				s: "DeepArchiveStorage",
			},
			want:    StorageTypeDeepArchiveStorage,
			wantErr: false,
		},
		{
			name: "DeepArchiveObjectOverhead",
			args: args{
				s: "DeepArchiveObjectOverhead",
			},
			want:    StorageTypeDeepArchiveObjectOverhead,
			wantErr: false,
		},
		{
			name: "DeepArchiveS3ObjectOverhead",
			args: args{
				s: "DeepArchiveS3ObjectOverhead",
			},
			want:    StorageTypeDeepArchiveS3ObjectOverhead,
			wantErr: false,
		},
		{
			name: "DeepArchiveStagingStorage",
			args: args{
				s: "DeepArchiveStagingStorage",
			},
			want:    StorageTypeDeepArchiveStagingStorage,
			wantErr: false,
		},
		{
			name: "AllStorageTypes",
			args: args{
				s: "AllStorageTypes",
			},
			want:    StorageTypeAllStorageTypes,
			wantErr: false,
		},
		{
			name: "UnsupportedStorageType",
			args: args{
				s: "UnsupportedStorageType",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "EmptyString",
			args: args{
				s: "",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "WhitespaceString",
			args: args{
				s: " ",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "NumericString",
			args: args{
				s: "123",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "PartialMatch",
			args: args{
				s: "Standard",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "CaseSensitiveMismatch",
			args: args{
				s: "standardstorage",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStorageType(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStorageType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseStorageType() = %v, want %v", got, tt.want)
			}
		})
	}
}
