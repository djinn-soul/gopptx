package chart

import (
	"math"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func trendlineChartXML() []byte {
	return []byte(`<c:chartSpace xmlns:c="x" xmlns:a="y"><c:chart><c:plotArea><c:barChart>` +
		`<c:ser><c:idx val="0"/><c:order val="0"/><c:dLbls><c:showVal val="1"/></c:dLbls>` +
		`<c:cat><c:strRef><c:f>Sheet1!$A$1</c:f></c:strRef></c:cat>` +
		`<c:val><c:numRef><c:f>Sheet1!$B$1</c:f></c:numRef></c:val></c:ser>` +
		`<c:ser><c:idx val="1"/><c:order val="1"/>` +
		`<c:cat><c:strRef><c:f>Sheet1!$A$1</c:f></c:strRef></c:cat>` +
		`<c:val><c:numRef><c:f>Sheet1!$C$1</c:f></c:numRef></c:val></c:ser>` +
		`</c:barChart></c:plotArea></c:chart></c:chartSpace>`)
}

func TestPatchChartFormatting_Trendline(t *testing.T) {
	name := "Fit"
	order := 3
	forward, backward := 2.0, 1.0
	showRSqr, showEq := true, true
	color, dash := "FF0000", lineDashPresetDash
	width := 19050

	got, err := PatchChartFormatting(trendlineChartXML(), common.ChartFormatUpdate{
		Trendlines: []common.ChartTrendline{{
			SeriesIndex:     0,
			Type:            "poly",
			Order:           &order,
			Name:            &name,
			Forward:         &forward,
			Backward:        &backward,
			DisplayRSquared: &showRSqr,
			DisplayEquation: &showEq,
			LineColor:       &color,
			LineWidthEMU:    &width,
			LineDash:        &dash,
		}},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	updated := string(got)
	for _, want := range []string{
		`<c:name>Fit</c:name>`,
		`<c:trendlineType val="poly"/>`,
		`<c:order val="3"/>`,
		`<c:forward val="2"/>`,
		`<c:backward val="1"/>`,
		`<c:dispRSqr val="1"/>`,
		`<c:dispEq val="1"/>`,
		`<a:srgbClr val="FF0000"/>`,
		`<a:prstDash val="dash"/>`,
		`w="19050"`,
	} {
		if !strings.Contains(updated, want) {
			t.Fatalf("updated XML missing %q: %s", want, updated)
		}
	}

	trendline := reTrendlineBlock.FindString(updated)
	assertXMLOrder(
		t, trendline,
		"<c:name", "<c:spPr", "<c:trendlineType", "<c:order",
		"<c:forward", "<c:backward", "<c:dispRSqr", "<c:dispEq",
	)

	series := strings.SplitAfter(updated, "</c:ser>")
	if !strings.Contains(series[0], "<c:trendline>") {
		t.Fatalf("first series missing trendline: %s", series[0])
	}
	if strings.Contains(series[1], "<c:trendline>") {
		t.Fatalf("trendline leaked onto the untargeted series: %s", series[1])
	}
	assertXMLOrder(t, series[0], "<c:dLbls", "<c:trendline", "<c:cat", "<c:val")

	state := ExtractChartState(got)
	if len(state.Trendlines) != 1 {
		t.Fatalf("expected one trendline in state, got %#v", state.Trendlines)
	}
	parsed := state.Trendlines[0]
	if parsed.Type != "poly" || parsed.Order == nil || *parsed.Order != order {
		t.Fatalf("trendline type/order not preserved in state: %#v", parsed)
	}
	if parsed.Name == nil || *parsed.Name != name {
		t.Fatalf("trendline name not preserved in state: %#v", parsed)
	}
	if parsed.LineColor == nil || *parsed.LineColor != color ||
		parsed.LineDash == nil || *parsed.LineDash != dash ||
		parsed.LineWidthEMU == nil || *parsed.LineWidthEMU != width {
		t.Fatalf("trendline line style not preserved in state: %#v", parsed)
	}

	again, err := PatchChartFormatting(got, common.ChartFormatUpdate{
		Trendlines: []common.ChartTrendline{{
			SeriesIndex:     0,
			Type:            "poly",
			Order:           &order,
			Name:            &name,
			Forward:         &forward,
			Backward:        &backward,
			DisplayRSquared: &showRSqr,
			DisplayEquation: &showEq,
			LineColor:       &color,
			LineWidthEMU:    &width,
			LineDash:        &dash,
		}},
	})
	if err != nil {
		t.Fatalf("re-patch error: %v", err)
	}
	if string(again) != updated {
		t.Fatalf("re-patch was not idempotent:\n%s\n%s", updated, string(again))
	}
}

func TestPatchChartFormattingWritesMovingAverageTrendlineOnSecondSeries(t *testing.T) {
	period := 3
	got, err := PatchChartFormatting(trendlineChartXML(), common.ChartFormatUpdate{
		Trendlines: []common.ChartTrendline{
			{SeriesIndex: 1, Type: "movingAvg", Period: &period},
		},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	series := strings.SplitAfter(string(got), "</c:ser>")
	if strings.Contains(series[0], "<c:trendline>") {
		t.Fatalf("trendline written to the wrong series: %s", series[0])
	}
	if !strings.Contains(series[1], `<c:period val="3"/>`) {
		t.Fatalf("second series missing moving-average trendline: %s", series[1])
	}
	assertXMLOrder(t, series[1], "<c:trendlineType", "<c:period", "<c:cat")
}

func TestPatchChartFormattingAppendsWithoutRebuildingExistingTrendline(t *testing.T) {
	existing := `<c:trendline><c:trendlineType val="linear"/><c:dispEq val="1"/>` +
		`<c:trendlineLbl><c:layout/><c:txPr><a:bodyPr/><a:lstStyle/>` +
		`<a:p><a:r><a:t>y = 2x</a:t></a:r></a:p></c:txPr></c:trendlineLbl>` +
		`<c:extLst><c:ext uri="preserve-me"/></c:extLst></c:trendline>`
	xml := strings.Replace(string(trendlineChartXML()), "<c:cat>", existing+"<c:cat>", 1)
	period := 2

	got, err := PatchChartFormatting([]byte(xml), common.ChartFormatUpdate{
		AppendTrendlines: []common.ChartTrendline{{
			SeriesIndex: 0,
			Type:        "movingAvg",
			Period:      &period,
		}},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	updated := string(got)
	if !strings.Contains(updated, existing) {
		t.Fatalf("append rebuilt or altered the existing trendline: %s", updated)
	}
	if count := strings.Count(updated, "<c:trendline>"); count != 2 {
		t.Fatalf("expected existing and appended trendlines, got %d: %s", count, updated)
	}
	assertXMLOrder(t, updated, existing, `<c:trendlineType val="movingAvg"/>`, "<c:cat>")
}

func TestValidateChartFormatUpdateRejectsInvalidTrendlines(t *testing.T) {
	order := 3
	period := 3
	bad := 9
	badColor := "nope"
	notANumber := math.NaN()
	infinity := math.Inf(1)
	cases := map[string]common.ChartTrendline{
		"unknown type":         {Type: "spline"},
		"order on non-poly":    {Type: "linear", Order: &order},
		"period on non-moving": {Type: "linear", Period: &period},
		"order out of range":   {Type: "poly", Order: &bad},
		"period below two":     {Type: "movingAvg", Period: new(int)},
		"negative series":      {Type: "linear", SeriesIndex: -1},
		"bad colour":           {Type: "linear", LineColor: &badColor},
		"bad dash":             {Type: "linear", LineDash: &badColor},
		"non-finite forward":   {Type: "linear", Forward: &notANumber},
		"non-finite backward":  {Type: "linear", Backward: &infinity},
		"non-finite intercept": {Type: "linear", Intercept: &notANumber},
	}
	for name, trendline := range cases {
		t.Run(name, func(t *testing.T) {
			err := ValidateChartFormatUpdate(common.ChartFormatUpdate{
				Trendlines: []common.ChartTrendline{trendline},
			})
			if err == nil {
				t.Fatalf("expected %s to be rejected", name)
			}
		})
	}
}

func TestPatchChartFormattingClearsTrendlines(t *testing.T) {
	seeded, err := PatchChartFormatting(trendlineChartXML(), common.ChartFormatUpdate{
		Trendlines: []common.ChartTrendline{{SeriesIndex: 0, Type: "linear"}},
	})
	if err != nil {
		t.Fatalf("seed error: %v", err)
	}
	cleared, err := PatchChartFormatting(seeded, common.ChartFormatUpdate{
		ClearTrendlineSeries: []int{0},
	})
	if err != nil {
		t.Fatalf("clear error: %v", err)
	}
	if strings.Contains(string(cleared), "<c:trendline>") {
		t.Fatalf("expected trendlines to be removed: %s", string(cleared))
	}
	if len(ExtractChartState(cleared).Trendlines) != 0 {
		t.Fatal("expected cleared state to report no trendlines")
	}
}
