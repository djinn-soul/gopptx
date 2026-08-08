package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// The reader records where a slide's title and body placeholders actually sit;
// the generator used to fall back to its built-in layout defaults and move them.
func TestGeneratorHonoursPlaceholderBounds(t *testing.T) {
	slide := elements.NewSlide("Bounded")
	slide = slide.AddBullet("Body")
	slide.TitleBoundsEMU = [4]int64{111111, 222222, 3333333, 444444}
	slide.ContentBoundsEMU = [4]int64{555555, 666666, 7777777, 888888}

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})
	slideXML := parts["ppt/slides/slide1.xml"]

	wants := []string{
		`<a:off x="111111" y="222222"/>`,
		`<a:ext cx="3333333" cy="444444"/>`,
		`<a:off x="555555" y="666666"/>`,
		`<a:ext cx="7777777" cy="888888"/>`,
	}
	for _, want := range wants {
		if !strings.Contains(slideXML, want) {
			t.Fatalf("slide1.xml missing %s: %s", want, slideXML)
		}
	}
}
