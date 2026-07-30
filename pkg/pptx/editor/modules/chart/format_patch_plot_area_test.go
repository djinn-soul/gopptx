package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestPatchChartFormattingPlotAreaLine(t *testing.T) {
	color := "FF0000"
	width := 38100
	dash := lineDashPresetDash
	input := dataTableChartXML()

	got, err := PatchChartFormatting(input, common.ChartFormatUpdate{
		PlotAreaLine: &common.ChartLineFormat{
			Color: &color, WidthEMU: &width, Dash: &dash,
		},
	})
	if err != nil {
		t.Fatalf("patch plot-area line: %v", err)
	}
	xml := string(got)
	expected := `<a:ln w="38100"><a:solidFill><a:srgbClr val="FF0000"/>` +
		`</a:solidFill><a:prstDash val="dash"/></a:ln>`
	if !strings.Contains(xml, `<c:spPr><a:noFill/>`+expected+`</c:spPr></c:plotArea>`) {
		t.Fatalf("expected exact line in plot-area c:spPr: %s", xml)
	}
	if strings.Count(xml, `val="FF0000"`) != 1 {
		t.Fatalf("plot-area line leaked into a nested c:spPr: %s", xml)
	}
	state := ExtractChartState(got)
	if state.PlotAreaLine == nil ||
		state.PlotAreaLine.Color == nil || *state.PlotAreaLine.Color != color ||
		state.PlotAreaLine.WidthEMU == nil || *state.PlotAreaLine.WidthEMU != width ||
		state.PlotAreaLine.Dash == nil || *state.PlotAreaLine.Dash != dash {
		t.Fatalf("unexpected plot-area line state: %#v", state.PlotAreaLine)
	}
}

func TestPatchChartFormattingAddsAndHidesPlotAreaLine(t *testing.T) {
	none := true
	input := []byte(
		`<c:chartSpace><c:chart><c:plotArea><c:barChart>` +
			`<c:ser><c:spPr><a:ln/></c:spPr></c:ser></c:barChart>` +
			`</c:plotArea></c:chart></c:chartSpace>`,
	)
	got, err := PatchChartFormatting(input, common.ChartFormatUpdate{
		PlotAreaLine: &common.ChartLineFormat{None: &none},
	})
	if err != nil {
		t.Fatalf("hide plot-area line: %v", err)
	}
	xml := string(got)
	if !strings.Contains(
		xml,
		`</c:barChart><c:spPr><a:ln><a:noFill/></a:ln></c:spPr></c:plotArea>`,
	) {
		t.Fatalf("expected new plot-area c:spPr after chart groups: %s", xml)
	}
	if strings.Count(xml, "<a:noFill/>") != 1 {
		t.Fatalf("nested series c:spPr was modified: %s", xml)
	}
}
