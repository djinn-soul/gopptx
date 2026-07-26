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
		`<c:numFmt formatCode="$#,##0" sourceLinked="1"/>`,
		`<c:grouping val="stacked"/>`,
		`<c:gapWidth val="35"/>`,
		`<c:overlap val="100"/>`,
	} {
		if !strings.Contains(xml, want) {
			t.Fatalf("expected XML to contain %s", want)
		}
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
