package templates

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestSimpleTemplate(t *testing.T) {
	tmpl := SimpleTemplate{
		Title:   "Simple Deck",
		Content: "Hello World",
	}

	slides, err := tmpl.Build()
	if err != nil {
		t.Fatalf("SimpleTemplate Build error: %v", err)
	}

	if len(slides) != 2 {
		t.Errorf("expected 2 slides, got %d", len(slides))
	}

	if slides[0].Title != "Simple Deck" {
		t.Errorf("expected first slide title 'Simple Deck', got %q", slides[0].Title)
	}
	if slides[1].Title != "Content" {
		t.Errorf("expected second slide title 'Content', got %q", slides[1].Title)
	}
	if slides[0].Background == nil {
		t.Fatalf("expected simple template cover slide to have background color")
	}
}

func TestProposalTemplate(t *testing.T) {
	tmpl := ProposalTemplate{
		Title:    "Big Project",
		Context:  "The problem",
		Solution: "The solution",
		Pricing: []PricingTier{
			{Name: "Basic", Price: "$100", Features: []string{"Feature 1"}},
			{Name: "Pro", Price: "$200", Features: []string{"Feature 1", "Feature 2"}},
		},
		Timeline: []Milestone{
			{Date: "2026-03-01", Task: "Phase 1", Status: "Planned"},
		},
	}

	slides, err := tmpl.Build()
	if err != nil {
		t.Fatalf("ProposalTemplate Build error: %v", err)
	}

	if len(slides) != 5 {
		t.Errorf("expected 5 slides, got %d", len(slides))
	}
}

func TestTrainingTemplate(t *testing.T) {
	tmpl := TrainingTemplate{
		Title:    "Go Basics",
		Agenda:   []string{"Syntax", "Types"},
		Concepts: []string{"Variables", "Functions", "Structs"},
		Summary:  "Done",
	}

	slides, err := tmpl.Build()
	if err != nil {
		t.Fatalf("TrainingTemplate Build error: %v", err)
	}

	// Title + Agenda + 3 Concepts + Summary = 6 slides
	expected := 6
	if len(slides) != expected {
		t.Errorf("expected %d slides, got %d", expected, len(slides))
	}
}

func TestStatusTemplate(t *testing.T) {
	tmpl := StatusTemplate{
		Project:   "Project X",
		OKRs:      []string{"Metric A", "Metric B"},
		Risks:     []string{"Risk 1"},
		NextSteps: []string{"Step 1"},
	}

	slides, err := tmpl.Build()
	if err != nil {
		t.Fatalf("StatusTemplate Build error: %v", err)
	}

	if len(slides) != 4 {
		t.Errorf("expected 4 slides, got %d", len(slides))
	}
	if slides[0].Background == nil || slides[1].Background == nil {
		t.Fatalf("expected status template slides to have color styling")
	}
}

func TestTechnicalTemplate(t *testing.T) {
	tmpl := TechnicalTemplate{
		Title:        "Microservices",
		Architecture: "Diagram...",
		DeepDive:     "Code...",
		Benchmarks:   "Fast",
	}

	slides, err := tmpl.Build()
	if err != nil {
		t.Fatalf("TechnicalTemplate Build error: %v", err)
	}

	if len(slides) != 4 {
		t.Errorf("expected 4 slides, got %d", len(slides))
	}
}

func TestTemplateValidation(t *testing.T) {
	tmpl := SimpleTemplate{} // missing title
	_, err := tmpl.Build()
	if err == nil {
		t.Error("expected error for empty title in SimpleTemplate")
	}

	tmpl2 := ProposalTemplate{}
	_, err = tmpl2.Build()
	if err == nil {
		t.Error("expected error for empty title in ProposalTemplate")
	}
}

func TestBrandingPresets(t *testing.T) {
	tests := []struct {
		preset BrandingPreset
		expect string
	}{
		{PresetCorporate, "Corporate"},
		{PresetModern, "Modern"},
		{PresetCreative, "Vibrant"},
	}

	for _, tt := range tests {
		theme := MapPreset(tt.preset)
		if theme.Name != tt.expect {
			t.Errorf("expected preset %q to map to theme %q, got %q", tt.preset, tt.expect, theme.Name)
		}
	}
}

func TestPresetToThemeMapping(t *testing.T) {
	tests := []struct {
		name   string
		preset BrandingPreset
		expect string
	}{
		{name: "corporate", preset: PresetCorporate, expect: "Corporate"},
		{name: "modern", preset: PresetModern, expect: "Modern"},
		{name: "creative", preset: PresetCreative, expect: "Vibrant"},
		{name: "default", preset: BrandingPreset("unknown"), expect: "Corporate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := presetToTheme(tt.preset)
			if theme.Name != tt.expect {
				t.Fatalf("expected %q, got %q", tt.expect, theme.Name)
			}
		})
	}
}

func TestRenderPricingTableStructuredRows(t *testing.T) {
	table := renderPricingTable([]PricingTier{
		{Name: "Starter", Price: "$9", Features: []string{"A", "B"}},
		{Name: "Pro", Price: "$19", Features: []string{"A", "B", "C"}},
	})

	if got, want := len(table.Rows), 3; got != want {
		t.Fatalf("expected %d total rows (header + data), got %d", want, got)
	}

	if table.Rows[0][0] != "Tier" || table.Rows[0][1] != "Price" || table.Rows[0][2] != "Features" {
		t.Fatalf("unexpected pricing header row: %#v", table.Rows[0])
	}

	if table.Rows[1][0] != "Starter" || table.Rows[1][1] != "$9" || table.Rows[1][2] != "A, B" {
		t.Fatalf("unexpected first pricing data row: %#v", table.Rows[1])
	}
}

func TestBrandingApplyAtAddsCoverAndBodyDynamics(t *testing.T) {
	branding := BrandingSpec{Preset: PresetModern, Footer: "Footer"}
	cover := branding.ApplyAt(elements.NewSlide("Cover"), 0)
	body := branding.ApplyAt(elements.NewSlide("Body"), 1)

	if cover.Background == nil || cover.Background.Type != elements.SlideBackgroundSolid {
		t.Fatalf("expected cover slide to have solid background")
	}
	if cover.TitleColor == "" {
		t.Fatalf("expected cover slide title color to be set")
	}
	if cover.TitleSize <= 44 {
		t.Fatalf("expected larger cover title size, got %d", cover.TitleSize)
	}
	if cover.FooterText != "Footer" {
		t.Fatalf("expected footer to be applied")
	}

	if body.Background == nil || body.Background.Type != elements.SlideBackgroundSolid {
		t.Fatalf("expected body slide to have solid background")
	}
	if body.TitleColor == "" {
		t.Fatalf("expected body title color to be set")
	}
	if body.ContentColor == "" {
		t.Fatalf("expected body content color to be set")
	}
}

func TestTemplateLayoutOverrides(t *testing.T) {
	tmpl := SimpleTemplate{
		Title:   "Simple Deck",
		Content: "Hello World",
		LayoutOverrides: LayoutOverrides{
			1: elements.SlideLayoutTwoColumn,
		},
	}

	slides, err := tmpl.Build()
	if err != nil {
		t.Fatalf("SimpleTemplate Build error: %v", err)
	}

	if slides[1].Layout != elements.SlideLayoutTwoColumn {
		t.Fatalf("expected slide 1 layout override to be %q, got %q", elements.SlideLayoutTwoColumn, slides[1].Layout)
	}
}

func TestTemplateLayoutOverridesValidation(t *testing.T) {
	t.Run("index out of range", func(t *testing.T) {
		tmpl := SimpleTemplate{
			Title:   "Simple Deck",
			Content: "Hello World",
			LayoutOverrides: LayoutOverrides{
				10: elements.SlideLayoutBlank,
			},
		}

		_, err := tmpl.Build()
		if err == nil {
			t.Fatal("expected index out of range error")
		}
	})

	t.Run("unsupported layout", func(t *testing.T) {
		tmpl := SimpleTemplate{
			Title:   "Simple Deck",
			Content: "Hello World",
			LayoutOverrides: LayoutOverrides{
				1: "custom-unsupported-layout",
			},
		}

		_, err := tmpl.Build()
		if err == nil {
			t.Fatal("expected unsupported layout error")
		}
	})
}
