package export

import (
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

type smartArtBox struct {
	X, Y, W, H float64
	Text       string
	ShapeType  string
	Fill       string
}

type smartArtLink struct {
	StartX, StartY float64
	EndX, EndY     float64
}

func renderPDFSmartArt(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	if renderPDFSmartArtSpecial(pdf, diagram) {
		return
	}
	boxes, links := layoutSmartArt(diagram)
	for _, link := range links {
		renderPDFConnector(pdf, smartArtLinkConnector(link))
	}
	for _, box := range boxes {
		renderPDFShape(pdf, smartArtNodeShape(box))
	}
}

// smartArtLinkConnector is the arrow PowerPoint draws between the nodes of a
// process or hierarchy: a solid accent-coloured line with a head on it, not the
// hairline this used to draw.
func smartArtLinkConnector(link smartArtLink) shapes.Connector {
	connector := shapes.NewElbowConnector(
		styling.Points(link.StartX),
		styling.Points(link.StartY),
		styling.Points(link.EndX),
		styling.Points(link.EndY),
	).WithLine(shapes.NewShapeLine(smartArtLinkColor, styling.Points(smartArtLinkWidthPt)))
	connector.EndArrow = shapes.ArrowTypeTriangle
	return connector
}

// smartArtNodeShape is one node drawn the way PowerPoint's default SmartArt
// style draws it: filled with the theme's first accent colour, with white text
// centred in it, rather than the pastel card with small dark text this used to
// produce.
func smartArtNodeShape(box smartArtBox) shapes.Shape {
	shape := shapes.NewShape(
		box.ShapeType,
		styling.Points(box.X),
		styling.Points(box.Y),
		styling.Points(box.W),
		styling.Points(box.H),
	).WithFill(shapes.NewShapeFill(box.Fill)).
		WithText(box.Text).
		WithVerticalAnchor(shapes.TextAnchorMiddle)
	shape.TextParagraphs = []text.Paragraph{{
		Runs: []text.Run{{
			Text:   box.Text,
			Color:  smartArtNodeTextColor,
			SizePt: smartArtNodeTextSizePt(box),
		}},
		Style: text.ParagraphStyle{Align: elements.TextAlignCenter},
	}}
	return shape
}

// smartArtNodeTextSizePt keeps a node's caption inside its box: PowerPoint
// shrinks SmartArt text to fit rather than letting it overflow.
func smartArtNodeTextSizePt(box smartArtBox) int {
	size := smartArtNodeMaxTextPt
	// Two lines plus the text insets have to fit the box height.
	for size > smartArtNodeMinTextPt && pdfLineHeight(size)*2 > box.H-2*defaultTextInsetTBPt {
		size--
	}
	return size
}

const (
	// smartArtNodeFill is the theme accent PowerPoint fills a default SmartArt
	// node with, and smartArtNodeTextColor the caption colour that goes on it.
	smartArtNodeFill      = "4472C4"
	smartArtNodeTextColor = "FFFFFF"
	smartArtLinkColor     = "4472C4"
	smartArtLinkWidthPt   = 2.0
	smartArtNodeMaxTextPt = 24
	smartArtNodeMinTextPt = 9
)

func layoutSmartArt(diagram smartart.SmartArt) ([]smartArtBox, []smartArtLink) {
	layoutURI := strings.ToLower(diagram.Layout.LayoutURI())
	switch {
	case strings.Contains(layoutURI, "orgchart"), strings.Contains(layoutURI, "hierarchy"):
		return layoutSmartArtHierarchy(diagram)
	case strings.Contains(layoutURI, "matrix"), strings.Contains(layoutURI, "picturegrid"):
		return layoutSmartArtGrid(diagram)
	case strings.Contains(layoutURI, "pyramid"):
		return layoutSmartArtPyramid(diagram)
	case strings.Contains(layoutURI, "cycle"),
		strings.Contains(layoutURI, "venn"),
		strings.Contains(layoutURI, "radial"):
		return layoutSmartArtRadial(diagram)
	default:
		return layoutSmartArtLinear(diagram)
	}
}

// smartArtPalette is the fill of a node the diagram gives no colour of its own.
func smartArtPalette(_ int) string {
	return smartArtNodeFill
}

// smartArtNodeColor is the fill a node is drawn with: its own colour when it has
// one — set on the node, or resolved from the diagram's colour style when the
// deck was read — and the default accent otherwise. Without this every diagram
// exported in the same blue whatever colour style it asked for.
func smartArtNodeColor(node smartart.Node, index int) string {
	if node.Color != "" {
		return node.Color
	}
	return smartArtPalette(index)
}

func flattenSmartArtNodes(nodes []smartart.Node) []smartart.Node {
	out := make([]smartart.Node, 0, len(nodes))
	queue := append([]smartart.Node(nil), nodes...)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		out = append(out, node)
		queue = append(queue, node.Children...)
	}
	return out
}
