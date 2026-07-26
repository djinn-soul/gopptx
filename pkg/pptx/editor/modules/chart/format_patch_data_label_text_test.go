package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestPatchChartFormattingAddsDataLabelWordWrap(t *testing.T) {
	wordWrap := false
	got, err := PatchChartFormatting(
		[]byte(sampleChartXML),
		common.ChartFormatUpdate{DataLabelWordWrap: &wordWrap},
	)
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	xml := string(got)
	want := `<c:txPr><a:bodyPr wrap="none"/><a:lstStyle/>` +
		`<a:p><a:endParaRPr lang="en-US"/></a:p></c:txPr>`
	mustContain(t, xml, want)
	if strings.Index(xml, "<c:txPr>") > strings.Index(xml, "<c:showVal") {
		t.Fatalf("expected c:txPr before data-label visibility properties")
	}
}

func TestPatchChartFormattingUpdatesExistingDataLabelWordWrap(t *testing.T) {
	wordWrap := true
	input := strings.Replace(
		sampleChartXML,
		`<c:dLbls>`,
		`<c:dLbls><c:txPr><a:bodyPr rot="0" wrap="none"/>`+
			`<a:lstStyle/><a:p/></c:txPr>`,
		1,
	)
	got, err := PatchChartFormatting(
		[]byte(input),
		common.ChartFormatUpdate{DataLabelWordWrap: &wordWrap},
	)
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	xml := string(got)
	mustContain(t, xml, `<a:bodyPr rot="0" wrap="square"/>`)
	if strings.Count(xml, "<c:txPr>") != 1 {
		t.Fatalf("expected existing c:txPr to be updated in place")
	}
}
