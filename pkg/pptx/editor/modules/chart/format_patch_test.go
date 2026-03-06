package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const sampleChartXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
<c:chart>
<c:title><c:tx><c:rich><a:bodyPr/><a:lstStyle/><a:p><a:r><a:t>Old</a:t></a:r></a:p></c:rich></c:tx><c:overlay val="0"/></c:title>
<c:autoTitleDeleted val="0"/>
<c:plotArea>
<c:layout/>
<c:barChart>
<c:ser></c:ser>
<c:dLbls><c:showVal val="1"/></c:dLbls>
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:barChart>
</c:plotArea>
<c:legend><c:legendPos val="r"/><c:overlay val="0"/></c:legend>
<c:plotVisOnly val="1"/>
</c:chart>
</c:chartSpace>`

func TestPatchChartFormatting_TitleLegendAndDataLabels(t *testing.T) {
	title := "Revenue 2026"
	titleOverlay := true
	showLegend := true
	legendPosition := "b"
	legendOverlay := true
	showLabels := true
	labelPosition := "outEnd"
	showSeriesName := true
	showPercent := false

	got, err := PatchChartFormatting([]byte(sampleChartXML), common.ChartFormatUpdate{
		Title:                   &title,
		TitleOverlay:            &titleOverlay,
		ShowLegend:              &showLegend,
		LegendPosition:          &legendPosition,
		LegendOverlay:           &legendOverlay,
		ShowDataLabels:          &showLabels,
		DataLabelPosition:       &labelPosition,
		DataLabelShowSeriesName: &showSeriesName,
		DataLabelShowPercent:    &showPercent,
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	xml := string(got)
	mustContain(t, xml, `<a:t>Revenue 2026</a:t>`)
	mustContain(t, xml, `<c:legendPos val="b"/>`)
	mustContain(t, xml, `<c:legend><c:legendPos val="b"/><c:overlay val="1"/></c:legend>`)
	mustContain(t, xml, `<c:dLblPos val="outEnd"/>`)
	mustContain(t, xml, `<c:showSerName val="1"/>`)
	if strings.Contains(xml, `<c:showPercent`) {
		t.Fatalf("expected showPercent removed when false")
	}
}

func TestPatchChartFormatting_HideLegendAndDataLabels(t *testing.T) {
	showLegend := false
	showLabels := false

	got, err := PatchChartFormatting([]byte(sampleChartXML), common.ChartFormatUpdate{
		ShowLegend:     &showLegend,
		ShowDataLabels: &showLabels,
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	xml := string(got)
	if strings.Contains(xml, "<c:legend>") {
		t.Fatalf("expected legend block removed")
	}
	if strings.Contains(xml, "<c:dLbls>") {
		t.Fatalf("expected data labels block removed")
	}
}

func TestValidateChartFormatUpdateRejectsInvalidPositions(t *testing.T) {
	legendPosition := "middle"
	if err := ValidateChartFormatUpdate(common.ChartFormatUpdate{LegendPosition: &legendPosition}); err == nil {
		t.Fatalf("expected invalid legend position error")
	}

	labelPosition := "diagonal"
	if err := ValidateChartFormatUpdate(common.ChartFormatUpdate{DataLabelPosition: &labelPosition}); err == nil {
		t.Fatalf("expected invalid data label position error")
	}
}

func TestPatchChartFormatting_ShowTitleAndPlotVisibleOnly(t *testing.T) {
	showTitle := false
	plotVisibleOnly := false
	got, err := PatchChartFormatting([]byte(sampleChartXML), common.ChartFormatUpdate{
		ShowTitle:       &showTitle,
		PlotVisibleOnly: &plotVisibleOnly,
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	xml := string(got)
	if strings.Contains(xml, "<c:title>") {
		t.Fatalf("expected title block removed")
	}
	mustContain(t, xml, `<c:autoTitleDeleted val="1"/>`)
	mustContain(t, xml, `<c:plotVisOnly val="0"/>`)
}

func mustContain(t *testing.T, xml string, want string) {
	t.Helper()
	if !strings.Contains(xml, want) {
		t.Fatalf("expected xml to contain %q", want)
	}
}
