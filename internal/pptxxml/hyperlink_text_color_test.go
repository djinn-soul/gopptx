package pptxxml

import (
	"strings"
	"testing"
)

// PowerPoint repaints a hyperlink run in the theme's hlink colour unless the
// link carries Office's hyperlink-colour extension, which is why setting the
// font colour on a linked run looked impossible (upstream issue #940).
func TestRichTextRunXMLKeepsColorOnHyperlinkRun(t *testing.T) {
	run := TextRunSpec{
		Text:      "linked",
		Color:     "FF0000",
		Hyperlink: &HyperlinkSpec{RelID: "rId2"},
	}

	xml := RichTextRunXML(run, ContentStyleSpec{})

	if !strings.Contains(xml, `<a:solidFill><a:srgbClr val="FF0000"/></a:solidFill>`) {
		t.Fatalf("expected the run colour:\n%s", xml)
	}
	if !strings.Contains(xml, `<ahyp:hlinkClr xmlns:ahyp="`+hyperlinkColorNS+`" val="tx"/>`) {
		t.Fatalf("expected the hyperlink colour override:\n%s", xml)
	}
	if !strings.Contains(xml, `<a:hlinkClick r:id="rId2">`) ||
		!strings.Contains(xml, `</a:hlinkClick>`) {
		t.Fatalf("expected the override nested inside hlinkClick:\n%s", xml)
	}
}

// A linked run with no colour of its own must keep the theme's hyperlink
// colour, so the element stays empty.
func TestRichTextRunXMLLeavesUncolouredHyperlinkAlone(t *testing.T) {
	run := TextRunSpec{Text: "linked", Hyperlink: &HyperlinkSpec{RelID: "rId2"}}

	xml := RichTextRunXML(run, ContentStyleSpec{})

	if strings.Contains(xml, "hlinkClr") {
		t.Fatalf("expected no colour override on an uncoloured run:\n%s", xml)
	}
	if !strings.Contains(xml, `<a:hlinkClick r:id="rId2"/>`) {
		t.Fatalf("expected a self-closing hlinkClick:\n%s", xml)
	}
}
