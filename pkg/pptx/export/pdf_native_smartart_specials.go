//nolint:mnd // SmartArt special renderers intentionally use fixed geometry constants per layout type.
package export

import (
	"math"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

func renderPDFSmartArtSpecial(pdf *gopdf.GoPdf, diagram smartart.SmartArt) bool {
	switch uri := strings.ToLower(smartArtLayoutURI(diagram)); {
	case strings.Contains(uri, "layout/default"):
		renderSmartArtBasicBlockList(pdf, diagram)
	case strings.Contains(uri, "vlist5"):
		renderSmartArtVerticalList(pdf, diagram)
	case strings.Contains(uri, "hlist1"):
		renderSmartArtHorizontalBulletList(pdf, diagram)
	case strings.Contains(uri, "hlist2"):
		renderSmartArtPictureAccentList(pdf, diagram)
	case strings.Contains(uri, "hprocess9"):
		renderSmartArtContinuousBlock(pdf, diagram)
	case strings.Contains(uri, "hierarchy2"):
		renderSmartArtHorizontalHierarchy(pdf, diagram)
	case strings.Contains(uri, "venn3"):
		renderSmartArtLinearVenn(pdf, diagram)
	case strings.Contains(uri, "venn2"):
		renderSmartArtStackedVenn(pdf, diagram)
	case strings.Contains(uri, "radial1"):
		renderSmartArtBasicRadial(pdf, diagram)
	case strings.Contains(uri, "matrix3"):
		renderSmartArtBasicMatrix(pdf, diagram)
	case strings.Contains(uri, "matrix1"):
		renderSmartArtTitledMatrix(pdf, diagram)
	case strings.Contains(uri, "pyramid1"):
		renderSmartArtPyramid(pdf, diagram, false)
	case strings.Contains(uri, "pyramid3"):
		renderSmartArtPyramid(pdf, diagram, true)
	case strings.Contains(uri, "picturegrid"):
		renderSmartArtPictureGrid(pdf, diagram)
	default:
		return false
	}
	return true
}

func renderSmartArtBasicBlockList(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	x, y, w, _ := smartArtBounds(diagram)
	boxW, boxH, gap := 158.0, 106.0, 22.0
	topY, bottomY := y+20, y+20+boxH+gap
	topLeft := x + (w-(3*boxW+2*gap))/2 - 10
	for i, node := range nodes {
		row, col := 0, i
		if i >= 3 {
			row, col = 1, i-3
		}
		left := topLeft + float64(col)*(boxW+gap)
		if row == 1 {
			left = x + (w-(2*boxW+gap))/2 + float64(col)*(boxW+gap)
		}
		top := topY
		if row == 1 {
			top = bottomY
		}
		drawSmartArtRect(pdf, left, top, boxW, boxH, smartArtBlueFill, smartArtBlueFill, 0)
		drawSmartArtCenteredText(pdf, node.Text, left+4, top, boxW-8, boxH, smartArtBlueText, 30)
	}
}

func renderSmartArtVerticalList(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	x, y, w, _ := smartArtBounds(diagram)
	boxW, boxH, gap := 220.0, 92.0, 12.0
	left := x + (w-boxW)/2
	for i, node := range nodes {
		top := y + 20 + float64(i)*(boxH+gap)
		drawSmartArtRect(pdf, left, top, boxW, boxH, smartArtBlueFill, smartArtBlueFill, 20)
		drawSmartArtCenteredText(pdf, node.Text, left, top, boxW, boxH, smartArtBlueText, 30)
	}
}

func renderSmartArtHorizontalBulletList(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	x, y, w, _ := smartArtBounds(diagram)
	boxW, boxH, gap := 184.0, 204.0, 26.0
	headerH := 76.0
	leftStart := x + (w-(3*boxW+2*gap))/2
	for i, node := range nodes {
		left := leftStart + float64(i)*(boxW+gap)
		drawSmartArtRect(pdf, left, y+34, boxW, boxH, smartArtPanelFill, smartArtPanelFill, 0)
		drawSmartArtRect(pdf, left, y+34, boxW, headerH, smartArtBlueFill, smartArtBlueFill, 0)
		drawSmartArtCenteredText(pdf, node.Text, left+10, y+44, boxW-20, headerH-10, smartArtBlueText, 26)
	}
}

func renderSmartArtPictureAccentList(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	x, y, w, _ := smartArtBounds(diagram)
	boxW, boxH, gap := 154.0, 248.0, 30.0
	leftStart := x + (w-(3*boxW+2*gap))/2 - 14
	for i, node := range nodes {
		left := leftStart + float64(i)*(boxW+gap)
		drawSmartArtRect(pdf, left, y+50, boxW, boxH, smartArtBlueFill, smartArtBlueFill, 0)
		drawSmartArtRect(pdf, left-34, y+18, 56, 56, smartArtLightFill, smartArtWhiteStroke, 0)
		drawSmartArtVerticalText(pdf, node.Text, left-14, y+96, smartArtInkText, 20)
	}
}

func renderSmartArtContinuousBlock(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	x, y, w, h := smartArtBounds(diagram)
	arrow := []gopdf.Point{
		{X: x + 46, Y: y + 78},
		{X: x + w - 160, Y: y + 78},
		{X: x + w - 160, Y: y + 28},
		{X: x + w - 38, Y: y + h/2},
		{X: x + w - 160, Y: y + h - 28},
		{X: x + w - 160, Y: y + h - 78},
		{X: x + 46, Y: y + h - 78},
	}
	drawSmartArtPolygon(pdf, arrow, smartArtLightFill, smartArtLightFill, 1)
	boxW, boxH, gap := 192.0, 124.0, 26.0
	for i, node := range nodes {
		left := x + float64(i)*(boxW+gap)
		top := y + (h-boxH)/2
		drawSmartArtRect(pdf, left, top, boxW, boxH, smartArtBlueFill, smartArtBlueFill, 18)
		drawSmartArtCenteredText(pdf, node.Text, left, top, boxW, boxH, smartArtBlueText, 28)
	}
}

func renderSmartArtHorizontalHierarchy(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes, levels := buildSmartArtHierarchy(diagram.Nodes)
	if len(nodes) == 0 || len(levels) == 0 {
		return
	}
	x, y, w, h := smartArtBounds(diagram)
	colGap, rowGap := 45.0, 16.0
	cols := len(levels)
	boxW := min(164.0, (w-colGap*float64(max(cols-1, 0)))/float64(cols))
	levelHeights := make([]float64, len(levels))
	for depth, level := range levels {
		levelHeights[depth] = float64(len(level))*80 + float64(max(len(level)-1, 0))*rowGap
	}
	centers := make([][2]float64, len(nodes))
	yOffset := -8.0
	for depth, level := range levels {
		left := x + float64(depth)*(boxW+colGap)
		totalH := levelHeights[depth]
		top := y + (h-totalH)/2 + yOffset
		for i, nodeIndex := range level {
			by := top + float64(i)*(80+rowGap)
			drawSmartArtRect(pdf, left, by, boxW, 80, smartArtBlueFill, smartArtBlueFill, 12)
			drawSmartArtCenteredText(pdf, nodes[nodeIndex].Node.Text, left, by, boxW, 80, smartArtBlueText, 28)
			centers[nodeIndex] = [2]float64{left + boxW/2, by + 40}
			if nodes[nodeIndex].Parent >= 0 {
				parent := centers[nodes[nodeIndex].Parent]
				drawSmartArtLine(pdf, parent[0]+boxW/2-6, parent[1], left, by+40)
			}
		}
	}
}

func renderSmartArtLinearVenn(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	x, y, w, _ := smartArtBounds(diagram)
	diameter, overlap := 188.0, 40.0
	totalW := diameter*float64(len(nodes)) - overlap*float64(max(len(nodes)-1, 0))
	left := x + (w-totalW)/2
	for i, node := range nodes {
		cx := left + float64(i)*(diameter-overlap)
		drawSmartArtEllipse(pdf, cx, y+34, diameter, diameter, "9FB7D8", smartArtWhiteStroke, 0.75)
		drawSmartArtCenteredText(pdf, node.Text, cx+22, y+84, diameter-44, 92, smartArtInkText, 24)
	}
}

func renderSmartArtStackedVenn(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	x, y, w, _ := smartArtBounds(diagram)
	cx := x + w/2
	baseY := y + 324
	for i, node := range nodes {
		diameter := 320.0 - float64(i)*68
		drawSmartArtEllipse(
			pdf,
			cx-diameter/2,
			baseY-diameter,
			diameter,
			diameter,
			smartArtBlueFill,
			smartArtWhiteStroke,
			0.7,
		)
		drawSmartArtCenteredText(
			pdf,
			node.Text,
			cx-diameter/4,
			baseY-diameter+10,
			diameter/2,
			diameter-20,
			smartArtBlueText,
			22,
		)
	}
}

func renderSmartArtBasicRadial(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	if len(nodes) == 0 {
		return
	}
	x, y, w, h := smartArtBounds(diagram)
	size := math.Min(w, h) * 0.72
	left := x + (w-size)/2
	top := y + (h-size)/2
	drawSmartArtEllipse(pdf, left, top, size, size, smartArtBlueFill, smartArtBlueFill, 1)
	drawSmartArtCenteredText(pdf, nodes[0].Text, left+42, top+80, size-84, size-160, smartArtBlueText, 40)
}

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
