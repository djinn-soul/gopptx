package pptxxml

// ChartSpec describes one chart written as a chart part.
type ChartSpec struct {
	Kind               string
	Title              string
	Categories         []string
	XValues            []float64
	Values             []float64
	BubbleSizes        []float64
	OpenValues         []float64
	HighValues         []float64
	LowValues          []float64
	CloseValues        []float64
	X                  int64
	Y                  int64
	CX                 int64
	CY                 int64
	Color              string
	SeriesName         string
	BarDir             string
	Grouping           string
	ShowMarkers        bool
	ScatterStyle       string
	RadarStyle         string
	BubbleScale        int
	BarSeries          []ChartSeries
	LineSeries         []ChartSeries
	ShowLegend         bool
	LegendPosition     string
	ShowDataLabels     bool
	HoleSize           int
	ShowMajorGridlines bool
	CategoryAxisTitle  string
	ValueAxisTitle     string
	ValueFormat        string
	MinValue           *float64
	MaxValue           *float64
	Smooth             bool
}

type ChartSeries struct {
	Name   string
	Values []float64
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
	ChartKindDoughnut       = "doughnut"
	ChartKindBubble         = "bubble"
	ChartKindRadar          = "radar"
	ChartKindRadarFilled    = "radarFilled"
	ChartKindStockHLC       = "stockHLC"
	ChartKindStockOHLC      = "stockOHLC"
	ChartKindCombo          = "combo"
)
