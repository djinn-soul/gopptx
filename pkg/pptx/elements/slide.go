package elements

import (
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

const (
	SlideLayoutTitleAndContent    = "title_and_content"
	SlideLayoutTitleOnly          = "title_only"
	SlideLayoutBlank              = "blank"
	SlideLayoutCenteredTitle      = "centered_title"
	SlideLayoutTitleAndBigContent = "title_and_big_content"
	SlideLayoutTwoColumn          = "two_column"

	// SlideLayoutTitle starts the legacy/descriptive layout aliases.
	SlideLayoutTitle          = "Title Slide"
	SlideLayoutSectionHeader  = "Section Header"
	SlideLayoutTwoContent     = "Two Content"
	SlideLayoutComparison     = "Comparison"
	SlideLayoutContentCaption = "Content with Caption"
	SlideLayoutPictureCaption = "Picture with Caption"

	placeholderTypeTitle = "title"
	placeholderTypeBody  = "body"
	placeholderTypePic   = "pic"

	defaultSlideTitleSizePt   = 44
	defaultSlideContentSizePt = 18
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
	DefaultBulletStyle   ParagraphStyle
	Bullets              []string
	BulletRuns           [][]Run
	BulletStyles         []ParagraphStyle
	ShowSlideNumber      bool
	FooterText           string
	Notes                string
	NotesBody            []Paragraph
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
	SmartArtDiagrams     []smartart.SmartArt
	PlaceholderOverrides []shapes.PlaceholderContent
	Comments             []SlideComment
	Hidden               bool
}

// SlideComment describes an author's comment on a slide.
type SlideComment struct {
	AuthorName string
	Text       string
	// Optional coordinates in EMU. If 0, defaults are applied locally.
	X int64
	Y int64
}

// NewSlide creates a new slide with default settings and a title.
func NewSlide(title string) SlideContent {
	return SlideContent{
		Title:       title,
		TitleSize:   defaultSlideTitleSizePt,
		ContentSize: defaultSlideContentSizePt,
		Layout:      SlideLayoutTitleAndContent,
	}
}

// Validate checks the slide for consistency.
func (s SlideContent) Validate(index int) error {
	return validateSlideContent(s, index)
}

// AddBullet appends one bullet item and returns the updated slide.
func (s SlideContent) AddBullet(text string) SlideContent {
	s.Bullets = append(s.Bullets, text)
	s.BulletRuns = append(s.BulletRuns, nil)
	s.BulletStyles = append(s.BulletStyles, s.DefaultBulletStyle)
	return s
}

// AddBulletWithStyle appends one bullet item with explicit styling.
func (s SlideContent) AddBulletWithStyle(text string, style ParagraphStyle) SlideContent {
	s.Bullets = append(s.Bullets, text)
	s.BulletRuns = append(s.BulletRuns, nil)
	s.BulletStyles = append(s.BulletStyles, style)
	return s
}

// AddBulletRuns appends one bullet item with rich text runs.
func (s SlideContent) AddBulletRuns(runs []Run) SlideContent {
	s.Bullets = append(s.Bullets, RunsToPlainText(runs))
	s.BulletRuns = append(s.BulletRuns, runs)
	s.BulletStyles = append(s.BulletStyles, s.DefaultBulletStyle)
	return s
}

// AddBulletRunsWithStyle appends one bullet item with rich text runs and paragraph styling.
func (s SlideContent) AddBulletRunsWithStyle(runs []Run, style ParagraphStyle) SlideContent {
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

// AddComment appends a comment from the specified author.
func (s SlideContent) AddComment(authorName, text string) SlideContent {
	s.Comments = append(s.Comments, SlideComment{
		AuthorName: authorName,
		Text:       text,
	})
	return s
}

// WithDefaultBulletStyle sets the base style for new bullets.
func (s SlideContent) WithDefaultBulletStyle(style ParagraphStyle) SlideContent {
	s.DefaultBulletStyle = style
	return s
}

// WithNotes sets the speaker notes for the slide.
func (s SlideContent) WithNotes(notes string) SlideContent {
	s.Notes = notes
	// Also populate NotesBody for internal consistency
	p := NewParagraph()
	p.Runs = append(p.Runs, NewRun(notes))
	s.NotesBody = []Paragraph{p}
	return s
}

// WithRichNotes sets the speaker notes using rich text paragraphs.
func (s SlideContent) WithRichNotes(body []Paragraph) SlideContent {
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
func (s SlideContent) AddNoteParagraph(p Paragraph) SlideContent {
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

// AddNoteBullet appends a bulleted paragraph to the speaker notes.
func (s SlideContent) AddNoteBullet(text string) SlideContent {
	p := NewParagraph()
	p.Runs = append(p.Runs, NewRun(text))
	p.Style.BulletStyle = BulletStyleBullet
	return s.AddNoteParagraph(p)
}

// AddNoteNumbered appends a numbered paragraph to the speaker notes.
func (s SlideContent) AddNoteNumbered(text string) SlideContent {
	p := NewParagraph()
	p.Runs = append(p.Runs, NewRun(text))
	p.Style.BulletStyle = BulletStyleNumber
	return s.AddNoteParagraph(p)
}

// AddNoteSubBullet appends an indented bullet paragraph to the speaker notes.
func (s SlideContent) AddNoteSubBullet(level int, text string) SlideContent {
	p := NewParagraph()
	p.Runs = append(p.Runs, NewRun(text))
	p.Style.BulletStyle = BulletStyleBullet
	p.Style.Level = level
	return s.AddNoteParagraph(p)
}

// WithTable sets the table for the slide.
func (s SlideContent) WithTable(t tables.Table) SlideContent {
	s.Table = &t
	return s
}

// WithBulletStyle sets the bullet style for all bullets on this slide.
func (s SlideContent) WithBulletStyle(style ParagraphStyle) SlideContent {
	s.DefaultBulletStyle = style
	for i := range s.BulletStyles {
		s.BulletStyles[i] = style
	}
	return s
}

// AddNumbered appends one numbered bullet item.
func (s SlideContent) AddNumbered(text string) SlideContent {
	return s.AddBulletWithStyle(text, DefaultParagraphStyle().WithNumbered())
}

// AddLettered appends one lettered bullet item.
func (s SlideContent) AddLettered(text string) SlideContent {
	return s.AddBulletWithStyle(text, DefaultParagraphStyle().WithLetteredLower())
}

// AddSubBullet adds a bullet at level index 1..8.
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

// WithMorphTransition sets a Morph transition using PowerPoint's default mode.
// We intentionally omit the option attribute to maximize compatibility.
func (s SlideContent) WithMorphTransition() SlideContent {
	s.Transition = transitions.TransitionOptions{
		Type: transitions.TransitionMorph,
	}
	return s
}

// WithMorphTransitionOptions sets a Morph transition with explicit options.
func (s SlideContent) WithMorphTransitionOptions(option transitions.MorphOption) SlideContent {
	s.Transition = transitions.TransitionOptions{
		Type:        transitions.TransitionMorph,
		MorphOption: option,
	}
	return s
}

// WithTransitionSound sets a sound file for the slide transition.
func (s SlideContent) WithTransitionSound(path string) SlideContent {
	// If transition is nil or not options, default to cut.
	opt, ok := s.Transition.(transitions.TransitionOptions)
	if !ok {
		opt = transitions.TransitionOptions{Type: transitions.TransitionCut}
	}
	if opt.Sound == nil {
		opt.Sound = &transitions.TransitionSound{}
	}
	// Store the path in RelID temporarily; it will be resolved to a relation ID
	// during package writing.
	opt.Sound.RelID = "file:" + path
	opt.Sound.Name = filepath.Base(path)
	s.Transition = opt
	return s
}

// WithBulletStyleName sets primary bullet style by name (e.g. BulletStyleNumber).
func (s SlideContent) WithBulletStyleName(styleName string) SlideContent {
	style := s.DefaultBulletStyle
	style.BulletStyle = NormalizeBulletStyle(styleName)
	return s.WithBulletStyle(style)
}

// WithLayout sets the slide layout (supports canonical and compatibility aliases).
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

// AddSmartArt adds a SmartArt diagram to the slide.
func (s SlideContent) AddSmartArt(sa smartart.SmartArt) SlideContent {
	s.SmartArtDiagrams = append(s.SmartArtDiagrams, sa)
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
	case "", strings.ToLower(SlideLayoutTitleAndContent), "titleandcontent", "title-and-content":
		return SlideLayoutTitleAndContent
	case strings.ToLower(SlideLayoutTitle), "title slide":
		return SlideLayoutTitle
	case strings.ToLower(SlideLayoutSectionHeader), "section header":
		return SlideLayoutSectionHeader
	case strings.ToLower(SlideLayoutTwoContent), "two content":
		return SlideLayoutTwoContent
	case strings.ToLower(SlideLayoutComparison), "comparison":
		return SlideLayoutComparison
	case strings.ToLower(SlideLayoutContentCaption), "content with caption":
		return SlideLayoutContentCaption
	case strings.ToLower(SlideLayoutPictureCaption), "picture with caption":
		return SlideLayoutPictureCaption
	case strings.ToLower(SlideLayoutTitleOnly), "titleonly", "title-only":
		return SlideLayoutTitleOnly
	case strings.ToLower(SlideLayoutBlank), "blank":
		return SlideLayoutBlank
	case strings.ToLower(SlideLayoutCenteredTitle), "centeredtitle", "centered-title":
		return SlideLayoutCenteredTitle
	case strings.ToLower(SlideLayoutTitleAndBigContent), "titleandbigcontent", "title-and-big-content", "big_content":
		return SlideLayoutTitleAndBigContent
	case strings.ToLower(SlideLayoutTwoColumn), "twocolumn", "two-column", "two column":
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

// WithPlaceholderOverride adds custom geometry or style overrides to a placeholder.
func (s SlideContent) WithPlaceholderOverride(
	target shapes.PlaceholderTarget,
	options shapes.PlaceholderOverrideOptions,
) SlideContent {
	// If no type/index specified, default to title (index 0)
	if target.Type == "" && target.Index == 0 && target.Name == "" {
		target.Index = 0
		target.Type = placeholderTypeTitle
	}

	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Index:    target.Index,
		Type:     target.Type,
		Target:   &target,
		Override: &options,
	})
	return s
}

func defaultPlaceholderTextType(index int) string {
	if index == 0 {
		return placeholderTypeTitle
	}
	return placeholderTypeBody
}

func defaultPlaceholderImageType(index int) string {
	if index == 0 {
		return placeholderTypeTitle
	}
	return placeholderTypePic
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
	clearCharts(&s)
	s.Chart = &chart
	return s
}

func (s SlideContent) WithBarHorizontalChart(chart charts.BarHorizontalChart) SlideContent {
	clearCharts(&s)
	s.BarHorizontal = &chart
	return s
}

func (s SlideContent) WithBarStackedChart(chart charts.BarStackedChart) SlideContent {
	clearCharts(&s)
	s.BarStacked = &chart
	return s
}

func (s SlideContent) WithBarStacked100Chart(chart charts.BarStacked100Chart) SlideContent {
	clearCharts(&s)
	s.BarStacked100 = &chart
	return s
}

// WithLineChart sets one line chart for the slide.
func (s SlideContent) WithLineChart(chart charts.LineChart) SlideContent {
	clearCharts(&s)
	s.Line = &chart
	return s
}

func (s SlideContent) WithLineMarkersChart(chart charts.LineMarkersChart) SlideContent {
	clearCharts(&s)
	s.LineMarkers = &chart
	return s
}

func (s SlideContent) WithLineStackedChart(chart charts.LineStackedChart) SlideContent {
	clearCharts(&s)
	s.LineStacked = &chart
	return s
}

// WithScatterChart sets one scatter chart for the slide.
func (s SlideContent) WithScatterChart(chart charts.ScatterChart) SlideContent {
	clearCharts(&s)
	s.Scatter = &chart
	return s
}

// WithAreaChart sets one area chart for the slide.
func (s SlideContent) WithAreaChart(chart charts.AreaChart) SlideContent {
	clearCharts(&s)
	s.Area = &chart
	return s
}

func (s SlideContent) WithAreaStackedChart(chart charts.AreaStackedChart) SlideContent {
	clearCharts(&s)
	s.AreaStacked = &chart
	return s
}

func (s SlideContent) WithAreaStacked100Chart(chart charts.AreaStacked100Chart) SlideContent {
	clearCharts(&s)
	s.AreaStacked100 = &chart
	return s
}

// WithPieChart sets one pie chart for the slide.
func (s SlideContent) WithPieChart(chart charts.PieChart) SlideContent {
	clearCharts(&s)
	s.Pie = &chart
	return s
}

// WithDoughnutChart sets one doughnut chart for the slide.
func (s SlideContent) WithDoughnutChart(chart charts.DoughnutChart) SlideContent {
	clearCharts(&s)
	s.Doughnut = &chart
	return s
}

func (s SlideContent) WithBubbleChart(chart charts.BubbleChart) SlideContent {
	clearCharts(&s)
	s.Bubble = &chart
	return s
}

func (s SlideContent) WithRadarChart(chart charts.RadarChart) SlideContent {
	clearCharts(&s)
	s.Radar = &chart
	return s
}

func (s SlideContent) WithRadarFilledChart(chart charts.RadarFilledChart) SlideContent {
	clearCharts(&s)
	s.RadarFilled = &chart
	return s
}

func (s SlideContent) WithStockHLCChart(chart charts.StockHLCChart) SlideContent {
	clearCharts(&s)
	s.StockHLC = &chart
	return s
}

func (s SlideContent) WithStockOHLCChart(chart charts.StockOHLCChart) SlideContent {
	clearCharts(&s)
	s.StockOHLC = &chart
	return s
}

func (s SlideContent) WithComboChart(chart charts.ComboChart) SlideContent {
	clearCharts(&s)
	s.Combo = &chart
	return s
}

func clearCharts(s *SlideContent) {
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
