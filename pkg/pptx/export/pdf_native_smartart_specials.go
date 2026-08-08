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
	if len(nodes) == 0 {
		return
	}
	x, y, w, h := smartArtBounds(diagram)
	topCols := min(len(nodes), 3)
	bottomCols := len(nodes) - topCols
	rows := 1
	if bottomCols > 0 {
		rows = 2
	}
	usableW := math.Max(1, w-20)
	usableH := math.Max(1, h-20)
	gapX := math.Max(8, math.Min(22, usableW*0.05))
	gapY := math.Max(8, math.Min(22, usableH*0.10))
	boxW := math.Min(158, (usableW-gapX*float64(max(topCols-1, 0)))/float64(topCols))
	if bottomCols > 0 {
		bottomBoxW := (usableW - gapX*float64(max(bottomCols-1, 0))) / float64(bottomCols)
		boxW = math.Min(boxW, bottomBoxW)
	}
	boxW = math.Max(1, boxW)
	boxH := math.Max(1, math.Min(106, (usableH-gapY*float64(rows-1))/float64(rows)))
	totalH := boxH * float64(rows)
	if rows > 1 {
		totalH += gapY
	}
	topY := y + (h-totalH)/2
	bottomY := topY + boxH + gapY
	topRowW := float64(topCols)*boxW + gapX*float64(max(topCols-1, 0))
	topLeft := x + (w-topRowW)/2
	bottomRowW := float64(bottomCols)*boxW + gapX*float64(max(bottomCols-1, 0))
	bottomLeft := x + (w-bottomRowW)/2
	for i, node := range nodes {
		row, col := 0, i
		if i >= 3 {
			row, col = 1, i-3
		}
		left := topLeft + float64(col)*(boxW+gapX)
		if row == 1 {
			left = bottomLeft + float64(col)*(boxW+gapX)
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
	if len(nodes) == 0 {
		return
	}
	x, y, w, h := smartArtBounds(diagram)
	boxW, boxH, gap, topPad := 220.0, 92.0, 12.0, 20.0
	// The stack is as tall as its entries make it, so a frame shorter than that
	// — or narrower than one box — scales the whole list down to fit.
	stackH := topPad + boxH*float64(len(nodes)) + gap*float64(len(nodes)-1)
	scale := smartArtFitScale(w, h, boxW, stackH)
	boxW, boxH, gap, topPad = boxW*scale, boxH*scale, gap*scale, topPad*scale

	left := x + (w-boxW)/2
	for i, node := range nodes {
		top := y + topPad + float64(i)*(boxH+gap)
		drawSmartArtRect(pdf, left, top, boxW, boxH, smartArtBlueFill, smartArtBlueFill, 20*scale)
		drawSmartArtCenteredText(pdf, node.Text, left, top, boxW, boxH, smartArtBlueText, 30)
	}
}

func renderSmartArtHorizontalBulletList(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtEntries(diagram)
	if len(nodes) == 0 {
		return
	}
	x, y, w, h := smartArtBounds(diagram)
	boxH, gap, headerH, topPad := 204.0, 26.0, 76.0, 34.0
	// The columns already share the frame's width; their height was fixed, so a
	// short frame had the panels hanging out of the bottom of it.
	scale := smartArtFitScale(w, h, w, topPad+boxH)
	boxH, gap, headerH, topPad = boxH*scale, gap*scale, headerH*scale, topPad*scale

	boxW := math.Max(1, (w-gap*float64(len(nodes)-1))/float64(len(nodes)))
	for i, node := range nodes {
		left := x + float64(i)*(boxW+gap)
		top := y + topPad
		drawSmartArtRect(pdf, left, top, boxW, boxH, smartArtPanelFill, smartArtPanelFill, 0)
		drawSmartArtRect(pdf, left, top, boxW, headerH, smartArtNodeColor(node), smartArtWhiteStroke, 0)
		drawSmartArtCenteredText(
			pdf, node.Text,
			left+10, top+10*scale,
			math.Max(1, boxW-20), math.Max(1, headerH-10*scale),
			smartArtBlueText, 26,
		)
		drawSmartArtChildLines(
			pdf, node.Children,
			left+12, top+headerH+10*scale,
			math.Max(1, boxW-24), smartArtInkText, 16,
		)
	}
}

func renderSmartArtPictureAccentList(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtEntries(diagram)
	if len(nodes) == 0 {
		return
	}
	x, y, w, h := smartArtBounds(diagram)
	boxH, gap, topPad, badge := 248.0, 30.0, 50.0, 56.0
	// The entries share the frame's width rather than assuming three of them,
	// which used to run a fourth entry off the edge of the slide. The height and
	// the picture badge were still fixed, so they are scaled to the frame too.
	scale := smartArtFitScale(w, h, w, topPad+boxH)
	boxH, gap, topPad, badge = boxH*scale, gap*scale, topPad*scale, badge*scale

	boxW := math.Max(1, (w-gap*float64(len(nodes)-1))/float64(len(nodes)))
	panelInset := math.Min(22*scale, boxW/2)
	for i, node := range nodes {
		left := x + float64(i)*(boxW+gap)
		top := y + topPad
		drawSmartArtRect(
			pdf, left+panelInset, top,
			math.Max(1, boxW-panelInset), boxH,
			smartArtNodeColor(node), smartArtWhiteStroke, 0,
		)
		badgeX, badgeY := left-12*scale, top-32*scale
		if !drawSmartArtNodeImage(pdf, node, badgeX, badgeY, badge, badge) {
			drawSmartArtRect(pdf, badgeX, badgeY, badge, badge, smartArtLightFill, smartArtWhiteStroke, 0)
		}
		drawSmartArtVerticalText(pdf, node.Text, left+8*scale, top+46*scale, smartArtInkText, 20)
		drawSmartArtChildLines(
			pdf, node.Children,
			left+36*scale, top+20*scale,
			math.Max(1, boxW-46*scale), smartArtBlueText, 14,
		)
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
	if len(nodes) == 0 {
		return
	}
	boxW, boxH, gap := 192.0, 124.0, 26.0
	// The row is as wide as its blocks make it. Four blocks at the calibrated
	// size already run past a 10in frame, so the row is scaled to the frame.
	rowW := boxW*float64(len(nodes)) + gap*float64(len(nodes)-1)
	scale := smartArtFitScale(w, h, rowW, boxH)
	boxW, boxH, gap = boxW*scale, boxH*scale, gap*scale
	rowW = boxW*float64(len(nodes)) + gap*float64(len(nodes)-1)

	for i, node := range nodes {
		left := x + (w-rowW)/2 + float64(i)*(boxW+gap)
		top := y + (h-boxH)/2
		drawSmartArtRect(pdf, left, top, boxW, boxH, smartArtBlueFill, smartArtBlueFill, 18*scale)
		drawSmartArtCenteredText(pdf, node.Text, left, top, boxW, boxH, smartArtBlueText, 28)
	}
}

func renderSmartArtHorizontalHierarchy(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes, levels := buildSmartArtHierarchy(diagram.Nodes)
	if len(nodes) == 0 || len(levels) == 0 {
		return
	}
	x, y, w, h := smartArtBounds(diagram)
	colGap, rowGap, boxH := 45.0, 16.0, 80.0
	cols := len(levels)
	// The widest level decides how tall the column of boxes has to be. Boxes were
	// a fixed 80pt, so a level with several siblings ran off the frame.
	widest := 1
	for _, level := range levels {
		widest = max(widest, len(level))
	}
	scale := smartArtFitScale(w, h, w, boxH*float64(widest)+rowGap*float64(widest-1))
	colGap, rowGap, boxH = colGap*scale, rowGap*scale, boxH*scale

	boxW := math.Max(1, min(164.0, (w-colGap*float64(max(cols-1, 0)))/float64(cols)))
	levelHeights := make([]float64, len(levels))
	for depth, level := range levels {
		levelHeights[depth] = float64(len(level))*boxH + float64(max(len(level)-1, 0))*rowGap
	}
	centers := make([][2]float64, len(nodes))
	yOffset := -8.0 * scale
	for depth, level := range levels {
		left := x + float64(depth)*(boxW+colGap)
		totalH := levelHeights[depth]
		top := y + (h-totalH)/2 + yOffset
		for i, nodeIndex := range level {
			by := top + float64(i)*(boxH+rowGap)
			drawSmartArtRect(pdf, left, by, boxW, boxH, smartArtBlueFill, smartArtBlueFill, 12*scale)
			drawSmartArtCenteredText(pdf, nodes[nodeIndex].Node.Text, left, by, boxW, boxH, smartArtBlueText, 28)
			centers[nodeIndex] = [2]float64{left + boxW/2, by + boxH/2}
			if nodes[nodeIndex].Parent >= 0 {
				parent := centers[nodes[nodeIndex].Parent]
				drawSmartArtLine(pdf, parent[0]+boxW/2-6*scale, parent[1], left, by+boxH/2)
			}
		}
	}
}

func renderSmartArtLinearVenn(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	if len(nodes) == 0 {
		return
	}
	x, y, w, h := smartArtBounds(diagram)
	diameter, overlap, topPad := 188.0, 40.0, 34.0
	totalW := diameter*float64(len(nodes)) - overlap*float64(max(len(nodes)-1, 0))
	// Circles of a fixed diameter run off both the sides and the bottom of a
	// small frame, so the whole row is scaled to fit it.
	scale := smartArtFitScale(w, h, totalW, topPad+diameter)
	diameter, overlap, topPad = diameter*scale, overlap*scale, topPad*scale
	totalW = diameter*float64(len(nodes)) - overlap*float64(max(len(nodes)-1, 0))

	left := x + (w-totalW)/2
	for i, node := range nodes {
		cx := left + float64(i)*(diameter-overlap)
		drawSmartArtEllipse(pdf, cx, y+topPad, diameter, diameter, "9FB7D8", smartArtWhiteStroke, 0.75)
		drawSmartArtCenteredText(
			pdf, node.Text,
			cx+diameter*vennTextInsetFraction, y+topPad+diameter*vennTextTopFraction,
			diameter*(1-2*vennTextInsetFraction), diameter*vennTextHeightFraction,
			smartArtInkText, 24,
		)
	}
}

// The Venn caption sits in the middle band of its circle, clear of the overlap
// on either side. The fractions are the calibrated 22/50/92 points over the
// calibrated 188pt diameter.
const (
	vennTextInsetFraction  = 22.0 / 188.0
	vennTextTopFraction    = 50.0 / 188.0
	vennTextHeightFraction = 92.0 / 188.0
)

func renderSmartArtStackedVenn(pdf *gopdf.GoPdf, diagram smartart.SmartArt) {
	nodes := smartArtNodes(diagram)
	if len(nodes) == 0 {
		return
	}
	x, y, w, h := smartArtBounds(diagram)
	// The stack is drawn from a 320pt outer circle down; it needs that much
	// room in both directions, and is scaled down when the frame has less.
	outer, step, base := 320.0, 68.0, 324.0
	scale := smartArtFitScale(w, h, outer, base)
	outer, step, base = outer*scale, step*scale, base*scale

	cx := x + w/2
	baseY := y + base
	for i, node := range nodes {
		diameter := math.Max(1, outer-float64(i)*step)
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
	cx, cy := x+w/2, y+h/2

	centerSize := math.Max(32, math.Min(220, math.Min(w, h)*0.42))
	centerLeft := cx - centerSize/2
	centerTop := cy - centerSize/2

	outerSize := math.Max(22, math.Min(110, math.Min(w, h)*0.22))
	maxRadius := math.Max(0, math.Min((w-outerSize)/2, (h-outerSize)/2)-4)
	radius := math.Min(maxRadius, math.Max(outerSize*1.3, math.Min(w, h)*0.32))
	if radius < 0 {
		radius = 0
	}

	if len(nodes) > 1 && radius > 0 {
		orbit := nodes[1:]
		for i, node := range orbit {
			angle := -math.Pi/2 + (2*math.Pi*float64(i))/float64(len(orbit))
			outerCX := cx + math.Cos(angle)*radius
			outerCY := cy + math.Sin(angle)*radius
			outerLeft := outerCX - outerSize/2
			outerTop := outerCY - outerSize/2

			drawSmartArtLine(pdf, cx, cy, outerCX, outerCY)
			drawSmartArtEllipse(
				pdf,
				outerLeft,
				outerTop,
				outerSize,
				outerSize,
				smartArtLightFill,
				smartArtWhiteStroke,
				1,
			)
			drawSmartArtCenteredText(
				pdf,
				node.Text,
				outerLeft+6,
				outerTop+6,
				outerSize-12,
				outerSize-12,
				smartArtInkText,
				22,
			)
		}
	}

	drawSmartArtEllipse(pdf, centerLeft, centerTop, centerSize, centerSize, smartArtBlueFill, smartArtBlueFill, 1)
	drawSmartArtCenteredText(
		pdf,
		nodes[0].Text,
		centerLeft+12,
		centerTop+12,
		centerSize-24,
		centerSize-24,
		smartArtBlueText,
		40,
	)
}
