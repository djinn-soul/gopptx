package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// PowerPoint writes these into every deck, and its own package validator lists
// them as required; gopptx wrote presProps only when print settings were set,
// and the other two never.
func TestGeneratedPackageCarriesStandardParts(t *testing.T) {
	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{elements.NewSlide("S1")})

	for _, want := range []string{"ppt/presProps.xml", "ppt/viewProps.xml", "ppt/tableStyles.xml"} {
		if _, ok := parts[want]; !ok {
			t.Fatalf("package is missing %s", want)
		}
		if !strings.Contains(parts["[Content_Types].xml"], `PartName="/`+want+`"`) {
			t.Fatalf("no content type declared for %s", want)
		}
		relTarget := `Target="` + strings.TrimPrefix(want, "ppt/") + `"`
		if !strings.Contains(parts["ppt/_rels/presentation.xml.rels"], relTarget) {
			t.Fatalf("no relationship declared for %s", want)
		}
	}

	// A table's <a:tableStyleId> now has a part to resolve against.
	if !strings.Contains(parts["ppt/tableStyles.xml"], "<a:tblStyleLst") {
		t.Fatalf("tableStyles.xml has no style list: %s", parts["ppt/tableStyles.xml"])
	}

	if !strings.Contains(parts["ppt/presentation.xml"], "<p:defaultTextStyle>") {
		t.Fatalf("presentation.xml has no defaultTextStyle: %s", parts["ppt/presentation.xml"])
	}
	if got := strings.Count(parts["ppt/presentation.xml"], "pPr marL="); got != 9 {
		t.Fatalf("defaultTextStyle level count = %d, want 9", got)
	}
}
