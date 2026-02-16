package charts_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

type chartVisualFixtureEntry struct {
	Chart       string `json:"chart"`
	Fingerprint string `json:"fingerprint"`
}

func TestChartVisualRegressionFixtures(t *testing.T) {
	slides := chartVisualRegressionSlides()
	charts := sortedSlideKeys(slides)

	got := make(map[string]string, len(charts))
	for _, chart := range charts {
		xml := chartXMLForSlide(t, slides[chart])
		got[chart] = chartVisualFingerprint(xml)
	}

	if os.Getenv("UPDATE_CHART_VISUAL_FIXTURES") == "1" {
		writeChartVisualFixtures(t, got)
		return
	}

	expected := loadChartVisualFixtures(t)

	if len(expected) != len(got) {
		t.Fatalf("chart fixture count mismatch: expected %d, got %d", len(expected), len(got))
	}

	for _, chart := range charts {
		want, ok := expected[chart]
		if !ok {
			t.Fatalf("missing chart fixture entry for %q", chart)
		}
		if got[chart] != want {
			t.Fatalf(
				"visual regression for %q: expected %s, got %s (run UPDATE_CHART_VISUAL_FIXTURES=1 go test ./pkg/pptx -run TestChartVisualRegressionFixtures -count=1 to refresh fixtures)",
				chart,
				want,
				got[chart],
			)
		}
	}

	for chart := range expected {
		if _, ok := got[chart]; !ok {
			t.Fatalf("fixture has unexpected chart entry %q", chart)
		}
	}
}

func chartVisualRegressionSlides() map[string]pptx.SlideContent {
	return map[string]pptx.SlideContent{
		"area": pptx.NewSlide("Area").WithAreaChart(
			charts.NewAreaChart([]string{"Q1", "Q2", "Q3"}, []float64{14, 17, 23}).WithTitle("Area"),
		),
		"areaStacked": pptx.NewSlide("Area Stacked").WithAreaStackedChart(
			charts.NewAreaStackedChart([]string{"Q1", "Q2", "Q3"}, []float64{14, 17, 23}).WithTitle("Area Stacked"),
		),
		"areaStacked100": pptx.NewSlide("Area Stacked 100").WithAreaStacked100Chart(
			charts.NewAreaStacked100Chart([]string{"Q1", "Q2", "Q3"}, []float64{14, 17, 23}).
				WithTitle("Area Stacked 100"),
		),
		"bar": pptx.NewSlide("Bar").WithBarChart(
			charts.NewBarChart([]string{"Q1", "Q2", "Q3"}, []float64{12, 18, 24}).WithTitle("Bar"),
		),
		"barHorizontal": pptx.NewSlide("Bar Horizontal").WithBarHorizontalChart(
			charts.NewBarHorizontalChart([]string{"Q1", "Q2", "Q3"}, []float64{12, 18, 24}).WithTitle("Bar Horizontal"),
		),
		"barStacked": pptx.NewSlide("Bar Stacked").WithBarStackedChart(
			charts.NewBarStackedChart([]string{"Q1", "Q2", "Q3"}, []float64{12, 18, 24}).WithTitle("Bar Stacked"),
		),
		"barStacked100": pptx.NewSlide("Bar Stacked 100").WithBarStacked100Chart(
			charts.NewBarStacked100Chart([]string{"Q1", "Q2", "Q3"}, []float64{12, 18, 24}).
				WithTitle("Bar Stacked 100"),
		),
		"bubble": pptx.NewSlide("Bubble").WithBubbleChart(
			charts.NewBubbleChart([]float64{1, 2, 3}, []float64{10, 20, 30}, []float64{10, 20, 30}).
				WithTitle("Bubble").WithSeriesName("Series 1").WithBubbleScale(100),
		),
		"combo": pptx.NewSlide("Combo").WithComboChart(
			charts.NewComboChart(
				[]string{"Q1", "Q2", "Q3"},
				[]charts.Series{{Name: "Bar A", Values: []float64{1, 2, 3}}},
				[]charts.Series{{Name: "Line A", Values: []float64{2, 3, 4}}},
			).WithTitle("Combo"),
		),
		"doughnut": pptx.NewSlide("Doughnut").WithDoughnutChart(
			charts.NewDoughnutChart([]string{"A", "B", "C"}, []float64{30, 45, 25}).WithTitle("Doughnut"),
		),
		"line": pptx.NewSlide("Line").WithLineChart(
			charts.NewLineChart([]string{"Q1", "Q2", "Q3"}, []float64{10, 16, 22}).WithTitle("Line"),
		),
		"lineMarkers": pptx.NewSlide("Line Markers").WithLineMarkersChart(
			charts.NewLineMarkersChart([]string{"Q1", "Q2", "Q3"}, []float64{10, 16, 22}).WithTitle("Line Markers"),
		),
		"lineStacked": pptx.NewSlide("Line Stacked").WithLineStackedChart(
			charts.NewLineStackedChart([]string{"Q1", "Q2", "Q3"}, []float64{10, 16, 22}).WithTitle("Line Stacked"),
		),
		"pie": pptx.NewSlide("Pie").WithPieChart(
			charts.NewPieChart([]string{"A", "B", "C"}, []float64{30, 45, 25}).WithTitle("Pie"),
		),
		"radar": pptx.NewSlide("Radar").WithRadarChart(
			charts.NewRadarChart([]string{"A", "B", "C"}, []float64{2, 3, 4}).WithTitle("Radar"),
		),
		"radarFilled": pptx.NewSlide("Radar Filled").WithRadarFilledChart(
			charts.NewRadarFilledChart([]string{"A", "B", "C"}, []float64{3, 4, 5}).WithTitle("Radar Filled"),
		),
		"scatter": pptx.NewSlide("Scatter Marker").WithScatterChart(
			charts.NewScatterChart([]float64{1, 2, 3}, []float64{10, 15, 20}).
				WithTitle("Scatter Marker").WithScatterStyle(charts.ScatterStyleMarker),
		),
		"scatterLines": pptx.NewSlide("Scatter Lines").WithScatterChart(
			charts.NewScatterChart([]float64{1, 2, 3}, []float64{10, 15, 20}).
				WithTitle("Scatter Lines").WithScatterStyle(charts.ScatterStyleLineMarker),
		),
		"scatterSmooth": pptx.NewSlide("Scatter Smooth").WithScatterChart(
			charts.NewScatterChart([]float64{1, 2, 3}, []float64{10, 15, 20}).
				WithTitle("Scatter Smooth").WithScatterStyle(charts.ScatterStyleSmoothMarker),
		),
		"stockHLC": pptx.NewSlide("StockHLC").WithStockHLCChart(
			charts.NewStockHLCChart(
				[]string{"D1", "D2", "D3"},
				[]float64{12, 13, 14},
				[]float64{8, 9, 10},
				[]float64{10, 11, 12},
			).WithTitle("StockHLC"),
		),
		"stockOHLC": pptx.NewSlide("StockOHLC").WithStockOHLCChart(
			charts.NewStockOHLCChart(
				[]string{"D1", "D2", "D3"},
				[]float64{9, 10, 11},
				[]float64{12, 13, 14},
				[]float64{8, 9, 10},
				[]float64{10, 11, 12},
			).WithTitle("StockOHLC"),
		),
	}
}

func loadChartVisualFixtures(t *testing.T) map[string]string {
	t.Helper()

	data, err := os.ReadFile(chartVisualFixturePath())
	if err != nil {
		t.Fatalf("read chart visual fixtures: %v", err)
	}

	var entries []chartVisualFixtureEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("decode chart visual fixtures: %v", err)
	}
	if len(entries) == 0 {
		t.Fatalf("chart visual fixtures file is empty")
	}

	out := make(map[string]string, len(entries))
	for _, entry := range entries {
		if entry.Chart == "" || entry.Fingerprint == "" {
			t.Fatalf("invalid chart fixture entry: %+v", entry)
		}
		if _, exists := out[entry.Chart]; exists {
			t.Fatalf("duplicate chart fixture entry for %q", entry.Chart)
		}
		out[entry.Chart] = entry.Fingerprint
	}
	return out
}

func writeChartVisualFixtures(t *testing.T, fixtures map[string]string) {
	t.Helper()

	keys := sortedSlideKeys(fixtures)
	entries := make([]chartVisualFixtureEntry, 0, len(keys))
	for _, key := range keys {
		entries = append(entries, chartVisualFixtureEntry{
			Chart:       key,
			Fingerprint: fixtures[key],
		})
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		t.Fatalf("encode chart visual fixtures: %v", err)
	}
	data = append(data, '\n')

	path := chartVisualFixturePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create chart visual fixtures dir: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write chart visual fixtures: %v", err)
	}
}

func chartVisualFingerprint(xml string) string {
	normalized := strings.TrimSpace(strings.ReplaceAll(xml, "\r\n", "\n"))
	normalized = strings.Join(strings.Fields(normalized), " ")
	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}

func chartVisualFixturePath() string {
	return filepath.Join("fixtures", "chart_visual_fingerprints.json")
}

func sortedSlideKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
