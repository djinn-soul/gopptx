package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// TestHyperlinkWritePathSanitizesScriptProtocols asserts that script protocols
// never reach the relationship target. Both the shape click-action path and the
// text-run path share getOrCreateHyperlinkRelID, so filtering there covers both.
func TestHyperlinkWritePathSanitizesScriptProtocols(t *testing.T) {
	cases := []struct {
		name    string
		address string
		want    string
	}{
		{name: "javascript", address: "javascript:alert(1)", want: "#"},
		{name: "javascript mixed case", address: "JavaScript:alert(1)", want: "#"},
		{name: "javascript leading space", address: "  javascript:alert(1)", want: "#"},
		{name: "vbscript", address: "vbscript:msgbox(1)", want: "#"},
		{name: "data html", address: "data:text/html;base64,PHNjcmlwdD4=", want: "#"},
		{name: "https preserved", address: "https://example.com/a?b=c", want: "https://example.com/a?b=c"},
		{name: "mailto preserved", address: "mailto:a@example.com", want: "mailto:a@example.com"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editor := newHyperlinkSanitizeEditor(t, "hyperlink-sanitize.pptx")
			defer func() { _ = editor.Close() }()

			const partPath = "ppt/slides/slide1.xml"
			relID, err := editor.getOrCreateHyperlinkRelID(partPath, tc.address)
			if err != nil {
				t.Fatalf("getOrCreateHyperlinkRelID(%q): %v", tc.address, err)
			}

			relsXML := readHyperlinkSanitizeRels(t, editor, partPath)
			if !strings.Contains(relsXML, `Id="`+relID+`"`) {
				t.Fatalf("relationship %s missing from rels:\n%s", relID, relsXML)
			}
			if !strings.Contains(relsXML, `Target="`+tc.want+`"`) {
				t.Fatalf("want Target=%q in rels, got:\n%s", tc.want, relsXML)
			}
			if tc.want == "#" && strings.Contains(strings.ToLower(relsXML), "alert") {
				t.Fatalf("script payload leaked into rels:\n%s", relsXML)
			}
		})
	}
}

// TestHyperlinkWritePathDedupesSanitizedTargets asserts sanitization happens
// before the reuse lookup, so two blocked addresses collapse onto a single
// relationship instead of accumulating duplicates.
func TestHyperlinkWritePathDedupesSanitizedTargets(t *testing.T) {
	editor := newHyperlinkSanitizeEditor(t, "hyperlink-sanitize-dedupe.pptx")
	defer func() { _ = editor.Close() }()

	const partPath = "ppt/slides/slide1.xml"
	first, err := editor.getOrCreateHyperlinkRelID(partPath, "javascript:alert(1)")
	if err != nil {
		t.Fatalf("first getOrCreateHyperlinkRelID: %v", err)
	}
	second, err := editor.getOrCreateHyperlinkRelID(partPath, "vbscript:msgbox(2)")
	if err != nil {
		t.Fatalf("second getOrCreateHyperlinkRelID: %v", err)
	}
	if first != second {
		t.Fatalf("sanitized targets should share one relationship, got %s and %s", first, second)
	}
}

func newHyperlinkSanitizeEditor(t *testing.T, name string) *PresentationEditor {
	t.Helper()

	base := writeDeckFixture(t, name, []elements.SlideContent{
		elements.NewSlide("Hyperlinks"),
	})
	editor, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	return editor
}

func readHyperlinkSanitizeRels(t *testing.T, editor *PresentationEditor, partPath string) string {
	t.Helper()

	relsPath := common.SlideRelsPartName(partPath)
	data, ok := editor.parts.Get(relsPath)
	if !ok {
		t.Fatalf("rels part %s not found", relsPath)
	}
	return string(data)
}
