package export

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestSlidesFromPPTX_PreservesSmartArtReaderData(t *testing.T) {
	diagram := smartart.NewSmartArt(smartart.BasicBlockList).
		WithAltText("SmartArt Alt").
		WithQuickStyle("urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1").
		WithColorStyle("urn:microsoft.com/office/officeart/2005/8/colors/accent1_2").
		Position(styling.Inches(1), styling.Inches(1.5)).
		Size(styling.Inches(5), styling.Inches(3)).
		AddItems([]string{"Plan", "Build", "Review"})
	slide := elements.NewSlide("SmartArt Reader").
		WithBlankLayout().
		AddSmartArt(diagram)

	data, err := pptx.CreateWithSlides("SmartArt Reader", []elements.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides failed: %v", err)
	}
	deckPath := filepath.Join(t.TempDir(), "reader_smartart.pptx")
	if err := os.WriteFile(deckPath, data, 0o600); err != nil {
		t.Fatalf("write deck: %v", err)
	}

	_, got, err := SlidesFromPPTX(deckPath)
	if err != nil {
		t.Fatalf("SlidesFromPPTX failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(got))
	}
	if len(got[0].SmartArtDiagrams) != 1 {
		t.Fatalf("expected 1 SmartArt diagram, got %d", len(got[0].SmartArtDiagrams))
	}
	if len(got[0].Shapes) != 0 {
		t.Fatalf("expected SmartArt graphicFrame to stay out of generic shapes, got %d shapes", len(got[0].Shapes))
	}
	read := got[0].SmartArtDiagrams[0]
	if read.AltText != "SmartArt Alt" {
		t.Fatalf("expected SmartArt alt text, got %q", read.AltText)
	}
	if read.Layout != smartart.BasicBlockList {
		t.Fatalf("expected BasicBlockList layout, got %q", read.Layout)
	}
	if read.QuickStyle != "urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1" {
		t.Fatalf("expected SmartArt quick style to round-trip, got %q", read.QuickStyle)
	}
	if read.ColorStyle != "urn:microsoft.com/office/officeart/2005/8/colors/accent1_2" {
		t.Fatalf("expected SmartArt color style to round-trip, got %q", read.ColorStyle)
	}
	if len(read.Nodes) != 3 || read.Nodes[0].Text != "Plan" || read.Nodes[2].Text != "Review" {
		t.Fatalf("expected SmartArt node texts to round-trip, got %+v", read.Nodes)
	}
}

func TestSlidesFromPPTX_PreservesSmartArtHierarchy(t *testing.T) {
	diagram := smartart.NewSmartArt(smartart.OrgChart).
		WithAltText("Org SmartArt Alt").
		Position(styling.Inches(1), styling.Inches(1.5)).
		Size(styling.Inches(5), styling.Inches(3)).
		AddNode(
			smartart.NewNode("CEO").
				WithChild(smartart.NewNode("Finance")).
				WithChild(smartart.NewNode("Engineering")),
		)
	slide := elements.NewSlide("SmartArt Hierarchy Reader").
		WithBlankLayout().
		AddSmartArt(diagram)

	data, err := pptx.CreateWithSlides("SmartArt Hierarchy Reader", []elements.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides failed: %v", err)
	}
	deckPath := filepath.Join(t.TempDir(), "reader_smartart_hierarchy.pptx")
	if err := os.WriteFile(deckPath, data, 0o600); err != nil {
		t.Fatalf("write deck: %v", err)
	}

	_, got, err := SlidesFromPPTX(deckPath)
	if err != nil {
		t.Fatalf("SlidesFromPPTX failed: %v", err)
	}
	if len(got) != 1 || len(got[0].SmartArtDiagrams) != 1 {
		t.Fatalf(
			"expected one slide with one SmartArt diagram, got slides=%d diagrams=%d",
			len(got),
			len(got[0].SmartArtDiagrams),
		)
	}
	if len(got[0].Shapes) != 0 {
		t.Fatalf(
			"expected hierarchy SmartArt graphicFrame to stay out of generic shapes, got %d shapes",
			len(got[0].Shapes),
		)
	}

	read := got[0].SmartArtDiagrams[0]
	if read.Layout != smartart.OrgChart {
		t.Fatalf("expected OrgChart layout, got %q", read.Layout)
	}
	if len(read.Nodes) != 1 {
		t.Fatalf("expected one root node, got %+v", read.Nodes)
	}
	root := read.Nodes[0]
	if root.Text != "CEO" {
		t.Fatalf("expected root CEO, got %+v", root)
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children under CEO, got %+v", root.Children)
	}
	if root.Children[0].Text != "Finance" || root.Children[1].Text != "Engineering" {
		t.Fatalf("expected ordered OrgChart children, got %+v", root.Children)
	}
}
