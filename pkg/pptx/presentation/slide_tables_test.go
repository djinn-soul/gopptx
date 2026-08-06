package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

func tableWith(label string, y styling.Length) tables.Table {
	return tables.NewTable([]styling.Length{styling.Inches(2)}).
		AddRow([]string{label}).
		Position(styling.Inches(1), y)
}

// SlideContent.Tables exists so extra tables survive a round trip, and the
// reader fills it — but the generator never read it, so three tables in gave
// one table out.
func TestGeneratorEmitsTableOverflow(t *testing.T) {
	slide := elements.NewSlide("Tables")
	first := tableWith("FIRST", styling.Inches(1))
	slide.Table = &first
	slide.Tables = []tables.Table{
		tableWith("SECOND", styling.Inches(3)),
		tableWith("THIRD", styling.Inches(5)),
	}

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})
	slideXML := parts["ppt/slides/slide1.xml"]

	for _, want := range []string{"FIRST", "SECOND", "THIRD"} {
		if !strings.Contains(slideXML, want) {
			t.Fatalf("slide1.xml is missing %q: %s", want, slideXML)
		}
	}
	if got := strings.Count(slideXML, "<a:tbl>"); got != 3 {
		t.Fatalf("a:tbl count = %d, want 3", got)
	}
}
