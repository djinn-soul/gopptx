package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// dataLabelPointChartXML has one bar series whose plot-level labels show the
// value and the category name, which a per-label patch has to carry over.
func dataLabelPointChartXML() []byte {
	return []byte(`<c:chartSpace xmlns:c="x" xmlns:a="y"><c:chart><c:plotArea><c:barChart>` +
		`<c:ser><c:idx val="0"/><c:order val="0"/>` +
		`<c:dLbls><c:showLegendKey val="0"/><c:showVal val="1"/>` +
		`<c:showCatName val="1"/><c:showSerName val="0"/><c:showPercent val="0"/>` +
		`</c:dLbls>` +
		`<c:cat><c:strLit><c:ptCount val="3"/></c:strLit></c:cat>` +
		`<c:val><c:numLit><c:ptCount val="3"/>` +
		`<c:pt idx="0"><c:v>4</c:v></c:pt><c:pt idx="1"><c:v>7</c:v></c:pt>` +
		`<c:pt idx="2"><c:v>9</c:v></c:pt></c:numLit></c:val></c:ser>` +
		`<c:ser><c:idx val="1"/><c:order val="1"/>` +
		`<c:val><c:numLit><c:ptCount val="1"/><c:pt idx="0"><c:v>1</c:v></c:pt></c:numLit></c:val>` +
		`</c:ser></c:barChart></c:plotArea></c:chart></c:chartSpace>`)
}

// Upstream #638: a single point takes a number format of its own, which lives
// on its c:dLbl rather than on the chart-wide c:dLbls.
func TestPatchChartFormatting_DataLabelPointNumberFormat(t *testing.T) {
	format := "0.0%"
	req := common.ChartFormatUpdate{
		DataLabelPoints: []common.ChartDataLabelPoint{
			{SeriesIndex: 0, PointIndex: 2, NumberFormat: &format},
		},
	}
	got, err := PatchChartFormatting(dataLabelPointChartXML(), req)
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	updated := string(got)

	series := strings.SplitAfter(updated, "</c:ser>")
	if strings.Contains(series[1], "<c:dLbl>") {
		t.Fatalf("label leaked onto the untargeted series: %s", series[1])
	}
	if !strings.Contains(updated, `<c:numFmt formatCode="0.0%" sourceLinked="0"/>`) {
		t.Fatalf("per-label number format missing: %s", updated)
	}
	// CT_DLbl orders idx, numFmt, then the display flags; CT_DLbls puts the
	// per-point labels before the series-wide ones.
	assertXMLOrder(t, series[0], `<c:dLbl><c:idx val="2"/>`, "<c:numFmt", "<c:showVal", "</c:dLbl>")

	// Upstream #803: the format reads back at the point level.
	state := ExtractChartState(got)
	if len(state.DataLabelPoints) != 1 {
		t.Fatalf("expected one label in state, got %#v", state.DataLabelPoints)
	}
	label := state.DataLabelPoints[0]
	if label.PointIndex != 2 || label.NumberFormat == nil || *label.NumberFormat != format {
		t.Fatalf("label number format not read back: %#v", label)
	}
	if label.FormatLinked == nil || *label.FormatLinked {
		t.Fatalf("expected sourceLinked false, got %#v", label.FormatLinked)
	}
}

// Upstream #650: colouring one label must not drop the category name it was
// already showing, because a c:dLbl inherits no flags from its parent.
func TestPatchChartFormatting_DataLabelPointFontKeepsFlags(t *testing.T) {
	color, size := "FFFFFF", 14
	bold := true
	req := common.ChartFormatUpdate{
		DataLabelPoints: []common.ChartDataLabelPoint{
			{SeriesIndex: 0, PointIndex: 0, FontColor: &color, FontSizePt: &size, FontBold: &bold},
		},
	}
	got, err := PatchChartFormatting(dataLabelPointChartXML(), req)
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	updated := string(got)

	label := dataLabelBlocks(strings.SplitAfter(updated, "</c:ser>")[0])
	if len(label) != 1 {
		t.Fatalf("expected exactly one per-point label, got %#v", label)
	}
	if !strings.Contains(label[0], `<a:defRPr sz="1400" b="1">`) {
		t.Fatalf("label font not applied: %s", label[0])
	}
	if !strings.Contains(label[0], `<a:srgbClr val="FFFFFF"/>`) {
		t.Fatalf("label font colour not applied: %s", label[0])
	}
	if !strings.Contains(label[0], `<c:showCatName val="1"/>`) {
		t.Fatalf("category flag dropped from the label: %s", label[0])
	}
	assertXMLOrder(t, label[0], "<c:idx", "<c:txPr", "<c:showVal", "<c:showCatName")

	state := ExtractChartState(got)
	if len(state.DataLabelPoints) != 1 {
		t.Fatalf("expected one label in state, got %#v", state.DataLabelPoints)
	}
	read := state.DataLabelPoints[0]
	if read.FontColor == nil || *read.FontColor != color {
		t.Fatalf("label font colour not read back: %#v", read)
	}
	if read.FontSizePt == nil || *read.FontSizePt != size {
		t.Fatalf("label font size not read back: %#v", read)
	}
	if read.ShowCategory == nil || !*read.ShowCategory {
		t.Fatalf("label category flag not read back: %#v", read)
	}
}

// A second patch on the same label merges: the number format survives a later
// font change, and the flags stay in schema order.
func TestPatchChartFormatting_DataLabelPointMergesAndDeletes(t *testing.T) {
	format, color := `#,##0.00 "kg"`, "FF0000"
	first, err := PatchChartFormatting(dataLabelPointChartXML(), common.ChartFormatUpdate{
		DataLabelPoints: []common.ChartDataLabelPoint{
			{SeriesIndex: 0, PointIndex: 1, NumberFormat: &format},
		},
	})
	if err != nil {
		t.Fatalf("first PatchChartFormatting error: %v", err)
	}
	second, err := PatchChartFormatting(first, common.ChartFormatUpdate{
		DataLabelPoints: []common.ChartDataLabelPoint{
			{SeriesIndex: 0, PointIndex: 1, FontColor: &color},
		},
	})
	if err != nil {
		t.Fatalf("second PatchChartFormatting error: %v", err)
	}
	updated := string(second)

	if strings.Count(updated, `<c:dLbl>`) != 1 {
		t.Fatalf("repeated patch duplicated the label: %s", updated)
	}
	if !strings.Contains(updated, `formatCode="#,##0.00 &quot;kg&quot;"`) {
		t.Fatalf("number format lost or unescaped on merge: %s", updated)
	}
	state := ExtractChartState(second)
	if len(state.DataLabelPoints) != 1 || state.DataLabelPoints[0].NumberFormat == nil ||
		*state.DataLabelPoints[0].NumberFormat != format {
		t.Fatalf("merged number format not read back: %#v", state.DataLabelPoints)
	}

	deleted := true
	third, err := PatchChartFormatting(second, common.ChartFormatUpdate{
		DataLabelPoints: []common.ChartDataLabelPoint{
			{SeriesIndex: 0, PointIndex: 1, Delete: &deleted},
		},
	})
	if err != nil {
		t.Fatalf("third PatchChartFormatting error: %v", err)
	}
	// CT_DLbl allows only idx and delete once the label is removed.
	label := dataLabelBlocks(string(third))[0]
	if label != `<c:dLbl><c:idx val="1"/><c:delete val="1"/></c:dLbl>` {
		t.Fatalf("deleted label carries more than idx and delete: %s", label)
	}
}

func TestValidateChartFormatUpdate_DataLabelPoints(t *testing.T) {
	badColor, negative := "not-a-color", -1
	oversized := 401
	cases := []struct {
		name string
		req  common.ChartFormatUpdate
	}{
		{"font colour", common.ChartFormatUpdate{DataLabelPoints: []common.ChartDataLabelPoint{
			{FontColor: &badColor},
		}}},
		{"point index", common.ChartFormatUpdate{DataLabelPoints: []common.ChartDataLabelPoint{
			{PointIndex: negative},
		}}},
		{"font size", common.ChartFormatUpdate{DataLabelPoints: []common.ChartDataLabelPoint{
			{FontSizePt: &oversized},
		}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateChartFormatUpdate(tc.req); err == nil {
				t.Fatalf("expected an error for %s", tc.name)
			}
		})
	}
}

// Hosting a per-point label requires a series-level <c:dLbls>, and that block
// shadows the plot-level one: the chart-wide number format has to be copied into
// it and onto the rebuilt label, or every other label of the series silently
// reverts to the general format.
func TestPatchChartFormatting_DataLabelPointKeepsChartWideFormat(t *testing.T) {
	format, color := "#,##0", "C00000"
	got, err := PatchChartFormatting(dataLabelPointChartXML(), common.ChartFormatUpdate{
		ShowDataLabels:        boolPointer(true),
		DataLabelShowValue:    boolPointer(true),
		DataLabelNumberFormat: &format,
		DataLabelPoints: []common.ChartDataLabelPoint{
			{PointIndex: 0, FontColor: &color},
		},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	updated := string(got)
	series := strings.SplitAfter(updated, "</c:ser>")[0]

	// The series block carries the chart-wide format for its untouched labels.
	seriesLabels := seriesDataLabelsBlock(series)
	if !strings.Contains(seriesLabels, `<c:numFmt formatCode="#,##0" sourceLinked="0"/>`) {
		t.Fatalf("series labels did not inherit the chart-wide format: %s", seriesLabels)
	}
	// The patched label carries it too, since a c:dLbl inherits nothing.
	label := dataLabelBlocks(series)[0]
	if !strings.Contains(label, `<c:numFmt formatCode="#,##0" sourceLinked="0"/>`) {
		t.Fatalf("patched label lost the chart-wide format: %s", label)
	}
	// A format of the label's own still wins.
	own := "0.0%"
	got, err = PatchChartFormatting(got, common.ChartFormatUpdate{
		DataLabelPoints: []common.ChartDataLabelPoint{{PointIndex: 0, NumberFormat: &own}},
	})
	if err != nil {
		t.Fatalf("second PatchChartFormatting error: %v", err)
	}
	label = dataLabelBlocks(strings.SplitAfter(string(got), "</c:ser>")[0])[0]
	if !strings.Contains(label, `formatCode="0.0%"`) {
		t.Fatalf("per-label format did not win: %s", label)
	}
}

func boolPointer(value bool) *bool { return &value }
