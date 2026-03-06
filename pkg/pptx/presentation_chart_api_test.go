package pptx

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func TestPresentation_ListSlideCharts(t *testing.T) {
	pptxPath := createPresentationWithBarChart(t)

	prs, err := Open(pptxPath)
	if err != nil {
		t.Fatalf("open presentation: %v", err)
	}
	defer prs.Close()

	refs, err := prs.ListSlideCharts(0)
	if err != nil {
		t.Fatalf("ListSlideCharts failed: %v", err)
	}
	if len(refs) != 1 {
		t.Fatalf("expected 1 chart ref, got %d", len(refs))
	}
	if refs[0].Index != 0 {
		t.Fatalf("expected chart index 0, got %d", refs[0].Index)
	}
	if refs[0].RelID == "" {
		t.Fatal("expected non-empty chart rel_id")
	}
	if refs[0].ChartPart == "" {
		t.Fatal("expected non-empty chart part")
	}
}

func TestPresentation_UpdateChartData(t *testing.T) {
	pptxPath := createPresentationWithBarChart(t)

	prs, err := Open(pptxPath)
	if err != nil {
		t.Fatalf("open presentation: %v", err)
	}

	err = prs.UpdateChartDataByIndex(0, 0, ChartDataUpdate{
		Categories: []string{"North", "South"},
		Series: []ChartSeriesData{
			{Values: []float64{33, 77}},
		},
	})
	if err != nil {
		_ = prs.Close()
		t.Fatalf("UpdateChartDataByIndex failed: %v", err)
	}
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save after UpdateChartDataByIndex failed: %v", err)
	}
	_ = prs.Close()

	chartXML := readZipEntry(t, pptxPath)
	if !strings.Contains(chartXML, "<c:v>North</c:v>") {
		t.Fatalf("updated chart XML missing category %q", "North")
	}
	if !strings.Contains(chartXML, "<c:v>77</c:v>") {
		t.Fatalf("updated chart XML missing value %q", "77")
	}
}

func TestPresentation_ReplaceChartData(t *testing.T) {
	pptxPath := createPresentationWithBarChart(t)

	prs, err := Open(pptxPath)
	if err != nil {
		t.Fatalf("open presentation: %v", err)
	}

	err = prs.ReplaceChartData(0, 0, []string{"Only"}, []float64{42})
	if err != nil {
		_ = prs.Close()
		t.Fatalf("ReplaceChartData failed: %v", err)
	}
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save after ReplaceChartData failed: %v", err)
	}
	_ = prs.Close()

	chartXML := readZipEntry(t, pptxPath)
	if !strings.Contains(chartXML, "<c:v>Only</c:v>") {
		t.Fatalf("replaced chart XML missing category %q", "Only")
	}
	if !strings.Contains(chartXML, "<c:v>42</c:v>") {
		t.Fatalf("replaced chart XML missing value %q", "42")
	}
}

func TestPresentation_UpdateChartDataByRelID(t *testing.T) {
	pptxPath := createPresentationWithBarChart(t)

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

	err = prs.UpdateChartDataByRelID(0, refs[0].RelID, ChartDataUpdate{
		Categories: []string{"West", "East"},
		Series:     []ChartSeriesData{{Values: []float64{12, 21}}},
	})
	if err != nil {
		_ = prs.Close()
		t.Fatalf("UpdateChartDataByRelID failed: %v", err)
	}
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save after UpdateChartDataByRelID failed: %v", err)
	}
	_ = prs.Close()

	chartXML := readZipEntry(t, pptxPath)
	if !strings.Contains(chartXML, "<c:v>West</c:v>") {
		t.Fatalf("updated chart XML missing category %q", "West")
	}
	if !strings.Contains(chartXML, "<c:v>21</c:v>") {
		t.Fatalf("updated chart XML missing value %q", "21")
	}
}

func TestPresentation_UpdateChartDataByRelID_Scatter(t *testing.T) {
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

	err = prs.UpdateChartDataByRelID(0, refs[0].RelID, ChartDataUpdate{
		Series: []ChartSeriesData{
			{XValues: []float64{3, 4}, YValues: []float64{30, 40}},
		},
	})
	if err != nil {
		_ = prs.Close()
		t.Fatalf("UpdateChartDataByRelID scatter failed: %v", err)
	}
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save after scatter update failed: %v", err)
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

func TestPresentation_UpdateChartDataByRelID_Bubble(t *testing.T) {
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

	err = prs.UpdateChartDataByRelID(0, refs[0].RelID, ChartDataUpdate{
		Series: []ChartSeriesData{
			{XValues: []float64{3, 4}, YValues: []float64{30, 40}, Sizes: []float64{7, 8}},
		},
	})
	if err != nil {
		_ = prs.Close()
		t.Fatalf("UpdateChartDataByRelID bubble failed: %v", err)
	}
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save after bubble update failed: %v", err)
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

func TestPresentation_ReplaceChartDataByRelID(t *testing.T) {
	pptxPath := createPresentationWithBarChart(t)

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

	err = prs.ReplaceChartDataByRelID(0, refs[0].RelID, []string{"Single"}, []float64{9})
	if err != nil {
		_ = prs.Close()
		t.Fatalf("ReplaceChartDataByRelID failed: %v", err)
	}
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save after ReplaceChartDataByRelID failed: %v", err)
	}
	_ = prs.Close()

	chartXML := readZipEntry(t, pptxPath)
	if !strings.Contains(chartXML, "<c:v>Single</c:v>") {
		t.Fatalf("replaced chart XML missing category %q", "Single")
	}
	if !strings.Contains(chartXML, "<c:v>9</c:v>") {
		t.Fatalf("replaced chart XML missing value %q", "9")
	}
}

func TestPresentation_ChartAPI_NilAndUninitialized(t *testing.T) {
	var nilPRS *Presentation

	if _, err := nilPRS.ListSlideCharts(0); err == nil {
		t.Fatal("expected ListSlideCharts error on nil presentation")
	}
	if err := nilPRS.UpdateChartData(0, ChartSelector{Index: intPtr(0)}, ChartDataUpdate{
		Categories: []string{"A"},
		Series:     []ChartSeriesData{{Values: []float64{1}}},
	}); err == nil {
		t.Fatal("expected UpdateChartData error on nil presentation")
	}
	if err := nilPRS.ReplaceChartData(0, 0, []string{"A"}, []float64{1}); err == nil {
		t.Fatal("expected ReplaceChartData error on nil presentation")
	}
	if err := nilPRS.UpdateChartDataByIndex(0, 0, ChartDataUpdate{
		Categories: []string{"A"},
		Series:     []ChartSeriesData{{Values: []float64{1}}},
	}); err == nil {
		t.Fatal("expected UpdateChartDataByIndex error on nil presentation")
	}
	if err := nilPRS.UpdateChartDataByRelID(0, "rId1", ChartDataUpdate{
		Categories: []string{"A"},
		Series:     []ChartSeriesData{{Values: []float64{1}}},
	}); err == nil {
		t.Fatal("expected UpdateChartDataByRelID error on nil presentation")
	}
	if err := nilPRS.ReplaceChartDataByRelID(0, "rId1", []string{"A"}, []float64{1}); err == nil {
		t.Fatal("expected ReplaceChartDataByRelID error on nil presentation")
	}
	if err := nilPRS.UpdateChartFormattingByIndex(0, 0, ChartFormatUpdate{}); err == nil {
		t.Fatal("expected UpdateChartFormattingByIndex error on nil presentation")
	}
	if err := nilPRS.UpdateChartFormattingByRelID(0, "rId1", ChartFormatUpdate{}); err == nil {
		t.Fatal("expected UpdateChartFormattingByRelID error on nil presentation")
	}

	uninit := &Presentation{}
	if _, err := uninit.ListSlideCharts(0); err == nil {
		t.Fatal("expected ListSlideCharts error on uninitialized presentation")
	}
	if err := uninit.UpdateChartData(0, ChartSelector{Index: intPtr(0)}, ChartDataUpdate{
		Categories: []string{"A"},
		Series:     []ChartSeriesData{{Values: []float64{1}}},
	}); err == nil {
		t.Fatal("expected UpdateChartData error on uninitialized presentation")
	}
	if err := uninit.ReplaceChartData(0, 0, []string{"A"}, []float64{1}); err == nil {
		t.Fatal("expected ReplaceChartData error on uninitialized presentation")
	}
	if err := uninit.UpdateChartDataByIndex(0, 0, ChartDataUpdate{
		Categories: []string{"A"},
		Series:     []ChartSeriesData{{Values: []float64{1}}},
	}); err == nil {
		t.Fatal("expected UpdateChartDataByIndex error on uninitialized presentation")
	}
	if err := uninit.UpdateChartDataByRelID(0, "rId1", ChartDataUpdate{
		Categories: []string{"A"},
		Series:     []ChartSeriesData{{Values: []float64{1}}},
	}); err == nil {
		t.Fatal("expected UpdateChartDataByRelID error on uninitialized presentation")
	}
	if err := uninit.ReplaceChartDataByRelID(0, "rId1", []string{"A"}, []float64{1}); err == nil {
		t.Fatal("expected ReplaceChartDataByRelID error on uninitialized presentation")
	}
	if err := uninit.UpdateChartFormattingByIndex(0, 0, ChartFormatUpdate{}); err == nil {
		t.Fatal("expected UpdateChartFormattingByIndex error on uninitialized presentation")
	}
	if err := uninit.UpdateChartFormattingByRelID(0, "rId1", ChartFormatUpdate{}); err == nil {
		t.Fatal("expected UpdateChartFormattingByRelID error on uninitialized presentation")
	}
}
