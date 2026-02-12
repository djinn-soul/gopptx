package elements

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

// Validate checks the slide for consistency.
func (s SlideContent) Validate(index int) error {
	for shapeIndex, shape := range s.Shapes {
		if err := shape.Validate(index, shapeIndex+1); err != nil {
			return err
		}
	}
	for connectorIndex, connector := range s.Connectors {
		if err := connector.ValidateWithShapes(s.Shapes, index, connectorIndex+1); err != nil {
			return err
		}
	}
	for imageIndex, image := range s.Images {
		if err := image.Validate(index, imageIndex+1); err != nil {
			return err
		}
	}

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
	if s.TitleColor != "" && !common.IsHexColor(s.TitleColor) {
		return fmt.Errorf("title color must be 6-digit RGB hex")
	}
	if (s.ContentSize != 0 && s.ContentSize < 1) || s.ContentSize > 400 {
		return fmt.Errorf("content size must be between 1 and 400 pt (or 0 for default)")
	}
	if s.ContentColor != "" && !common.IsHexColor(s.ContentColor) {
		return fmt.Errorf("content color must be 6-digit RGB hex")
	}
	if s.Background != nil {
		if err := s.Background.Validate(); err != nil {
			return fmt.Errorf("invalid background: %w", err)
		}
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

	if s.TitleAlign != "" {
		switch s.TitleAlign {
		case "l", "ctr", "r", "just":
		default:
			return fmt.Errorf("invalid title alignment: %q (expected l|ctr|r|just)", s.TitleAlign)
		}
	}

	if s.ContentVAlign != "" {
		switch s.ContentVAlign {
		case "t", "ctr", "b":
		default:
			return fmt.Errorf("invalid content vertical alignment: %q (expected t|ctr|b)", s.ContentVAlign)
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
		if i == 0 && (anim.Trigger == animations.AnimationWithPrevious || anim.Trigger == animations.AnimationAfterPrevious) {
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
	TitleAlign           string
	TitleFont            string
	ContentSize          int
	ContentColor         string
	ContentBold          bool
	ContentItalic        bool
	ContentUnderline     bool
	ContentVAlign        string
	Layout               string
	Background           *SlideBackground
	Transition           transitions.SlideTransition
	DefaultBulletStyle   TextParagraphStyle
	Bullets              []string
	BulletRuns           [][]TextRun
	BulletStyles         []TextParagraphStyle
	ShowSlideNumber      bool
	Notes                string
	NotesBody            []TextParagraph
	Images               []shapes.Image
	Shapes               []shapes.Shape
	Connectors           []shapes.Connector
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
	Animations           []animations.Animation
	PlaceholderOverrides []shapes.PlaceholderContent
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
func (s SlideContent) AddShape(sd shapes.ShapeDefinition) SlideContent {
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
	// Also populate NotesBody for internal consistency
	p := NewTextParagraph()
	p.Runs = append(p.Runs, NewTextRun(notes))
	s.NotesBody = []TextParagraph{p}
	return s
}

// WithRichNotes sets the speaker notes using rich text paragraphs.
func (s SlideContent) WithRichNotes(body []TextParagraph) SlideContent {
	s.NotesBody = body
	// Sync to plain text Notes
	var sb strings.Builder
	for i, p := range body {
		for _, r := range p.Runs {
			sb.WriteString(r.Text)
		}
		if i < len(body)-1 {
			sb.WriteString("\n")
		}
	}
	s.Notes = sb.String()
	return s
}

// AddNoteParagraph appends a rich text paragraph to the speaker notes.
func (s SlideContent) AddNoteParagraph(p TextParagraph) SlideContent {
	s.NotesBody = append(s.NotesBody, p)
	// Sync to plain text Notes
	if s.Notes != "" {
		s.Notes += "\n"
	}
	for _, r := range p.Runs {
		s.Notes += r.Text
	}
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
func (s SlideContent) WithTransition(t transitions.SlideTransition) SlideContent {
	s.Transition = t
	return s
}

// WithTransitionOptions sets built-in transition options.
func (s SlideContent) WithTransitionOptions(opt transitions.TransitionOptions) SlideContent {
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

// WithBackgroundColor sets the slide background as RGB hex.
func (s SlideContent) WithBackgroundColor(color string) SlideContent {
	normalized := common.NormalizeHexColor(color)
	bg := NewSolidBackground(normalized)
	s.Background = &bg
	return s
}

// WithBackground sets a complex background for the slide.
func (s SlideContent) WithBackground(bg SlideBackground) SlideContent {
	s.Background = &bg
	return s
}

// WithGradientBackground sets a gradient background for the slide.
func (s SlideContent) WithGradientBackground(gradient shapes.ShapeGradientFill) SlideContent {
	bg := NewGradientBackground(gradient)
	return s.WithBackground(bg)
}

// WithPictureBackground sets a picture background for the slide using image data.
func (s SlideContent) WithPictureBackground(img shapes.Image) SlideContent {
	bg := NewPictureBackground(img)
	return s.WithBackground(bg)
}

// WithTitleSize sets the title font size in points.
func (s SlideContent) WithTitleSize(size int) SlideContent {
	s.TitleSize = size
	return s
}

// WithTitleColor sets the title color as RGB hex.
func (s SlideContent) WithTitleColor(color string) SlideContent {
	s.TitleColor = common.NormalizeHexColor(color)
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

// WithTitleAlign sets the horizontal alignment of the title (l|ctr|r|just).
func (s SlideContent) WithTitleAlign(align string) SlideContent {
	s.TitleAlign = strings.ToLower(strings.TrimSpace(align))
	return s
}

// WithTitleFont sets the typeface for the slide title (e.g., "Consolas").
func (s SlideContent) WithTitleFont(font string) SlideContent {
	s.TitleFont = strings.TrimSpace(font)
	return s
}

// WithContentSize sets the content font size in points.
func (s SlideContent) WithContentSize(size int) SlideContent {
	s.ContentSize = size
	return s
}

// WithContentColor sets the content color as RGB hex.
func (s SlideContent) WithContentColor(color string) SlideContent {
	s.ContentColor = common.NormalizeHexColor(color)
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

// WithContentVAlign sets the vertical alignment of the main content (t|ctr|b).
func (s SlideContent) WithContentVAlign(align string) SlideContent {
	s.ContentVAlign = strings.ToLower(strings.TrimSpace(align))
	return s
}

// AddImage adds an image to the slide.
func (s SlideContent) AddImage(img shapes.Image) SlideContent {
	s.Images = append(s.Images, img)
	return s
}

// WithSlideNumber enables or disables slide number display.
func (s SlideContent) WithSlideNumber(show bool) SlideContent {
	s.ShowSlideNumber = show
	return s
}

// AddAnimation adds an animation to the slide.
func (s SlideContent) AddAnimation(anim animations.AnimationDefinition) SlideContent {
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

// WithPlaceholderText overrides a placeholder with text using the default placeholder type.
func (s SlideContent) WithPlaceholderText(index int, text string) SlideContent {
	return s.WithPlaceholderTextAs(index, defaultPlaceholderTextType(index), text)
}

// WithPlaceholderTextAs overrides a placeholder with text and explicit placeholder type.
func (s SlideContent) WithPlaceholderTextAs(index int, placeholderType, text string) SlideContent {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Index: index,
		Type:  placeholderType,
		Text:  text,
	})
	return s
}

// WithPlaceholderImage overrides a placeholder with an image using the default placeholder type.
func (s SlideContent) WithPlaceholderImage(index int, img shapes.Image) SlideContent {
	return s.WithPlaceholderImageAs(index, defaultPlaceholderImageType(index), img)
}

// WithPlaceholderImageAs overrides a placeholder with an image and explicit placeholder type.
func (s SlideContent) WithPlaceholderImageAs(index int, placeholderType string, img shapes.Image) SlideContent {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Index: index,
		Type:  placeholderType,
		Image: &img,
	})
	return s
}

// WithPlaceholderTable overrides a placeholder with a table using the default placeholder type.
func (s SlideContent) WithPlaceholderTable(index int, table tables.Table) SlideContent {
	return s.WithPlaceholderTableAs(index, defaultPlaceholderTextType(index), table)
}

// WithPlaceholderTableAs overrides a placeholder with a table and explicit placeholder type.
func (s SlideContent) WithPlaceholderTableAs(index int, placeholderType string, table tables.Table) SlideContent {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Index: index,
		Type:  placeholderType,
		Table: &table,
	})
	return s
}

// WithPlaceholderChart overrides a placeholder with a chart using the default placeholder type.
func (s SlideContent) WithPlaceholderChart(index int, chart ChartDefinition) SlideContent {
	return s.WithPlaceholderChartAs(index, defaultPlaceholderTextType(index), chart)
}

// WithPlaceholderChartAs overrides a placeholder with a chart and explicit placeholder type.
func (s SlideContent) WithPlaceholderChartAs(index int, placeholderType string, chart ChartDefinition) SlideContent {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Index: index,
		Type:  placeholderType,
		Chart: chart,
	})
	return s
}

func defaultPlaceholderTextType(index int) string {
	if index == 0 {
		return "title"
	}
	return "body"
}

func defaultPlaceholderImageType(index int) string {
	if index == 0 {
		return "title"
	}
	return "pic"
}

// AddConnector adds a connector to the slide.
func (s SlideContent) AddConnector(c shapes.Connector) SlideContent {
	s.Connectors = append(s.Connectors, c)
	return s
}

// AutoRerouteConnectors recalculates connector sites from current shape positions.
func (s SlideContent) AutoRerouteConnectors() SlideContent {
	rerouted := make([]shapes.Connector, 0, len(s.Connectors))
	for _, connector := range s.Connectors {
		rerouted = append(rerouted, connector.AutoReroute(s.Shapes))
	}
	s.Connectors = rerouted
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
