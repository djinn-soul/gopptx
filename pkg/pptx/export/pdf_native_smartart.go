package export

import (
	"math"
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
	// Image is the node's picture, in the picture layouts that carry one. The
	// generic layouts used to drop it, so a picture diagram exported as a row of
	// empty coloured boxes.
	Image []byte
}

type smartArtLink struct {
	StartX, StartY float64
	EndX, EndY     float64
}

func renderPDFSmartArt(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	// The layout PowerPoint itself computed beats anything guessed from the
	// layout URI, so the cache is drawn whenever the deck carries one.
	if renderPDFSmartArtFromCache(pdf, diagram) {
		return
	}
	if renderPDFSmartArtSpecial(pdf, diagram) {
		return
	}
	boxes, links := layoutSmartArt(diagram)
	for _, link := range links {
		renderPDFConnector(pdf, smartArtLinkConnector(link))
	}
	for _, box := range boxes {
		renderPDFSmartArtBox(pdf, box)
	}
}

// renderPDFSmartArtBox draws one node. A node with a picture gets the picture in
// the top of its box and its caption underneath, which is how the picture
// layouts arrange the two; a node without one is drawn as a plain captioned
// shape.
func renderPDFSmartArtBox(pdf *gopdf.GoPdf, box smartArtBox) {
	if len(box.Image) == 0 {
		renderPDFShape(pdf, smartArtNodeShape(box))
		return
	}
	frame := box
	frame.Text = ""
	renderPDFShape(pdf, smartArtNodeShape(frame))

	inset := math.Min(smartArtPictureInsetPt, math.Min(box.W, box.H)*smartArtPictureMaxInsetFraction)
	pictureH := (box.H - 2*inset) * smartArtPictureHeightFraction
	drawn := drawSmartArtImageBytes(
		pdf, box.Image,
		box.X+inset, box.Y+inset,
		box.W-2*inset, pictureH,
	)
	captionY := box.Y + inset
	captionH := box.H - 2*inset
	if drawn {
		captionY += pictureH + inset
		captionH -= pictureH + inset
	}
	if box.Text == "" || captionH <= 0 {
		return
	}
	drawSmartArtCenteredText(
		pdf, box.Text,
		box.X+inset, captionY,
		box.W-2*inset, captionH,
		smartArtNodeTextColor, smartArtNodeMaxTextPt,
	)
}

const (
	// smartArtPictureInsetPt is the margin a node keeps around its picture.
	smartArtPictureInsetPt = 6.0
	// smartArtPictureHeightFraction is how much of a node's inner height the
	// picture takes, leaving the rest for the caption.
	smartArtPictureHeightFraction = 0.6
	// smartArtPictureMaxInsetFraction caps the margin on a small node, so the
	// picture never loses more of the box to its margin than it keeps.
	smartArtPictureMaxInsetFraction = 0.25
)

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

// smartArtNodeColor is the fill a node is drawn with: its own colour when it has
// one — set on the node, or resolved from the diagram's colour style when the
// deck was read — and the default accent otherwise. Without this every diagram
// exported in the same blue whatever colour style it asked for.
//
// The fallback deliberately does not vary by node. PowerPoint's default SmartArt
// colour style is accent1-based and paints every node the same accent; only the
// "colorful" styles spread the palette, and those are resolved per node when the
// deck is read (see resolveSmartArtNodeColors). Cycling here would tint the
// default style's nodes a colour the diagram never asked for.
func smartArtNodeColor(node smartart.Node) string {
	if node.Color != "" {
		return node.Color
	}
	return smartArtNodeFill
}

// smartArtLayoutNodes are the entries a flat layout draws a box for: the
// diagram's own top-level nodes, each carrying its descendants' text as further
// lines of its caption.
//
// This used to breadth-first flatten the tree, which promoted every child to an
// entry of its own: a three-topic diagram with two sub-points each was drawn as
// nine boxes rather than three. PowerPoint puts a node's children inside the
// node, as the second and later lines of its text, and only the hierarchy
// layouts give a child a box of its own.
func smartArtLayoutNodes(nodes []smartart.Node) []smartart.Node {
	out := make([]smartart.Node, 0, len(nodes))
	for _, node := range nodes {
		entry := node
		entry.Text = strings.Join(smartArtNodeTextLines(node), "\n")
		entry.Children = nil
		out = append(out, entry)
	}
	return out
}

// smartArtNodeTextLines is a node's own text followed by its descendants', in
// the order the diagram lists them. Empty captions are dropped so a node with no
// text of its own does not open its box with a blank line.
func smartArtNodeTextLines(node smartart.Node) []string {
	lines := make([]string, 0, 1+len(node.Children))
	if text := strings.TrimSpace(node.Text); text != "" {
		lines = append(lines, text)
	}
	for _, child := range node.Children {
		lines = append(lines, smartArtNodeTextLines(child)...)
	}
	return lines
}
