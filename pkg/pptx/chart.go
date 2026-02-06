package pptx

import (
	"fmt"
	"regexp"
	"strings"
)

var hexColorPattern = regexp.MustCompile(`^[0-9A-F]{6}$`)

// BarChart is a simple categorical bar chart.
type BarChart struct {
	Title      string
	Categories []string
	Values     []float64
	X          int64
	Y          int64
	CX         int64
	CY         int64
	BarColor   string
}

// NewBarChart creates a bar chart with default layout and style.
func NewBarChart(categories []string, values []float64) BarChart {
	cats := make([]string, len(categories))
	copy(cats, categories)
	vals := make([]float64, len(values))
	copy(vals, values)

	return BarChart{
		Title:      "Chart",
		Categories: cats,
		Values:     vals,
		X:          685800,
		Y:          1800000,
		CX:         7772400,
		CY:         4114800,
		BarColor:   "4F81BD",
	}
}

// Position sets chart position in EMU.
func (c BarChart) Position(x int64, y int64) BarChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c BarChart) Size(cx int64, cy int64) BarChart {
	c.CX = cx
	c.CY = cy
	return c
}

// WithTitle sets the chart title.
func (c BarChart) WithTitle(title string) BarChart {
	c.Title = title
	return c
}

// WithBarColor sets the bar fill color using RGB hex.
func (c BarChart) WithBarColor(color string) BarChart {
	c.BarColor = normalizeHexColor(color)
	return c
}

func validateBarChart(chart BarChart, slideIndex int) error {
	if chart.X < 0 || chart.Y < 0 {
		return fmt.Errorf("slide %d chart position cannot be negative", slideIndex)
	}
	if chart.CX <= 0 || chart.CY <= 0 {
		return fmt.Errorf("slide %d chart size must be > 0", slideIndex)
	}
	if strings.TrimSpace(chart.Title) == "" {
		return fmt.Errorf("slide %d chart title cannot be empty", slideIndex)
	}
	if len(chart.Categories) == 0 {
		return fmt.Errorf("slide %d chart must define at least one category", slideIndex)
	}
	if len(chart.Categories) != len(chart.Values) {
		return fmt.Errorf(
			"slide %d chart category/value length mismatch (%d vs %d)",
			slideIndex,
			len(chart.Categories),
			len(chart.Values),
		)
	}

	hasPositive := false
	for i := range chart.Categories {
		if strings.TrimSpace(chart.Categories[i]) == "" {
			return fmt.Errorf("slide %d chart category %d cannot be empty", slideIndex, i+1)
		}
		if chart.Values[i] < 0 {
			return fmt.Errorf("slide %d chart value %d cannot be negative", slideIndex, i+1)
		}
		if chart.Values[i] > 0 {
			hasPositive = true
		}
	}
	if !hasPositive {
		return fmt.Errorf("slide %d chart requires at least one positive value", slideIndex)
	}
	if !hexColorPattern.MatchString(normalizeHexColor(chart.BarColor)) {
		return fmt.Errorf("slide %d chart color must be 6-digit RGB hex", slideIndex)
	}
	return nil
}

func normalizeHexColor(color string) string {
	clean := strings.TrimPrefix(strings.TrimSpace(color), "#")
	return strings.ToUpper(clean)
}
