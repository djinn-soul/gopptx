package editor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

func TestChangeSmartArtLayoutPreservesHierarchyAndStyles(t *testing.T) {
	diagram := smartart.NewSmartArt(smartart.OrgChart).
		WithQuickStyle("urn:test:smartart:quickstyle").
		WithColorStyle("urn:test:smartart:colorstyle").
		AddNode(
			smartart.NewNode("CEO").
				WithChild(smartart.NewNode("Finance")).
				WithChild(smartart.NewNode("Engineering")),
		)
	slide := elements.NewSlide("SmartArt Edit").WithBlankLayout().AddSmartArt(diagram)

	data, err := pptx.CreateWithSlides("SmartArt Edit", []elements.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides failed: %v", err)
	}
	path := filepath.Join(t.TempDir(), "smartart_edit_roundtrip.pptx")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write deck: %v", err)
	}

	ed, err := editor.OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("OpenPresentationEditor failed: %v", err)
	}
	defer func() { _ = ed.Close() }()

	shapes, err := ed.GetShapes(0)
	if err != nil {
		t.Fatalf("GetShapes failed: %v", err)
	}
	shapeID := 0
	for _, shape := range shapes {
		if shape.Type == "graphicFrame" {
			shapeID = shape.ID
			break
		}
	}
	if shapeID == 0 {
		t.Fatal("expected SmartArt graphicFrame shape")
	}

	if err := ed.ChangeSmartArtLayout(0, shapeID, smartart.Hierarchy); err != nil {
		t.Fatalf("ChangeSmartArtLayout failed: %v", err)
	}

	updatedPath := filepath.Join(t.TempDir(), "smartart_edit_roundtrip_updated.pptx")
	if err := ed.Save(updatedPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	_, slides, err := export.SlidesFromPPTX(updatedPath)
	if err != nil {
		t.Fatalf("SlidesFromPPTX failed: %v", err)
	}
	if len(slides) != 1 || len(slides[0].SmartArtDiagrams) != 1 {
		t.Fatalf("expected one SmartArt diagram after save, got slides=%d diagrams=%d", len(slides), len(slides[0].SmartArtDiagrams))
	}

	read := slides[0].SmartArtDiagrams[0]
	if read.Layout != smartart.Hierarchy {
		t.Fatalf("expected hierarchy layout after edit, got %q", read.Layout)
	}
	if len(read.Nodes) != 1 || read.Nodes[0].Text != "CEO" {
		t.Fatalf("expected CEO root after layout change, got %+v", read.Nodes)
	}
	if len(read.Nodes[0].Children) != 2 {
		t.Fatalf("expected 2 child nodes after layout change, got %+v", read.Nodes[0].Children)
	}
	if read.Nodes[0].Children[0].Text != "Finance" || read.Nodes[0].Children[1].Text != "Engineering" {
		t.Fatalf("expected ordered child nodes after layout change, got %+v", read.Nodes[0].Children)
	}
	if read.QuickStyle != "urn:test:smartart:quickstyle" {
		t.Fatalf("expected quick style to be preserved, got %q", read.QuickStyle)
	}
	if read.ColorStyle != "urn:test:smartart:colorstyle" {
		t.Fatalf("expected color style to be preserved, got %q", read.ColorStyle)
	}
}
