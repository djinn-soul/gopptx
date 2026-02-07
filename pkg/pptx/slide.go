package pptx

import (
	"fmt"
)

// SlideContent describes the user-visible content of a slide.
type SlideContent struct {
	Title          string
	Bullets        []string
	Images         []Image
	Table          *Table
	Chart          *BarChart
	BarHorizontal  *BarHorizontalChart
	BarStacked     *BarStackedChart
	BarStacked100  *BarStacked100Chart
	Line           *LineChart
	LineMarkers    *LineMarkersChart
	LineStacked    *LineStackedChart
	Scatter        *ScatterChart
	Area           *AreaChart
	AreaStacked    *AreaStackedChart
	AreaStacked100 *AreaStacked100Chart
	Pie            *PieChart
	Dough          *DoughnutChart
	Bubble         *BubbleChart
	Radar          *RadarChart
	RadarFilled    *RadarFilledChart
	StockHLC       *StockHLCChart
	StockOHLC      *StockOHLCChart
	Combo          *ComboChart
}

// NewSlide creates a new slide with a title.
func NewSlide(title string) SlideContent {
	return SlideContent{Title: title}
}

// AddBullet appends one bullet item and returns the updated slide.
func (s SlideContent) AddBullet(text string) SlideContent {
	s.Bullets = append(s.Bullets, text)
	return s
}

// AddImage appends one image and returns the updated slide.
func (s SlideContent) AddImage(image Image) SlideContent {
	s.Images = append(s.Images, image)
	return s
}

// WithTable sets one table for the slide.
func (s SlideContent) WithTable(table Table) SlideContent {
	s.Table = &table
	return s
}

func validateSlide(s SlideContent, index int) error {
	if s.Title == "" {
		return fmt.Errorf("slide %d title cannot be empty", index)
	}
	for bulletIndex, bullet := range s.Bullets {
		if bullet == "" {
			return fmt.Errorf("slide %d bullet %d cannot be empty", index, bulletIndex+1)
		}
	}
	for imageIndex, image := range s.Images {
		if err := validateImage(image, index, imageIndex+1); err != nil {
			return err
		}
	}
	if s.Table != nil {
		if err := validateTable(*s.Table, index); err != nil {
			return err
		}
	}
	if err := validateSlideCharts(s, index); err != nil {
		return err
	}
	if chartKindCount(s) > 1 {
		return fmt.Errorf("slide %d cannot have more than one chart kind", index)
	}
	return nil
}
