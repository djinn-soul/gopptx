package pptx

import (
	"fmt"
)

// SlideContent describes the user-visible content of a slide.
type SlideContent struct {
	Title                string
	TitleSize            int
	TitleColor           string
	TitleBold            bool
	TitleItalic          bool
	TitleUnderline       bool
	ContentSize          int
	ContentColor         string
	ContentBold          bool
	ContentItalic        bool
	ContentUnderline     bool
	Layout               string
	Transition           SlideTransition
	DefaultBulletStyle   TextParagraphStyle
	Bullets              []string
	BulletRuns           [][]TextRun
	BulletStyles         []TextParagraphStyle
	Notes                string
	Images               []Image
	Shapes               []Shape
	Connectors           []Connector
	Table                *Table
	Chart                *BarChart
	BarHorizontal        *BarHorizontalChart
	BarStacked           *BarStackedChart
	BarStacked100        *BarStacked100Chart
	Line                 *LineChart
	LineMarkers          *LineMarkersChart
	LineStacked          *LineStackedChart
	Scatter              *ScatterChart
	Area                 *AreaChart
	AreaStacked          *AreaStackedChart
	AreaStacked100       *AreaStacked100Chart
	Pie                  *PieChart
	Dough                *DoughnutChart
	Bubble               *BubbleChart
	Radar                *RadarChart
	RadarFilled          *RadarFilledChart
	StockHLC             *StockHLCChart
	StockOHLC            *StockOHLCChart
	Combo                *ComboChart
	Animations           []Animation
	PlaceholderOverrides []PlaceholderContent
}

// WithTitleSize sets the font size for the slide title.
func (s SlideContent) WithTitleSize(size int) SlideContent {
	s.TitleSize = size
	return s
}

// WithTitleColor sets the RGB hex color for the slide title.
func (s SlideContent) WithTitleColor(color string) SlideContent {
	s.TitleColor = normalizeHexColor(color)
	return s
}

// WithTitleBold sets whether the slide title is bold.
func (s SlideContent) WithTitleBold(bold bool) SlideContent {
	s.TitleBold = bold
	return s
}

// WithTitleItalic sets whether the slide title is italic.
func (s SlideContent) WithTitleItalic(italic bool) SlideContent {
	s.TitleItalic = italic
	return s
}

// WithTitleUnderline sets whether the slide title is underlined.
func (s SlideContent) WithTitleUnderline(underline bool) SlideContent {
	s.TitleUnderline = underline
	return s
}

// WithContentSize sets the default font size for bullet content.
func (s SlideContent) WithContentSize(size int) SlideContent {
	s.ContentSize = size
	return s
}

// WithContentColor sets the default RGB hex color for bullet content.
func (s SlideContent) WithContentColor(color string) SlideContent {
	s.ContentColor = normalizeHexColor(color)
	return s
}

// WithContentBold sets whether the bullet content is bold by default.
func (s SlideContent) WithContentBold(bold bool) SlideContent {
	s.ContentBold = bold
	return s
}

// WithContentItalic sets whether the bullet content is italic by default.
func (s SlideContent) WithContentItalic(italic bool) SlideContent {
	s.ContentItalic = italic
	return s
}

// WithContentUnderline sets whether the bullet content is underlined by default.
func (s SlideContent) WithContentUnderline(underline bool) SlideContent {
	s.ContentUnderline = underline
	return s
}

// PlaceholderContent overrides content for a layout placeholder.
type PlaceholderContent struct {
	Index int
	Type  string
	Text  string
	Image *Image
	Table *Table
	Chart ChartDefinition
}

// WithPlaceholderChart adds a chart to a specific placeholder index.
func (s SlideContent) WithPlaceholderChart(idx int, chart ChartDefinition) SlideContent {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, PlaceholderContent{
		Index: idx,
		Chart: chart,
	})
	return s
}

// WithPlaceholderImage adds an image to a specific placeholder index.
func (s SlideContent) WithPlaceholderImage(idx int, image Image) SlideContent {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, PlaceholderContent{
		Index: idx,
		Image: &image,
	})
	return s
}

// WithPlaceholderText adds text to a specific placeholder index.
func (s SlideContent) WithPlaceholderText(idx int, text string) SlideContent {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, PlaceholderContent{
		Index: idx,
		Text:  text,
	})
	return s
}

// WithPlaceholderTable adds a table to a specific placeholder index.
func (s SlideContent) WithPlaceholderTable(idx int, table Table) SlideContent {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, PlaceholderContent{
		Index: idx,
		Table: &table,
	})
	return s
}

// AddAnimation appends one animation effect and returns the updated slide.
func (s SlideContent) AddAnimation(def AnimationDefinition) SlideContent {
	s.Animations = append(s.Animations, def.ToAnimation())
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

// Validate checks the slide content for OOXML compliance and common constraints.
func (s SlideContent) Validate(slideIndex int) error {
	if err := validateSlideLayout(s, slideIndex); err != nil {
		return err
	}
	if err := validateSlideTransition(s, slideIndex); err != nil {
		return err
	}
	if err := validateSlideStyle(s, slideIndex); err != nil {
		return err
	}
	for bulletIndex, bullet := range s.Bullets {
		if bullet == "" {
			return fmt.Errorf("slide %d bullet %d cannot be empty", slideIndex, bulletIndex+1)
		}
	}
	if err := validateSlideTextRuns(s, slideIndex); err != nil {
		return err
	}
	if err := validateSlideTextParagraphStyles(s, slideIndex); err != nil {
		return err
	}
	for imageIndex, image := range s.Images {
		if err := image.Validate(slideIndex, imageIndex+1); err != nil {
			return err
		}
	}
	if err := validateSlideDrawings(s, slideIndex); err != nil {
		return err
	}
	if err := validateSlideAnimations(s, slideIndex); err != nil {
		return err
	}
	if s.Table != nil {
		if err := s.Table.Validate(slideIndex); err != nil {
			return err
		}
	}
	if err := validateSlideCharts(s, slideIndex); err != nil {
		return err
	}
	return nil
}
