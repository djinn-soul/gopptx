// Package smartart provides SmartArt (DiagramML) generation for presentations.
//
// SmartArt diagrams consist of 5 OOXML parts (data, layout, colors, quickStyle,
// drawing) embedded via a p:graphicFrame on the slide. PowerPoint resolves the
// actual shape layout at render time from the layout URI we reference.
package smartart

import (
	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Layout identifies a SmartArt layout type by its OOXML URN.
type Layout string

// Phase-1 supported SmartArt layouts.
const (
	// List layouts.
	BasicBlockList      Layout = "urn:microsoft.com/office/officeart/2005/8/layout/vList1"
	VerticalBlockList   Layout = "urn:microsoft.com/office/officeart/2005/8/layout/vList2"
	HorizontalBulletLst Layout = "urn:microsoft.com/office/officeart/2005/8/layout/hList1"
	SquareAccentList    Layout = "urn:microsoft.com/office/officeart/2005/8/layout/vList3"
	PictureAccentList   Layout = "urn:microsoft.com/office/officeart/2005/8/layout/vList5"

	// Process layouts.
	BasicProcess           Layout = "urn:microsoft.com/office/officeart/2005/8/layout/process1"
	AccentProcess          Layout = "urn:microsoft.com/office/officeart/2005/8/layout/process2"
	AlternatingFlow        Layout = "urn:microsoft.com/office/officeart/2005/8/layout/process3"
	ContinuousBlockProcess Layout = "urn:microsoft.com/office/officeart/2005/8/layout/process4"

	// Cycle layouts.
	BasicCycle Layout = "urn:microsoft.com/office/officeart/2005/8/layout/cycle1"
	TextCycle  Layout = "urn:microsoft.com/office/officeart/2005/8/layout/cycle2"
	BlockCycle Layout = "urn:microsoft.com/office/officeart/2005/8/layout/cycle3"

	// Hierarchy layouts.
	OrgChart            Layout = "urn:microsoft.com/office/officeart/2005/8/layout/orgChart1"
	Hierarchy           Layout = "urn:microsoft.com/office/officeart/2005/8/layout/hierarchy1"
	HorizontalHierarchy Layout = "urn:microsoft.com/office/officeart/2005/8/layout/hierarchy2"

	// Relationship layouts.
	BasicVenn   Layout = "urn:microsoft.com/office/officeart/2005/8/layout/venn1"
	LinearVenn  Layout = "urn:microsoft.com/office/officeart/2005/8/layout/venn2"
	StackedVenn Layout = "urn:microsoft.com/office/officeart/2005/8/layout/venn3"
	BasicRadial Layout = "urn:microsoft.com/office/officeart/2005/8/layout/radial1"

	// Matrix layouts.
	BasicMatrix  Layout = "urn:microsoft.com/office/officeart/2005/8/layout/matrix1"
	TitledMatrix Layout = "urn:microsoft.com/office/officeart/2005/8/layout/matrix2"

	// Pyramid layouts.
	BasicPyramid    Layout = "urn:microsoft.com/office/officeart/2005/8/layout/pyramid1"
	InvertedPyramid Layout = "urn:microsoft.com/office/officeart/2005/8/layout/pyramid2"

	// Picture layouts.
	PictureStrips Layout = "urn:microsoft.com/office/officeart/2005/8/layout/picture1"
	PictureGrid   Layout = "urn:microsoft.com/office/officeart/2005/8/layout/picture2"
)

// Name returns the human-readable name of the layout.
func (l Layout) Name() string {
	if n, ok := layoutNames[l]; ok {
		return n
	}
	return string(l)
}

//nolint:gochecknoglobals // lookup table
var layoutNames = map[Layout]string{
	BasicBlockList:         "Basic Block List",
	VerticalBlockList:      "Vertical Block List",
	HorizontalBulletLst:    "Horizontal Bullet List",
	SquareAccentList:       "Square Accent List",
	PictureAccentList:      "Picture Accent List",
	BasicProcess:           "Basic Process",
	AccentProcess:          "Accent Process",
	AlternatingFlow:        "Alternating Flow",
	ContinuousBlockProcess: "Continuous Block Process",
	BasicCycle:             "Basic Cycle",
	TextCycle:              "Text Cycle",
	BlockCycle:             "Block Cycle",
	OrgChart:               "Organization Chart",
	Hierarchy:              "Hierarchy",
	HorizontalHierarchy:    "Horizontal Hierarchy",
	BasicVenn:              "Basic Venn",
	LinearVenn:             "Linear Venn",
	StackedVenn:            "Stacked Venn",
	BasicRadial:            "Basic Radial",
	BasicMatrix:            "Basic Matrix",
	TitledMatrix:           "Titled Matrix",
	BasicPyramid:           "Basic Pyramid",
	InvertedPyramid:        "Inverted Pyramid",
	PictureStrips:          "Picture Strips",
	PictureGrid:            "Picture Grid",
}

const (
	defaultX  styling.Length = 914400  // 1 inch
	defaultY  styling.Length = 1828800 // 2 inches
	defaultCX styling.Length = 7315200 // 8 inches
	defaultCY styling.Length = 3657600 // 4 inches
)

// Node represents a single data point in a SmartArt diagram.
// Nodes form a recursive tree for hierarchy layouts.
type Node struct {
	Text     string
	Children []Node
	Color    string // optional RGB hex override, e.g. "FF0000"
}

// NewNode creates a Node with the given text.
func NewNode(text string) Node {
	return Node{Text: text}
}

// WithChild appends a child node (for hierarchy layouts).
func (n Node) WithChild(child Node) Node {
	n.Children = append(n.Children, child)
	return n
}

// WithColor sets an optional RGB hex color on the node.
func (n Node) WithColor(color string) Node {
	n.Color = color
	return n
}

// SmartArt describes a SmartArt diagram to be placed on a slide.
type SmartArt struct {
	Layout     Layout
	Nodes      []Node
	X          styling.Length
	Y          styling.Length
	CX         styling.Length
	CY         styling.Length
	ColorStyle string // optional csTypeId, e.g. "colorful1"
	QuickStyle string // optional qsTypeId
}

// NewSmartArt creates a SmartArt diagram with the given layout and default size.
func NewSmartArt(layout Layout) SmartArt {
	return SmartArt{
		Layout: layout,
		X:      defaultX,
		Y:      defaultY,
		CX:     defaultCX,
		CY:     defaultCY,
	}
}

// AddNode appends a top-level node to the diagram.
func (sa SmartArt) AddNode(node Node) SmartArt {
	sa.Nodes = append(sa.Nodes, node)
	return sa
}

// AddItems appends multiple simple text nodes at once.
func (sa SmartArt) AddItems(items []string) SmartArt {
	for _, item := range items {
		sa.Nodes = append(sa.Nodes, NewNode(item))
	}
	return sa
}

// Position sets the diagram position in EMU.
func (sa SmartArt) Position(x, y styling.Length) SmartArt {
	sa.X = x
	sa.Y = y
	return sa
}

// Size sets the diagram size in EMU.
func (sa SmartArt) Size(cx, cy styling.Length) SmartArt {
	sa.CX = cx
	sa.CY = cy
	return sa
}

// WithColorStyle sets the color style identifier (e.g. "colorful1").
func (sa SmartArt) WithColorStyle(cs string) SmartArt {
	sa.ColorStyle = cs
	return sa
}

// WithQuickStyle sets the quick style identifier.
func (sa SmartArt) WithQuickStyle(qs string) SmartArt {
	sa.QuickStyle = qs
	return sa
}

// ToSpec converts the SmartArt to an internal XML specification.
func (sa SmartArt) ToSpec() pptxxml.SmartArtSpec {
	nodes := make([]pptxxml.SmartArtNodeSpec, len(sa.Nodes))
	for i, n := range sa.Nodes {
		nodes[i] = n.toSpec()
	}
	return pptxxml.SmartArtSpec{
		LayoutURI:    string(sa.Layout),
		ColorStyleID: sa.ColorStyle,
		QuickStyleID: sa.QuickStyle,
		Nodes:        nodes,
		X:            int64(sa.X),
		Y:            int64(sa.Y),
		CX:           int64(sa.CX),
		CY:           int64(sa.CY),
	}
}

func (n Node) toSpec() pptxxml.SmartArtNodeSpec {
	children := make([]pptxxml.SmartArtNodeSpec, len(n.Children))
	for i, c := range n.Children {
		children[i] = c.toSpec()
	}
	return pptxxml.SmartArtNodeSpec{
		Text:     n.Text,
		Children: children,
	}
}
