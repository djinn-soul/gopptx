//nolint:mnd // SmartArt seed geometry uses fixed template-calibrated dimensions.
package export

import (
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
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

func renderNativePDFSlideSmartArt(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	for _, diagram := range slide.SmartArtDiagrams {
		renderPDFSmartArt(pdf, diagram)
	}
}

func renderPDFSmartArt(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	if renderPDFSmartArtSpecial(pdf, diagram) {
		return
	}
	boxes, links := layoutSmartArt(diagram)
	for _, link := range links {
		renderPDFConnector(pdf, shapes.NewElbowConnector(
			styling.Points(link.StartX),
			styling.Points(link.StartY),
			styling.Points(link.EndX),
			styling.Points(link.EndY),
		).WithLine(shapes.NewShapeLine("7A869A", styling.Points(1.1))))
	}
	for _, box := range boxes {
		renderPDFShape(pdf, shapes.NewShape(
			box.ShapeType,
			styling.Points(box.X),
			styling.Points(box.Y),
			styling.Points(box.W),
			styling.Points(box.H),
		).WithFill(shapes.NewShapeFill(box.Fill)).
			WithLine(shapes.NewShapeLine("5B6578", styling.Points(0.9))).
			WithText(box.Text).
			WithTextFrame(shapes.NewTextFrame().WithRotation(0)))
	}
}

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

func smartArtPalette(depth int) string {
	colors := []string{"DCEBFA", "E7F4E8", "FDE8D7", "F3E6FA", "FFF3CD", "E6F5F7"}
	return colors[depth%len(colors)]
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
