package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestPatchChartFormatting_ErrorBars(t *testing.T) {
	value := 1.5
	noEndCap := true
	color := "0000FF"

	got, err := PatchChartFormatting(trendlineChartXML(), common.ChartFormatUpdate{
		ErrorBars: []common.ChartErrorBars{{
			SeriesIndex: 0,
			BarType:     "both",
			ValueType:   "fixedVal",
			Value:       &value,
			NoEndCap:    &noEndCap,
			LineColor:   &color,
		}},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	updated := string(got)
	for _, want := range []string{
		`<c:errBarType val="both"/>`,
		`<c:errValType val="fixedVal"/>`,
		`<c:noEndCap val="1"/>`,
		`<c:val val="1.5"/>`,
		`<a:srgbClr val="0000FF"/>`,
	} {
		if !strings.Contains(updated, want) {
			t.Fatalf("updated XML missing %q: %s", want, updated)
		}
	}
	assertXMLOrder(
		t, reErrBarsBlock.FindString(updated),
		"<c:errBarType", "<c:errValType", "<c:noEndCap", "<c:val", "<c:spPr",
	)

	series := strings.SplitAfter(updated, "</c:ser>")
	assertXMLOrder(t, series[0], "<c:dLbls", "<c:errBars", "<c:cat", "</c:ser>")
	if strings.Contains(series[1], "<c:errBars>") {
		t.Fatalf("error bars leaked onto the untargeted series: %s", series[1])
	}

	state := ExtractChartState(got)
	if len(state.ErrorBars) != 1 {
		t.Fatalf("expected one error bar set in state, got %#v", state.ErrorBars)
	}
	parsed := state.ErrorBars[0]
	if parsed.BarType != "both" || parsed.ValueType != "fixedVal" {
		t.Fatalf("error bar types not preserved in state: %#v", parsed)
	}
	if parsed.Value == nil || *parsed.Value != value ||
		parsed.NoEndCap == nil || !*parsed.NoEndCap ||
		parsed.LineColor == nil || *parsed.LineColor != color {
		t.Fatalf("error bar details not preserved in state: %#v", parsed)
	}

	again, err := PatchChartFormatting(got, common.ChartFormatUpdate{
		ErrorBars: []common.ChartErrorBars{{
			SeriesIndex: 0,
			BarType:     "both",
			ValueType:   "fixedVal",
			Value:       &value,
			NoEndCap:    &noEndCap,
			LineColor:   &color,
		}},
	})
	if err != nil {
		t.Fatalf("re-patch error: %v", err)
	}
	if string(again) != updated {
		t.Fatalf("re-patch was not idempotent:\n%s\n%s", updated, string(again))
	}
}

func TestPatchChartFormattingWritesCustomAndDirectionalErrorBars(t *testing.T) {
	plus := "Sheet1!$D$2:$D$5"
	minus := "Sheet1!$E$2:$E$5"
	dirX, dirY := "x", "y"

	got, err := PatchChartFormatting(trendlineChartXML(), common.ChartFormatUpdate{
		ErrorBars: []common.ChartErrorBars{
			{
				SeriesIndex: 0, BarType: "both", ValueType: "cust",
				Direction: &dirX, PlusReference: &plus, MinusReference: &minus,
			},
			{SeriesIndex: 0, BarType: "plus", ValueType: "stdDev", Direction: &dirY},
		},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	updated := string(got)
	if !strings.Contains(updated, `<c:plus><c:numRef><c:f>Sheet1!$D$2:$D$5</c:f></c:numRef></c:plus>`) {
		t.Fatalf("custom plus reference missing: %s", updated)
	}
	blocks := reErrBarsBlock.FindAllString(updated, -1)
	if len(blocks) != 2 {
		t.Fatalf("expected separate X and Y error bars, got %d: %s", len(blocks), updated)
	}
	assertXMLOrder(t, blocks[0], "<c:errDir", "<c:errBarType", "<c:errValType", "<c:plus", "<c:minus")

	state := ExtractChartState(got)
	if len(state.ErrorBars) != 2 {
		t.Fatalf("expected two error bar sets in state, got %#v", state.ErrorBars)
	}
	if state.ErrorBars[0].PlusReference == nil || *state.ErrorBars[0].PlusReference != plus {
		t.Fatalf("custom reference not preserved in state: %#v", state.ErrorBars[0])
	}
	if state.ErrorBars[1].Direction == nil || *state.ErrorBars[1].Direction != dirY {
		t.Fatalf("error bar direction not preserved in state: %#v", state.ErrorBars[1])
	}
}

func TestPatchChartFormattingClearsErrorBars(t *testing.T) {
	value := 2.0
	seeded, err := PatchChartFormatting(trendlineChartXML(), common.ChartFormatUpdate{
		ErrorBars: []common.ChartErrorBars{
			{SeriesIndex: 0, BarType: "both", ValueType: "fixedVal", Value: &value},
		},
	})
	if err != nil {
		t.Fatalf("seed error: %v", err)
	}
	cleared, err := PatchChartFormatting(seeded, common.ChartFormatUpdate{
		ClearErrorBarSeries: []int{0},
	})
	if err != nil {
		t.Fatalf("clear error: %v", err)
	}
	if strings.Contains(string(cleared), "<c:errBars>") {
		t.Fatalf("expected error bars to be removed: %s", string(cleared))
	}
}

func TestPatchChartFormattingOrdersTrendlineBeforeErrorBars(t *testing.T) {
	value := 1.0
	got, err := PatchChartFormatting(trendlineChartXML(), common.ChartFormatUpdate{
		Trendlines: []common.ChartTrendline{{SeriesIndex: 0, Type: "linear"}},
		ErrorBars: []common.ChartErrorBars{
			{SeriesIndex: 0, BarType: "both", ValueType: "fixedVal", Value: &value},
		},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	series := strings.SplitAfter(string(got), "</c:ser>")
	assertXMLOrder(t, series[0], "<c:dLbls", "<c:trendline", "<c:errBars", "<c:cat", "<c:val>")
}

func TestValidateChartFormatUpdateRejectsInvalidErrorBars(t *testing.T) {
	negative := -1.0
	reference := "Sheet1!$D$2"
	badDirection := "z"
	badColor := "nope"
	cases := map[string]common.ChartErrorBars{
		"unknown bar type":     {BarType: "up", ValueType: "fixedVal"},
		"unknown value type":   {BarType: "both", ValueType: "guess"},
		"unknown direction":    {BarType: "both", ValueType: "stdDev", Direction: &badDirection},
		"custom without refs":  {BarType: "both", ValueType: "cust"},
		"refs without custom":  {BarType: "both", ValueType: "stdDev", PlusReference: &reference},
		"negative fixed value": {BarType: "both", ValueType: "fixedVal", Value: &negative},
		"negative series":      {BarType: "both", ValueType: "stdDev", SeriesIndex: -1},
		"bad colour":           {BarType: "both", ValueType: "stdDev", LineColor: &badColor},
	}
	for name, bars := range cases {
		t.Run(name, func(t *testing.T) {
			err := ValidateChartFormatUpdate(common.ChartFormatUpdate{
				ErrorBars: []common.ChartErrorBars{bars},
			})
			if err == nil {
				t.Fatalf("expected %s to be rejected", name)
			}
		})
	}
}

// PowerPoint only applies the c:spPr line style to every data point when
// c:noEndCap is present, so it is written even at its schema default.
func TestPatchChartFormattingAlwaysWritesNoEndCap(t *testing.T) {
	value := 1.5
	color := "00B050"
	got, err := PatchChartFormatting(trendlineChartXML(), common.ChartFormatUpdate{
		ErrorBars: []common.ChartErrorBars{
			{SeriesIndex: 0, BarType: "both", ValueType: "fixedVal", Value: &value, LineColor: &color},
		},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	block := reErrBarsBlock.FindString(string(got))
	if !strings.Contains(block, `<c:noEndCap val="0"/>`) {
		t.Fatalf("noEndCap must be written even when unset: %s", block)
	}
	assertXMLOrder(t, block, "<c:errValType", "<c:noEndCap", "<c:val", "<c:spPr")

	// An explicit true must still win.
	capless := true
	got, err = PatchChartFormatting(trendlineChartXML(), common.ChartFormatUpdate{
		ErrorBars: []common.ChartErrorBars{
			{SeriesIndex: 0, BarType: "both", ValueType: "stdDev", NoEndCap: &capless},
		},
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	if !strings.Contains(reErrBarsBlock.FindString(string(got)), `<c:noEndCap val="1"/>`) {
		t.Fatalf("explicit noEndCap not honoured: %s", string(got))
	}
}
