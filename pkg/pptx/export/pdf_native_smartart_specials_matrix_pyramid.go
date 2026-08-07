//nolint:mnd // SmartArt special renderers intentionally use fixed geometry constants per layout type.
package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

func renderSmartArtBasicMatrix(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	x, y, w, h := smartArtBounds(diagram)
	cx, cy := x+w/2, y+h/2
	// The backing diamond is inset from the frame's sides by a fraction of its
	// width rather than a flat 120pt, which used to invert the diamond in any
	// frame under 240pt wide.
	inset := w * matrixDiamondInsetFraction
	drawSmartArtPolygon(
		pdf,
		[]gopdf.Point{
			{X: cx, Y: y + 12}, {X: x + w - inset, Y: cy},
			{X: cx, Y: y + h - 12}, {X: x + inset, Y: cy},
		},
		smartArtLightFill,
		smartArtLightFill,
		1,
	)
	boxW, boxH, gap := 124.0, 124.0, 10.0
	// The four quadrants have to fit inside the diamond, not just the frame.
	scale := smartArtFitScale(w-2*inset, h, 2*boxW+gap, 2*boxH+gap)
	boxW, boxH, gap = boxW*scale, boxH*scale, gap*scale
	for i := range min(len(nodes), 4) {
		col, row := i%2, i/2
		left := cx - boxW - gap/2 + float64(col)*(boxW+gap)
		top := cy - boxH - gap/2 + float64(row)*(boxH+gap)
		drawSmartArtRect(pdf, left, top, boxW, boxH, smartArtBlueFill, smartArtWhiteStroke, 20*scale)
		drawSmartArtCenteredText(pdf, nodes[i].Text, left, top, boxW, boxH, smartArtBlueText, 28)
	}
}

// matrixDiamondInsetFraction is how far in from each side the backing diamond's
// points sit, as a fraction of the frame width. It is the calibrated 120pt over
// the ~9in frame the layout was measured in.
const matrixDiamondInsetFraction = 120.0 / 648.0

// matrixCentreMaxFraction is the most of the frame the centre label may take, so
// it stays clear of the four quadrant captions around it.
const matrixCentreMaxFraction = 0.45

// pyramidTextInsetPt is the margin a pyramid tier's caption keeps from the
// sloped edges either side of it.
const pyramidTextInsetPt = 20.0

// pictureGridCaptionPt is the height reserved above each tile for its caption.
const pictureGridCaptionPt = 22.0

func renderSmartArtTitledMatrix(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	x, y, w, h := smartArtBounds(diagram)
	drawSmartArtRect(pdf, x, y, w, h, smartArtBlueFill, smartArtWhiteStroke, 28)
	drawSmartArtLine(pdf, x+w/2, y, x+w/2, y+h)
	drawSmartArtLine(pdf, x, y+h/2, x+w, y+h/2)
	if len(nodes) > 1 {
		drawSmartArtCenteredText(pdf, nodes[1].Text, x+14, y+16, w/2-28, h/2-32, smartArtBlueText, 28)
	}
	if len(nodes) > 2 {
		drawSmartArtCenteredText(pdf, nodes[2].Text, x+w/2+14, y+16, w/2-28, h/2-32, smartArtBlueText, 28)
	}
	if len(nodes) > 3 {
		drawSmartArtCenteredText(pdf, nodes[3].Text, x+14, y+h/2+16, w/2-28, h/2-32, smartArtBlueText, 28)
	}
	// The centre label is fixed at 188x92pt, which covered a small frame edge to
	// edge; it is scaled to stay a label rather than a lid.
	centerW, centerH := 188.0, 92.0
	centerScale := smartArtFitScale(w*matrixCentreMaxFraction, h*matrixCentreMaxFraction, centerW, centerH)
	centerW, centerH = centerW*centerScale, centerH*centerScale
	centerX, centerY := x+(w-centerW)/2, y+(h-centerH)/2
	drawSmartArtRect(pdf, centerX, centerY, centerW, centerH, "A6B8D9", smartArtWhiteStroke, 16*centerScale)
	if len(nodes) > 0 {
		drawSmartArtCenteredText(pdf, nodes[0].Text, centerX, centerY, centerW, centerH, smartArtInkText, 30)
	}
}

func renderSmartArtPyramid(pdf *gopdf.GoPdf, diagram smartart.SmartArt, inverted bool) {
	nodes := smartArtNodes(diagram)
	x, y, w, h := smartArtBounds(diagram)
	for i := range nodes {
		topRatio := float64(i) / float64(len(nodes))
		bottomRatio := float64(i+1) / float64(len(nodes))
		if inverted {
			topRatio, bottomRatio = 1-bottomRatio, 1-topRatio
		}
		topY := y + float64(i)*(h/float64(len(nodes)))
		bottomY := y + float64(i+1)*(h/float64(len(nodes)))
		topHalf := (w / 2) * topRatio
		bottomHalf := (w / 2) * bottomRatio
		poly := []gopdf.Point{
			{X: x + w/2 - topHalf, Y: topY},
			{X: x + w/2 + topHalf, Y: topY},
			{X: x + w/2 + bottomHalf, Y: bottomY},
			{X: x + w/2 - bottomHalf, Y: bottomY},
		}
		drawSmartArtPolygon(pdf, poly, smartArtBlueFill, smartArtWhiteStroke, 1)
		// The caption is inset from the tier's own edges. A narrow top tier is
		// narrower than the flat 20pt inset used to be, which gave the text box a
		// negative width.
		textInset := math.Min(pyramidTextInsetPt, bottomHalf/2)
		drawSmartArtCenteredText(
			pdf,
			nodes[i].Text,
			x+w/2-bottomHalf+textInset,
			topY+4,
			math.Max(1, bottomHalf*2-2*textInset),
			math.Max(1, bottomY-topY-8),
			smartArtInkText,
			32,
		)
	}
}

func renderSmartArtPictureGrid(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	x, y, w, h := smartArtBounds(diagram)
	boxW, boxH, gapX, gapY, topPad := 132.0, 132.0, 18.0, 42.0, 34.0
	// Two rows of captioned tiles: the grid is scaled to the frame rather than
	// drawn at a fixed 132pt per tile, which ran the second row off the bottom.
	scale := smartArtFitScale(w, h, 2*boxW+gapX, topPad+2*boxH+gapY+pictureGridCaptionPt)
	boxW, boxH, gapX, gapY, topPad = boxW*scale, boxH*scale, gapX*scale, gapY*scale, topPad*scale
	captionH := pictureGridCaptionPt * scale

	left := x + (w-(2*boxW+gapX))/2
	for i := range min(len(nodes), 4) {
		col, row := i%2, i/2
		bx := left + float64(col)*(boxW+gapX)
		by := y + topPad + float64(row)*(boxH+gapY+captionH)
		drawSmartArtTopText(pdf, nodes[i].Text, bx, by-captionH, boxW, smartArtInkText, 20)
		if !drawSmartArtNodeImage(pdf, nodes[i], bx, by, boxW, boxH) {
			drawSmartArtRect(pdf, bx, by, boxW, boxH, smartArtLightFill, smartArtWhiteStroke, 0)
		}
	}
}
