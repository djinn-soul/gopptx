package elements

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestToXMLBulletParagraphStyles_MapsLineSpacingPtsAndTabStops(t *testing.T) {
	styles := []ParagraphStyle{
		NewParagraphStyle().
			WithLineSpacingPts(18).
			WithTabStops(styling.Emu(914400)),
	}

	specs := ToXMLBulletParagraphStyles(styles)
	if len(specs) != 1 {
		t.Fatalf("expected 1 paragraph spec, got %d", len(specs))
	}
	if specs[0].LineSpacingPts != 18 {
		t.Fatalf("expected line spacing points 18, got %d", specs[0].LineSpacingPts)
	}
	if len(specs[0].TabStops) != 1 || specs[0].TabStops[0] != 914400 {
		t.Fatalf("expected tab stop 914400, got %+v", specs[0].TabStops)
	}
}
