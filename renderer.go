package s3bytes

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/nekrassov01/mintab"
)

// Renderer is a renderer struct for the s3bytes package.
type Renderer struct {
	Metrics    []Metric
	MetricName MetricName
	OutputType OutputType

	w io.Writer
}

// NewRenderer creates a new renderer.
func NewRenderer(w io.Writer, metrics []Metric, metricName MetricName, outputType OutputType) *Renderer {
	return &Renderer{
		Metrics:    metrics,
		MetricName: metricName,
		OutputType: outputType,
		w:          w,
	}
}

// String returns a string representation of the renderer.
func (ren *Renderer) String() string {
	b, _ := json.MarshalIndent(ren, "", "  ")
	return string(b)
}

// Render renders the metrics.
func (ren *Renderer) Render() error {
	switch ren.OutputType {
	case OutputTypeJSON:
		return ren.toJSON()
	case OutputTypeText, OutputTypeMarkdown, OutputTypeBacklog:
		return ren.toTable()
	case OutputTypeTSV:
		return ren.toTSV()
	default:
		return nil
	}
}

func (ren *Renderer) toJSON() error {
	b := json.NewEncoder(ren.w)
	b.SetIndent("", "  ")
	return b.Encode(ren.Metrics)
}

func (ren *Renderer) toTable() error {
	var opt mintab.Option
	switch ren.OutputType {
	case OutputTypeText:
		opt = mintab.WithFormat(mintab.TextFormat)
	case OutputTypeMarkdown:
		opt = mintab.WithFormat(mintab.MarkdownFormat)
	case OutputTypeBacklog:
		opt = mintab.WithFormat(mintab.BacklogFormat)
	}
	table := mintab.New(ren.w, opt)
	if err := table.Load(ren.toInput()); err != nil {
		return err
	}
	table.Render()
	return nil
}

func (ren *Renderer) toTSV() error {
	buf := &bytes.Buffer{}
	switch ren.MetricName {
	case MetricNameBucketSizeBytes:
		buf.WriteString("BucketName\tRegion\tStorageType\tBytes\tReadableBytes\n")
		for _, metric := range ren.Metrics {
			m := metric.(*SizeMetric)
			s := []string{
				m.BucketName,
				m.Region,
				m.StorageType.String(),
				strconv.FormatFloat(m.Bytes, 'f', -1, 64),
				m.ReadableBytes,
			}
			buf.WriteString(strings.Join(s, "\t"))
			buf.WriteString("\n")
		}
	case MetricNameNumberOfObjects:
		buf.WriteString("BucketName\tRegion\tStorageType\tObjects\n")
		for _, metric := range ren.Metrics {
			m := metric.(*ObjectMetric)
			s := []string{
				m.BucketName,
				m.Region,
				m.StorageType.String(),
				strconv.FormatFloat(m.Objects, 'f', -1, 64),
			}
			buf.WriteString(strings.Join(s, "\t"))
			buf.WriteString("\n")
		}
	default:
	}
	if _, err := ren.w.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (ren *Renderer) toInput() mintab.Input {
	var (
		header = ([]string)(nil)
		data   = make([][]any, len(ren.Metrics))
	)
	switch ren.MetricName {
	case MetricNameBucketSizeBytes:
		header = []string{
			"BucketName",
			"Region",
			"StorageType",
			"Bytes",
			"ReadableBytes",
		}
		for i, metric := range ren.Metrics {
			m := metric.(*SizeMetric)
			data[i] = []any{
				m.BucketName,
				m.Region,
				m.StorageType,
				m.Bytes,
				m.ReadableBytes,
			}
		}
	case MetricNameNumberOfObjects:
		header = []string{
			"BucketName",
			"Region",
			"StorageType",
			"Objects",
		}
		for i, metric := range ren.Metrics {
			m := metric.(*ObjectMetric)
			data[i] = []any{
				m.BucketName,
				m.Region,
				m.StorageType,
				m.Objects,
			}
		}
	default:
	}
	return mintab.Input{
		Header: header,
		Data:   data,
	}
}
