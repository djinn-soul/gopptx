package elements

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

// A link on a run inside a shape's rich text needs a slide relationship like
// any other. Without one the renderer still writes <a:hlinkClick>, but with no
// r:id to resolve, so the link is inert.
func TestBuildSlideHyperlinkRelsRegistersRichTextRuns(t *testing.T) {
	link := action.NewHyperlink(action.HyperlinkURL("https://example.com/rich"))
	hover := action.NewHyperlink(action.HyperlinkURL("https://example.com/hover"))

	shape := shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 100, 100)
	shape.TextParagraphs = []text.Paragraph{{
		Runs: []text.Run{
			{Text: "click", Hyperlink: &link},
			{Text: "hover", HoverAction: &hover},
		},
	}}

	slide := NewSlide("Rich").AddShape(shape)
	rids, rels, _ := BuildSlideHyperlinkRels(slide, 2)

	if _, ok := rids[&link]; !ok {
		t.Error("the run's click link got no relationship ID")
	}
	if _, ok := rids[&hover]; !ok {
		t.Error("the run's hover link got no relationship ID")
	}
	if len(rels) != 2 {
		t.Fatalf("got %d relationships, want one per link", len(rels))
	}
}

// The same target reached from a shape and from one of its runs is one
// relationship, as it already was for shapes and bullets.
func TestBuildSlideHyperlinkRelsDedupesAcrossShapeAndRuns(t *testing.T) {
	shapeLink := action.NewHyperlink(action.HyperlinkURL("https://example.com/same"))
	runLink := action.NewHyperlink(action.HyperlinkURL("https://example.com/same"))

	shape := shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 100, 100)
	shape.ClickAction = &shapeLink
	shape.TextParagraphs = []text.Paragraph{{
		Runs: []text.Run{{Text: "click", Hyperlink: &runLink}},
	}}

	slide := NewSlide("Dedup").AddShape(shape)
	rids, rels, _ := BuildSlideHyperlinkRels(slide, 2)

	if len(rels) != 1 {
		t.Fatalf("got %d relationships for one target, want 1", len(rels))
	}
	if rids[&shapeLink] != rids[&runLink] {
		t.Errorf("same target got two IDs: %q and %q", rids[&shapeLink], rids[&runLink])
	}
}
