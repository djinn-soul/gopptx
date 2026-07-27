package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Regexp.ReplaceAllString expands "$1" and friends in the *replacement*, so any
// site embedding caller text has to replace literally. These cover the two
// non-chart features that took user text into a replacement.

func TestDefineTableStyleKeepsDollarInRedefinedName(t *testing.T) {
	e := newTableEditorFixture()
	e.parts.Set("[Content_Types].xml", []byte(contentTypesFixtureXML))
	e.parts.Set("ppt/_rels/presentation.xml.rels", []byte(presentationRelsFixtureXML))

	styleID, err := e.DefineTableStyle(common.TableStyleDefinition{Name: "plain"})
	if err != nil {
		t.Fatalf("DefineTableStyle (create) failed: %v", err)
	}
	// The create path inserts; only redefining takes the regexp update path.
	if _, err = e.DefineTableStyle(common.TableStyleDefinition{
		StyleID: styleID,
		Name:    "Q1 $1 Report",
	}); err != nil {
		t.Fatalf("DefineTableStyle (update) failed: %v", err)
	}

	styles, err := e.ListTableStyles()
	if err != nil {
		t.Fatalf("ListTableStyles failed: %v", err)
	}
	for _, style := range styles {
		if style.StyleID == styleID {
			if style.Name != "Q1 $1 Report" {
				t.Fatalf("redefined style name = %q, want %q", style.Name, "Q1 $1 Report")
			}
			return
		}
	}
	t.Fatalf("style %s missing from %#v", styleID, styles)
}

func TestRewriteFontBlockKeepsDollarInTypeface(t *testing.T) {
	source := `<a:fontScheme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
		`<a:majorFont><a:latin typeface="Calibri Light"/></a:majorFont>` +
		`<a:minorFont><a:latin typeface="Calibri"/></a:minorFont></a:fontScheme>`

	out, ok := rewriteFontBlock(source, majorFontPattern, "Gill Sans $1 Pro")
	if !ok {
		t.Fatal("rewriteFontBlock did not match the majorFont block")
	}
	if !strings.Contains(out, `typeface="Gill Sans $1 Pro"`) {
		t.Errorf("typeface lost its $-reference:\n%s", out)
	}
}
