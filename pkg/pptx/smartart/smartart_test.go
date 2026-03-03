package smartart

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestLayout_Name(t *testing.T) {
	tests := []struct {
		layout   Layout
		expected string
	}{
		{BasicBlockList, "Basic Block List"},
		{BasicProcess, "Basic Process"},
		{BasicCycle, "Basic Cycle"},
		{OrgChart, "Organization Chart"},
		{BasicVenn, "Basic Venn"},
		{PictureStrips, "Picture Strips"},
		{Layout("unknown"), "unknown"},
	}

	for _, tt := range tests {
		if tt.layout.Name() != tt.expected {
			t.Errorf("Layout(%q).Name() = %q; want %q", tt.layout, tt.layout.Name(), tt.expected)
		}
	}
}

func TestNode(t *testing.T) {
	node := NewNode("parent").
		WithChild(NewNode("child")).
		WithColor("FF0000")

	if node.Text != "parent" {
		t.Errorf("expected text 'parent', got '%s'", node.Text)
	}
	if len(node.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(node.Children))
	}
	if node.Children[0].Text != "child" {
		t.Errorf("expected child text 'child', got '%s'", node.Children[0].Text)
	}
	if node.Color != "FF0000" {
		t.Errorf("expected color 'FF0000', got '%s'", node.Color)
	}
}

func TestSmartArt(t *testing.T) {
	sa := NewSmartArt(BasicBlockList).
		WithAltText("alt text").
		WithDecorative(true).
		AddNode(NewNode("node1")).
		AddItems([]string{"node2", "node3"}).
		Position(styling.Emu(100), styling.Emu(200)).
		Size(styling.Emu(300), styling.Emu(400)).
		WithColorStyle("colorful1").
		WithQuickStyle("simple")

	if sa.Layout != BasicBlockList {
		t.Errorf("expected layout %s, got %s", BasicBlockList, sa.Layout)
	}
	if sa.AltText != "alt text" {
		t.Errorf("expected alt text 'alt text', got '%s'", sa.AltText)
	}
	if !sa.IsDecorative {
		t.Error("expected decorative to be true")
	}
	if len(sa.Nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(sa.Nodes))
	}
	if sa.X != 100 || sa.Y != 200 {
		t.Errorf("expected position (100, 200), got (%d, %d)", sa.X, sa.Y)
	}
	if sa.CX != 300 || sa.CY != 400 {
		t.Errorf("expected size (300, 400), got (%d, %d)", sa.CX, sa.CY)
	}
	if sa.ColorStyle != "colorful1" {
		t.Errorf("expected color style 'colorful1', got '%s'", sa.ColorStyle)
	}
	if sa.QuickStyle != "simple" {
		t.Errorf("expected quick style 'simple', got '%s'", sa.QuickStyle)
	}
}

func TestSmartArt_ToSpec(t *testing.T) {
	sa := NewSmartArt(BasicBlockList).
		AddNode(NewNode("node1").WithChild(NewNode("child1")))

	spec := sa.ToSpec()
	if spec.LayoutURI != string(BasicBlockList) {
		t.Errorf("expected layout URI %s, got %s", BasicBlockList, spec.LayoutURI)
	}
	if len(spec.Nodes) != 1 {
		t.Errorf("expected 1 node spec, got %d", len(spec.Nodes))
	}
	if len(spec.Nodes[0].Children) != 1 {
		t.Errorf("expected 1 child spec, got %d", len(spec.Nodes[0].Children))
	}
}

func TestSmartArt_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sa      SmartArt
		wantErr bool
	}{
		{"Valid", NewSmartArt(BasicBlockList).AddNode(NewNode("A")), false},
		{"Empty Layout", SmartArt{Nodes: []Node{{Text: "A"}}}, true},
		{"No Nodes", NewSmartArt(BasicBlockList), true},
		{"Invalid Node", NewSmartArt(BasicBlockList).AddNode(Node{}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sa.Validate(1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSmartArt_LayoutNameHelpers(t *testing.T) {
	name, ok := processLayoutName(BasicProcess)
	if !ok || name != "Basic Process" {
		t.Error("process failed")
	}

	name, ok = relationshipLayoutName(BasicRadial)
	if !ok || name != "Basic Radial" {
		t.Error("rel failed")
	}

	name, ok = matrixPictureLayoutName(PictureStrips)
	if !ok || name != "Picture Strips" {
		t.Error("matrix failed")
	}
}
