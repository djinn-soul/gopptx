package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Chart XML with a stacked bar plot, one series and both axes, used by the
// gridline, series-line, series-format and data-label-box patches.
const lineStyleChartXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
<c:chart><c:plotArea><c:barChart><c:barDir val="col"/><c:grouping val="stacked"/>
<c:ser><c:idx val="0"/><c:order val="0"/><c:cat><c:strRef><c:f>Sheet1!$A$2</c:f></c:strRef></c:cat>
<c:val><c:numRef><c:f>Sheet1!$B$2</c:f></c:numRef></c:val></c:ser>
<c:gapWidth val="150"/><c:overlap val="100"/><c:axId val="1"/><c:axId val="2"/></c:barChart>
<c:catAx><c:axId val="1"/><c:crosses val="autoZero"/></c:catAx>
<c:valAx><c:axId val="2"/><c:majorGridlines/><c:crosses val="autoZero"/></c:valAx>
</c:plotArea></c:chart></c:chartSpace>`

func patchLineStyleChart(t *testing.T, req common.ChartFormatUpdate) string {
	t.Helper()
	got, err := PatchChartFormatting([]byte(lineStyleChartXML), req)
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	return string(got)
}

func intPtr(value int) *int { return &value }

// Upstream #984: the gridline style lives in the gridline's own c:spPr.
func TestPatchChartFormatting_GridlineFormats(t *testing.T) {
	xml := patchLineStyleChart(t, common.ChartFormatUpdate{
		ValueAxisMajorGridFormat: &common.ChartLineFormat{
			Color: stringPtr("D9D9D9"), WidthEMU: intPtr(9525), Dash: stringPtr("sysDot"),
		},
		CategoryAxisMinorGridFormat: &common.ChartLineFormat{Color: stringPtr("#ff0000")},
	})
	mustContain(t, xml, `<c:majorGridlines><c:spPr><a:ln w="9525">`+
		`<a:solidFill><a:srgbClr val="D9D9D9"/></a:solidFill><a:prstDash val="sysDot"/></a:ln></c:spPr></c:majorGridlines>`)
	// Styling a gridline the axis lacks draws it: there is nowhere else for the
	// c:spPr to live.
	mustContain(t, xml, `<c:catAx><c:axId val="1"/><c:minorGridlines><c:spPr>`+
		`<a:ln><a:solidFill><a:srgbClr val="FF0000"/></a:solidFill></a:ln></c:spPr></c:minorGridlines>`)
}

func TestPatchChartFormatting_GridlineFormatIsIdempotent(t *testing.T) {
	format := &common.ChartLineFormat{Color: stringPtr("112233")}
	first := patchLineStyleChart(t, common.ChartFormatUpdate{ValueAxisMajorGridFormat: format})
	second, err := PatchChartFormatting([]byte(first), common.ChartFormatUpdate{
		ValueAxisMajorGridFormat: &common.ChartLineFormat{Color: stringPtr("445566")},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	xml := string(second)
	if strings.Count(xml, "<c:majorGridlines>") != 1 {
		t.Fatalf("expected exactly one majorGridlines element, got %q", xml)
	}
	mustContain(t, xml, `<a:srgbClr val="445566"/>`)
	if strings.Contains(xml, `val="112233"`) {
		t.Fatalf("expected the first colour replaced, got %q", xml)
	}
}

func TestPatchChartFormatting_GridlineFormatRejectsBadColor(t *testing.T) {
	_, err := PatchChartFormatting([]byte(lineStyleChartXML), common.ChartFormatUpdate{
		ValueAxisMajorGridFormat: &common.ChartLineFormat{Color: stringPtr("nope")},
	})
	if err == nil {
		t.Fatal("expected an error for an unparsable gridline colour")
	}
}

// Upstream #846: c:serLines sits after c:overlap and before the axis ids.
func TestPatchChartFormatting_SeriesLines(t *testing.T) {
	show := true
	xml := patchLineStyleChart(t, common.ChartFormatUpdate{
		SeriesLines: &common.ChartSeriesLines{
			Show: &show,
			Line: &common.ChartLineFormat{Color: stringPtr("404040"), WidthEMU: intPtr(12700)},
		},
	})
	mustContain(t, xml, `<c:overlap val="100"/><c:serLines><c:spPr><a:ln w="12700">`+
		`<a:solidFill><a:srgbClr val="404040"/></a:solidFill></a:ln></c:spPr></c:serLines><c:axId val="1"/>`)
}

func TestPatchChartFormatting_SeriesLinesRemoved(t *testing.T) {
	show := true
	withLines := patchLineStyleChart(t, common.ChartFormatUpdate{
		SeriesLines: &common.ChartSeriesLines{Show: &show},
	})
	mustContain(t, withLines, `<c:serLines/>`)

	hide := false
	got, err := PatchChartFormatting([]byte(withLines), common.ChartFormatUpdate{
		SeriesLines: &common.ChartSeriesLines{Show: &hide},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	if strings.Contains(string(got), "<c:serLines") {
		t.Fatalf("expected the series lines removed, got %q", got)
	}
}

// Upstream #872: the markers of a series carry their own fill and outline.
func TestPatchChartFormatting_SeriesFormat(t *testing.T) {
	xml := patchLineStyleChart(t, common.ChartFormatUpdate{
		SeriesFormats: []common.ChartSeriesFormat{{
			FillColor:       stringPtr("1F77B4"),
			LineColor:       stringPtr("0A3A5A"),
			LineWidthEMU:    intPtr(19050),
			MarkerSymbol:    stringPtr("circle"),
			MarkerSize:      intPtr(7),
			MarkerFillColor: stringPtr("FFFFFF"),
			MarkerLineColor: stringPtr("1F77B4"),
		}},
	})
	// CT_Ser puts c:spPr and c:marker before the category and value references.
	mustContain(t, xml, `<c:order val="0"/><c:spPr><a:solidFill><a:srgbClr val="1F77B4"/></a:solidFill>`+
		`<a:ln w="19050"><a:solidFill><a:srgbClr val="0A3A5A"/></a:solidFill></a:ln></c:spPr>`+
		`<c:marker><c:symbol val="circle"/><c:size val="7"/><c:spPr>`+
		`<a:solidFill><a:srgbClr val="FFFFFF"/></a:solidFill>`+
		`<a:ln><a:solidFill><a:srgbClr val="1F77B4"/></a:solidFill></a:ln></c:spPr></c:marker><c:cat>`)
}

// Recolouring a marker must not reset the symbol and size it already had.
func TestPatchChartFormatting_SeriesMarkerKeepsShape(t *testing.T) {
	first := patchLineStyleChart(t, common.ChartFormatUpdate{
		SeriesFormats: []common.ChartSeriesFormat{{
			MarkerSymbol: stringPtr("diamond"), MarkerSize: intPtr(9),
			MarkerFillColor: stringPtr("112233"),
		}},
	})
	got, err := PatchChartFormatting([]byte(first), common.ChartFormatUpdate{
		SeriesFormats: []common.ChartSeriesFormat{{MarkerFillColor: stringPtr("445566")}},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	xml := string(got)
	mustContain(t, xml, `<c:symbol val="diamond"/><c:size val="9"/>`)
	mustContain(t, xml, `<a:srgbClr val="445566"/>`)
	if strings.Count(xml, "<c:marker>") != 1 {
		t.Fatalf("expected exactly one marker element, got %q", xml)
	}
}

// A series format must not disturb the per-point run written before it.
func TestPatchChartFormatting_SeriesFormatKeepsDataPoints(t *testing.T) {
	xml := patchLineStyleChart(t, common.ChartFormatUpdate{
		DataPoints: []common.ChartDataPoint{{PointIndex: 1, FillColor: stringPtr("AABBCC")}},
		SeriesFormats: []common.ChartSeriesFormat{{
			LineColor: stringPtr("102030"),
		}},
	})
	mustContain(t, xml, `<c:dPt><c:idx val="1"/><c:spPr>`+
		`<a:solidFill><a:srgbClr val="AABBCC"/></a:solidFill></c:spPr></c:dPt>`)
	mustContain(t, xml, `<c:spPr><a:ln><a:solidFill><a:srgbClr val="102030"/></a:solidFill></a:ln></c:spPr><c:dPt>`)
}

func TestPatchChartFormatting_SeriesFormatRejectsBadMarker(t *testing.T) {
	_, err := PatchChartFormatting([]byte(lineStyleChartXML), common.ChartFormatUpdate{
		SeriesFormats: []common.ChartSeriesFormat{{MarkerSymbol: stringPtr("blob")}},
	})
	if err == nil {
		t.Fatal("expected an error for a symbol outside CT_MarkerStyle")
	}
}

// Upstream #662 and #716: the label box is the c:spPr of the c:dLbls.
func TestPatchChartFormatting_DataLabelBox(t *testing.T) {
	show := true
	xml := patchLineStyleChart(t, common.ChartFormatUpdate{
		ShowDataLabels:     &show,
		DataLabelFillColor: stringPtr("FFF2CC"),
		DataLabelBorder: &common.ChartLineFormat{
			Color: stringPtr("BF8F00"), WidthEMU: intPtr(9525), Dash: stringPtr("dash"),
		},
	})
	mustContain(t, xml, `<c:spPr><a:solidFill><a:srgbClr val="FFF2CC"/></a:solidFill>`+
		`<a:ln w="9525"><a:solidFill><a:srgbClr val="BF8F00"/></a:solidFill>`+
		`<a:prstDash val="dash"/></a:ln></c:spPr>`)
	if strings.Count(xml, "<c:dLbls>") != 1 {
		t.Fatalf("expected exactly one dLbls element, got %q", xml)
	}
}

func TestPatchChartFormatting_DataLabelPointBox(t *testing.T) {
	xml := patchLineStyleChart(t, common.ChartFormatUpdate{
		DataLabelPoints: []common.ChartDataLabelPoint{{
			PointIndex: 0,
			FillColor:  stringPtr("E2EFDA"),
			Border:     &common.ChartLineFormat{Color: stringPtr("375623")},
		}},
	})
	mustContain(t, xml, `<c:dLbl><c:idx val="0"/><c:spPr>`+
		`<a:solidFill><a:srgbClr val="E2EFDA"/></a:solidFill>`+
		`<a:ln><a:solidFill><a:srgbClr val="375623"/></a:solidFill></a:ln></c:spPr>`)
}

// The state snapshot has to report back what the patches wrote.
func TestExtractChartState_LineStyles(t *testing.T) {
	show := true
	xml := patchLineStyleChart(t, common.ChartFormatUpdate{
		ValueAxisMajorGridFormat: &common.ChartLineFormat{
			Color: stringPtr("D9D9D9"), WidthEMU: intPtr(9525),
		},
		SeriesLines: &common.ChartSeriesLines{Show: &show},
		SeriesFormats: []common.ChartSeriesFormat{{
			FillColor: stringPtr("1F77B4"), MarkerSymbol: stringPtr("circle"),
		}},
	})
	state := ExtractChartState([]byte(xml))

	gridline := state.ValueAx.MajorGridlineFormat
	if gridline == nil || gridline.Color == nil || *gridline.Color != "D9D9D9" {
		t.Fatalf("expected the value axis gridline colour read back, got %+v", gridline)
	}
	if gridline.WidthEMU == nil || *gridline.WidthEMU != 9525 {
		t.Fatalf("expected the gridline width read back, got %+v", gridline)
	}
	if state.SeriesLines == nil || state.SeriesLines.Show == nil || !*state.SeriesLines.Show {
		t.Fatalf("expected the series lines read back, got %+v", state.SeriesLines)
	}
	if len(state.SeriesFormats) != 1 {
		t.Fatalf("expected one series format, got %+v", state.SeriesFormats)
	}
	format := state.SeriesFormats[0]
	if format.FillColor == nil || *format.FillColor != "1F77B4" {
		t.Fatalf("expected the series fill read back, got %+v", format)
	}
	if format.MarkerSymbol == nil || *format.MarkerSymbol != "circle" {
		t.Fatalf("expected the marker symbol read back, got %+v", format)
	}
}
