package elements_test

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
	"github.com/djinn-soul/gopptx/pkg/pptx/validation/structural"
)

func TestNewSlide(t *testing.T) {
	s := elements.NewSlide("Test Title")
	if s.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got %q", s.Title)
	}
	if s.Layout != elements.SlideLayoutTitleAndContent {
		t.Errorf("expected layout %q, got %q", elements.SlideLayoutTitleAndContent, s.Layout)
	}
}

func TestSlideContent_Mutations(t *testing.T) {
	s := elements.NewSlide("Title")
	s = s.WithTitleSize(50).
		WithTitleColor("#FF0000").
		WithTitleBold(true).
		WithTitleItalic(true).
		WithTitleUnderline(true).
		WithTitleAlign("ctr").
		WithTitleFont("Arial").
		WithContentSize(20).
		WithContentColor("#00FF00").
		WithContentBold(true).
		WithContentItalic(true).
		WithContentUnderline(true).
		WithContentVAlign("ctr").
		WithSlideNumber(true).
		WithLayout(elements.SlideLayoutBlank)

	if s.TitleSize != 50 {
		t.Errorf("expected TitleSize 50, got %d", s.TitleSize)
	}
	if s.TitleColor != "FF0000" {
		t.Errorf("expected TitleColor FF0000, got %q", s.TitleColor)
	}
	if !s.TitleBold || !s.TitleItalic || !s.TitleUnderline {
		t.Error("expected TitleBold, TitleItalic, TitleUnderline to be true")
	}
	if s.TitleAlign != "ctr" {
		t.Errorf("expected TitleAlign ctr, got %q", s.TitleAlign)
	}
	if s.TitleFont != "Arial" {
		t.Errorf("expected TitleFont Arial, got %q", s.TitleFont)
	}
	if s.ContentSize != 20 {
		t.Errorf("expected ContentSize 20, got %d", s.ContentSize)
	}
	if s.ContentColor != "00FF00" {
		t.Errorf("expected ContentColor 00FF00, got %q", s.ContentColor)
	}
	if !s.ContentBold || !s.ContentItalic || !s.ContentUnderline {
		t.Error("expected ContentBold, ContentItalic, ContentUnderline to be true")
	}
	if s.ContentVAlign != "ctr" {
		t.Errorf("expected ContentVAlign ctr, got %q", s.ContentVAlign)
	}
	if !s.ShowSlideNumber {
		t.Error("expected ShowSlideNumber to be true")
	}
	if s.Layout != elements.SlideLayoutBlank {
		t.Errorf("expected Layout %q, got %q", elements.SlideLayoutBlank, s.Layout)
	}
}

func TestSlideContent_Bullets(t *testing.T) {
	s := elements.NewSlide("Bullets")
	s = s.AddBullet("Bullet 1").
		AddNumbered("Numbered 1").
		AddLettered("Lettered 1").
		AddSubBullet(1, "SubBullet 1")

	if len(s.Bullets) != 4 {
		t.Errorf("expected 4 bullets, got %d", len(s.Bullets))
	}
	if s.BulletStyles[1].BulletStyle != elements.BulletStyleNumber {
		t.Errorf("expected BulletStyleNumber, got %q", s.BulletStyles[1].BulletStyle)
	}
	if s.BulletStyles[2].BulletStyle != elements.BulletStyleLetterLower {
		t.Errorf("expected BulletStyleLetterLower, got %q", s.BulletStyles[2].BulletStyle)
	}
	if s.BulletStyles[3].Level != 1 {
		t.Errorf("expected level 1, got %d", s.BulletStyles[3].Level)
	}

	s = s.WithBulletStyleName("number")
	if s.DefaultBulletStyle.BulletStyle != elements.BulletStyleNumber {
		t.Errorf("expected default style number, got %q", s.DefaultBulletStyle.BulletStyle)
	}
}

func TestSlideContent_Notes(t *testing.T) {
	s := elements.NewSlide("Notes")
	s = s.WithNotes("Plain notes")
	if s.Notes != "Plain notes" || len(s.NotesBody) != 1 {
		t.Error("WithNotes failed")
	}

	s = s.AddNoteBullet("Bullet note")
	if len(s.NotesBody) != 2 || s.NotesBody[1].Style.BulletStyle != elements.BulletStyleBullet {
		t.Error("AddNoteBullet failed")
	}

	s = s.AddNoteNumbered("Numbered note")
	if len(s.NotesBody) != 3 || s.NotesBody[2].Style.BulletStyle != elements.BulletStyleNumber {
		t.Error("AddNoteNumbered failed")
	}

	s = s.AddNoteSubBullet(1, "Subbullet note")
	if len(s.NotesBody) != 4 || s.NotesBody[3].Style.Level != 1 {
		t.Error("AddNoteSubBullet failed")
	}
}

func TestSlideContent_Elements(t *testing.T) {
	s := elements.NewSlide("Elements")
	s = s.AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, styling.Inches(1), styling.Inches(1), styling.Inches(2), styling.Inches(2))).
		AddImage(shapes.Image{Data: testutil.TinyPNG(), Format: "png", CX: styling.Inches(1), CY: styling.Inches(1)}).
		WithTable(tables.NewTable([]styling.Length{styling.Inches(1), styling.Inches(1)}).AddRow([]string{"A", "B"})).
		AddComment("Author", "Comment text").
		AddAnimation(animations.NewAnimation(1, animations.AnimationEntranceAppear)).
		AddSmartArt(smartart.NewSmartArt(smartart.BasicBlockList).AddNode(smartart.NewNode("Node"))).
		AddConnector(shapes.NewStraightConnector(0, 0, styling.Inches(1), styling.Inches(1)))

	if len(s.Shapes) != 1 {
		t.Error("AddShape failed")
	}
	if len(s.Images) != 1 {
		t.Error("AddImage failed")
	}
	if s.Table == nil {
		t.Error("WithTable failed")
	}
	if len(s.Comments) != 1 {
		t.Error("AddComment failed")
	}
	if len(s.Animations) != 1 {
		t.Error("AddAnimation failed")
	}
	if len(s.SmartArtDiagrams) != 1 {
		t.Error("AddSmartArt failed")
	}
	if len(s.Connectors) != 1 {
		t.Error("AddConnector failed")
	}

	// Full generation and validation
	data, err := pptx.CreateWithSlides("Elements Test", []elements.SlideContent{s})
	if err != nil {
		t.Fatalf("CreateWithSlides failed: %v", err)
	}
	validatePPTX(t, data)
}

func TestSlideContent_Background(t *testing.T) {
	s := elements.NewSlide("Background")
	s = s.WithBackgroundColor("#FFFFFF")
	if s.Background == nil || s.Background.Type != elements.SlideBackgroundSolid {
		t.Error("WithBackgroundColor failed")
	}

	s = s.WithGradientBackground(shapes.ShapeGradientFill{})
	if s.Background.Type != elements.SlideBackgroundGradient {
		t.Error("WithGradientBackground failed")
	}

	s = s.WithPictureBackground(shapes.Image{Data: testutil.TinyPNG(), Format: "png"})
	if s.Background.Type != elements.SlideBackgroundPicture {
		t.Error("WithPictureBackground failed")
	}
}

func TestSlideContent_Transitions(t *testing.T) {
	s := elements.NewSlide("Transitions")
	s = s.WithTransition(transitions.TransitionOptions{Type: transitions.TransitionFade})
	if s.Transition == nil {
		t.Error("WithTransition failed")
	}

	s = s.WithMorphTransition()
	opt := s.Transition.(transitions.TransitionOptions)
	if opt.Type != transitions.TransitionMorph {
		t.Error("WithMorphTransition failed")
	}

	s = s.WithTransitionSound("sound.wav")
	opt = s.Transition.(transitions.TransitionOptions)
	if opt.Sound == nil || opt.Sound.Name != "sound.wav" {
		t.Error("WithTransitionSound failed")
	}
}

func TestSlideContent_Layouts(t *testing.T) {
	s := elements.NewSlide("Layouts")
	s = s.WithTitleOnlyLayout()
	if s.Layout != elements.SlideLayoutTitleOnly {
		t.Error("WithTitleOnlyLayout failed")
	}
	s = s.WithBlankLayout()
	if s.Layout != elements.SlideLayoutBlank {
		t.Error("WithBlankLayout failed")
	}
	s = s.WithCenteredTitleLayout()
	if s.Layout != elements.SlideLayoutCenteredTitle {
		t.Error("WithCenteredTitleLayout failed")
	}
	s = s.WithTitleAndBigContentLayout()
	if s.Layout != elements.SlideLayoutTitleAndBigContent {
		t.Error("WithTitleAndBigContentLayout failed")
	}
	s = s.WithTwoColumnLayout()
	if s.Layout != elements.SlideLayoutTwoColumn {
		t.Error("WithTwoColumnLayout failed")
	}
	s = s.WithTitleAndContentLayout()
	if s.Layout != elements.SlideLayoutTitleAndContent {
		t.Error("WithTitleAndContentLayout failed")
	}
}

func TestSlideContent_Placeholders(t *testing.T) {
	s := elements.NewSlide("Placeholders")
	s = s.WithPlaceholderText(1, "Text").
		WithPlaceholderImage(2, shapes.Image{Data: testutil.TinyPNG(), Format: "png"}).
		WithPlaceholderTable(3, tables.NewTable([]styling.Length{styling.Inches(1)}).AddRow([]string{"A"})).
		WithPlaceholderChart(4, &charts.BarChart{})

	if len(s.PlaceholderOverrides) != 4 {
		t.Errorf("expected 4 placeholder overrides, got %d", len(s.PlaceholderOverrides))
	}
}

func TestSlideContent_Charts(t *testing.T) {
	s := elements.NewSlide("Charts")
	s = s.WithBarChart(charts.BarChart{}).
		WithLineChart(charts.LineChart{}).
		WithPieChart(charts.PieChart{}).
		WithScatterChart(charts.ScatterChart{}).
		WithAreaChart(charts.AreaChart{}).
		WithDoughnutChart(charts.DoughnutChart{}).
		WithBubbleChart(charts.BubbleChart{}).
		WithRadarChart(charts.RadarChart{}).
		WithStockHLCChart(charts.StockHLCChart{}).
		WithComboChart(charts.ComboChart{})

	if s.Combo == nil {
		t.Error("WithComboChart failed")
	}
	if s.Chart != nil {
		t.Error("clearCharts failed to clear previous chart")
	}
}

func TestNormalizeSlideLayout(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", elements.SlideLayoutTitleAndContent},
		{"Title Slide", elements.SlideLayoutTitle},
		{"title-only", elements.SlideLayoutTitleOnly},
		{"BIG_CONTENT", elements.SlideLayoutTitleAndBigContent},
		{"Two Column", elements.SlideLayoutTwoColumn},
		{"unknown", "unknown"},
	}
	for _, tt := range tests {
		if got := elements.NormalizeSlideLayout(tt.input); got != tt.expected {
			t.Errorf("NormalizeSlideLayout(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

type mockPartStore struct {
	parts map[string][]byte
}

func (m *mockPartStore) Has(path string) bool {
	_, ok := m.parts[path]
	return ok
}

func (m *mockPartStore) Get(path string) ([]byte, bool) {
	data, ok := m.parts[path]
	return data, ok
}

func (m *mockPartStore) Keys() []string {
	keys := make([]string, 0, len(m.parts))
	for k := range m.parts {
		keys = append(keys, k)
	}
	return keys
}

func validatePPTX(t *testing.T, blob []byte) {
	r, err := zip.NewReader(bytes.NewReader(blob), int64(len(blob)))
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}

	m := &mockPartStore{parts: make(map[string][]byte)}
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("failed to open file %s: %v", f.Name, err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatalf("failed to read file %s: %v", f.Name, err)
		}
		m.parts[f.Name] = data
	}

	v := structural.NewValidator(m)
	issues := v.Validate()
	for _, issue := range issues {
		if issue.Severity == structural.SeverityError {
			t.Errorf("validation error: %v", issue)
		}
	}
}
