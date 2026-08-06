package editor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

func TestGetSmartArtReadsBackTheNodeTree(t *testing.T) {
	const quickStyle = "urn:microsoft.com/office/officeart/2005/8/quickstyle/3d1"
	const colorStyle = "urn:microsoft.com/office/officeart/2005/8/colors/colorful1"

	diagram := smartart.NewSmartArt(smartart.PictureAccentList).
		WithQuickStyle(quickStyle).
		WithColorStyle(colorStyle).
		AddNode(smartart.NewNode("Topic A").WithChild(smartart.NewNode("Detail A"))).
		AddNode(smartart.NewNode("Topic B"))
	slide := elements.NewSlide("SmartArt Read").WithBlankLayout().AddSmartArt(diagram)

	data, err := pptx.CreateWithSlides("SmartArt Read", []elements.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides failed: %v", err)
	}
	path := filepath.Join(t.TempDir(), "smartart_read.pptx")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write deck: %v", err)
	}

	ed, err := editor.OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("OpenPresentationEditor failed: %v", err)
	}
	defer func() { _ = ed.Close() }()

	diagrams, err := ed.ListSmartArt(0)
	if err != nil {
		t.Fatalf("ListSmartArt failed: %v", err)
	}
	if len(diagrams) != 1 {
		t.Fatalf("expected one diagram on the slide, got %d", len(diagrams))
	}

	info := diagrams[0]
	if info.QuickStyleID != quickStyle {
		t.Errorf("quick style = %q, want %q", info.QuickStyleID, quickStyle)
	}
	if info.ColorStyleID != colorStyle {
		t.Errorf("color style = %q, want %q", info.ColorStyleID, colorStyle)
	}
	if info.LayoutURI != smartart.PictureAccentList.LayoutURI() {
		t.Errorf("layout = %q, want %q", info.LayoutURI, smartart.PictureAccentList.LayoutURI())
	}
	if len(info.Nodes) != 2 {
		t.Fatalf("expected two top-level nodes, got %d", len(info.Nodes))
	}
	if info.Nodes[0].Text != "Topic A" || info.Nodes[1].Text != "Topic B" {
		t.Errorf("node order = %q, %q", info.Nodes[0].Text, info.Nodes[1].Text)
	}
	if len(info.Nodes[0].Children) != 1 || info.Nodes[0].Children[0].Text != "Detail A" {
		t.Errorf("expected Detail A under Topic A, got %+v", info.Nodes[0].Children)
	}

	single, err := ed.GetSmartArt(0, info.ShapeID)
	if err != nil {
		t.Fatalf("GetSmartArt failed: %v", err)
	}
	if len(single.Nodes) != len(info.Nodes) {
		t.Errorf("GetSmartArt returned %d nodes, ListSmartArt %d", len(single.Nodes), len(info.Nodes))
	}
}
