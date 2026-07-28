package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// dataPointChartXML has one bar series with a cached negative value at index 2.
func dataPointChartXML() []byte {
	return []byte(`<c:chartSpace xmlns:c="x" xmlns:a="y"><c:chart><c:plotArea><c:barChart>` +
		`<c:ser><c:idx val="0"/><c:order val="0"/>` +
		`<c:spPr><a:solidFill><a:srgbClr val="4F81BD"/></a:solidFill></c:spPr>` +
		`<c:dLbls><c:showVal val="1"/></c:dLbls>` +
		`<c:cat><c:strLit><c:ptCount val="4"/></c:strLit></c:cat>` +
		`<c:val><c:numLit><c:ptCount val="4"/>` +
		`<c:pt idx="0"><c:v>4</c:v></c:pt><c:pt idx="1"><c:v>7</c:v></c:pt>` +
		`<c:pt idx="2"><c:v>-3</c:v></c:pt><c:pt idx="3"><c:v>9</c:v></c:pt>` +
		`</c:numLit></c:val></c:ser>` +
		`<c:ser><c:idx val="1"/><c:order val="1"/>` +
		`<c:val><c:numLit><c:ptCount val="1"/><c:pt idx="0"><c:v>1</c:v></c:pt></c:numLit></c:val>` +
		`</c:ser></c:barChart></c:plotArea></c:chart></c:chartSpace>`)
}

func TestPatchChartFormatting_DataPoints(t *testing.T) {
	fill, line := "FF0000", "00FF00"
	width := 12700
	explosion := 25

	req := common.ChartFormatUpdate{
		DataPoints: []common.ChartDataPoint{
			{SeriesIndex: 0, PointIndex: 3, FillColor: &fill, Explosion: &explosion},
			{SeriesIndex: 0, PointIndex: 1, LineColor: &line, LineWidthEMU: &width},
		},
	}
	got, err := PatchChartFormatting(dataPointChartXML(), req)
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	updated := string(got)

	series := strings.SplitAfter(updated, "</c:ser>")
	if strings.Contains(series[1], "<c:dPt>") {
		t.Fatalf("data points leaked onto the untargeted series: %s", series[1])
	}
	// The run must be sorted by idx and sit before the labels.
	assertXMLOrder(t, series[0], `<c:dPt><c:idx val="1"/>`, `<c:dPt><c:idx val="3"/>`, "<c:dLbls")
	assertXMLOrder(t, series[0], "<c:spPr", "<c:dPt>", "<c:dLbls", "<c:cat")
	if !strings.Contains(updated, `<a:ln w="12700"><a:solidFill><a:srgbClr val="00FF00"/></a:solidFill></a:ln>`) {
		t.Fatalf("point line formatting missing: %s", updated)
	}
	assertXMLOrder(t, reDataPointBlock.FindAllString(updated, -1)[1], "<c:idx", "<c:explosion", "<c:spPr")

	state := ExtractChartState(got)
	if len(state.DataPoints) != 2 {
		t.Fatalf("expected two data points in state, got %#v", state.DataPoints)
	}
	if state.DataPoints[1].FillColor == nil || *state.DataPoints[1].FillColor != fill {
		t.Fatalf("point fill not preserved in state: %#v", state.DataPoints[1])
	}
	if state.DataPoints[0].LineColor == nil || *state.DataPoints[0].LineColor != line ||
		state.DataPoints[0].LineWidthEMU == nil || *state.DataPoints[0].LineWidthEMU != width {
		t.Fatalf("point line not preserved in state: %#v", state.DataPoints[0])
	}

	again, err := PatchChartFormatting(got, req)
	if err != nil {
		t.Fatalf("re-patch error: %v", err)
	}
	if string(again) != updated {
		t.Fatalf("re-patch was not idempotent:\n%s\n%s", updated, string(again))
	}
}

func TestPatchChartFormattingMergesAndClearsDataPoints(t *testing.T) {
	fill, line := "FF0000", "00FF00"
	seeded, err := PatchChartFormatting(dataPointChartXML(), common.ChartFormatUpdate{
		DataPoints: []common.ChartDataPoint{{SeriesIndex: 0, PointIndex: 1, FillColor: &fill}},
	})
	if err != nil {
		t.Fatalf("seed error: %v", err)
	}
	merged, err := PatchChartFormatting(seeded, common.ChartFormatUpdate{
		DataPoints: []common.ChartDataPoint{{SeriesIndex: 0, PointIndex: 1, LineColor: &line}},
	})
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	block := reDataPointBlock.FindString(string(merged))
	if !strings.Contains(block, `<a:srgbClr val="FF0000"/>`) ||
		!strings.Contains(block, `<a:srgbClr val="00FF00"/>`) {
		t.Fatalf("second patch dropped the earlier point formatting: %s", block)
	}

	cleared, err := PatchChartFormatting(merged, common.ChartFormatUpdate{
		ClearDataPointSeries: []int{0},
	})
	if err != nil {
		t.Fatalf("clear error: %v", err)
	}
	if strings.Contains(string(cleared), "<c:dPt>") {
		t.Fatalf("expected data points to be removed: %s", string(cleared))
	}
}

func TestPatchChartFormattingInvertIfNegativeWritesNegativePointFill(t *testing.T) {
	negative := "C00000"
	got, err := PatchChartFormatting(dataPointChartXML(), common.ChartFormatUpdate{
		SeriesInverts: []common.ChartSeriesInvert{
			{SeriesIndex: 0, InvertIfNegative: true, NegativeFillColor: &negative},
		},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	updated := string(got)
	series := strings.SplitAfter(updated, "</c:ser>")
	assertXMLOrder(t, series[0], "<c:spPr", "<c:invertIfNegative", "<c:dPt>", "<c:dLbls")

	points := reDataPointBlock.FindAllString(updated, -1)
	if len(points) != 1 {
		t.Fatalf("expected one negative data point, got %d: %s", len(points), updated)
	}
	if !strings.Contains(points[0], `<c:idx val="2"/>`) {
		t.Fatalf("negative point index wrong: %s", points[0])
	}
	// The point turns its own inversion off so the explicit fill is drawn.
	if !strings.Contains(points[0], `<c:invertIfNegative val="0"/>`) ||
		!strings.Contains(points[0], `<a:srgbClr val="C00000"/>`) {
		t.Fatalf("negative point fill not written: %s", points[0])
	}

	again, err := PatchChartFormatting(got, common.ChartFormatUpdate{
		SeriesInverts: []common.ChartSeriesInvert{
			{SeriesIndex: 0, InvertIfNegative: true, NegativeFillColor: &negative},
		},
	})
	if err != nil {
		t.Fatalf("re-patch error: %v", err)
	}
	if string(again) != updated {
		t.Fatalf("re-patch was not idempotent:\n%s\n%s", updated, string(again))
	}
}

func TestPatchChartFormattingInvertIfNegativeFlagOnly(t *testing.T) {
	got, err := PatchChartFormatting(dataPointChartXML(), common.ChartFormatUpdate{
		SeriesInverts: []common.ChartSeriesInvert{{SeriesIndex: 0, InvertIfNegative: true}},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	if !strings.Contains(string(got), `<c:invertIfNegative val="1"/>`) {
		t.Fatalf("series invert flag missing: %s", string(got))
	}
	if strings.Contains(string(got), "<c:dPt>") {
		t.Fatalf("no negative fill was requested, so no data point should be written: %s", string(got))
	}
}

func TestValidateChartFormatUpdateRejectsInvalidDataPointsAndInverts(t *testing.T) {
	badColor := "nope"
	badExplosion := 500
	negativeWidth := -1
	pointCases := map[string]common.ChartDataPoint{
		"negative series":   {SeriesIndex: -1},
		"negative point":    {PointIndex: -1},
		"bad fill":          {FillColor: &badColor},
		"bad line":          {LineColor: &badColor},
		"negative width":    {LineWidthEMU: &negativeWidth},
		"explosion too big": {Explosion: &badExplosion},
	}
	for name, point := range pointCases {
		t.Run(name, func(t *testing.T) {
			err := ValidateChartFormatUpdate(common.ChartFormatUpdate{
				DataPoints: []common.ChartDataPoint{point},
			})
			if err == nil {
				t.Fatalf("expected %s to be rejected", name)
			}
		})
	}

	if err := ValidateChartFormatUpdate(common.ChartFormatUpdate{
		SeriesInverts: []common.ChartSeriesInvert{{SeriesIndex: -1}},
	}); err == nil {
		t.Fatal("expected a negative invert series index to be rejected")
	}
	if err := ValidateChartFormatUpdate(common.ChartFormatUpdate{
		SeriesInverts: []common.ChartSeriesInvert{{NegativeFillColor: &badColor}},
	}); err == nil {
		t.Fatal("expected an invalid negative fill colour to be rejected")
	}
}

// Upstream #825: recolouring a scatter point means recolouring its marker. The
// point's own c:spPr formats the segment leading to it, so a marker-less patch
// leaves the marker in the series colour.
func TestPatchChartFormatting_DataPointMarker(t *testing.T) {
	fill, line, symbol := "C00000", "7F0000", "circle"
	size := 9

	got, err := PatchChartFormatting(dataPointChartXML(), common.ChartFormatUpdate{
		DataPoints: []common.ChartDataPoint{{
			PointIndex:      2,
			MarkerFillColor: &fill,
			MarkerLineColor: &line,
			MarkerSymbol:    &symbol,
			MarkerSize:      &size,
		}},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	point := reDataPointBlock.FindString(string(got))
	// CT_Marker orders symbol, size, spPr; CT_DPt puts the marker after idx.
	assertXMLOrder(t, point, "<c:idx", "<c:marker>", "<c:symbol", "<c:size", "<c:spPr>")
	if !strings.Contains(point, `<a:solidFill><a:srgbClr val="C00000"/></a:solidFill>`) {
		t.Fatalf("marker fill missing: %s", point)
	}

	state := ExtractChartState(got)
	if len(state.DataPoints) != 1 {
		t.Fatalf("expected one data point in state, got %#v", state.DataPoints)
	}
	read := state.DataPoints[0]
	if read.MarkerFillColor == nil || *read.MarkerFillColor != fill {
		t.Fatalf("marker fill not read back: %#v", read)
	}
	if read.MarkerSize == nil || *read.MarkerSize != size {
		t.Fatalf("marker size not read back: %#v", read)
	}
	// The marker's colours must not be mistaken for the point's own.
	if read.FillColor != nil || read.LineColor != nil {
		t.Fatalf("marker colours leaked into the point shape: %#v", read)
	}
}

func TestValidateChartFormatUpdateRejectsInvalidMarkers(t *testing.T) {
	badColor, badSymbol := "nope", "hexagon"
	badSize := 100
	cases := map[string]common.ChartDataPoint{
		"marker fill":   {MarkerFillColor: &badColor},
		"marker line":   {MarkerLineColor: &badColor},
		"marker symbol": {MarkerSymbol: &badSymbol},
		"marker size":   {MarkerSize: &badSize},
	}
	for name, point := range cases {
		t.Run(name, func(t *testing.T) {
			err := ValidateChartFormatUpdate(common.ChartFormatUpdate{
				DataPoints: []common.ChartDataPoint{point},
			})
			if err == nil {
				t.Fatalf("expected %s to be rejected", name)
			}
		})
	}
}

// Upstream #450: recolouring a point keeps the effects it already carried. The
// c:dPt run is rebuilt on every patch, so the shadow has to be carried across.
func TestPatchChartFormatting_DataPointKeepsEffects(t *testing.T) {
	shadow := `<a:effectLst><a:outerShdw blurRad="50800" dist="38100" dir="2700000">` +
		`<a:srgbClr val="000000"><a:alpha val="40000"/></a:srgbClr></a:outerShdw></a:effectLst>`
	base := strings.Replace(
		string(dataPointChartXML()),
		`<c:dLbls><c:showVal val="1"/></c:dLbls>`,
		`<c:dPt><c:idx val="1"/><c:spPr><a:solidFill><a:srgbClr val="4F81BD"/></a:solidFill>`+
			shadow+`</c:spPr></c:dPt><c:dLbls><c:showVal val="1"/></c:dLbls>`,
		1,
	)

	fill := "FF0000"
	got, err := PatchChartFormatting([]byte(base), common.ChartFormatUpdate{
		DataPoints: []common.ChartDataPoint{{PointIndex: 1, FillColor: &fill}},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	point := reDataPointBlock.FindString(string(got))
	if !strings.Contains(point, "<a:outerShdw") {
		t.Fatalf("shadow dropped when recolouring the point: %s", point)
	}
	if !strings.Contains(point, `<a:solidFill><a:srgbClr val="FF0000"/></a:solidFill>`) {
		t.Fatalf("new fill missing: %s", point)
	}
	// CT_ShapeProperties orders the fill before the effects.
	assertXMLOrder(t, point, "<a:solidFill>", "<a:effectLst>")
}
