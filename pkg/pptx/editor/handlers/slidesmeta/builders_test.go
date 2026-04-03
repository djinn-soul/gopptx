package slidesmeta

import (
	"errors"
	"strings"
	"testing"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestBuildChartDefinition(t *testing.T) {
	base := editorcommand.AddChartRequest{
		ChartType:  "bar",
		Title:      "Revenue",
		Categories: []string{"Q1", "Q2"},
		Values:     []float64{10, 20},
		X:          10,
		Y:          20,
		W:          30,
		H:          40,
	}
	def, err := BuildChartDefinition(base)
	if err != nil {
		t.Fatalf("BuildChartDefinition(bar) failed: %v", err)
	}
	if def == nil {
		t.Fatal("expected non-nil chart definition")
	}

	for _, chartType := range []string{
		"line",
		"pie",
		"barHorizontal",
		"barStacked",
		"barStacked100",
		"lineMarkers",
		"lineStacked",
		"scatter",
		"area",
		"areaStacked",
		"areaStacked100",
		"doughnut",
		"bubble",
		"radar",
		"radarFilled",
		"stockHLC",
		"stockOHLC",
	} {
		req := base
		req.ChartType = chartType
		if _, err = BuildChartDefinition(req); err != nil {
			t.Fatalf("BuildChartDefinition(%s) failed: %v", chartType, err)
		}
	}

	comboReq := base
	comboReq.ChartType = "combo"
	comboReq.BarSeries = []editorcommand.ChartSeriesRequest{
		{Name: "Revenue", Values: []float64{10, 20}},
	}
	comboReq.LineSeries = []editorcommand.ChartSeriesRequest{
		{Name: "Growth", Values: []float64{2, 3}},
	}
	if _, err = BuildChartDefinition(comboReq); err != nil {
		t.Fatalf("BuildChartDefinition(combo) failed: %v", err)
	}

	_, err = BuildChartDefinition(editorcommand.AddChartRequest{
		ChartType:  "notAChart",
		Categories: []string{"A"},
		Values:     []float64{1},
	})
	if err == nil {
		t.Fatal("expected unsupported chart type error")
	}
	if !errors.Is(err, ErrUnsupportedChartType) {
		t.Fatalf("expected ErrUnsupportedChartType, got %v", err)
	}
	if !strings.Contains(err.Error(), `"notAChart"`) {
		t.Fatalf("expected chart type in wrapped error message, got %v", err)
	}
}

func TestBuildSlideContent(t *testing.T) {
	req := editorcommand.UpdateSlideRequest{
		Title:   "New Title",
		Layout:  elements.SlideLayoutTwoColumn,
		Bullets: []string{"One", "Two"},
	}
	slide := BuildSlideContent(req, "Current Title")
	if slide.Title != "New Title" {
		t.Fatalf("expected explicit title override, got %q", slide.Title)
	}
	if slide.Layout != elements.SlideLayoutTwoColumn {
		t.Fatalf("expected requested layout, got %q", slide.Layout)
	}
	if len(slide.Bullets) != 2 || slide.Bullets[1] != "Two" {
		t.Fatalf("expected bullets to be copied, got %+v", slide.Bullets)
	}

	req = editorcommand.UpdateSlideRequest{}
	slide = BuildSlideContent(req, "Current Title")
	if slide.Title != "Current Title" {
		t.Fatalf("expected fallback current title, got %q", slide.Title)
	}
}

func TestResolveThemeByName(t *testing.T) {
	cases := map[string]styling.Theme{
		"Corporate": styling.ThemeCorporate,
		"Modern":    styling.ThemeModern,
		"Vibrant":   styling.ThemeVibrant,
		"Dark":      styling.ThemeDark,
		"Nature":    styling.ThemeNature,
		"Tech":      styling.ThemeTech,
		"Carbon":    styling.ThemeCarbon,
	}

	for name, want := range cases {
		got, err := ResolveThemeByName(name)
		if err != nil {
			t.Fatalf("ResolveThemeByName(%q) failed: %v", name, err)
		}
		if got.Name != want.Name {
			t.Fatalf("ResolveThemeByName(%q) returned wrong theme name %q != %q", name, got.Name, want.Name)
		}
	}

	_, err := ResolveThemeByName("Unknown")
	if err == nil {
		t.Fatal("expected unknown theme error")
	}
	if !errors.Is(err, ErrUnknownThemeName) {
		t.Fatalf("expected ErrUnknownThemeName, got %v", err)
	}
}
