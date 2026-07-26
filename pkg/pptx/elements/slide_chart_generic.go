package elements

import (
	"reflect"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

// WithChart sets any chart on the slide, replacing a chart already set. It is
// the type-generic counterpart to the WithBarChart / WithPieChart / ... family
// and dispatches to the same fields, so the two styles are interchangeable.
//
// Both a chart value and a pointer to one are accepted. A nil chart, or a nil
// pointer, leaves the slide unchanged.
func (s SlideContent) WithChart(chart charts.Chart) SlideContent {
	chart, ok := derefChart(chart)
	if !ok {
		return s
	}

	switch c := chart.(type) {
	case charts.BarChart:
		return s.WithBarChart(c)
	case charts.BarHorizontalChart:
		return s.WithBarHorizontalChart(c)
	case charts.BarStackedChart:
		return s.WithBarStackedChart(c)
	case charts.BarStacked100Chart:
		return s.WithBarStacked100Chart(c)
	case charts.LineChart:
		return s.WithLineChart(c)
	case charts.LineMarkersChart:
		return s.WithLineMarkersChart(c)
	case charts.LineStackedChart:
		return s.WithLineStackedChart(c)
	case charts.ScatterChart:
		return s.WithScatterChart(c)
	case charts.AreaChart:
		return s.WithAreaChart(c)
	case charts.AreaStackedChart:
		return s.WithAreaStackedChart(c)
	case charts.AreaStacked100Chart:
		return s.WithAreaStacked100Chart(c)
	case charts.PieChart:
		return s.WithPieChart(c)
	case charts.DoughnutChart:
		return s.WithDoughnutChart(c)
	case charts.BubbleChart:
		return s.WithBubbleChart(c)
	case charts.RadarChart:
		return s.WithRadarChart(c)
	case charts.RadarFilledChart:
		return s.WithRadarFilledChart(c)
	case charts.StockHLCChart:
		return s.WithStockHLCChart(c)
	case charts.StockOHLCChart:
		return s.WithStockOHLCChart(c)
	case charts.ComboChart:
		return s.WithComboChart(c)
	}
	return s
}

// derefChart reduces a pointer-to-chart to the chart value it points at, so the
// dispatch above needs one case per chart type rather than one per type and one
// per pointer to it. It reports false for a nil chart or a nil pointer, and for
// a pointer whose element does not itself satisfy charts.Chart.
func derefChart(chart charts.Chart) (charts.Chart, bool) {
	if chart == nil {
		return nil, false
	}

	value := reflect.ValueOf(chart)
	if value.Kind() != reflect.Pointer {
		return chart, true
	}
	if value.IsNil() {
		return nil, false
	}

	elem, ok := value.Elem().Interface().(charts.Chart)
	if !ok {
		return nil, false
	}
	return elem, true
}
