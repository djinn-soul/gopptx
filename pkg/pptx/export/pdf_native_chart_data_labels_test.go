package export

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func TestFormatChartValue(t *testing.T) {
	cases := []struct {
		name   string
		value  float64
		format string
		want   string
	}{
		{"general trims zeros", 12.5, "", "12.5"},
		{"general whole number", 12, "General", "12"},
		{"general keeps precision", 12.345, "General", "12.345"},
		{"fixed decimals", 12.345, "0.00", "12.35"},
		{"no decimals rounds", 12.6, "0", "13"},
		{"percent scales", 0.256, "0.0%", "25.6%"},
		{"currency prefix", 1234.5, "$#,##0.00", "$1,234.50"},
		{"thousands grouping", 1234567, "#,##0", "1,234,567"},
		{"negative grouping", -1234567, "#,##0", "-1,234,567"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatChartValue(tc.value, tc.format); got != tc.want {
				t.Fatalf("formatChartValue(%v, %q)=%q want %q", tc.value, tc.format, got, tc.want)
			}
		})
	}
}

func TestChartDataLabelTextDefaultsToValueOnly(t *testing.T) {
	opts := chartSeriesOpts{}
	got := chartDataLabelText(opts, chartDataLabel{category: "Q1", value: 42})
	if got != "42" {
		t.Fatalf("default label=%q want %q", got, "42")
	}
}

func TestChartDataLabelTextComposesCustomSettings(t *testing.T) {
	opts := chartSeriesOpts{
		valueFormat: "0.0",
		seriesName:  "Revenue",
		dataLabels: charts.DataLabelSettings{
			UseCustom:      true,
			ShowSeriesName: true,
			ShowCategory:   true,
			ShowValue:      true,
			ShowPercent:    true,
		},
	}
	label := chartDataLabel{category: "Q1", seriesName: "Revenue", value: 25, total: 100}
	want := "Revenue, Q1, 25.0, 25%"
	if got := chartDataLabelText(opts, label); got != want {
		t.Fatalf("custom label=%q want %q", got, want)
	}
}

func TestChartDataLabelTextOmitsPercentWithoutTotal(t *testing.T) {
	opts := chartSeriesOpts{
		dataLabels: charts.DataLabelSettings{UseCustom: true, ShowValue: true, ShowPercent: true},
	}
	if got := chartDataLabelText(opts, chartDataLabel{value: 7}); got != "7" {
		t.Fatalf("label=%q want %q", got, "7")
	}
}

func TestPieSliceLabelTextKeepsLegacyDefaults(t *testing.T) {
	categories := []string{"A", "B"}
	withCat := chartSeriesOpts{showCatName: true}
	if got := pieSliceLabelText(withCat, categories, 0, 25, 100, 0.25); got != "A" {
		t.Fatalf("category label=%q want %q", got, "A")
	}
	withoutCat := chartSeriesOpts{}
	if got := pieSliceLabelText(withoutCat, categories, 0, 25, 100, 0.25); got != "25%" {
		t.Fatalf("percent label=%q want %q", got, "25%")
	}
}

func TestPieSliceLabelTextUsesCustomSettings(t *testing.T) {
	opts := chartSeriesOpts{
		dataLabels: charts.DataLabelSettings{UseCustom: true, ShowCategory: true, ShowPercent: true},
	}
	want := "A, 25%"
	if got := pieSliceLabelText(opts, []string{"A"}, 0, 25, 100, 0.25); got != want {
		t.Fatalf("custom pie label=%q want %q", got, want)
	}
}

func TestBarLabelYHonoursPosition(t *testing.T) {
	const barTop, barBottom = 100.0, 180.0
	custom := func(position string) chartSeriesOpts {
		return chartSeriesOpts{dataLabels: charts.DataLabelSettings{UseCustom: true, Position: position}}
	}

	if got := barLabelY(chartSeriesOpts{}, barTop, barBottom, false); got >= barTop {
		t.Fatalf("default label y=%v want above bar top %v", got, barTop)
	}
	if got := barLabelY(custom(charts.DataLabelPositionCenter), barTop, barBottom, false); got != 140 {
		t.Fatalf("centre label y=%v want 140", got)
	}
	if got := barLabelY(custom(charts.DataLabelPositionInsideEnd), barTop, barBottom, false); got <= barTop {
		t.Fatalf("inside-end label y=%v want below bar top %v", got, barTop)
	}
	if got := barLabelY(custom(charts.DataLabelPositionInsideBase), barTop, barBottom, false); got >= barBottom {
		t.Fatalf("inside-base label y=%v want above bar bottom %v", got, barBottom)
	}
}

func TestBarLabelYFlipsForNegativeBars(t *testing.T) {
	const barTop, barBottom = 100.0, 180.0
	got := barLabelY(chartSeriesOpts{}, barTop, barBottom, true)
	if got <= barBottom {
		t.Fatalf("negative bar label y=%v want below bar bottom %v", got, barBottom)
	}
}

func TestDataLabelWrapsUnlessTurnedOff(t *testing.T) {
	if !dataLabelWraps(chartSeriesOpts{}) {
		t.Fatal("a label with no c:dLbls should wrap, matching PowerPoint's default")
	}
	off := false
	opts := chartSeriesOpts{dataLabels: charts.DataLabelSettings{WordWrap: &off}}
	if dataLabelWraps(opts) {
		t.Fatal("wrap=0 not honoured")
	}
}

func TestDataLabelKeyWidthOnlyWithLegendKey(t *testing.T) {
	if got := dataLabelKeyWidth(chartSeriesOpts{}); got != 0 {
		t.Fatalf("key width=%v want 0 without c:showLegendKey", got)
	}
	// The flag only counts as part of a custom c:dLbls.
	ignored := chartSeriesOpts{dataLabels: charts.DataLabelSettings{ShowLegendKey: true}}
	if got := dataLabelKeyWidth(ignored); got != 0 {
		t.Fatalf("key width=%v want 0 when the settings are not custom", got)
	}
	withKey := chartSeriesOpts{
		dataLabels: charts.DataLabelSettings{UseCustom: true, ShowLegendKey: true},
	}
	if got := dataLabelKeyWidth(withKey); got <= 0 {
		t.Fatalf("key width=%v want room for the swatch", got)
	}
}

func TestClampDataLabelToPlot(t *testing.T) {
	unbounded := chartSeriesOpts{}
	if x, y := clampDataLabelToPlot(unbounded, -50, -50, 10, 10); x != -50 || y != -50 {
		t.Fatalf("clamp without a plot rect moved the label to %v,%v", x, y)
	}

	opts := withChartPlotArea(chartSeriesOpts{}, 100, 200, 300, 150)
	if x, _ := clampDataLabelToPlot(opts, 40, 250, 60, 10); x != 100 {
		t.Fatalf("left overflow x=%v want 100", x)
	}
	if x, _ := clampDataLabelToPlot(opts, 380, 250, 60, 10); x != 340 {
		t.Fatalf("right overflow x=%v want 340", x)
	}
	if _, y := clampDataLabelToPlot(opts, 150, 190, 20, 10); y != 200 {
		t.Fatalf("top overflow y=%v want 200", y)
	}
	if _, y := clampDataLabelToPlot(opts, 150, 348, 20, 10); y != 340 {
		t.Fatalf("bottom overflow y=%v want 340", y)
	}
	// A label that already fits is left alone.
	if x, y := clampDataLabelToPlot(opts, 150, 250, 20, 10); x != 150 || y != 250 {
		t.Fatalf("fitting label moved to %v,%v", x, y)
	}
}

func TestWithChartPlotAreaSetsWrapWidth(t *testing.T) {
	opts := withChartLabelData(chartSeriesOpts{}, []string{"a", "b", "c", "d"}, []float64{1, 2, 3, 4})
	opts = withChartPlotArea(opts, 0, 0, 400, 100)
	if !opts.hasPlot {
		t.Fatal("plot rect not recorded")
	}
	if opts.labelWrapWidth != 100 {
		t.Fatalf("wrap width=%v want one category slot (100)", opts.labelWrapWidth)
	}

	// Many categories must not squeeze the wrap width to nothing.
	narrow := withChartLabelData(chartSeriesOpts{}, make([]string, 100), nil)
	narrow = withChartPlotArea(narrow, 0, 0, 400, 100)
	if narrow.labelWrapWidth != minDataLabelWrapWidthPt {
		t.Fatalf("wrap width=%v want the floor %v", narrow.labelWrapWidth, minDataLabelWrapWidthPt)
	}
}

func TestDataLabelPointAnchorHonoursPosition(t *testing.T) {
	const x, y = 100.0, 100.0
	custom := func(position string) chartSeriesOpts {
		return chartSeriesOpts{dataLabels: charts.DataLabelSettings{UseCustom: true, Position: position}}
	}

	if _, ay := dataLabelPointAnchor(chartSeriesOpts{}, x, y); ay >= y {
		t.Fatalf("default anchor y=%v want above the point", ay)
	}
	if _, ay := dataLabelPointAnchor(custom(charts.DataLabelPositionBottom), x, y); ay <= y {
		t.Fatalf("bottom anchor y=%v want below the point", ay)
	}
	if ax, _ := dataLabelPointAnchor(custom(charts.DataLabelPositionLeft), x, y); ax >= x {
		t.Fatalf("left anchor x=%v want left of the point", ax)
	}
	if ax, _ := dataLabelPointAnchor(custom(charts.DataLabelPositionRight), x, y); ax <= x {
		t.Fatalf("right anchor x=%v want right of the point", ax)
	}
}

func TestSumChartValuesUsesAbsolutes(t *testing.T) {
	if got := sumChartValues([]float64{3, -4, 5}); got != 12 {
		t.Fatalf("sumChartValues=%v want 12", got)
	}
}
