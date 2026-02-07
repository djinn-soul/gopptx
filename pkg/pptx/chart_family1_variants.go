package pptx

// BarHorizontalChart is a horizontal clustered bar chart variant.
type BarHorizontalChart struct {
	BarChart
}

func NewBarHorizontalChart(categories []string, values []float64) BarHorizontalChart {
	return BarHorizontalChart{BarChart: NewBarChart(categories, values)}
}

func validateBarHorizontalChart(chart BarHorizontalChart, slideIndex int) error {
	return validateBarChart(chart.BarChart, slideIndex)
}

// BarStackedChart is a stacked bar chart variant.
type BarStackedChart struct {
	BarChart
}

func NewBarStackedChart(categories []string, values []float64) BarStackedChart {
	return BarStackedChart{BarChart: NewBarChart(categories, values)}
}

func validateBarStackedChart(chart BarStackedChart, slideIndex int) error {
	return validateBarChart(chart.BarChart, slideIndex)
}

// BarStacked100Chart is a 100%% stacked bar chart variant.
type BarStacked100Chart struct {
	BarChart
}

func NewBarStacked100Chart(categories []string, values []float64) BarStacked100Chart {
	return BarStacked100Chart{BarChart: NewBarChart(categories, values)}
}

func validateBarStacked100Chart(chart BarStacked100Chart, slideIndex int) error {
	return validateBarChart(chart.BarChart, slideIndex)
}

// LineMarkersChart is a line-with-markers chart variant.
type LineMarkersChart struct {
	LineChart
}

func NewLineMarkersChart(categories []string, values []float64) LineMarkersChart {
	return LineMarkersChart{LineChart: NewLineChart(categories, values)}
}

func validateLineMarkersChart(chart LineMarkersChart, slideIndex int) error {
	return validateLineChart(chart.LineChart, slideIndex)
}

// LineStackedChart is a stacked line chart variant.
type LineStackedChart struct {
	LineChart
}

func NewLineStackedChart(categories []string, values []float64) LineStackedChart {
	return LineStackedChart{LineChart: NewLineChart(categories, values)}
}

func validateLineStackedChart(chart LineStackedChart, slideIndex int) error {
	return validateLineChart(chart.LineChart, slideIndex)
}

// AreaStackedChart is a stacked area chart variant.
type AreaStackedChart struct {
	AreaChart
}

func NewAreaStackedChart(categories []string, values []float64) AreaStackedChart {
	return AreaStackedChart{AreaChart: NewAreaChart(categories, values)}
}

func validateAreaStackedChart(chart AreaStackedChart, slideIndex int) error {
	return validateAreaChart(chart.AreaChart, slideIndex)
}

// AreaStacked100Chart is a 100%% stacked area chart variant.
type AreaStacked100Chart struct {
	AreaChart
}

func NewAreaStacked100Chart(categories []string, values []float64) AreaStacked100Chart {
	return AreaStacked100Chart{AreaChart: NewAreaChart(categories, values)}
}

func validateAreaStacked100Chart(chart AreaStacked100Chart, slideIndex int) error {
	return validateAreaChart(chart.AreaChart, slideIndex)
}
