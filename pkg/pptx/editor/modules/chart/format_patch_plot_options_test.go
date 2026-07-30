package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestPatchChartFormatting_DataLabelsAndBarOptions(t *testing.T) {
	format, grouping := "$#,##0", "stacked"
	gapWidth, overlap := 35, 100
	got, err := PatchChartFormatting([]byte(sampleChartXML), common.ChartFormatUpdate{
		DataLabelNumberFormat: &format,
		ChartGrouping:         &grouping,
		GapWidth:              &gapWidth,
		Overlap:               &overlap,
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	xml := string(got)
	for _, want := range []string{
		// Asking for a code unlinks the labels: a linked label takes its format
		// from the source cells, so PowerPoint drew the old numbers unchanged.
		`<c:numFmt formatCode="$#,##0" sourceLinked="0"/>`,
		`<c:grouping val="stacked"/>`,
		`<c:gapWidth val="35"/>`,
		`<c:overlap val="100"/>`,
	} {
		if !strings.Contains(xml, want) {
			t.Fatalf("expected XML to contain %s", want)
		}
	}
	assertXMLOrder(t, xml, "<c:grouping", "<c:ser>", "<c:dLbls", "<c:gapWidth", "<c:overlap", "<c:axId")
}

func TestPatchChartFormatting_DataLabelFormatLinkedWins(t *testing.T) {
	format, linked := "0.00", true
	got, err := PatchChartFormatting([]byte(sampleChartXML), common.ChartFormatUpdate{
		DataLabelNumberFormat: &format,
		DataLabelFormatLinked: &linked,
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	if !strings.Contains(string(got), `<c:numFmt formatCode="0.00" sourceLinked="1"/>`) {
		t.Fatalf("explicit format_linked not honoured: %s", got)
	}
}

func TestValidateChartFormatUpdateRejectsInvalidBarOptions(t *testing.T) {
	grouping := "diagonal"
	if err := ValidateChartFormatUpdate(common.ChartFormatUpdate{ChartGrouping: &grouping}); err == nil {
		t.Fatal("expected invalid grouping error")
	}
	gapWidth := 501
	if err := ValidateChartFormatUpdate(common.ChartFormatUpdate{GapWidth: &gapWidth}); err == nil {
		t.Fatal("expected invalid gap width error")
	}
}

func assertXMLOrder(t *testing.T, xml string, nodes ...string) {
	t.Helper()
	previous := -1
	for _, node := range nodes {
		index := strings.Index(xml, node)
		if index < 0 {
			t.Fatalf("XML missing ordered node %q: %s", node, xml)
		}
		if index <= previous {
			t.Fatalf("node %q is out of schema order in %s", node, xml)
		}
		previous = index
	}
}
