package s3bytes

import (
	"fmt"
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/pkg/browser"
)

func getPieItems(data *MetricData) (string, []opts.PieData) {
	var (
		othersTotal = 0.0
		items       = make([]opts.PieData, 0, MaxChartItems)
		title       = ""
	)
	for i, metric := range data.Metrics {
		if metric.Value == 0 {
			continue
		}
		if title == "" {
			switch metric.MetricName {
			case MetricNameBucketSizeBytes:
				title = "Bucket Size Bytes"
			case MetricNameNumberOfObjects:
				title = "Number Of Objects"
			}
		}
		if i < MaxChartItems-1 {
			item := opts.PieData{
				Name:  metric.BucketName,
				Value: metric.Value,
			}
			items = append(items, item)
		} else {
			othersTotal += metric.Value
		}
	}
	if othersTotal > 0 {
		item := opts.PieData{
			Name:  "others",
			Value: othersTotal,
		}
		items = append(items, item)
	}
	return title, items
}

func newPie(title string, items []opts.PieData) *charts.Pie {
	if len(items) == 0 {
		return nil
	}
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  "light",
			Width:  "1280px",
			Height: "720px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: title,
			Left:  "center",
		}),
		charts.WithLegendOpts(opts.Legend{
			Orient: "vertical",
			X:      "right",
			Y:      "bottom",
		}),
	)
	pie.AddSeries("", items)
	pie.SetSeriesOptions(
		charts.WithLabelOpts(opts.Label{
			Show:      opts.Bool(true),
			Position:  "inside",
			Formatter: "{d}%",
		}),
	)
	return pie
}

func render(pie *charts.Pie) error {
	if pie == nil {
		return nil
	}
	title := "s3bytes"
	page := components.NewPage()
	page.SetPageTitle(title)
	page.AddCharts(pie)
	fname := title + ".html"
	i := 1
	for {
		if _, err := os.Stat(fname); err != nil {
			if os.IsNotExist(err) {
				break
			}
			return err
		}
		fname = fmt.Sprintf("%s%d.html", title, i)
		i++
	}
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	if err := page.Render(io.MultiWriter(f)); err != nil {
		return err
	}
	browser.OpenFile(fname) //nolint:errcheck
	return nil
}
