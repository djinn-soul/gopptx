package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestExcelRowHelpers(t *testing.T) {
	req := common.ChartDataUpdate{
		MultiLevelCategories: [][]string{{"A", "B"}, {"A1", "B1"}},
		Series:               []common.ChartSeriesData{{Values: []float64{1, 2}}},
	}
	headers := categoryHeaders(req)
	if len(headers) != 2 || headers[0] != "Category Level 1" {
		t.Fatalf("unexpected category headers: %+v", headers)
	}

	rows := buildCategoryRows(nil, req.MultiLevelCategories, req.Series)
	if len(rows) != 2 || len(rows[0]) != 3 {
		t.Fatalf("unexpected category rows shape: %+v", rows)
	}
	if rows[0][0] != "A" || rows[0][1] != "A1" || rows[0][2] != "1" {
		t.Fatalf("unexpected first category row: %+v", rows[0])
	}

	series := []common.ChartSeriesData{
		{XValues: []float64{1, 2}, YValues: []float64{3, 4}, Sizes: []float64{5, 6}},
	}
	scatterHeaders, scatterRows := buildScatterSheet(series, true)
	if len(scatterHeaders) != 3 || scatterHeaders[2] != "S1" {
		t.Fatalf("unexpected scatter headers: %+v", scatterHeaders)
	}
	if len(scatterRows) != 2 || len(scatterRows[0]) != 3 {
		t.Fatalf("unexpected scatter rows shape: %+v", scatterRows)
	}
	if scatterRows[1][0] != "2" || scatterRows[1][1] != "4" || scatterRows[1][2] != "6" {
		t.Fatalf("unexpected second scatter row: %+v", scatterRows[1])
	}

	if got := ColumnName(27); got != "AA" {
		t.Fatalf("ColumnName(27)=%q, want AA", got)
	}
	if !isNumberLiteral("3.14") || isNumberLiteral("x") {
		t.Fatal("isNumberLiteral classification mismatch")
	}
}

func TestSelectorAndPayloadValidation(t *testing.T) {
	refs := []common.SlideChartRef{
		{Index: 0, RelID: "rId1", ChartPart: "ppt/charts/chart1.xml"},
		{Index: 1, RelID: "rId2", ChartPart: "ppt/charts/chart2.xml"},
	}
	idx := 1
	selected, err := ResolveChartSelector(refs, common.ChartSelector{Index: &idx}, 0)
	if err != nil || selected.RelID != "rId2" {
		t.Fatalf("ResolveChartSelector(index) failed: selected=%+v err=%v", selected, err)
	}
	selected, err = ResolveChartSelector(refs, common.ChartSelector{RelID: "rId1"}, 0)
	if err != nil || selected.Index != 0 {
		t.Fatalf("ResolveChartSelector(relID) failed: selected=%+v err=%v", selected, err)
	}
	rel := "rId1"
	selected, err = ResolveChartSelector(refs, common.ChartSelector{Index: &[]int{0}[0], RelID: rel}, 0)
	if err != nil || selected.RelID != rel {
		t.Fatalf("ResolveChartSelector(index+relID) failed: selected=%+v err=%v", selected, err)
	}
	if _, err = ResolveChartSelector(refs, common.ChartSelector{}, 0); err == nil {
		t.Fatal("expected selector missing criteria error")
	}
	badIdx := 4
	if _, err = ResolveChartSelector(refs, common.ChartSelector{Index: &badIdx}, 0); err == nil {
		t.Fatal("expected selector index out-of-range error")
	}

	if kind := DetectChartKind([]byte("<c:scatterChart/>")); kind != KindScatter {
		t.Fatalf("DetectChartKind scatter=%v", kind)
	}
	if kind := DetectChartKind([]byte("<c:bubbleChart/>")); kind != KindBubble {
		t.Fatalf("DetectChartKind bubble=%v", kind)
	}
	if kind := DetectChartKind([]byte("<c:barChart/>")); kind != KindCategory {
		t.Fatalf("DetectChartKind category=%v", kind)
	}

	if err := ValidateChartUpdatePayload(KindCategory, common.ChartDataUpdate{
		Categories: []string{"A"},
		Series:     []common.ChartSeriesData{{Values: []float64{1}}},
	}); err != nil {
		t.Fatalf("category payload should be valid: %v", err)
	}
	if err := ValidateChartUpdatePayload(KindScatter, common.ChartDataUpdate{
		Series: []common.ChartSeriesData{{XValues: []float64{1}, YValues: []float64{2}}},
	}); err != nil {
		t.Fatalf("scatter payload should be valid: %v", err)
	}
	if err := ValidateChartUpdatePayload(KindBubble, common.ChartDataUpdate{
		Series: []common.ChartSeriesData{{XValues: []float64{1}, YValues: []float64{2}, Sizes: []float64{3}}},
	}); err != nil {
		t.Fatalf("bubble payload should be valid: %v", err)
	}
	if err := ValidateChartUpdatePayload(KindScatter, common.ChartDataUpdate{
		Series: []common.ChartSeriesData{{XValues: []float64{1}, YValues: []float64{}}},
	}); err == nil {
		t.Fatal("expected invalid scatter payload")
	}
}

func TestPlaceholderDefinitionAndXMLHelpers(t *testing.T) {
	def, err := PlaceholderChartDefinition(
		map[string]any{
			"categories": []any{"A", "B"},
			"values":     []any{1.0, 2.0},
		},
		"bar",
		"Sales",
		10, 20, 100, 200,
	)
	if err != nil || def == nil {
		t.Fatalf("PlaceholderChartDefinition(bar) failed: def=%T err=%v", def, err)
	}
	def, err = PlaceholderChartDefinition(
		map[string]any{
			"x_values": []any{1.0, 2.0},
			"y_values": []any{3.0, 4.0},
		},
		"scatter",
		"Scatter",
		0, 0, 100, 200,
	)
	if err != nil || def == nil {
		t.Fatalf("PlaceholderChartDefinition(scatter) failed: def=%T err=%v", def, err)
	}
	if _, err = PlaceholderChartDefinition(map[string]any{}, "bar", "", 0, 0, 0, 0); err == nil {
		t.Fatal("expected category parsing error")
	}
	if _, err = PlaceholderChartDefinition(
		map[string]any{"categories": []any{"A"}, "values": []any{1.0}},
		"unknown",
		"",
		0, 0, 0, 0,
	); err == nil {
		t.Fatal("expected unsupported chart type error")
	}

	if got := boolToOneZero(true); got != "1" {
		t.Fatalf("boolToOneZero(true)=%q", got)
	}
	if got := boolToOneZero(false); got != "0" {
		t.Fatalf("boolToOneZero(false)=%q", got)
	}
	escaped := xmlEscape(`<a&"b">`)
	if !strings.Contains(escaped, "&lt;a&amp;&quot;b&quot;&gt;") {
		t.Fatalf("unexpected xmlEscape output: %q", escaped)
	}
}
