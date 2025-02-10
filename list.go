package s3bytes

import (
	"context"
	"sync"
	"sync/atomic"
)

func (man *Manager) List() (*MetricData, error) {
	var (
		total       int64
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(man.ctx)
		size        = MaxQueries * 2 * len(man.regions)
		metricsChan = make(chan []*Metric, size)
		errorChan   = make(chan error, 1)
		data        = &MetricData{
			Header:  header,
			Metrics: make([]*Metric, 0, size),
		}
	)
	defer cancel()
	errorFunc := func(err error) {
		select {
		case errorChan <- err:
		default:
		}
	}
	for _, region := range man.regions {
		region := region
		if err := man.sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer man.sem.Release(1)
			buckets, err := man.getBuckets(region)
			if err != nil {
				errorFunc(err)
				return
			}
			m, n, err := man.getMetrics(man.buildQueries(buckets), region)
			if err != nil {
				errorFunc(err)
				return
			}
			atomic.AddInt64(&total, n)
			select {
			case metricsChan <- m:
			case <-ctx.Done():
				return
			}
		}()
	}
	go func() {
		wg.Wait()
		close(metricsChan)
	}()
	for {
		select {
		case m, ok := <-metricsChan:
			if !ok {
				data.Total = atomic.LoadInt64(&total)
				return data, nil
			}
			data.Metrics = append(data.Metrics, m...)
		case err := <-errorChan:
			cancel()
			return nil, err
		}
	}
}
