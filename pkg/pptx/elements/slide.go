package elements

//go:generate go run ../../../cmd/gen_slide_builder . slide_builder_gen.go

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/ink"
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
	// SlideLayoutTitleAndVerticalText and SlideLayoutVerticalTitleAndText are the
	// two ST_SlideLayoutType layouts PowerPoint offers for vertically written
	// text.
	SlideLayoutTitleAndVerticalText = "title_and_vertical_text"
	SlideLayoutVerticalTitleAndText = "vertical_title_and_text"

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
	Title            string
	TitleSize        int
	TitleColor       string
	TitleBold        bool
	TitleItalic      bool
	TitleUnderline   bool
	TitleAlign       string
	TitleFont        string
	ContentSize      int
	ContentColor     string
	ContentBold      bool
	ContentItalic    bool
	ContentUnderline bool
	ContentVAlign    string
	Layout           string
	// TitleBoundsEMU holds the actual title placeholder geometry read from an
	// existing PPTX (x, y, cx, cy in EMU).  All four values are zero when the
	// slide was constructed programmatically and the renderer should fall back
	// to its built-in layout defaults.
	TitleBoundsEMU [4]int64
	// ContentBoundsEMU holds the actual body/content placeholder geometry read
	// from an existing PPTX (x, y, cx, cy in EMU).  Same zero-means-default
	// convention as TitleBoundsEMU.
	ContentBoundsEMU   [4]int64
	Background         *SlideBackground
	Transition         transitions.SlideTransition
	DefaultBulletStyle ParagraphStyle
	Bullets            []string
	BulletRuns         [][]Run
	BulletStyles       []ParagraphStyle
	ShowSlideNumber    bool
	FooterText         string
	Notes              string
	NotesBody          []Paragraph
	Images             []shapes.Image
	Shapes             []shapes.Shape
	Connectors         []shapes.Connector
	Table              *tables.Table
	// Tables holds any additional tables beyond Table. Slides read from a PPTX
	// can carry several; Table keeps the first for backwards compatibility and
	// the rest land here so exporters do not silently drop them.
	Tables         []tables.Table
	Chart          *charts.BarChart
	BarHorizontal  *charts.BarHorizontalChart
	BarStacked     *charts.BarStackedChart
	BarStacked100  *charts.BarStacked100Chart
	Line           *charts.LineChart
	LineMarkers    *charts.LineMarkersChart
	LineStacked    *charts.LineStackedChart
	Scatter        *charts.ScatterChart
	Area           *charts.AreaChart
	AreaStacked    *charts.AreaStackedChart
	AreaStacked100 *charts.AreaStacked100Chart
	Pie            *charts.PieChart
	Pie3D          *charts.Pie3DChart
	Doughnut       *charts.DoughnutChart
	Bubble         *charts.BubbleChart
	Radar          *charts.RadarChart
	RadarFilled    *charts.RadarFilledChart
	StockHLC       *charts.StockHLCChart
	StockOHLC      *charts.StockOHLCChart
	Combo          *charts.ComboChart
	// Charts holds any additional charts beyond the typed fields above, the way
	// Tables holds the tables beyond Table. A slide can carry several charts;
	// the typed pointers are one per kind, so this is how a caller places two
	// of the same kind, or more than one chart at all.
	Charts               []charts.ChartDefinition
	Animations           []animations.Animation
	SmartArtDiagrams     []smartart.SmartArt
	PlaceholderOverrides []shapes.PlaceholderContent
	Comments             []SlideComment
	// InkAnnotations holds pen strokes drawn on the slide. Each annotation
	// becomes one InkML part plus a content-part reference in the shape tree.
	InkAnnotations []*ink.Annotation
	Hidden         bool
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

// Validate checks the slide for consistency, returning the first problem found.
func (s SlideContent) Validate(index int) error {
	return validateSlideContent(s, index)
}

// ValidateAll reports every logical problem on the slide.
//
// Per-object rules are independent, so stopping at the first defect forces a
// caller to run validate → fix → validate once per problem. Slide-level rules
// still short-circuit, and contribute at most one error.
func (s SlideContent) ValidateAll(index int) []error {
	errs := validateSlideObjectsAll(s, index)
	if err := validateSlideContentBeyondObjects(s, index); err != nil {
		errs = append(errs, err)
	}
	return errs
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
