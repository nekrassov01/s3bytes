package s3bytes

import (
	"encoding/csv"
	"encoding/json"
	"io"

	"github.com/nekrassov01/mintab"
)

// Renderer is a renderer struct for the s3bytes package.
// OutputType represents the type of the output.
type Renderer struct {
	Data       *MetricData
	OutputType OutputType
	w          io.Writer
}

// NewRenderer creates a new renderer with the specified parameters.
func NewRenderer(w io.Writer, data *MetricData, outputType OutputType) *Renderer {
	return &Renderer{
		Data:       data,
		OutputType: outputType,
		w:          w,
	}
}

// String returns the string representation of the renderer.
func (ren *Renderer) String() string {
	b, _ := json.MarshalIndent(ren, "", "  ")
	return string(b)
}

// Render renders the output.
func (ren *Renderer) Render() error {
	switch ren.OutputType {
	case OutputTypeJSON, OutputTypePrettyJSON:
		return ren.toJSON()
	case OutputTypeText, OutputTypeCompressedText, OutputTypeMarkdown, OutputTypeBacklog:
		return ren.toTable()
	case OutputTypeTSV:
		return ren.toTSV()
	case OutputTypeChart:
		return ren.toChart()
	default:
		return nil
	}
}

func (ren *Renderer) toJSON() error {
	b := json.NewEncoder(ren.w)
	if ren.OutputType == OutputTypePrettyJSON {
		b.SetIndent("", "  ")
	}
	return b.Encode(ren.Data.Metrics)
}

func (ren *Renderer) toTable() error {
	var opt mintab.Option
	switch ren.OutputType {
	case OutputTypeText:
		opt = mintab.WithFormat(mintab.TextFormat)
	case OutputTypeCompressedText:
		opt = mintab.WithFormat(mintab.CompressedTextFormat)
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

func (ren *Renderer) toInput() mintab.Input {
	data := make([][]any, len(ren.Data.Metrics))
	for i, row := range ren.Data.Metrics {
		data[i] = row.toInput()
	}
	return mintab.Input{
		Header: ren.Data.Header,
		Data:   data,
	}
}

func (ren *Renderer) toTSV() error {
	w := csv.NewWriter(ren.w)
	w.Comma = '\t'
	if err := w.Write(ren.Data.Header); err != nil {
		return err
	}
	for _, metric := range ren.Data.Metrics {
		if err := w.Write(metric.toTSV()); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func (ren *Renderer) toChart() error {
	title, items := getPieItems(ren.Data)
	pie := newPie(title, items)
	return render(pie)
}
