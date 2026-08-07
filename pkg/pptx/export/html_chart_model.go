package export

import (
	"strconv"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// The HTML export drew nothing for a chart, whichever of the twenty kinds a
// slide carried. These normalise every kind into one shape — a title, category
// labels and named series — so the writer can draw the ones that map onto a
// simple plot and tabulate the rest. Nothing is dropped silently.

// htmlChartKind is how a chart is drawn in HTML.
type htmlChartKind string

const (
	htmlChartBar        htmlChartKind = "bar"
	htmlChartBarStacked htmlChartKind = "bar-stacked"
	htmlChartLine       htmlChartKind = "line"
	htmlChartArea       htmlChartKind = "area"
	htmlChartPie        htmlChartKind = "pie"
	// htmlChartTable is the fallback: a chart whose shape has no simple plot —
	// scatter, bubble, stock, radar — is written out as its data.
	htmlChartTable htmlChartKind = "table"
)

// The series a stock chart is made of, named as both the PDF legend and the
// HTML table label them.
const (
	stockSeriesOpen  = "Open"
	stockSeriesHigh  = "High"
	stockSeriesLow   = "Low"
	stockSeriesClose = "Close"
)

// stockAndComboChartCount is how many charts stockAndComboChartsToHTML can
// return: the two stock kinds and the combo.
const stockAndComboChartCount = 3

type htmlChartSeries struct {
	Name   string
	Values []float64
}

type htmlChart struct {
	Kind       htmlChartKind
	Title      string
	Categories []string
	Series     []htmlChartSeries
}

// hasData reports whether there is anything to draw or tabulate.
func (c htmlChart) hasData() bool {
	for _, series := range c.Series {
		if len(series.Values) > 0 {
			return true
		}
	}
	return false
}

// maxValue is the largest value across the series, used to scale a plot. A
// stacked chart is scaled by its column totals instead.
func (c htmlChart) maxValue() float64 {
	if c.Kind == htmlChartBarStacked {
		return c.maxStackTotal()
	}
	maximum := 0.0
	for _, series := range c.Series {
		for _, value := range series.Values {
			if value > maximum {
				maximum = value
			}
		}
	}
	return maximum
}

func (c htmlChart) maxStackTotal() float64 {
	maximum := 0.0
	for index := range c.categoryCount() {
		total := 0.0
		for _, series := range c.Series {
			if index < len(series.Values) {
				total += series.Values[index]
			}
		}
		if total > maximum {
			maximum = total
		}
	}
	return maximum
}

// categoryCount is the number of points along the category axis, taken from the
// longest series when the labels run short.
func (c htmlChart) categoryCount() int {
	count := len(c.Categories)
	for _, series := range c.Series {
		if len(series.Values) > count {
			count = len(series.Values)
		}
	}
	return count
}

// categoryLabel is the label at index, or the position when the chart states
// fewer labels than points.
func (c htmlChart) categoryLabel(index int) string {
	if index < len(c.Categories) {
		return c.Categories[index]
	}
	return strconv.Itoa(index + 1)
}

// slideHTMLCharts collects every chart on the slide, in the order the fields are
// declared.
func slideHTMLCharts(slide elements.SlideContent) []htmlChart {
	out := make([]htmlChart, 0, 2)
	add := func(chart htmlChart) {
		if chart.hasData() {
			out = append(out, chart)
		}
	}
	single := func(kind htmlChartKind, title, name string, categories []string, values []float64) htmlChart {
		return htmlChart{
			Kind:       kind,
			Title:      title,
			Categories: categories,
			Series:     []htmlChartSeries{{Name: name, Values: values}},
		}
	}

	if c := slide.Chart; c != nil {
		add(single(htmlChartBar, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.BarHorizontal; c != nil {
		add(single(htmlChartBar, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.BarStacked; c != nil {
		add(single(htmlChartBarStacked, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.BarStacked100; c != nil {
		add(single(htmlChartBarStacked, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.Line; c != nil {
		add(single(htmlChartLine, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.LineMarkers; c != nil {
		add(single(htmlChartLine, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.LineStacked; c != nil {
		add(single(htmlChartLine, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.Area; c != nil {
		add(single(htmlChartArea, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.AreaStacked; c != nil {
		add(single(htmlChartArea, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.AreaStacked100; c != nil {
		add(single(htmlChartArea, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.Pie; c != nil {
		add(single(htmlChartPie, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.Pie3D; c != nil {
		add(single(htmlChartPie, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.Doughnut; c != nil {
		add(single(htmlChartPie, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.Radar; c != nil {
		add(single(htmlChartTable, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	if c := slide.RadarFilled; c != nil {
		add(single(htmlChartTable, c.Title, c.SeriesName, c.Categories, c.Values))
	}
	out = append(out, pointChartsToHTML(slide)...)
	out = append(out, stockAndComboChartsToHTML(slide)...)
	return out
}

// pointChartsToHTML tabulates the charts plotted on two value axes, which have
// no category labels to plot against.
func pointChartsToHTML(slide elements.SlideContent) []htmlChart {
	out := make([]htmlChart, 0, 2)
	if c := slide.Scatter; c != nil && len(c.YValues) > 0 {
		out = append(out, htmlChart{
			Kind:       htmlChartTable,
			Title:      c.Title,
			Categories: floatsAsLabels(c.XValues),
			Series:     []htmlChartSeries{{Name: "Y", Values: c.YValues}},
		})
	}
	if c := slide.Bubble; c != nil && len(c.YValues) > 0 {
		series := []htmlChartSeries{{Name: "Y", Values: c.YValues}}
		if len(c.BubbleSizes) > 0 {
			series = append(series, htmlChartSeries{Name: "Size", Values: c.BubbleSizes})
		}
		out = append(out, htmlChart{
			Kind:       htmlChartTable,
			Title:      c.Title,
			Categories: floatsAsLabels(c.XValues),
			Series:     series,
		})
	}
	return out
}

func stockAndComboChartsToHTML(slide elements.SlideContent) []htmlChart {
	out := make([]htmlChart, 0, stockAndComboChartCount)
	if c := slide.StockHLC; c != nil {
		out = append(out, htmlChart{
			Kind:       htmlChartTable,
			Title:      c.Title,
			Categories: c.Categories,
			Series: []htmlChartSeries{
				{Name: stockSeriesHigh, Values: c.HighValues},
				{Name: stockSeriesLow, Values: c.LowValues},
				{Name: stockSeriesClose, Values: c.CloseValues},
			},
		})
	}
	if c := slide.StockOHLC; c != nil {
		out = append(out, htmlChart{
			Kind:       htmlChartTable,
			Title:      c.Title,
			Categories: c.Categories,
			Series: []htmlChartSeries{
				{Name: stockSeriesOpen, Values: c.OpenValues},
				{Name: stockSeriesHigh, Values: c.HighValues},
				{Name: stockSeriesLow, Values: c.LowValues},
				{Name: stockSeriesClose, Values: c.CloseValues},
			},
		})
	}
	if c := slide.Combo; c != nil {
		out = append(out, htmlChart{
			Kind:       htmlChartBar,
			Title:      c.Title,
			Categories: c.Categories,
			Series:     append(namedSeries(c.BarSeries), namedSeries(c.LineSeries)...),
		})
	}
	return out
}

func namedSeries(series []charts.Series) []htmlChartSeries {
	out := make([]htmlChartSeries, 0, len(series))
	for _, s := range series {
		out = append(out, htmlChartSeries{Name: s.Name, Values: s.Values})
	}
	return out
}

func floatsAsLabels(values []float64) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, formatHTMLChartValue(value))
	}
	return out
}

// formatHTMLChartValue prints a value without a trailing ".0" on whole numbers,
// which is how the deck's own labels read.
func formatHTMLChartValue(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}
