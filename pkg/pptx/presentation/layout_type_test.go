package presentation

import (
	"regexp"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// ST_SlideLayoutType on p:sldLayout/@type is how PowerPoint classifies a layout
// in the New Slide gallery and matches placeholders when one is applied. Every
// layout gopptx generated was untyped.
func TestGeneratedLayoutsAreTyped(t *testing.T) {
	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{elements.NewSlide("S1")})

	wantTypes := map[string]string{
		"ppt/slideLayouts/slideLayout1.xml": "obj",
		"ppt/slideLayouts/slideLayout2.xml": "titleOnly",
		"ppt/slideLayouts/slideLayout3.xml": "blank",
		"ppt/slideLayouts/slideLayout4.xml": "title",
		"ppt/slideLayouts/slideLayout5.xml": "obj",
		"ppt/slideLayouts/slideLayout6.xml": "twoObj",
		"ppt/slideLayouts/slideLayout7.xml": "vertTx",
		"ppt/slideLayouts/slideLayout8.xml": "vertTitleAndTx",
	}
	for part, want := range wantTypes {
		xml, ok := parts[part]
		if !ok {
			t.Fatalf("package is missing %s", part)
		}
		if !strings.Contains(xml, ` type="`+want+`"`) {
			t.Fatalf("%s is not typed %q: %s", part, want, xml)
		}
		if !strings.Contains(parts["[Content_Types].xml"], `PartName="/`+part+`"`) {
			t.Fatalf("no content type declared for %s", part)
		}
	}

	// The master has to name all eight, and the rels file has to reach them.
	master := parts["ppt/slideMasters/slideMaster1.xml"]
	if got := strings.Count(master, "<p:sldLayoutId "); got != len(wantTypes) {
		t.Fatalf("master declares %d layouts, want %d", got, len(wantTypes))
	}
	rels := parts["ppt/slideMasters/_rels/slideMaster1.xml.rels"]
	layoutRels := regexp.MustCompile(`Target="\.\./slideLayouts/slideLayout\d+\.xml"`).FindAllString(rels, -1)
	if len(layoutRels) != len(wantTypes) {
		t.Fatalf("master rels reach %d layouts, want %d", len(layoutRels), len(wantTypes))
	}
}

// The two vertical-text layouts are selectable, and their body is written down
// the page.
func TestVerticalTextLayoutWritesVerticalBody(t *testing.T) {
	slide := elements.NewSlide("Vertical")
	slide.Layout = elements.SlideLayoutTitleAndVerticalText
	slide = slide.AddBullet("Body")

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})
	if !strings.Contains(parts["ppt/slides/slide1.xml"], `vert="eaVert"`) {
		t.Fatalf("body is not vertical: %s", parts["ppt/slides/slide1.xml"])
	}
	if !strings.Contains(parts["ppt/slides/_rels/slide1.xml.rels"], "slideLayout7.xml") {
		t.Fatalf("slide does not use the vertical-text layout: %s", parts["ppt/slides/_rels/slide1.xml.rels"])
	}
}
