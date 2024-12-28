package s3bytes

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/dustin/go-humanize"
)

// SetQueries sets the metric data queries for the cloudwatch client.
// https://docs.aws.amazon.com/AmazonS3/latest/userguide/metrics-dimensions.html
func (man *Manager) SetQueries() error {
	if man.MetricName == MetricNameBucketSizeBytes && man.StorageType == StorageTypeAllStorageTypes {
		return errors.New("BucketSizeBytes metric does not support AllStorageTypes")
	}
	if man.MetricName == MetricNameNumberOfObjects && man.StorageType != StorageTypeAllStorageTypes {
		return errors.New("NumberOfObjects metric only supports AllStorageTypes")
	}
	var (
		batch          = make([]types.MetricDataQuery, 0, man.MaxQueries)
		namespace      = aws.String("AWS/S3")
		metricName     = aws.String(man.MetricName.String())
		bucketNameKey  = aws.String("BucketName")
		storageTypeKey = aws.String("StorageType")
		storageType    = aws.String(man.StorageType.String())
		period         = aws.Int32(86400)
		stat           = aws.String("Average")
	)
	for i, bucket := range man.Buckets {
		query := types.MetricDataQuery{
			Id:    aws.String(fmt.Sprintf("m%d", i)),
			Label: bucket.Name,
			MetricStat: &types.MetricStat{
				Metric: &types.Metric{
					Namespace:  namespace,
					MetricName: metricName,
					Dimensions: []types.Dimension{
						{
							Name:  bucketNameKey,
							Value: bucket.Name,
						},
						{
							Name:  storageTypeKey,
							Value: storageType,
						},
					},
				},
				Period: period,
				Stat:   stat,
			},
		}
		batch = append(batch, query)
		if len(batch) == man.MaxQueries {
			man.Batches = append(man.Batches, batch)
			batch = make([]types.MetricDataQuery, 0, man.MaxQueries)
		}
	}
	if len(batch) > 0 {
		man.Batches = append(man.Batches, batch)
	}
	return nil
}

// SetData sets the metric data for the cloudwatch client.
func (man *Manager) SetData() error {
	var (
		startTime = aws.Time(time.Now().Add(-48 * time.Hour))
		endTime   = aws.Time(time.Now())
		opt       = func(o *cloudwatch.Options) { o.Region = man.Region }
	)
	for _, batch := range man.Batches {
		batch := batch
		var token *string
		for {
			in := &cloudwatch.GetMetricDataInput{
				StartTime:         startTime,
				EndTime:           endTime,
				MetricDataQueries: batch,
				NextToken:         token,
			}
			out, err := man.cw.GetMetricData(man.ctx, in, opt)
			if err != nil {
				return err
			}
			for _, result := range out.MetricDataResults {
				var (
					metric Metric
					size   float64
				)
				if len(result.Values) == 0 {
					size = 0
				} else {
					size = slices.Max(result.Values)
				}
				if ok := man.filterFunc(size); !ok {
					continue
				}
				switch man.MetricName {
				case MetricNameBucketSizeBytes:
					metric = &SizeMetric{
						BucketName:    aws.ToString(result.Label),
						Region:        man.Region,
						StorageType:   man.StorageType,
						Bytes:         size,
						ReadableBytes: humanize.IBytes(uint64(size)),
					}
				case MetricNameNumberOfObjects:
					metric = &ObjectMetric{
						BucketName:  aws.ToString(result.Label),
						Region:      man.Region,
						StorageType: man.StorageType,
						Objects:     size,
					}
				default:
				}
				man.Metrics = append(man.Metrics, metric)
			}
			token = out.NextToken
			if token == nil {
				break
			}
		}
	}
	return nil
}
