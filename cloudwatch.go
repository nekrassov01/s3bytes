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

var (
	namespace      = aws.String("AWS/S3")
	bucketNameKey  = aws.String("BucketName")
	storageTypeKey = aws.String("StorageType")
	period         = aws.Int32(86400)
	stat           = aws.String("Average")
	startTime      = aws.Time(time.Now().Add(-48 * time.Hour))
	endTime        = aws.Time(time.Now())
)

func (man *Manager) getMetrics(buckets []s3types.Bucket, region string) ([]*Metric, int64, error) {
	var (
		total       int64
		metricName  = aws.String(man.metricName.String())
		storageType = aws.String(man.storageType.String())
		queries     = make([]cwtypes.MetricDataQuery, 0, MaxQueries)
		metrics     = make([]*Metric, 0, MaxQueries*2)
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
		queries = append(queries, query)
		if len(queries) < MaxQueries {
			continue
		}
		m, n, err := man.getMetricsFromQueries(queries, region)
		if err != nil {
			return nil, 0, err
		}
		metrics = append(metrics, m...)
		total += n
		queries = make([]cwtypes.MetricDataQuery, 0, MaxQueries)
	}
	if len(queries) > 0 {
		m, n, err := man.getMetricsFromQueries(queries, region)
		if err != nil {
			return nil, 0, err
		}
		metrics = append(metrics, m...)
		total += n
	}
	return metrics, total, nil
}

func (man *Manager) getMetricsFromQueries(queries []cwtypes.MetricDataQuery, region string) ([]*Metric, int64, error) {
	var (
		total   int64
		token   *string
		metrics = make([]*Metric, 0, MaxQueries)
		opt     = func(o *cloudwatch.Options) { o.Region = region }
	)
	for {
		in := &cloudwatch.GetMetricDataInput{
			StartTime:         startTime,
			EndTime:           endTime,
			MetricDataQueries: queries,
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
	return metrics, total, nil
}
