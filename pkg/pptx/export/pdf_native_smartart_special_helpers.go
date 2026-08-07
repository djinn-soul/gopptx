//nolint:mnd // SmartArt special helper geometry and text paddings are template-calibrated constants.
package export

import (
	"math"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

// The per-layout renderers draw from the same accent as the generic layouts.
// They used to carry their own 4F81BD — the Office 2007 accent1 — so one deck
// exported two different blues depending on which layout a diagram used.
const (
	smartArtBlueFill    = smartArtNodeFill
	smartArtBlueText    = smartArtNodeTextColor
	smartArtInkText     = "000000"
	smartArtLightFill   = "C2CDE1"
	smartArtPanelFill   = "D5DCEA"
	smartArtWhiteStroke = "FFFFFF"
	smartArtLineStroke  = smartArtNodeFill
)

// The per-layout renderers below were calibrated against PowerPoint's render of
// each layout in a large frame, and state their geometry in absolute points. A
// diagram placed in a smaller frame therefore drew straight past its edges — a
// stacked Venn 320pt across still drew 320pt across in a 2in frame.
//
// smartArtFitScale is the factor that brings a renderer's calibrated content
// back inside the frame it was given. It never scales up: at or above the size
// the layout was calibrated for, the geometry is used exactly as measured.
func smartArtFitScale(frameW, frameH, contentW, contentH float64) float64 {
	scale := 1.0
	if contentW > 0 && frameW > 0 {
		scale = math.Min(scale, frameW/contentW)
	}
	if contentH > 0 && frameH > 0 {
		scale = math.Min(scale, frameH/contentH)
	}
	if scale <= 0 {
		return 1
	}
	return scale
}

func smartArtBounds(diagram smartart.SmartArt) (float64, float64, float64, float64) {
	return emuToPt(int64(diagram.X)), emuToPt(int64(diagram.Y)), emuToPt(int64(diagram.CX)), emuToPt(int64(diagram.CY))
}

func smartArtNodes(diagram smartart.SmartArt) []smartart.Node {
	return smartArtLayoutNodes(diagram.Nodes)
}

// smartArtEntries are the diagram's own entries, with their children left in
// place. Layouts that draw a body under each entry need the tree: flattening it
// turned every child into another entry, so a three-topic diagram was drawn with
// six boxes running off the edge of its frame.
func smartArtEntries(diagram smartart.SmartArt) []smartart.Node {
	if len(diagram.Nodes) == 0 {
		return nil
	}
	return diagram.Nodes
}

// drawSmartArtNodeImage paints a node's picture into the given box, and reports
// whether there was one to paint.
func drawSmartArtNodeImage(pdf *gopdf.GoPdf, node smartart.Node, x, y, w, h float64) bool {
	return drawSmartArtImageBytes(pdf, node.ImageData, x, y, w, h)
}

// drawSmartArtImageBytes paints picture bytes into the given box, and reports
// whether there was a picture to paint.
func drawSmartArtImageBytes(pdf *gopdf.GoPdf, data []byte, x, y, w, h float64) bool {
	if len(data) == 0 || w <= 0 || h <= 0 {
		return false
	}
	holder, err := gopdf.ImageHolderByBytes(data)
	if err != nil {
		return false
	}
	if err := pdf.ImageByHolder(holder, x, y, &gopdf.Rect{W: w, H: h}); err != nil {
		return false
	}
	return true
}

// drawSmartArtChildLines writes an entry's children as bullet lines inside its
// body, which is where the layout puts them.
func drawSmartArtChildLines(
	pdf *gopdf.GoPdf,
	children []smartart.Node,
	x, y, w float64,
	color string,
	sizePt int,
) {
	lineHeight := pdfLineHeight(sizePt)
	for i, child := range children {
		drawSmartArtTopText(pdf, "• "+child.Text, x, y+float64(i)*lineHeight, w, color, sizePt)
	}
}

func smartArtLayoutURI(diagram smartart.SmartArt) string {
	return strings.ToLower(diagram.Layout.LayoutURI())
}

func drawSmartArtRect(pdf *gopdf.GoPdf, x, y, w, h float64, fill, stroke string, radius float64) {
	pdf.SetFillColor(hexToRGB(fill))
	pdf.SetStrokeColor(hexToRGB(stroke))
	pdf.SetLineWidth(1)
	if radius > 0 {
		_ = pdf.Rectangle(x, y, x+w, y+h, "DF", radius, 0)
		return
	}
	pdf.RectFromUpperLeftWithStyle(x, y, w, h, "DF")
}

func drawSmartArtEllipse(pdf *gopdf.GoPdf, x, y, w, h float64, fill, stroke string, alpha float64) {
	cx := x + w/2
	cy := y + h/2
	rx := w / 2
	ry := h / 2
	points := make([]gopdf.Point, 0, 40)
	for i := range 40 {
		angle := (2 * math.Pi * float64(i)) / 40
		points = append(points, gopdf.Point{
			X: cx + math.Cos(angle)*rx,
			Y: cy + math.Sin(angle)*ry,
		})
	}
	drawSmartArtPolygon(pdf, points, fill, stroke, alpha)
}

func drawSmartArtPolygon(pdf *gopdf.GoPdf, points []gopdf.Point, fill, stroke string, alpha float64) {
	pdf.SetFillColor(hexToRGB(fill))
	pdf.SetStrokeColor(hexToRGB(stroke))
	pdf.SetLineWidth(1)
	if alpha > 0 && alpha < 1 {
		transparency, err := gopdf.NewTransparency(alpha, shapeEffectsBlendMode)
		if err == nil {
			_ = pdf.SetTransparency(transparency)
		}
	}
	pdf.Polygon(points, "DF")
	if alpha > 0 && alpha < 1 {
		pdf.ClearTransparency()
	}
}

func drawSmartArtCenteredText(pdf *gopdf.GoPdf, text string, x, y, w, h float64, color string, maxSize int) {
	fontSize := fitPDFTextToBoxWithMetrics(pdf, text, maxSize, minTextAutoFitSize, false, false, w-8, h-8, "")
	setPDFTextFontWithHint(pdf, fontSize, false, false, "")
	lines := wrapPDFTextWithMetrics(pdf, text, w-8)
	lineH := pdfLineHeight(fontSize)
	totalH := lineH * float64(len(lines))
	startY := y + max((h-totalH)/2, 0)
	pdf.SetTextColor(hexToRGB(color))
	for i, line := range lines {
		lineW := measuredWidth(pdf, line)
		pdf.SetX(x + max((w-lineW)/2, 0))
		pdf.SetY(startY + float64(i)*lineH + fontBaselineShift(pdf, "", fontSize))
		_ = pdf.Cell(nil, line)
	}
	setPDFTextFontWithHint(pdf, defaultFontSize, false, false, "")
}

func drawSmartArtTopText(pdf *gopdf.GoPdf, text string, x, y, w float64, color string, fontSize int) {
	setPDFTextFontWithHint(pdf, fontSize, false, false, "")
	pdf.SetTextColor(hexToRGB(color))
	lineW := measuredWidth(pdf, text)
	pdf.SetX(x + max((w-lineW)/2, 0))
	pdf.SetY(y + fontBaselineShift(pdf, "", fontSize))
	_ = pdf.Cell(nil, text)
	setPDFTextFontWithHint(pdf, defaultFontSize, false, false, "")
}

func drawSmartArtVerticalText(pdf *gopdf.GoPdf, text string, cx, cy float64, color string, fontSize int) {
	setPDFTextFontWithHint(pdf, fontSize, false, false, "")
	pdf.SetTextColor(hexToRGB(color))
	lineW := measuredWidth(pdf, text)
	pdf.Rotate(90, cx, cy)
	pdf.SetX(cx - lineW/2)
	pdf.SetY(cy + fontBaselineShift(pdf, "", fontSize))
	_ = pdf.Cell(nil, text)
	pdf.RotateReset()
	setPDFTextFontWithHint(pdf, defaultFontSize, false, false, "")
}

func drawSmartArtLine(pdf *gopdf.GoPdf, x1, y1, x2, y2 float64) {
	pdf.SetStrokeColor(hexToRGB(smartArtLineStroke))
	pdf.SetLineWidth(1.4)
	pdf.Line(x1, y1, x2, y2)
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(1)
}
