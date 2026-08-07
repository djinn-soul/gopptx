package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Code could only reach a slide as bullets in the body placeholder, so there
// was no way to put a listing at chosen coordinates.
func TestCodeBlockIsPlacedAtItsOwnCoordinates(t *testing.T) {
	block := elements.NewCodeBlock(
		"func main() {\n\tprintln(\"hi\")\n}",
		"go",
		styling.Inches(1), styling.Inches(2), styling.Inches(6), styling.Inches(3),
	).WithLanguageLabel(true)

	slide := elements.NewSlide("Code").AddCodeBlock(block)
	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})
	slideXML := parts["ppt/slides/slide1.xml"]

	if !strings.Contains(slideXML, `<a:off x="914400" y="1828800"/>`) {
		t.Fatalf("code block is not at its stated position: %s", slideXML)
	}
	if !strings.Contains(slideXML, "func") || !strings.Contains(slideXML, "main") {
		t.Fatalf("code text missing: %s", slideXML)
	}
	if !strings.Contains(slideXML, "[GO]") {
		t.Fatalf("language label missing: %s", slideXML)
	}
	// The listing sits on the dark background the palette is designed against.
	if !strings.Contains(slideXML, elements.DefaultCodeBackground) {
		t.Fatalf("code background missing: %s", slideXML)
	}
}

func TestCodeBlockKeepsOneParagraphPerLine(t *testing.T) {
	block := elements.NewCodeBlock(
		"one\n\nthree",
		"text",
		styling.Inches(1), styling.Inches(1), styling.Inches(4), styling.Inches(2),
	)
	slide := elements.NewSlide("Code").AddCodeBlock(block)

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})
	slideXML := parts["ppt/slides/slide1.xml"]

	for _, want := range []string{"one", "three"} {
		if !strings.Contains(slideXML, want) {
			t.Fatalf("line %q missing: %s", want, slideXML)
		}
	}
	// Three source lines, so three paragraphs — the blank one included.
	shapeStart := strings.Index(slideXML, elements.DefaultCodeBackground)
	if shapeStart < 0 {
		t.Fatal("code shape not found")
	}
	if got := strings.Count(slideXML[shapeStart:], "<a:p>"); got < 3 {
		t.Fatalf("paragraph count = %d, want at least 3", got)
	}
}
