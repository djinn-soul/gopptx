package pptx

import (
	"fmt"
)

// SlideContent describes the user-visible content of a slide.
type SlideContent struct {
	Title   string
	Bullets []string
	Images  []Image
	Table   *Table
	Chart   *BarChart
	Line    *LineChart
	Pie     *PieChart
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

// WithBarChart sets one bar chart for the slide.
func (s SlideContent) WithBarChart(chart BarChart) SlideContent {
	s.Chart = &chart
	s.Line = nil
	s.Pie = nil
	return s
}

// WithLineChart sets one line chart for the slide.
func (s SlideContent) WithLineChart(chart LineChart) SlideContent {
	s.Line = &chart
	s.Chart = nil
	s.Pie = nil
	return s
}

// WithPieChart sets one pie chart for the slide.
func (s SlideContent) WithPieChart(chart PieChart) SlideContent {
	s.Pie = &chart
	s.Chart = nil
	s.Line = nil
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
	if s.Chart != nil {
		if err := validateBarChart(*s.Chart, index); err != nil {
			return err
		}
	}
	if s.Line != nil {
		if err := validateLineChart(*s.Line, index); err != nil {
			return err
		}
	}
	if s.Pie != nil {
		if err := validatePieChart(*s.Pie, index); err != nil {
			return err
		}
	}
	chartKinds := 0
	if s.Chart != nil {
		chartKinds++
	}
	if s.Line != nil {
		chartKinds++
	}
	if s.Pie != nil {
		chartKinds++
	}
	if chartKinds > 1 {
		return fmt.Errorf("slide %d cannot have more than one chart kind", index)
	}
	return nil
}
