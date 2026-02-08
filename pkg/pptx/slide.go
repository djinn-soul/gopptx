package pptx

import (
	"fmt"
)

// SlideContent describes the user-visible content of a slide.
type SlideContent struct {
	Title              string
	Layout             string
	Transition         SlideTransition
	DefaultBulletStyle TextParagraphStyle
	Bullets            []string
	BulletRuns         [][]TextRun
	BulletStyles       []TextParagraphStyle
	Notes              string
	Images             []Image
	Shapes             []Shape
	Connectors         []Connector
	Table              *Table
	Chart              *BarChart
	BarHorizontal      *BarHorizontalChart
	BarStacked         *BarStackedChart
	BarStacked100      *BarStacked100Chart
	Line               *LineChart
	LineMarkers        *LineMarkersChart
	LineStacked        *LineStackedChart
	Scatter            *ScatterChart
	Area               *AreaChart
	AreaStacked        *AreaStackedChart
	AreaStacked100     *AreaStacked100Chart
	Pie                *PieChart
	Dough              *DoughnutChart
	Bubble             *BubbleChart
	Radar              *RadarChart
	RadarFilled        *RadarFilledChart
	StockHLC           *StockHLCChart
	StockOHLC          *StockOHLCChart
	Combo              *ComboChart
	Animations         []Animation
}

// AddAnimation appends one animation effect and returns the updated slide.
func (s SlideContent) AddAnimation(animation Animation) SlideContent {
	s.Animations = append(s.Animations, animation)
	return s
}

// NewSlide creates a new slide with a title.
func NewSlide(title string) SlideContent {
	return SlideContent{
		Title:              title,
		Layout:             SlideLayoutTitleAndContent,
		DefaultBulletStyle: defaultTextParagraphStyle(),
	}
}

// AddBullet appends one bullet item and returns the updated slide.
func (s SlideContent) AddBullet(text string) SlideContent {
	s.Bullets = append(s.Bullets, text)
	s.BulletRuns = append(s.BulletRuns, nil)
	s.BulletStyles = append(s.BulletStyles, s.DefaultBulletStyle)
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
	if err := validateSlideLayout(s, index); err != nil {
		return err
	}
	if err := validateSlideTransition(s, index); err != nil {
		return err
	}
	for bulletIndex, bullet := range s.Bullets {
		if bullet == "" {
			return fmt.Errorf("slide %d bullet %d cannot be empty", index, bulletIndex+1)
		}
	}
	if err := validateSlideTextRuns(s, index); err != nil {
		return err
	}
	if err := validateSlideTextParagraphStyles(s, index); err != nil {
		return err
	}
	for imageIndex, image := range s.Images {
		if err := validateImage(image, index, imageIndex+1); err != nil {
			return err
		}
	}
	if err := validateSlideDrawings(s, index); err != nil {
		return err
	}
	if err := validateSlideAnimations(s, index); err != nil {
		return err
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
