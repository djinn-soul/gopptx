package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const titledChartXML = `<?xml version="1.0"?><c:chartSpace xmlns:c="http://x" xmlns:a="http://y">` +
	`<c:chart><c:title><c:tx><c:rich><a:bodyPr/><a:p><a:r><a:t>Revenue</a:t></a:r></a:p></c:rich></c:tx>` +
	`<c:overlay val="0"/></c:title><c:autoTitleDeleted val="0"/>` +
	`<c:plotArea><c:layout/></c:plotArea></c:chart></c:chartSpace>`

func floatPtrTL(v float64) *float64 { return &v }

// Issue #1030: a chart title needs a manual x/y so it can be pinned left.
func TestPatchChartFormattingSetsManualTitleLayout(t *testing.T) {
	out, err := PatchChartFormatting([]byte(titledChartXML), common.ChartFormatUpdate{
		TitleX: floatPtrTL(0.02),
		TitleY: floatPtrTL(0.03),
	})
	if err != nil {
		t.Fatalf("patch: %v", err)
	}

	got := string(out)
	for _, want := range []string{
		`<c:manualLayout>`, `<c:xMode val="edge"/>`, `<c:yMode val="edge"/>`,
		`<c:x val="0.02"/>`, `<c:y val="0.03"/>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %s in %s", want, got)
		}
	}

	// CT_Title orders tx, layout, overlay.
	title := got[strings.Index(got, "<c:title>"):strings.Index(got, "</c:title>")]
	tx := strings.Index(title, "</c:tx>")
	layout := strings.Index(title, "<c:layout>")
	overlay := strings.Index(title, "<c:overlay")
	if tx < 0 || layout < 0 || overlay < 0 {
		t.Fatalf("missing title children: %s", title)
	}
	if tx > layout || layout > overlay {
		t.Fatalf("title children out of schema order: %s", title)
	}
}

func TestPatchChartTitleLayoutKeepsUnsetAxis(t *testing.T) {
	first, err := PatchChartFormatting([]byte(titledChartXML), common.ChartFormatUpdate{
		TitleX: floatPtrTL(0.1),
		TitleY: floatPtrTL(0.2),
	})
	if err != nil {
		t.Fatalf("first patch: %v", err)
	}
	second, err := PatchChartFormatting(first, common.ChartFormatUpdate{
		TitleX: floatPtrTL(0.5),
	})
	if err != nil {
		t.Fatalf("second patch: %v", err)
	}

	got := string(second)
	if !strings.Contains(got, `<c:x val="0.5"/>`) {
		t.Fatalf("x not updated: %s", got)
	}
	if !strings.Contains(got, `<c:y val="0.2"/>`) {
		t.Fatalf("y should be preserved: %s", got)
	}
	if strings.Count(got, "<c:manualLayout>") != 1 {
		t.Fatalf("layout duplicated: %s", got)
	}
}

func TestPatchChartTitleLayoutWithTitleText(t *testing.T) {
	title := "Left aligned"
	out, err := PatchChartFormatting([]byte(titledChartXML), common.ChartFormatUpdate{
		Title:  &title,
		TitleX: floatPtrTL(0.01),
	})
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "<a:t>Left aligned</a:t>") {
		t.Fatalf("title text missing: %s", got)
	}
	if !strings.Contains(got, `<c:x val="0.01"/>`) {
		t.Fatalf("layout missing: %s", got)
	}
}
