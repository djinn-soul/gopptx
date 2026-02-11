package elements

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

// Validate checks the slide for consistency.
func (s SlideContent) Validate(index int) error {
	for shapeIndex, shape := range s.Shapes {
		if err := shape.Validate(index, shapeIndex+1); err != nil {
			return err
		}
	}
	for connectorIndex, connector := range s.Connectors {
		if err := connector.Validate(len(s.Shapes), index, connectorIndex+1); err != nil {
			return err
		}
	}
	for imageIndex, image := range s.Images {
		if err := image.Validate(index, imageIndex+1); err != nil {
			return err
		}
	}

	// Validate charts
	// Validate charts
	var slideCharts []ChartDefinition
	if s.Chart != nil {
		slideCharts = append(slideCharts, s.Chart)
	}
	if s.BarHorizontal != nil {
		slideCharts = append(slideCharts, s.BarHorizontal)
	}
	if s.BarStacked != nil {
		slideCharts = append(slideCharts, s.BarStacked)
	}
	if s.BarStacked100 != nil {
		slideCharts = append(slideCharts, s.BarStacked100)
	}
	if s.Line != nil {
		slideCharts = append(slideCharts, s.Line)
	}
	if s.LineMarkers != nil {
		slideCharts = append(slideCharts, s.LineMarkers)
	}
	if s.LineStacked != nil {
		slideCharts = append(slideCharts, s.LineStacked)
	}
	if s.Scatter != nil {
		slideCharts = append(slideCharts, s.Scatter)
	}
	if s.Area != nil {
		slideCharts = append(slideCharts, s.Area)
	}
	if s.AreaStacked != nil {
		slideCharts = append(slideCharts, s.AreaStacked)
	}
	if s.AreaStacked100 != nil {
		slideCharts = append(slideCharts, s.AreaStacked100)
	}
	if s.Pie != nil {
		slideCharts = append(slideCharts, s.Pie)
	}
	if s.Doughnut != nil {
		slideCharts = append(slideCharts, s.Doughnut)
	}
	if s.Bubble != nil {
		slideCharts = append(slideCharts, s.Bubble)
	}
	if s.Radar != nil {
		slideCharts = append(slideCharts, s.Radar)
	}
	if s.RadarFilled != nil {
		slideCharts = append(slideCharts, s.RadarFilled)
	}
	if s.StockHLC != nil {
		slideCharts = append(slideCharts, s.StockHLC)
	}
	if s.StockOHLC != nil {
		slideCharts = append(slideCharts, s.StockOHLC)
	}
	if s.Combo != nil {
		slideCharts = append(slideCharts, s.Combo)
	}

	for _, c := range slideCharts {
		if err := c.Validate(index); err != nil {
			return err
		}
	}

	// Validate text styles and runs
	if err := s.DefaultBulletStyle.Validate(); err != nil {
		return err
	}
	for _, style := range s.BulletStyles {
		if err := style.Validate(); err != nil {
			return err
		}
	}
	for _, runs := range s.BulletRuns {
		for _, run := range runs {
			if err := run.Validate(); err != nil {
				return err
			}
		}
	}

	if (s.TitleSize != 0 && s.TitleSize < 1) || s.TitleSize > 400 {
		return fmt.Errorf("title size must be between 1 and 400 pt (or 0 for default)")
	}
	if s.TitleColor != "" && !IsHexColor(s.TitleColor) {
		return fmt.Errorf("title color must be 6-digit RGB hex")
	}
	if (s.ContentSize != 0 && s.ContentSize < 1) || s.ContentSize > 400 {
		return fmt.Errorf("content size must be between 1 and 400 pt (or 0 for default)")
	}
	if s.ContentColor != "" && !IsHexColor(s.ContentColor) {
		return fmt.Errorf("content color must be 6-digit RGB hex")
	}

	// Validate transition
	if err := ValidateSlideTransition(s, index); err != nil {
		return err
	}

	// Validate table
	if s.Table != nil {
		if err := s.Table.Validate(index); err != nil {
			return err
		}
	}

	// Validate bullets
	for _, b := range s.Bullets {
		if b == "" {
			return fmt.Errorf("bullet cannot be empty")
		}
	}

	if s.Title == "" && s.Layout != SlideLayoutBlank {
		return fmt.Errorf("title cannot be empty")
	}

	for i, anim := range s.Animations {
		if err := anim.Validate(); err != nil {
			return err
		}
		// First animation in a sequence (on a slide) cannot be WithPrevious/AfterPrevious
		if i == 0 && (anim.Trigger == AnimationWithPrevious || anim.Trigger == AnimationAfterPrevious) {
			return fmt.Errorf("first animation trigger cannot be with/after previous")
		}
	}

	return nil
}

const (
	SlideLayoutTitleAndContent    = "title_and_content"
	SlideLayoutTitleOnly          = "title_only"
	SlideLayoutBlank              = "blank"
	SlideLayoutCenteredTitle      = "centered_title"
	SlideLayoutTitleAndBigContent = "title_and_big_content"
	SlideLayoutTwoColumn          = "two_column"

	// Legacy or descriptive aliases
	SlideLayoutTitle          = "Title Slide"
	SlideLayoutSectionHeader  = "Section Header"
	SlideLayoutTwoContent     = "Two Content"
	SlideLayoutComparison     = "Comparison"
	SlideLayoutContentCaption = "Content with Caption"
	SlideLayoutPictureCaption = "Picture with Caption"
)

// ChartDefinition describes the public interface for all chart types.
type ChartDefinition = charts.ChartDefinition

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
	Table                *tables.Table
	Chart                *charts.BarChart
	BarHorizontal        *charts.BarHorizontalChart
	BarStacked           *charts.BarStackedChart
	BarStacked100        *charts.BarStacked100Chart
	Line                 *charts.LineChart
	LineMarkers          *charts.LineMarkersChart
	LineStacked          *charts.LineStackedChart
	Scatter              *charts.ScatterChart
	Area                 *charts.AreaChart
	AreaStacked          *charts.AreaStackedChart
	AreaStacked100       *charts.AreaStacked100Chart
	Pie                  *charts.PieChart
	Doughnut             *charts.DoughnutChart
	Bubble               *charts.BubbleChart
	Radar                *charts.RadarChart
	RadarFilled          *charts.RadarFilledChart
	StockHLC             *charts.StockHLCChart
	StockOHLC            *charts.StockOHLCChart
	Combo                *charts.ComboChart
	Animations           []Animation
	PlaceholderOverrides []PlaceholderContent
}

// PlaceholderContent describes overridden content for a slide layout placeholder.
type PlaceholderContent struct {
	Index int
	Type  string
	Text  string
	Image *Image
	Table *tables.Table
	Chart charts.ChartDefinition
}

// AddBullet appends one bullet item and returns the updated slide.
func (s SlideContent) AddBullet(text string) SlideContent {
	s.Bullets = append(s.Bullets, text)
	s.BulletRuns = append(s.BulletRuns, nil)
	s.BulletStyles = append(s.BulletStyles, s.DefaultBulletStyle)
	return s
}

// AddBulletWithStyle appends one bullet item with explicit styling.
func (s SlideContent) AddBulletWithStyle(text string, style TextParagraphStyle) SlideContent {
	s.Bullets = append(s.Bullets, text)
	s.BulletRuns = append(s.BulletRuns, nil)
	s.BulletStyles = append(s.BulletStyles, style)
	return s
}

// AddBulletRuns appends one bullet item with rich text runs.
func (s SlideContent) AddBulletRuns(runs []TextRun) SlideContent {
	s.Bullets = append(s.Bullets, RunsToPlainText(runs))
	s.BulletRuns = append(s.BulletRuns, runs)
	s.BulletStyles = append(s.BulletStyles, s.DefaultBulletStyle)
	return s
}

// AddBulletRunsWithStyle appends one bullet item with rich text runs and paragraph styling.
func (s SlideContent) AddBulletRunsWithStyle(runs []TextRun, style TextParagraphStyle) SlideContent {
	s.Bullets = append(s.Bullets, RunsToPlainText(runs))
	s.BulletRuns = append(s.BulletRuns, runs)
	s.BulletStyles = append(s.BulletStyles, style)
	return s
}

// AddShape appends one shape and returns the updated slide content.
func (s SlideContent) AddShape(sd ShapeDefinition) SlideContent {
	s.Shapes = append(s.Shapes, sd.ToShape())
	return s
}

// WithDefaultBulletStyle sets the base style for new bullets.
func (s SlideContent) WithDefaultBulletStyle(style TextParagraphStyle) SlideContent {
	s.DefaultBulletStyle = style
	return s
}

// NewSlide creates a new slide with default settings and a title.
func NewSlide(title string) SlideContent {
	return SlideContent{
		Title:       title,
		TitleSize:   44,
		ContentSize: 18,
		Layout:      SlideLayoutTitleAndContent,
	}
}

// WithNotes sets the speaker notes for the slide.
func (s SlideContent) WithNotes(notes string) SlideContent {
	s.Notes = notes
	return s
}

// WithTable sets the table for the slide.
func (s SlideContent) WithTable(t tables.Table) SlideContent {
	s.Table = &t
	return s
}

// WithBulletStyle sets the bullet style for all bullets on this slide.
func (s SlideContent) WithBulletStyle(style TextParagraphStyle) SlideContent {
	s.DefaultBulletStyle = style
	for i := range s.BulletStyles {
		s.BulletStyles[i] = style
	}
	return s
}

// AddNumbered appends one numbered bullet item.
func (s SlideContent) AddNumbered(text string) SlideContent {
	return s.AddBulletWithStyle(text, DefaultTextParagraphStyle().WithNumbered())
}

// AddLettered appends one lettered bullet item.
func (s SlideContent) AddLettered(text string) SlideContent {
	return s.AddBulletWithStyle(text, DefaultTextParagraphStyle().WithLetteredLower())
}

// Adding bullet at index 1..8
func (s SlideContent) AddSubBullet(level int, text string) SlideContent {
	s.Bullets = append(s.Bullets, text)
	s.BulletRuns = append(s.BulletRuns, nil)
	style := s.DefaultBulletStyle
	style.Level = level
	s.BulletStyles = append(s.BulletStyles, style)
	return s
}

// WithTransition sets the transition for the slide.
func (s SlideContent) WithTransition(t SlideTransition) SlideContent {
	s.Transition = t
	return s
}

// WithTransitionOptions sets built-in transition options.
func (s SlideContent) WithTransitionOptions(opt TransitionOptions) SlideContent {
	s.Transition = opt
	return s
}

// WithBulletStyleName sets primary bullet style by name (e.g. BulletStyleNumber).
func (s SlideContent) WithBulletStyleName(styleName string) SlideContent {
	style := s.DefaultBulletStyle
	style.BulletStyle = NormalizeBulletStyle(styleName)
	return s.WithBulletStyle(style)
}

// WithSubtitleLayout sets the layout to Title and Content. (Convenience)
func (s SlideContent) WithLayout(layout string) SlideContent {
	s.Layout = NormalizeSlideLayout(layout)
	return s
}

// WithTitleSize sets the title font size in points.
func (s SlideContent) WithTitleSize(size int) SlideContent {
	s.TitleSize = size
	return s
}

// WithTitleColor sets the title color as RGB hex.
func (s SlideContent) WithTitleColor(color string) SlideContent {
	s.TitleColor = NormalizeHexColor(color)
	return s
}

// WithTitleBold sets whether the title is bold.
func (s SlideContent) WithTitleBold(bold bool) SlideContent {
	s.TitleBold = bold
	return s
}

// WithTitleItalic sets whether the title is italic.
func (s SlideContent) WithTitleItalic(italic bool) SlideContent {
	s.TitleItalic = italic
	return s
}

// WithTitleUnderline sets whether the title is underlined.
func (s SlideContent) WithTitleUnderline(underline bool) SlideContent {
	s.TitleUnderline = underline
	return s
}

// WithContentSize sets the content font size in points.
func (s SlideContent) WithContentSize(size int) SlideContent {
	s.ContentSize = size
	return s
}

// WithContentColor sets the content color as RGB hex.
func (s SlideContent) WithContentColor(color string) SlideContent {
	s.ContentColor = NormalizeHexColor(color)
	return s
}

// WithContentBold sets whether the content is bold.
func (s SlideContent) WithContentBold(bold bool) SlideContent {
	s.ContentBold = bold
	return s
}

// WithContentItalic sets whether the content is italic.
func (s SlideContent) WithContentItalic(italic bool) SlideContent {
	s.ContentItalic = italic
	return s
}

// WithContentUnderline sets whether the content is underlined.
func (s SlideContent) WithContentUnderline(underline bool) SlideContent {
	s.ContentUnderline = underline
	return s
}

// AddImage adds an image to the slide.
func (s SlideContent) AddImage(img Image) SlideContent {
	s.Images = append(s.Images, img)
	return s
}

// AddAnimation adds an animation to the slide.
func (s SlideContent) AddAnimation(anim AnimationDefinition) SlideContent {
	s.Animations = append(s.Animations, anim.ToAnimation())
	return s
}

// WithTitleOnlyLayout sets the layout to title_only.
func (s SlideContent) WithTitleOnlyLayout() SlideContent {
	s.Layout = SlideLayoutTitleOnly
	return s
}

// WithBlankLayout sets the layout to blank.
func (s SlideContent) WithBlankLayout() SlideContent {
	s.Layout = SlideLayoutBlank
	return s
}

// WithCenteredTitleLayout sets the layout to centered_title.
func (s SlideContent) WithCenteredTitleLayout() SlideContent {
	s.Layout = SlideLayoutCenteredTitle
	return s
}

// WithTitleAndBigContentLayout sets the layout to title_and_big_content.
func (s SlideContent) WithTitleAndBigContentLayout() SlideContent {
	s.Layout = SlideLayoutTitleAndBigContent
	return s
}

// WithTwoColumnLayout sets the layout to two_column.
func (s SlideContent) WithTwoColumnLayout() SlideContent {
	return s.WithLayout(SlideLayoutTwoColumn)
}

// WithTitleAndContentLayout sets the layout to title_and_content.
func (s SlideContent) WithTitleAndContentLayout() SlideContent {
	return s.WithLayout(SlideLayoutTitleAndContent)
}

func NormalizeSlideLayout(layout string) string {
	normalized := strings.ToLower(strings.TrimSpace(layout))
	switch normalized {
	case "", SlideLayoutTitleAndContent:
		return SlideLayoutTitleAndContent
	case "titleandcontent", "title-and-content":
		return SlideLayoutTitleAndContent
	case SlideLayoutTitleOnly:
		return SlideLayoutTitleOnly
	case "titleonly", "title-only":
		return SlideLayoutTitleOnly
	case SlideLayoutBlank:
		return SlideLayoutBlank
	case SlideLayoutCenteredTitle:
		return SlideLayoutCenteredTitle
	case "centeredtitle", "centered-title":
		return SlideLayoutCenteredTitle
	case SlideLayoutTitleAndBigContent:
		return SlideLayoutTitleAndBigContent
	case "titleandbigcontent", "title-and-big-content", "big_content":
		return SlideLayoutTitleAndBigContent
	case SlideLayoutTwoColumn:
		return SlideLayoutTwoColumn
	case "twocolumn", "two-column":
		return SlideLayoutTwoColumn
	default:
		return normalized
	}
}

func SlideLayoutXMLMode(layout string) string {
	switch NormalizeSlideLayout(layout) {
	case SlideLayoutTitleOnly:
		return "titleOnly"
	case SlideLayoutBlank:
		return "blank"
	case SlideLayoutCenteredTitle:
		return "centeredTitle"
	case SlideLayoutTitleAndBigContent:
		return "titleAndBigContent"
	case SlideLayoutTwoColumn:
		return "twoColumn"
	default:
		return "titleAndContent"
	}
}

func SlideLayoutTarget(layout string) string {
	switch NormalizeSlideLayout(layout) {
	case SlideLayoutTitleOnly:
		return "../slideLayouts/slideLayout2.xml"
	case SlideLayoutBlank:
		return "../slideLayouts/slideLayout3.xml"
	case SlideLayoutCenteredTitle:
		return "../slideLayouts/slideLayout4.xml"
	case SlideLayoutTitleAndBigContent:
		return "../slideLayouts/slideLayout5.xml"
	case SlideLayoutTwoColumn:
		return "../slideLayouts/slideLayout6.xml"
	default:
		return "../slideLayouts/slideLayout1.xml"
	}
}

// WithPlaceholderText overrides a placeholder with text.
// Supported call forms:
//   - WithPlaceholderText(index, text)
//   - WithPlaceholderText(index, placeholderType, text)
func (s SlideContent) WithPlaceholderText(index int, args ...any) SlideContent {
	phType, text := parsePlaceholderTextArgs(args...)
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, PlaceholderContent{
		Index: index,
		Type:  phType,
		Text:  text,
	})
	return s
}

// WithPlaceholderImage overrides a placeholder with an image.
// Supported call forms:
//   - WithPlaceholderImage(index, image)
//   - WithPlaceholderImage(index, placeholderType, image)
func (s SlideContent) WithPlaceholderImage(index int, args ...any) SlideContent {
	phType, img := parsePlaceholderImageArgs(args...)
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, PlaceholderContent{
		Index: index,
		Type:  phType,
		Image: &img,
	})
	return s
}

// WithPlaceholderTable overrides a placeholder with a table.
// Supported call forms:
//   - WithPlaceholderTable(index, table)
//   - WithPlaceholderTable(index, placeholderType, table)
func (s SlideContent) WithPlaceholderTable(index int, args ...any) SlideContent {
	phType, table := parsePlaceholderTableArgs(args...)
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, PlaceholderContent{
		Index: index,
		Type:  phType,
		Table: &table,
	})
	return s
}

// WithPlaceholderChart overrides a placeholder with a chart.
// Supported call forms:
//   - WithPlaceholderChart(index, chart)
//   - WithPlaceholderChart(index, placeholderType, chart)
func (s SlideContent) WithPlaceholderChart(index int, args ...any) SlideContent {
	phType, chart := parsePlaceholderChartArgs(args...)
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, PlaceholderContent{
		Index: index,
		Type:  phType,
		Chart: chart,
	})
	return s
}

func parsePlaceholderTextArgs(args ...any) (string, string) {
	switch len(args) {
	case 1:
		text, ok := args[0].(string)
		if !ok {
			panic("WithPlaceholderText(index, text): text must be string")
		}
		return "", text
	case 2:
		phType, ok := args[0].(string)
		if !ok {
			panic("WithPlaceholderText(index, type, text): type must be string")
		}
		text, ok := args[1].(string)
		if !ok {
			panic("WithPlaceholderText(index, type, text): text must be string")
		}
		return phType, text
	default:
		panic("WithPlaceholderText requires (index, text) or (index, type, text)")
	}
}

func parsePlaceholderImageArgs(args ...any) (string, Image) {
	switch len(args) {
	case 1:
		img, ok := args[0].(Image)
		if !ok {
			panic("WithPlaceholderImage(index, image): image must be Image")
		}
		return "", img
	case 2:
		phType, ok := args[0].(string)
		if !ok {
			panic("WithPlaceholderImage(index, type, image): type must be string")
		}
		img, ok := args[1].(Image)
		if !ok {
			panic("WithPlaceholderImage(index, type, image): image must be Image")
		}
		return phType, img
	default:
		panic("WithPlaceholderImage requires (index, image) or (index, type, image)")
	}
}

func parsePlaceholderTableArgs(args ...any) (string, tables.Table) {
	switch len(args) {
	case 1:
		table, ok := args[0].(tables.Table)
		if !ok {
			panic("WithPlaceholderTable(index, table): table must be tables.Table")
		}
		return "", table
	case 2:
		phType, ok := args[0].(string)
		if !ok {
			panic("WithPlaceholderTable(index, type, table): type must be string")
		}
		table, ok := args[1].(tables.Table)
		if !ok {
			panic("WithPlaceholderTable(index, type, table): table must be tables.Table")
		}
		return phType, table
	default:
		panic("WithPlaceholderTable requires (index, table) or (index, type, table)")
	}
}

func parsePlaceholderChartArgs(args ...any) (string, ChartDefinition) {
	switch len(args) {
	case 1:
		chart, ok := args[0].(ChartDefinition)
		if !ok {
			panic("WithPlaceholderChart(index, chart): chart must implement ChartDefinition")
		}
		return "", chart
	case 2:
		phType, ok := args[0].(string)
		if !ok {
			panic("WithPlaceholderChart(index, type, chart): type must be string")
		}
		chart, ok := args[1].(ChartDefinition)
		if !ok {
			panic("WithPlaceholderChart(index, type, chart): chart must implement ChartDefinition")
		}
		return phType, chart
	default:
		panic("WithPlaceholderChart requires (index, chart) or (index, type, chart)")
	}
}

// AddConnector adds a connector to the slide.
func (s SlideContent) AddConnector(c Connector) SlideContent {
	s.Connectors = append(s.Connectors, c)
	return s
}

// WithBarChart sets one bar chart for the slide.
func (s SlideContent) WithBarChart(chart charts.BarChart) SlideContent {
	s.clearCharts()
	s.Chart = &chart
	return s
}

func (s SlideContent) WithBarHorizontalChart(chart charts.BarHorizontalChart) SlideContent {
	s.clearCharts()
	s.BarHorizontal = &chart
	return s
}

func (s SlideContent) WithBarStackedChart(chart charts.BarStackedChart) SlideContent {
	s.clearCharts()
	s.BarStacked = &chart
	return s
}

func (s SlideContent) WithBarStacked100Chart(chart charts.BarStacked100Chart) SlideContent {
	s.clearCharts()
	s.BarStacked100 = &chart
	return s
}

// WithLineChart sets one line chart for the slide.
func (s SlideContent) WithLineChart(chart charts.LineChart) SlideContent {
	s.clearCharts()
	s.Line = &chart
	return s
}

func (s SlideContent) WithLineMarkersChart(chart charts.LineMarkersChart) SlideContent {
	s.clearCharts()
	s.LineMarkers = &chart
	return s
}

func (s SlideContent) WithLineStackedChart(chart charts.LineStackedChart) SlideContent {
	s.clearCharts()
	s.LineStacked = &chart
	return s
}

// WithScatterChart sets one scatter chart for the slide.
func (s SlideContent) WithScatterChart(chart charts.ScatterChart) SlideContent {
	s.clearCharts()
	s.Scatter = &chart
	return s
}

// WithAreaChart sets one area chart for the slide.
func (s SlideContent) WithAreaChart(chart charts.AreaChart) SlideContent {
	s.clearCharts()
	s.Area = &chart
	return s
}

func (s SlideContent) WithAreaStackedChart(chart charts.AreaStackedChart) SlideContent {
	s.clearCharts()
	s.AreaStacked = &chart
	return s
}

func (s SlideContent) WithAreaStacked100Chart(chart charts.AreaStacked100Chart) SlideContent {
	s.clearCharts()
	s.AreaStacked100 = &chart
	return s
}

// WithPieChart sets one pie chart for the slide.
func (s SlideContent) WithPieChart(chart charts.PieChart) SlideContent {
	s.clearCharts()
	s.Pie = &chart
	return s
}

// WithDoughnutChart sets one doughnut chart for the slide.
func (s SlideContent) WithDoughnutChart(chart charts.DoughnutChart) SlideContent {
	s.clearCharts()
	s.Doughnut = &chart
	return s
}

func (s SlideContent) WithBubbleChart(chart charts.BubbleChart) SlideContent {
	s.clearCharts()
	s.Bubble = &chart
	return s
}

func (s SlideContent) WithRadarChart(chart charts.RadarChart) SlideContent {
	s.clearCharts()
	s.Radar = &chart
	return s
}

func (s SlideContent) WithRadarFilledChart(chart charts.RadarFilledChart) SlideContent {
	s.clearCharts()
	s.RadarFilled = &chart
	return s
}

func (s SlideContent) WithStockHLCChart(chart charts.StockHLCChart) SlideContent {
	s.clearCharts()
	s.StockHLC = &chart
	return s
}

func (s SlideContent) WithStockOHLCChart(chart charts.StockOHLCChart) SlideContent {
	s.clearCharts()
	s.StockOHLC = &chart
	return s
}

func (s SlideContent) WithComboChart(chart charts.ComboChart) SlideContent {
	s.clearCharts()
	s.Combo = &chart
	return s
}

func (s *SlideContent) clearCharts() {
	s.Chart = nil
	s.BarHorizontal = nil
	s.BarStacked = nil
	s.BarStacked100 = nil
	s.Line = nil
	s.LineMarkers = nil
	s.LineStacked = nil
	s.Scatter = nil
	s.Area = nil
	s.AreaStacked = nil
	s.AreaStacked100 = nil
	s.Pie = nil
	s.Doughnut = nil
	s.Bubble = nil
	s.Radar = nil
	s.RadarFilled = nil
	s.StockHLC = nil
	s.StockOHLC = nil
	s.Combo = nil
}
