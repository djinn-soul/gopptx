package charts

import "github.com/djinn-soul/gopptx/internal/pptxxml"

// ChartDefinition is the interface for all chart types supported by gopptx.
type ChartDefinition interface {
	ToChartSpec() *pptxxml.ChartSpec
	Validate(slideIndex int) error
	GetCategories() []string
	GetValues() []float64
}
