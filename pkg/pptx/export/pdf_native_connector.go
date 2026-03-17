//nolint:mnd // Connector geometry uses fixed ratios/segment defaults chosen for PPT-like rendering.
package export

import (
	"math"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

const arrowFillDrawStyle = "DF"

const connectorAdjustmentPrimary = "adj1"

//nolint:funlen // Connector rendering keeps geometry/routing/arrow/label logic together for predictable output.
func renderPDFConnector(pdf *gopdf.GoPdf, c shapes.Connector) {
	x1 := emuToPt(c.StartX.Emu())
	y1 := emuToPt(c.StartY.Emu())
	x2 := emuToPt(c.EndX.Emu())
	y2 := emuToPt(c.EndY.Emu())
	startTipX, startTipY := x1, y1
	endTipX, endTipY := x2, y2

	strokeW := emuToPt(c.Line.Width.Emu())
	if strokeW < minStrokeWidth {
		strokeW = minStrokeWidth
	}
	pdf.SetLineWidth(strokeW)
	r, g, b := hexToRGB(c.Line.Color)
	pdf.SetStrokeColor(r, g, b)
	pdf.SetFillColor(r, g, b)

	startAngle := math.Atan2(y2-y1, x2-x1)
	endAngle := startAngle
	labelX, labelY := (x1+x2)/2, (y1+y2)/2
	startInset := connectorArrowInset(c.StartArrowLen, c.StartArrow)
	endInset := connectorArrowInset(c.EndArrowLen, c.EndArrow)
	x1, y1 = insetConnectorPoint(x1, y1, x2, y2, startInset)
	x2, y2 = insetConnectorPoint(x2, y2, x1, y1, endInset)

	switch shapes.NormalizeConnectorType(c.Type) {
	case shapes.ConnectorTypeElbow:
		midX := x1 + (x2-x1)/2
		pdf.Line(x1, y1, midX, y1)
		pdf.Line(midX, y1, midX, y2)
		pdf.Line(midX, y2, x2, y2)
		startAngle = math.Atan2(0, midX-x1)
		endAngle = math.Atan2(0, x2-midX)
		labelX, labelY = connectorLabelPosition(midX, y1+(y2-y1)/2, endAngle, c.Label)
	case shapes.ConnectorTypeCurved:
		cx, cy := connectorControlPoint(x1, y1, x2, y2, c.Adjustments)
		pts := quadraticBezierPoints(x1, y1, cx, cy, x2, y2, 20)
		for i := range len(pts) - 1 {
			pdf.Line(pts[i].X, pts[i].Y, pts[i+1].X, pts[i+1].Y)
		}
		startAngle = math.Atan2(pts[1].Y-pts[0].Y, pts[1].X-pts[0].X)
		last := len(pts) - 1
		endAngle = math.Atan2(pts[last].Y-pts[last-1].Y, pts[last].X-pts[last-1].X)
		labelX, labelY = connectorLabelPosition(pts[len(pts)/2].X, pts[len(pts)/2].Y, endAngle, c.Label)
	default:
		pdf.Line(x1, y1, x2, y2)
		labelX, labelY = connectorLabelPosition(labelX, labelY, endAngle, c.Label)
	}

	if shapes.NormalizeArrowType(c.StartArrow) != shapes.ArrowTypeNone {
		drawPDFArrowhead(
			pdf,
			startTipX,
			startTipY,
			startAngle+math.Pi,
			c.StartArrow,
			c.StartArrowWidth,
			c.StartArrowLen,
		)
	}
	if shapes.NormalizeArrowType(c.EndArrow) != shapes.ArrowTypeNone {
		drawPDFArrowhead(pdf, endTipX, endTipY, endAngle, c.EndArrow, c.EndArrowWidth, c.EndArrowLen)
	}

	if c.Label != "" {
		pdf.SetTextColor(0, 0, 0)
		pdf.SetX(labelX)
		pdf.SetY(labelY)
		_ = pdf.Cell(nil, c.Label)
	}

	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(1)
}

func connectorArrowInset(sizeToken string, arrowType string) float64 {
	if shapes.NormalizeArrowType(arrowType) == shapes.ArrowTypeNone {
		return 0
	}
	_, length := arrowSizePt(shapes.ArrowSizeMedium, sizeToken)
	return length * 0.6
}

func insetConnectorPoint(x, y, towardX, towardY, inset float64) (float64, float64) {
	if inset <= 0 {
		return x, y
	}
	dx := towardX - x
	dy := towardY - y
	dist := math.Hypot(dx, dy)
	if dist <= 0.01 {
		return x, y
	}
	return x + (dx/dist)*inset, y + (dy/dist)*inset
}

func connectorLabelPosition(x, y, lineAngle float64, label string) (float64, float64) {
	if label == "" {
		return x, y
	}
	normal := lineAngle + math.Pi/2
	offset := 7.0
	return x + math.Cos(normal)*offset, y + math.Sin(normal)*offset - 6
}

func connectorControlPoint(x1, y1, x2, y2 float64, adjustments []shapes.ConnectorAdjustment) (float64, float64) {
	// Default control point bows perpendicular to connector direction.
	mx := (x1 + x2) / 2
	my := (y1 + y2) / 2
	dx := x2 - x1
	dy := y2 - y1
	norm := math.Hypot(dx, dy)
	if norm == 0 {
		return mx, my
	}
	perpX := -dy / norm
	perpY := dx / norm
	bendPct := 0.15

	for _, adj := range adjustments {
		if strings.TrimSpace(adj.Name) != connectorAdjustmentPrimary {
			continue
		}
		formula := strings.TrimSpace(adj.Formula)
		if !strings.HasPrefix(formula, "val ") {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(formula, "val ")))
		if err != nil {
			continue
		}
		// OOXML connector adjustments are typically in 1/100000 ranges.
		bendPct = math.Max(-0.5, math.Min(0.5, float64(n)/100000.0))
	}

	mag := norm * bendPct
	return mx + perpX*mag, my + perpY*mag
}

func quadraticBezierPoints(
	x1, y1, cx, cy, x2, y2 float64,
	segments int,
) []gopdf.Point {
	if segments < 4 {
		segments = 4
	}
	pts := make([]gopdf.Point, 0, segments+1)
	for i := 0; i <= segments; i++ {
		t := float64(i) / float64(segments)
		mt := 1 - t
		x := mt*mt*x1 + 2*mt*t*cx + t*t*x2
		y := mt*mt*y1 + 2*mt*t*cy + t*t*y2
		pts = append(pts, gopdf.Point{X: x, Y: y})
	}
	return pts
}

func drawPDFArrowhead(
	pdf *gopdf.GoPdf,
	tipX, tipY, angle float64,
	arrowType string,
	widthToken string,
	lengthToken string,
) {
	halfWidth, length := arrowSizePt(widthToken, lengthToken)
	left := gopdf.Point{
		X: tipX - length*math.Cos(angle) + halfWidth*math.Sin(angle),
		Y: tipY - length*math.Sin(angle) - halfWidth*math.Cos(angle),
	}
	right := gopdf.Point{
		X: tipX - length*math.Cos(angle) - halfWidth*math.Sin(angle),
		Y: tipY - length*math.Sin(angle) + halfWidth*math.Cos(angle),
	}
	back := gopdf.Point{
		X: tipX - (length * 0.55 * math.Cos(angle)),
		Y: tipY - (length * 0.55 * math.Sin(angle)),
	}

	switch shapes.NormalizeArrowType(arrowType) {
	case shapes.ArrowTypeOpen:
		pdf.Line(tipX, tipY, left.X, left.Y)
		pdf.Line(tipX, tipY, right.X, right.Y)
	case shapes.ArrowTypeDiamond:
		rear := gopdf.Point{
			X: tipX - length*math.Cos(angle),
			Y: tipY - length*math.Sin(angle),
		}
		pdf.Polygon([]gopdf.Point{{X: tipX, Y: tipY}, left, rear, right}, arrowFillDrawStyle)
	case shapes.ArrowTypeOval:
		pdf.Oval(
			back.X-halfWidth,
			back.Y-halfWidth,
			back.X+halfWidth,
			back.Y+halfWidth,
		)
	case shapes.ArrowTypeStealth:
		pdf.Polygon([]gopdf.Point{{X: tipX, Y: tipY}, left, back, right}, arrowFillDrawStyle)
	default:
		pdf.Polygon([]gopdf.Point{{X: tipX, Y: tipY}, left, right}, arrowFillDrawStyle)
	}
}

func arrowSizePt(widthToken string, lengthToken string) (float64, float64) {
	halfWidth := 3.0
	length := 6.0
	switch shapes.NormalizeArrowSize(widthToken) {
	case shapes.ArrowSizeSmall:
		halfWidth = 2.0
	case shapes.ArrowSizeLarge:
		halfWidth = 4.0
	}
	switch shapes.NormalizeArrowSize(lengthToken) {
	case shapes.ArrowSizeSmall:
		length = 4.0
	case shapes.ArrowSizeLarge:
		length = 8.0
	}
	return halfWidth, length
}
