package pptxxml

// ChartSpec describes a categorical chart rendered as native slide shapes.
type ChartSpec struct {
	Kind       string
	Title      string
	Categories []string
	Values     []float64
	X          int64
	Y          int64
	CX         int64
	CY         int64
	Color      string
}

const (
	ChartKindBar  = "bar"
	ChartKindLine = "line"
)

func chartShape(chart *ChartSpec, baseID int) (string, int) {
	switch chart.Kind {
	case ChartKindLine:
		return lineChartShape(chart, baseID)
	case ChartKindBar:
		fallthrough
	default:
		return barChartShape(chart, baseID)
	}
}

func chartMaxValue(values []float64) float64 {
	if len(values) == 0 {
		return 1
	}

	maxValue := values[0]
	for _, value := range values[1:] {
		if value > maxValue {
			maxValue = value
		}
	}
	if maxValue <= 0 {
		return 1
	}
	return maxValue
}
