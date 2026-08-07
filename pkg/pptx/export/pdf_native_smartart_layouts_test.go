package export

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func layoutDiagram(layout smartart.Layout, nodes ...smartart.Node) smartart.SmartArt {
	diagram := smartart.NewSmartArt(layout).
		Position(styling.Inches(0.5), styling.Inches(1)).
		Size(styling.Inches(9), styling.Inches(2))
	for _, node := range nodes {
		diagram = diagram.AddNode(node)
	}
	return diagram
}

func TestSmartArtLayoutNodesFoldChildrenIntoTheirParent(t *testing.T) {
	nodes := smartArtLayoutNodes([]smartart.Node{
		smartart.NewNode("Discover").
			WithChild(smartart.NewNode("Interviews")).
			WithChild(smartart.NewNode("Surveys")),
		smartart.NewNode("Deliver"),
	})

	if len(nodes) != 2 {
		t.Fatalf("got %d entries, want 2 — children must not become entries", len(nodes))
	}
	want := "Discover\nInterviews\nSurveys"
	if nodes[0].Text != want {
		t.Errorf("first entry text = %q, want %q", nodes[0].Text, want)
	}
	if len(nodes[0].Children) != 0 {
		t.Errorf("children left on the folded entry: %d", len(nodes[0].Children))
	}
	if nodes[1].Text != "Deliver" {
		t.Errorf("second entry text = %q, want Deliver", nodes[1].Text)
	}
}

func TestSmartArtLayoutNodesSkipEmptyCaptions(t *testing.T) {
	nodes := smartArtLayoutNodes([]smartart.Node{
		smartart.NewNode("  ").WithChild(smartart.NewNode("Only child")),
	})
	if len(nodes) != 1 || nodes[0].Text != "Only child" {
		t.Errorf("got %+v, want a single entry reading %q", nodes, "Only child")
	}
}

func TestLinearLayoutKeepsBoxesInsideACrowdedFrame(t *testing.T) {
	// The inter-node gap is a fraction of the frame width, so enough nodes used
	// to consume the whole frame and leave every box a negative width.
	nodes := make([]smartart.Node, 0, 40)
	for i := range 40 {
		nodes = append(nodes, smartart.NewNode(strings.Repeat("x", i%5+1)))
	}
	diagram := layoutDiagram(smartart.BasicProcess, nodes...)

	boxes, _ := layoutSmartArtLinear(diagram)
	if len(boxes) != len(nodes) {
		t.Fatalf("got %d boxes, want %d", len(boxes), len(nodes))
	}
	frameRight := emuToPt(int64(diagram.X)) + emuToPt(int64(diagram.CX))
	for i, box := range boxes {
		if box.W <= 0 || box.H <= 0 {
			t.Fatalf("box %d has non-positive size %vx%v", i, box.W, box.H)
		}
		if box.X+box.W > frameRight+1 {
			t.Fatalf("box %d ends at %v, past the frame's right edge %v", i, box.X+box.W, frameRight)
		}
	}
}

func TestPyramidTiersKeepAPositiveHeight(t *testing.T) {
	nodes := make([]smartart.Node, 0, 60)
	for range 60 {
		nodes = append(nodes, smartart.NewNode("tier"))
	}
	boxes, _ := layoutSmartArtPyramid(layoutDiagram(smartart.BasicPyramid, nodes...))
	if len(boxes) != len(nodes) {
		t.Fatalf("got %d boxes, want %d", len(boxes), len(nodes))
	}
	for i, box := range boxes {
		if box.H <= 0 {
			t.Fatalf("tier %d has height %v", i, box.H)
		}
	}
}

func TestGridLayoutKeepsPositiveCells(t *testing.T) {
	nodes := make([]smartart.Node, 0, 200)
	for range 200 {
		nodes = append(nodes, smartart.NewNode("cell"))
	}
	boxes, _ := layoutSmartArtGrid(layoutDiagram(smartart.BasicMatrix, nodes...))
	for i, box := range boxes {
		if box.W <= 0 || box.H <= 0 {
			t.Fatalf("cell %d has size %vx%v", i, box.W, box.H)
		}
	}
}

func TestGenericLayoutCarriesNodePictures(t *testing.T) {
	picture := []byte{0x89, 'P', 'N', 'G'}
	node := smartart.NewNode("With picture")
	node.ImageData = picture

	boxes, _ := layoutSmartArtLinear(layoutDiagram(smartart.BasicProcess, node))
	if len(boxes) != 1 {
		t.Fatalf("got %d boxes, want 1", len(boxes))
	}
	if len(boxes[0].Image) != len(picture) {
		t.Errorf("box carries %d picture bytes, want %d", len(boxes[0].Image), len(picture))
	}
}

func TestNodeColorPrefersItsOwnColour(t *testing.T) {
	own := smartart.NewNode("Coloured")
	own.Color = "FF0000"
	if got := smartArtNodeColor(own); got != "FF0000" {
		t.Errorf("got %q, want the node's own FF0000", got)
	}
	if got := smartArtNodeColor(smartart.NewNode("Plain")); got != smartArtNodeFill {
		t.Errorf("got %q, want the default accent %q", got, smartArtNodeFill)
	}
}
