//nolint:mnd // Connector geometry uses fixed ratios/segment defaults chosen for PPT-like rendering.
package export

import (
	"math"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func renderPDFConnector(pdf *gopdf.GoPdf, c shapes.Connector) {
	x1 := emuToPt(c.StartX.Emu())
	y1 := emuToPt(c.StartY.Emu())
	x2 := emuToPt(c.EndX.Emu())
	y2 := emuToPt(c.EndY.Emu())

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

	switch shapes.NormalizeConnectorType(c.Type) {
	case shapes.ConnectorTypeElbow:
		midX := x1 + (x2-x1)/2
		pdf.Line(x1, y1, midX, y1)
		pdf.Line(midX, y1, midX, y2)
		pdf.Line(midX, y2, x2, y2)
		startAngle = math.Atan2(0, midX-x1)
		endAngle = math.Atan2(0, x2-midX)
		labelX, labelY = midX, y1+(y2-y1)/2
	case shapes.ConnectorTypeCurved:
		cx, cy := connectorControlPoint(x1, y1, x2, y2, c.Adjustments)
		pts := quadraticBezierPoints(x1, y1, cx, cy, x2, y2, 20)
		for i := range len(pts) - 1 {
			pdf.Line(pts[i].X, pts[i].Y, pts[i+1].X, pts[i+1].Y)
		}
		startAngle = math.Atan2(pts[1].Y-pts[0].Y, pts[1].X-pts[0].X)
		last := len(pts) - 1
		endAngle = math.Atan2(pts[last].Y-pts[last-1].Y, pts[last].X-pts[last-1].X)
		labelX, labelY = pts[len(pts)/2].X, pts[len(pts)/2].Y
	default:
		pdf.Line(x1, y1, x2, y2)
	}

	if shapes.NormalizeArrowType(c.StartArrow) != shapes.ArrowTypeNone {
		drawPDFArrowhead(pdf, x1, y1, startAngle+math.Pi, c.StartArrowLen)
	}
	if shapes.NormalizeArrowType(c.EndArrow) != shapes.ArrowTypeNone {
		drawPDFArrowhead(pdf, x2, y2, endAngle, c.EndArrowLen)
	}

	if c.Label != "" {
		pdf.SetTextColor(0, 0, 0)
		pdf.SetX(labelX + 4)
		pdf.SetY(labelY - 10)
		_ = pdf.Cell(nil, c.Label)
	}

	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(1)
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
		if strings.TrimSpace(adj.Name) != "adj1" {
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

func drawPDFArrowhead(pdf *gopdf.GoPdf, tipX, tipY, angle float64, sizeToken string) {
	size := arrowSizePt(sizeToken)
	left := gopdf.Point{
		X: tipX - size*math.Cos(angle-math.Pi/6),
		Y: tipY - size*math.Sin(angle-math.Pi/6),
	}
	right := gopdf.Point{
		X: tipX - size*math.Cos(angle+math.Pi/6),
		Y: tipY - size*math.Sin(angle+math.Pi/6),
	}
	pdf.Polygon([]gopdf.Point{{X: tipX, Y: tipY}, left, right}, "DF")
}

func arrowSizePt(sizeToken string) float64 {
	switch shapes.NormalizeArrowSize(sizeToken) {
	case shapes.ArrowSizeSmall:
		return 4
	case shapes.ArrowSizeLarge:
		return 8
	default:
		return 6
	}
}
