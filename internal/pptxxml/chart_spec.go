package pptxxml

// ChartSpec describes one chart written as a chart part.
type ChartSpec struct {
	Kind                       string
	Title                      string
	TitleOverlay               bool
	Categories                 []string
	XValues                    []float64
	Values                     []float64
	BubbleSizes                []float64
	OpenValues                 []float64
	HighValues                 []float64
	LowValues                  []float64
	CloseValues                []float64
	X                          int64
	Y                          int64
	CX                         int64
	CY                         int64
	Color                      string
	SeriesName                 string
	BarDir                     string
	Grouping                   string
	ShowMarkers                bool
	ScatterStyle               string
	RadarStyle                 string
	BubbleScale                int
	BarSeries                  []ChartSeries
	LineSeries                 []ChartSeries
	ShowLegend                 bool
	LegendPosition             string
	LegendOverlay              bool
	ShowDataLabels             bool
	DataLabelPosition          string
	DataLabelShowLegendKey     *bool
	DataLabelShowValue         *bool
	DataLabelShowCategoryName  *bool
	DataLabelShowSeriesName    *bool
	DataLabelShowPercent       *bool
	DataLabelShowBubbleSize    *bool
	DataLabelWordWrap          *bool
	HoleSize                   int
	ShowMajorGridlines         bool
	ShowCategoryMajorGridlines bool
	CategoryAxisTitle          string
	ValueAxisTitle             string
	CategoryTickLabelPosition  string
	ValueTickLabelPosition     string
	CategoryAxisCrosses        string
	ValueAxisCrosses           string
	ValueFormat                string
	ValueAxisCrossBetween      string
	MinValue                   *float64
	MaxValue                   *float64
	// DisplayBlanksAs decides how a missing value is drawn: "gap", "zero" or
	// "span". Empty leaves the element out, which PowerPoint reads as "gap".
	DisplayBlanksAs string

	// Secondary value axis. SecondaryAxis declares the axis; a plot only draws
	// against it when its c:axId pair references the secondary ids, which is
	// what puts two series on independent scales.
	SecondaryAxis                   bool
	SecondaryValueAxisTitle         string
	SecondaryValueTickLabelPosition string
	SecondaryValueFormat            string
	ShowSecondaryMajorGridlines     bool
	SecondaryMinValue               *float64
	SecondaryMaxValue               *float64
	Smooth                          bool
	ExternalDataID                  string

	// Data table drawn under the plot area (<c:dTable>). Only chart kinds with
	// a category axis can carry one; the border and legend-key flags default to
	// true when the table is shown and nothing else is stated.
	ShowDataTable           bool
	DataTableShowHorzBorder *bool
	DataTableShowVertBorder *bool
	DataTableShowOutline    *bool
	DataTableShowLegendKeys *bool

	// Accessibility
	AltText      string
	IsDecorative bool
}

type ChartSeries struct {
	Name   string
	Values []float64
	// SecondaryAxis draws this series against the secondary value axis. Only
	// combo line series honour it; a combo with a mix emits one line plot per
	// axis so an unmarked series keeps the primary scale.
	SecondaryAxis bool
}

const (
	ChartKindBar            = "bar"
	ChartKindBarHorizontal  = "barHorizontal"
	ChartKindBarStacked     = "barStacked"
	ChartKindBarStacked100  = "barStacked100"
	ChartKindLine           = "line"
	ChartKindLineMarkers    = "lineMarkers"
	ChartKindLineStacked    = "lineStacked"
	ChartKindScatter        = "scatter"
	ChartKindArea           = "area"
	ChartKindAreaStacked    = "areaStacked"
	ChartKindAreaStacked100 = "areaStacked100"
	ChartKindPie            = "pie"
	ChartKindThreeDPie      = "pie3D"
	ChartKindThreeDColumn   = "column3D"
	ChartKindThreeDBar      = "bar3D"
	ChartKindThreeDLine     = "line3D"
	ChartKindThreeDArea     = "area3D"
	ChartKindDoughnut       = "doughnut"
	ChartKindBubble         = "bubble"
	ChartKindRadar          = "radar"
	ChartKindRadarFilled    = "radarFilled"
	ChartKindStockHLC       = "stockHLC"
	ChartKindStockOHLC      = "stockOHLC"
	ChartKindCombo          = "combo"
)
