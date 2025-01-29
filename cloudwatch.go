package s3bytes

import (
	"fmt"
	"slices"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// buildQueries builds the metric data queries.
// See: https://docs.aws.amazon.com/AmazonS3/latest/userguide/metrics-dimensions.html
func (man *Manager) buildQueries(buckets []s3types.Bucket) [][]cwtypes.MetricDataQuery {
	var (
		batches        = make([][]cwtypes.MetricDataQuery, 0, 2)
		batch          = make([]cwtypes.MetricDataQuery, 0, MaxQueries)
		namespace      = aws.String("AWS/S3")
		metricName     = aws.String(man.metricName.String())
		bucketNameKey  = aws.String("BucketName")
		storageTypeKey = aws.String("StorageType")
		storageType    = aws.String(man.storageType.String())
		period         = aws.Int32(86400)
		stat           = aws.String("Average")
	)
	for i, bucket := range buckets {
		query := cwtypes.MetricDataQuery{
			Id:    aws.String(fmt.Sprintf("m%d", i)),
			Label: bucket.Name,
			MetricStat: &cwtypes.MetricStat{
				Metric: &cwtypes.Metric{
					Namespace:  namespace,
					MetricName: metricName,
					Dimensions: []cwtypes.Dimension{
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
		if len(batch) == MaxQueries {
			batches = append(batches, batch)
			batch = make([]cwtypes.MetricDataQuery, 0, MaxQueries)
		}
	}
	if len(batch) > 0 {
		batches = append(batches, batch)
	}
	return batches
}

// getMetrics gets the metrics and the total bytes.
func (man *Manager) getMetrics(batches [][]cwtypes.MetricDataQuery, region string) ([]*Metric, int64, error) {
	var (
		total     int64
		metrics   = make([]*Metric, 0, MaxQueries*2)
		startTime = aws.Time(time.Now().Add(-48 * time.Hour))
		endTime   = aws.Time(time.Now())
		opt       = func(o *cloudwatch.Options) { o.Region = region }
	)
	for _, batch := range batches {
		batch := batch
		var token *string
		for {
			in := &cloudwatch.GetMetricDataInput{
				StartTime:         startTime,
				EndTime:           endTime,
				MetricDataQueries: batch,
				NextToken:         token,
			}
			out, err := man.GetMetricData(man.ctx, in, opt)
			if err != nil {
				return nil, 0, err
			}
			for _, result := range out.MetricDataResults {
				var value float64
				if len(result.Values) == 0 {
					value = 0
				} else {
					value = slices.Max(result.Values)
				}
				if ok := man.filterFunc(value); !ok {
					continue
				}
				metric := &Metric{
					BucketName:  aws.ToString(result.Label),
					Region:      region,
					MetricName:  man.metricName,
					StorageType: man.storageType,
					Value:       value,
				}
				metrics = append(metrics, metric)
				atomic.AddInt64(&total, int64(metric.Value))
			}
			token = out.NextToken
			if token == nil {
				break
			}
		}
	}
	return metrics, total, nil
}
