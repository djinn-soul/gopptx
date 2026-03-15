package elements

import (
	"testing"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func TestChartHelpers_Counts(t *testing.T) {
	s := NewSlide("Test")
	// Test direct chart counts
	s.Chart = &charts.BarChart{}
	s.BarHorizontal = &charts.BarHorizontalChart{}
	s.Pie = &charts.PieChart{}

	if count := s.directChartCount(); count != 3 {
		t.Errorf("expected 3 direct charts, got %d", count)
	}

	// Test placeholder chart counts
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Chart: &charts.BarChart{},
	})

	if count := s.placeholderChartCount(); count != 1 {
		t.Errorf("expected 1 placeholder chart, got %d", count)
	}

	if count := s.ChartKindCount(); count != 4 {
		t.Errorf("expected 4 total charts, got %d", count)
	}
}

func TestText_Normalizers(t *testing.T) {
	// Simple New* tests
	_ = NewParagraphStyle()
	_ = NewParagraph()

	// Normalizers
	if align := NormalizeTextAlign("  Center  "); align != "center" {
		t.Errorf("expected normalized lowercase center, got %q", align)
	}

	style := ParagraphStyle{
		Level: 10, // MaxBulletLevel is 8
	}
	normalized := NormalizeParagraphStyle(style)
	if normalized.Level != 10 {
		t.Errorf("expected level passthrough 10, got %d", normalized.Level)
	}

	runs := []Run{{Text: "Hello"}, {Text: "World"}}
	normalizedRuns := NormalizeRuns(runs)
	if len(normalizedRuns) != 1 || normalizedRuns[0].Text != "HelloWorld" {
		t.Errorf("expected merged runs, got %+v", normalizedRuns)
	}
}

func TestSlideMaster_Metadata(t *testing.T) {
	m := NewMaster()
	m.WithFooter("Footer Content")
	if m.FooterText != "Footer Content" {
		t.Errorf("expected Footer Content, got %q", m.FooterText)
	}

	ph := shapes.Placeholder{Type: "body"}
	m.AddPlaceholder(ph)
	if len(m.Placeholders) != 1 {
		t.Errorf("expected 1 placeholder")
	}

	m.WithColorMapping("lt1", "dk1")
	if m.ColorMapping == nil || m.ColorMapping.BG1 != "lt1" {
		t.Errorf("expected color mapping")
	}

	m.WithTitleStyle([]TextLevelStyle{{Level: 0, SizePt: 24}})
	if len(m.TxStyles.TitleStyle) != 1 {
		t.Errorf("expected title style")
	}

	m.WithBodyStyle([]TextLevelStyle{{Level: 1, SizePt: 18}})
	if len(m.TxStyles.BodyStyle) != 1 {
		t.Errorf("expected body style")
	}

	m.WithOtherStyle([]TextLevelStyle{{Level: 2, SizePt: 12}})
	if len(m.TxStyles.OtherStyle) != 1 {
		t.Errorf("expected other style")
	}
}

type mockTransition struct {
	xml string
}
func (m mockTransition) XML() string { return m.xml }
func (m mockTransition) Validate() error { return nil }

func TestTransition_Coverage(t *testing.T) {
	s := NewSlide("Test")
	if SlideTransitionXML(s) != "" {
		t.Errorf("expected empty transition XML")
	}

	s.Transition = mockTransition{xml: "<p:transition>...</p:transition>"}
	if xml := SlideTransitionXML(s); xml != "<p:transition>...</p:transition>" {
		t.Errorf("got %q", xml)
	}
}
