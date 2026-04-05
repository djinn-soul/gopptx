//nolint:mnd // SmartArt special renderers intentionally use fixed geometry constants per layout type.
package export

import (
	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

func renderSmartArtBasicMatrix(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	x, y, w, h := smartArtBounds(diagram)
	cx, cy := x+w/2, y+h/2
	drawSmartArtPolygon(
		pdf,
		[]gopdf.Point{{X: cx, Y: y + 12}, {X: x + w - 120, Y: cy}, {X: cx, Y: y + h - 12}, {X: x + 120, Y: cy}},
		smartArtLightFill,
		smartArtLightFill,
		1,
	)
	boxW, boxH, gap := 124.0, 124.0, 10.0
	for i := range min(len(nodes), 4) {
		col, row := i%2, i/2
		left := cx - boxW - gap/2 + float64(col)*(boxW+gap)
		top := cy - boxH - gap/2 + float64(row)*(boxH+gap)
		drawSmartArtRect(pdf, left, top, boxW, boxH, smartArtBlueFill, smartArtWhiteStroke, 20)
		drawSmartArtCenteredText(pdf, nodes[i].Text, left, top, boxW, boxH, smartArtBlueText, 28)
	}
}

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
	centerW, centerH := 188.0, 92.0
	centerX, centerY := x+(w-centerW)/2, y+(h-centerH)/2
	drawSmartArtRect(pdf, centerX, centerY, centerW, centerH, "A6B8D9", smartArtWhiteStroke, 16)
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
		drawSmartArtCenteredText(
			pdf,
			nodes[i].Text,
			x+w/2-bottomHalf+20,
			topY+4,
			bottomHalf*2-40,
			bottomY-topY-8,
			smartArtInkText,
			32,
		)
	}
}

func renderSmartArtPictureGrid(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	x, y, w, _ := smartArtBounds(diagram)
	boxW, boxH, gapX, gapY := 132.0, 132.0, 18.0, 42.0
	left := x + (w-(2*boxW+gapX))/2
	for i := range min(len(nodes), 4) {
		col, row := i%2, i/2
		bx := left + float64(col)*(boxW+gapX)
		by := y + 34 + float64(row)*(boxH+gapY+22)
		drawSmartArtTopText(pdf, nodes[i].Text, bx, by-20, boxW, smartArtInkText, 20)
		drawSmartArtRect(pdf, bx, by, boxW, boxH, smartArtLightFill, smartArtWhiteStroke, 0)
	}
}
