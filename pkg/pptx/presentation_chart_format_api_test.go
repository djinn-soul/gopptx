package pptx

import (
	"strings"
	"testing"
)

func TestPresentation_UpdateChartFormattingByIndex(t *testing.T) {
	pptxPath := createPresentationWithBarChart(t)

	prs, err := Open(pptxPath)
	if err != nil {
		t.Fatalf("open presentation: %v", err)
	}

	title := "Updated Title"
	showTitle := true
	plotVisibleOnly := false
	showLegend := true
	legendPos := "b"
	showLabels := true
	labelPos := "outEnd"
	showSeriesName := true
	showPercent := false
	if err := prs.UpdateChartFormattingByIndex(0, 0, ChartFormatUpdate{
		ShowTitle:               &showTitle,
		Title:                   &title,
		PlotVisibleOnly:         &plotVisibleOnly,
		ShowLegend:              &showLegend,
		LegendPosition:          &legendPos,
		ShowDataLabels:          &showLabels,
		DataLabelPosition:       &labelPos,
		DataLabelShowSeriesName: &showSeriesName,
		DataLabelShowPercent:    &showPercent,
	}); err != nil {
		_ = prs.Close()
		t.Fatalf("UpdateChartFormattingByIndex failed: %v", err)
	}
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save after UpdateChartFormattingByIndex failed: %v", err)
	}
	_ = prs.Close()

	chartXML := readZipEntry(t, pptxPath)
	if !strings.Contains(chartXML, "<a:t>Updated Title</a:t>") {
		t.Fatalf("updated chart XML missing updated title")
	}
	if !strings.Contains(chartXML, `<c:plotVisOnly val="0"/>`) {
		t.Fatalf("updated chart XML missing plot visibility override")
	}
	if !strings.Contains(chartXML, `<c:legendPos val="b"/>`) {
		t.Fatalf("updated chart XML missing legend position")
	}
	if !strings.Contains(chartXML, `<c:dLblPos val="outEnd"/>`) {
		t.Fatalf("updated chart XML missing data label position")
	}
	if !strings.Contains(chartXML, `<c:showSerName val="1"/>`) {
		t.Fatalf("updated chart XML missing showSerName")
	}
	if strings.Contains(chartXML, `<c:showPercent`) {
		t.Fatalf("expected showPercent removed when false")
	}
}

func TestPresentation_UpdateChartFormattingByRelID(t *testing.T) {
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

	showLegend := false
	showLabels := false
	showTitle := false
	if err := prs.UpdateChartFormattingByRelID(0, refs[0].RelID, ChartFormatUpdate{
		ShowTitle:      &showTitle,
		ShowLegend:     &showLegend,
		ShowDataLabels: &showLabels,
	}); err != nil {
		_ = prs.Close()
		t.Fatalf("UpdateChartFormattingByRelID failed: %v", err)
	}
	if err := prs.Save(); err != nil {
		_ = prs.Close()
		t.Fatalf("save after UpdateChartFormattingByRelID failed: %v", err)
	}
	_ = prs.Close()

	chartXML := readZipEntry(t, pptxPath)
	if strings.Contains(chartXML, "<c:title>") {
		t.Fatalf("expected title block removed")
	}
	if strings.Contains(chartXML, "<c:legend>") {
		t.Fatalf("expected legend block removed")
	}
	if strings.Contains(chartXML, "<c:dLbls>") {
		t.Fatalf("expected data label block removed")
	}
}
