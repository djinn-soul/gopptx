package pptx

import (
	"reflect"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/enums"
)

func TestNewChartBuildsEveryCategoryChartType(t *testing.T) {
	categories := []string{"Q1", "Q2", "Q3"}
	values := []float64{10, 20, 30}

	tests := []struct {
		chartType enums.XLChartType
		wantKind  string
		wantType  Chart
	}{
		{enums.XLChartTypeBar, "bar", charts.BarChart{}},
		{enums.XLChartTypeBarHorizontal, "barHorizontal", charts.BarHorizontalChart{}},
		{enums.XLChartTypeBarStacked, "barStacked", charts.BarStackedChart{}},
		{enums.XLChartTypeBarStacked100, "barStacked100", charts.BarStacked100Chart{}},
		{enums.XLChartTypeLine, "line", charts.LineChart{}},
		{enums.XLChartTypeLineMarkers, "lineMarkers", charts.LineMarkersChart{}},
		{enums.XLChartTypeLineStacked, "lineStacked", charts.LineStackedChart{}},
		{enums.XLChartTypeArea, "area", charts.AreaChart{}},
		{enums.XLChartTypeAreaStacked, "areaStacked", charts.AreaStackedChart{}},
		{enums.XLChartTypeAreaStacked100, "areaStacked100", charts.AreaStacked100Chart{}},
		{enums.XLChartTypePie, "pie", charts.PieChart{}},
		{enums.XLChartTypeDoughnut, "doughnut", charts.DoughnutChart{}},
		{enums.XLChartTypeRadar, "radar", charts.RadarChart{}},
		{enums.XLChartTypeRadarFilled, "radarFilled", charts.RadarFilledChart{}},
	}

	for _, tt := range tests {
		t.Run(string(tt.chartType), func(t *testing.T) {
			chart, err := NewChart(tt.chartType, categories, values)
			if err != nil {
				t.Fatalf("NewChart(%q): %v", tt.chartType, err)
			}
			if got := chart.ChartKind(); got != tt.wantKind {
				t.Errorf("ChartKind() = %q, want %q", got, tt.wantKind)
			}
			if gotType, wantType := typeName(chart), typeName(tt.wantType); gotType != wantType {
				t.Errorf("concrete type = %s, want %s", gotType, wantType)
			}
		})
	}
}

func TestNewChartMatchesDedicatedConstructor(t *testing.T) {
	categories := []string{"A", "B"}
	values := []float64{1, 2}

	chart, err := NewChart(enums.XLChartTypeBar, categories, values)
	if err != nil {
		t.Fatalf("NewChart: %v", err)
	}
	got, ok := chart.(charts.BarChart)
	if !ok {
		t.Fatalf("NewChart returned %s, want charts.BarChart", typeName(chart))
	}
	want := charts.NewBarChart(categories, values)

	if gotCats, wantCats := got.GetCategories(), want.GetCategories(); !equalStrings(gotCats, wantCats) {
		t.Errorf("categories = %v, want %v", gotCats, wantCats)
	}
	if gotVals, wantVals := got.GetValues(), want.GetValues(); !equalFloats(gotVals, wantVals) {
		t.Errorf("values = %v, want %v", gotVals, wantVals)
	}
}

func TestNewChartRejectsNonCategoryChartTypes(t *testing.T) {
	tests := []struct {
		chartType    enums.XLChartType
		wantCtorHint string
	}{
		{enums.XLChartTypeScatter, "NewScatterChart"},
		{enums.XLChartTypeBubble, "NewBubbleChart"},
		{enums.XLChartTypeStockHLC, "NewStockHLCChart"},
		{enums.XLChartTypeStockOHLC, "NewStockOHLCChart"},
		{enums.XLChartTypeCombo, "NewComboChart"},
	}

	for _, tt := range tests {
		t.Run(string(tt.chartType), func(t *testing.T) {
			chart, err := NewChart(tt.chartType, []string{"A"}, []float64{1})
			if err == nil {
				t.Fatalf("NewChart(%q) = %v, want error", tt.chartType, chart)
			}
			if chart != nil {
				t.Errorf("chart = %v, want nil alongside error", chart)
			}
			if !strings.Contains(err.Error(), tt.wantCtorHint) {
				t.Errorf("error %q does not point at %s", err, tt.wantCtorHint)
			}
		})
	}
}

func TestNewChartRejectsUnknownChartType(t *testing.T) {
	chart, err := NewChart(enums.XLChartType("spiral"), []string{"A"}, []float64{1})
	if err == nil {
		t.Fatalf("NewChart = %v, want error", chart)
	}
	if !strings.Contains(err.Error(), "spiral") {
		t.Errorf("error %q does not name the rejected type", err)
	}
}

// WithChart must land in the same slide field as the type-specific setter, so
// that the generic and the explicit style produce identical slides.
func TestWithChartMatchesTypeSpecificSetter(t *testing.T) {
	categories := []string{"A", "B"}
	values := []float64{3, 4}

	pie := charts.NewPieChart(categories, values).WithTitle("Share")

	generic := NewSlide("Chart").WithChart(pie)
	explicit := NewSlide("Chart").WithPieChart(pie)

	if generic.Pie == nil {
		t.Fatal("WithChart did not set the Pie field")
	}
	if !reflect.DeepEqual(*generic.Pie, *explicit.Pie) {
		t.Errorf("WithChart produced %+v, want %+v", *generic.Pie, *explicit.Pie)
	}
	if generic.Chart != nil || generic.Line != nil {
		t.Error("WithChart set a field belonging to another chart type")
	}
}

// WithChart replaces a chart already on the slide rather than leaving two set.
func TestWithChartReplacesExistingChart(t *testing.T) {
	bar := charts.NewBarChart([]string{"A"}, []float64{1})
	pie := charts.NewPieChart([]string{"B"}, []float64{2})

	slide := NewSlide("Chart").WithChart(bar).WithChart(pie)

	if slide.Chart != nil {
		t.Error("bar chart survived being replaced")
	}
	if slide.Pie == nil {
		t.Error("pie chart was not set")
	}
}

func TestWithChartIgnoresNilChart(t *testing.T) {
	bar := charts.NewBarChart([]string{"A"}, []float64{1})
	slide := NewSlide("Chart").WithBarChart(bar)

	if got := slide.WithChart(nil); got.Chart == nil {
		t.Error("WithChart(nil) cleared the existing chart, want it left unchanged")
	}
}

func typeName(v any) string {
	if v == nil {
		return "<nil>"
	}
	return reflect.TypeOf(v).String()
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalFloats(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
