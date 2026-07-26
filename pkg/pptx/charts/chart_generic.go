package charts

import "github.com/djinn-soul/gopptx/internal/pptxxml"

// Chart is satisfied by every chart type in this package. The interface is
// sealed via an unexported method, so it can be used as a type-safe "any chart"
// parameter without accepting arbitrary values.
//
// Chart carries no fluent option methods: those return their concrete type
// (BarChart, PieChart, ...) so that chaining stays type-checked. Configure a
// chart through its concrete type first, then pass it where a Chart is wanted.
type Chart interface {
	// ChartKind reports the pptxxml chart kind this chart renders as.
	ChartKind() string

	sealedChart()
}

func (BarChart) ChartKind() string { return pptxxml.ChartKindBar }
func (BarChart) sealedChart()      {}

func (BarHorizontalChart) ChartKind() string { return pptxxml.ChartKindBarHorizontal }
func (BarHorizontalChart) sealedChart()      {}

func (BarStackedChart) ChartKind() string { return pptxxml.ChartKindBarStacked }
func (BarStackedChart) sealedChart()      {}

func (BarStacked100Chart) ChartKind() string { return pptxxml.ChartKindBarStacked100 }
func (BarStacked100Chart) sealedChart()      {}

func (LineChart) ChartKind() string { return pptxxml.ChartKindLine }
func (LineChart) sealedChart()      {}

func (LineMarkersChart) ChartKind() string { return pptxxml.ChartKindLineMarkers }
func (LineMarkersChart) sealedChart()      {}

func (LineStackedChart) ChartKind() string { return pptxxml.ChartKindLineStacked }
func (LineStackedChart) sealedChart()      {}

func (ScatterChart) ChartKind() string { return pptxxml.ChartKindScatter }
func (ScatterChart) sealedChart()      {}

func (AreaChart) ChartKind() string { return pptxxml.ChartKindArea }
func (AreaChart) sealedChart()      {}

func (AreaStackedChart) ChartKind() string { return pptxxml.ChartKindAreaStacked }
func (AreaStackedChart) sealedChart()      {}

func (AreaStacked100Chart) ChartKind() string { return pptxxml.ChartKindAreaStacked100 }
func (AreaStacked100Chart) sealedChart()      {}

func (PieChart) ChartKind() string { return pptxxml.ChartKindPie }
func (PieChart) sealedChart()      {}

func (DoughnutChart) ChartKind() string { return pptxxml.ChartKindDoughnut }
func (DoughnutChart) sealedChart()      {}

func (BubbleChart) ChartKind() string { return pptxxml.ChartKindBubble }
func (BubbleChart) sealedChart()      {}

func (RadarChart) ChartKind() string { return pptxxml.ChartKindRadar }
func (RadarChart) sealedChart()      {}

func (RadarFilledChart) ChartKind() string { return pptxxml.ChartKindRadarFilled }
func (RadarFilledChart) sealedChart()      {}

func (StockHLCChart) ChartKind() string { return pptxxml.ChartKindStockHLC }
func (StockHLCChart) sealedChart()      {}

func (StockOHLCChart) ChartKind() string { return pptxxml.ChartKindStockOHLC }
func (StockOHLCChart) sealedChart()      {}

func (ComboChart) ChartKind() string { return pptxxml.ChartKindCombo }
func (ComboChart) sealedChart()      {}
