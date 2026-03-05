package pptx

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func TestPresentation_UpdateChartDataByIndexFromBuilder_Category(t *testing.T) {
	pptxPath := createPresentationWithBarChart(t)

	prs, err := Open(pptxPath)
	if err != nil {
		t.Fatalf("open presentation: %v", err)
	}

	builder := NewCategoryChartData([]string{"Q1", "Q2"}).
		AddSeries("Revenue", []float64{12, 24})
	if err := prs.UpdateChartDataByIndexFromBuilder(0, 0, builder); err != nil {
		_ = prs.Close()
		t.Fatalf("UpdateChartDataByIndexFromBuilder failed: %v", err)
	}
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save after category builder update failed: %v", err)
	}
	_ = prs.Close()

	chartXML := readZipEntry(t, pptxPath)
	if !strings.Contains(chartXML, "<c:v>Q1</c:v>") {
		t.Fatalf("updated chart XML missing category %q", "Q1")
	}
	if !strings.Contains(chartXML, "<c:v>24</c:v>") {
		t.Fatalf("updated chart XML missing value %q", "24")
	}
}

func TestPresentation_UpdateChartDataByRelIDFromBuilder_XY(t *testing.T) {
	pptxPath := createPresentationWithChart(t, charts.NewScatterChart(
		[]float64{1, 2},
		[]float64{10, 20},
	).WithTitle("Scatter Fixture"))

	prs, err := Open(pptxPath)
	if err != nil {
		t.Fatalf("open presentation: %v", err)
	}

	refs, err := prs.ListSlideCharts(0)
	if err != nil {
		_ = prs.Close()
		t.Fatalf("ListSlideCharts failed: %v", err)
	}
	if len(refs) != 1 {
		_ = prs.Close()
		t.Fatalf("expected 1 chart ref, got %d", len(refs))
	}

	builder := NewXyChartData().AddSeries("S1", []float64{3, 4}, []float64{30, 40})
	if err := prs.UpdateChartDataByRelIDFromBuilder(0, refs[0].RelID, builder); err != nil {
		_ = prs.Close()
		t.Fatalf("UpdateChartDataByRelIDFromBuilder (xy) failed: %v", err)
	}
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save after xy builder update failed: %v", err)
	}
	_ = prs.Close()

	chartXML := readZipEntry(t, pptxPath)
	if !strings.Contains(chartXML, "<c:scatterChart>") {
		t.Fatal("expected scatter chart XML")
	}
	if !strings.Contains(chartXML, "<c:v>40</c:v>") {
		t.Fatalf("updated scatter chart XML missing value %q", "40")
	}
}

func TestPresentation_UpdateChartDataByRelIDFromBuilder_Bubble(t *testing.T) {
	pptxPath := createPresentationWithChart(t, charts.NewBubbleChart(
		[]float64{1, 2},
		[]float64{10, 20},
		[]float64{5, 6},
	).WithTitle("Bubble Fixture"))

	prs, err := Open(pptxPath)
	if err != nil {
		t.Fatalf("open presentation: %v", err)
	}

	refs, err := prs.ListSlideCharts(0)
	if err != nil {
		_ = prs.Close()
		t.Fatalf("ListSlideCharts failed: %v", err)
	}
	if len(refs) != 1 {
		_ = prs.Close()
		t.Fatalf("expected 1 chart ref, got %d", len(refs))
	}

	builder := NewBubbleChartData().AddSeries(
		"B1",
		[]float64{3, 4},
		[]float64{30, 40},
		[]float64{7, 8},
	)
	if err := prs.UpdateChartDataByRelIDFromBuilder(0, refs[0].RelID, builder); err != nil {
		_ = prs.Close()
		t.Fatalf("UpdateChartDataByRelIDFromBuilder (bubble) failed: %v", err)
	}
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save after bubble builder update failed: %v", err)
	}
	_ = prs.Close()

	chartXML := readZipEntry(t, pptxPath)
	if !strings.Contains(chartXML, "<c:bubbleChart>") {
		t.Fatal("expected bubble chart XML")
	}
	if !strings.Contains(chartXML, "<c:v>8</c:v>") {
		t.Fatalf("updated bubble chart XML missing bubble size %q", "8")
	}
}

func TestChartDataBuilderValidationErrors(t *testing.T) {
	pptxPath := createPresentationWithBarChart(t)

	prs, err := Open(pptxPath)
	if err != nil {
		t.Fatalf("open presentation: %v", err)
	}
	defer prs.Close()

	if err := prs.UpdateChartDataByIndexFromBuilder(0, 0, nil); err == nil {
		t.Fatal("expected error for nil builder")
	}

	emptyCategory := NewCategoryChartData([]string{"A"})
	if err := prs.UpdateChartDataByIndexFromBuilder(0, 0, emptyCategory); err == nil {
		t.Fatal("expected error for category builder without series")
	}

	invalidXY := NewXyChartData().AddSeries("bad", []float64{1}, []float64{2, 3})
	if err := prs.UpdateChartDataByIndexFromBuilder(0, 0, invalidXY); err == nil {
		t.Fatal("expected error for mismatched xy series lengths")
	}

	invalidBubble := NewBubbleChartData().AddSeries(
		"bad",
		[]float64{1},
		[]float64{2},
		[]float64{3, 4},
	)
	if err := prs.UpdateChartDataByIndexFromBuilder(0, 0, invalidBubble); err == nil {
		t.Fatal("expected error for mismatched bubble series lengths")
	}

	multiLevelMismatch := NewCategoryChartData(nil).
		AddCategoryLevel([]string{"2026", "2026"}).
		AddCategoryLevel([]string{"Q1"}).
		AddSeries("Revenue", []float64{1, 2})
	if err := prs.UpdateChartDataByIndexFromBuilder(0, 0, multiLevelMismatch); err == nil {
		t.Fatal("expected error for mismatched multi-level category lengths")
	}
}

func TestCategoryChartDataBuilderMultiLevelPayload(t *testing.T) {
	builder := NewCategoryChartData(nil).
		AddCategoryLevel([]string{"2026", "2026"}).
		AddCategoryLevel([]string{"Q1", "Q2"}).
		AddSeries("Revenue", []float64{10, 20})

	update, err := builder.chartDataUpdate()
	if err != nil {
		t.Fatalf("chartDataUpdate failed: %v", err)
	}
	if len(update.MultiLevelCategories) != 2 {
		t.Fatalf("expected 2 category levels, got %d", len(update.MultiLevelCategories))
	}
	if update.MultiLevelCategories[1][1] != "Q2" {
		t.Fatalf("unexpected multi-level category value: %+v", update.MultiLevelCategories)
	}
	if len(update.Series) != 1 || len(update.Series[0].Values) != 2 {
		t.Fatalf("unexpected series payload: %+v", update.Series)
	}
}
