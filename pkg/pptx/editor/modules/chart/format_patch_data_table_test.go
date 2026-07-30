package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func dataTableChartXML() []byte {
	// The series and the axes carry c:spPr of their own: only the plot area's
	// own one may anchor the data table.
	return []byte(`<c:chartSpace xmlns:c="x" xmlns:a="y"><c:chart><c:plotArea><c:layout/>` +
		`<c:barChart><c:ser><c:idx val="0"/>` +
		`<c:spPr><a:solidFill><a:srgbClr val="4F81BD"/></a:solidFill></c:spPr>` +
		`<c:val><c:numLit><c:ptCount val="1"/></c:numLit></c:val></c:ser>` +
		`<c:axId val="1"/><c:axId val="2"/></c:barChart>` +
		`<c:catAx><c:axId val="1"/><c:spPr><a:ln/></c:spPr></c:catAx>` +
		`<c:valAx><c:axId val="2"/><c:spPr><a:ln/></c:spPr></c:valAx>` +
		`<c:spPr><a:noFill/></c:spPr></c:plotArea></c:chart></c:chartSpace>`)
}

func TestPatchChartFormatting_DataTable(t *testing.T) {
	showKeys := true
	noVertical := false
	fontSize := 9

	req := common.ChartFormatUpdate{
		DataTable: &common.ChartDataTable{
			Show:               true,
			ShowKeys:           &showKeys,
			ShowVerticalBorder: &noVertical,
			FontSizePt:         &fontSize,
		},
	}
	got, err := PatchChartFormatting(dataTableChartXML(), req)
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	updated := string(got)
	for _, want := range []string{
		`<c:showHorzBorder val="1"/>`,
		`<c:showVertBorder val="0"/>`,
		`<c:showOutline val="1"/>`,
		`<c:showKeys val="1"/>`,
		`<a:defRPr sz="900"/>`,
	} {
		if !strings.Contains(updated, want) {
			t.Fatalf("updated XML missing %q: %s", want, updated)
		}
	}
	assertXMLOrder(
		t, reDataTableBlock.FindString(updated),
		"<c:showHorzBorder", "<c:showVertBorder", "<c:showOutline", "<c:showKeys", "<c:txPr",
	)
	// CT_PlotArea puts c:dTable after the axes and before the plot-area c:spPr,
	// never inside a series or an axis.
	assertXMLOrder(t, updated, "<c:catAx", "<c:valAx", "<c:dTable", "<a:noFill")
	if strings.Contains(reSerBlocks.FindString(updated), "<c:dTable") {
		t.Fatalf("data table was written inside the series: %s", updated)
	}
	if strings.Index(updated, "</c:barChart>") > strings.Index(updated, "<c:dTable") {
		t.Fatalf("data table must follow the chart group: %s", updated)
	}

	state := ExtractChartState(got)
	if state.DataTable == nil || !state.DataTable.Show {
		t.Fatalf("data table missing from state: %#v", state.DataTable)
	}
	if state.DataTable.ShowVerticalBorder == nil || *state.DataTable.ShowVerticalBorder {
		t.Fatalf("data table border flag not preserved: %#v", state.DataTable)
	}
	if state.DataTable.FontSizePt == nil || *state.DataTable.FontSizePt != fontSize {
		t.Fatalf("data table font size not preserved: %#v", state.DataTable)
	}

	again, err := PatchChartFormatting(got, req)
	if err != nil {
		t.Fatalf("re-patch error: %v", err)
	}
	if string(again) != updated {
		t.Fatalf("re-patch was not idempotent:\n%s\n%s", updated, string(again))
	}
}

func TestPatchChartFormattingDataTableKeepsUntouchedFlags(t *testing.T) {
	off := false
	seeded, err := PatchChartFormatting(dataTableChartXML(), common.ChartFormatUpdate{
		DataTable: &common.ChartDataTable{Show: true, ShowOutline: &off},
	})
	if err != nil {
		t.Fatalf("seed error: %v", err)
	}
	keys := false
	updated, err := PatchChartFormatting(seeded, common.ChartFormatUpdate{
		DataTable: &common.ChartDataTable{Show: true, ShowKeys: &keys},
	})
	if err != nil {
		t.Fatalf("second patch error: %v", err)
	}
	block := reDataTableBlock.FindString(string(updated))
	if !strings.Contains(block, `<c:showOutline val="0"/>`) {
		t.Fatalf("second patch dropped the earlier outline flag: %s", block)
	}
	if !strings.Contains(block, `<c:showKeys val="0"/>`) {
		t.Fatalf("second patch did not apply the keys flag: %s", block)
	}
}

func TestPatchChartFormattingRemovesDataTable(t *testing.T) {
	seeded, err := PatchChartFormatting(dataTableChartXML(), common.ChartFormatUpdate{
		DataTable: &common.ChartDataTable{Show: true},
	})
	if err != nil {
		t.Fatalf("seed error: %v", err)
	}
	cleared, err := PatchChartFormatting(seeded, common.ChartFormatUpdate{
		DataTable: &common.ChartDataTable{Show: false},
	})
	if err != nil {
		t.Fatalf("clear error: %v", err)
	}
	if strings.Contains(string(cleared), "<c:dTable") {
		t.Fatalf("expected the data table to be removed: %s", string(cleared))
	}
	if ExtractChartState(cleared).DataTable != nil {
		t.Fatal("expected cleared state to report no data table")
	}
}

func TestValidateChartFormatUpdateRejectsInvalidDataTableFont(t *testing.T) {
	tooBig := 500
	if err := ValidateChartFormatUpdate(common.ChartFormatUpdate{
		DataTable: &common.ChartDataTable{Show: true, FontSizePt: &tooBig},
	}); err == nil {
		t.Fatal("expected an out-of-range data table font size to be rejected")
	}
}
