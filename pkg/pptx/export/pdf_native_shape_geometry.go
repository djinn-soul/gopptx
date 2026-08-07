package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// Shape geometry constants.
const (
	// Polygon vertex counts for regular polygons.
	pentagonSides = 5
	octagonSides  = 8

	// hexagonAdjustment is the OOXML hexagon preset's default adj value (25000
	// of 100000), the fraction of the shorter side each end is cut back by.
	hexagonAdjustment = 0.25

	// Star point counts.
	starPoints4  = 4
	starPoints5  = 5
	starPoints6  = 6
	starPoints7  = 7
	starPoints8  = 8
	starPoints10 = 10
	starPoints12 = 12
	starPoints16 = 16
	starPoints24 = 24
	starPoints32 = 32

	// Polygon resolution for callout, heart, and cloud shape approximations.
	calloutPolyRes = 32

	// shapeOffsetRatio is the fraction of width used as the side offset in
	// parallelogram and trapezoid shapes.
	shapeOffsetRatio = 0.2

	// tailHeightRatio is the callout tail height as a fraction of the shape height.
	tailHeightRatio = 0.25

	// Star inner/outer radius ratios (inner ÷ outer).
	starInnerRatio4  = 0.35
	starInnerRatio5  = 0.382
	starInnerRatio6  = 0.5
	starInnerRatio7  = 0.45
	starInnerRatio8  = 0.414
	starInnerRatio10 = 0.38
	starInnerRatio12 = 0.45
	starInnerRatio16 = 0.5
	starInnerRatio24 = 0.6
	starInnerRatio32 = 0.65

	// Heart curve (x=16·sin³t, y=13·cos(t)−5·cos(2t)−2·cos(3t)−cos(4t)) constants.
	heartSinExponent = 3    // exponent for the sin³(t) term
	heartXScale      = 16.0 // x-axis scale factor; also the shift for normalisation
	heartCos3Freq    = 3    // frequency of the 3rd cosine harmonic
	heartCos4Freq    = 4    // frequency of the 4th cosine harmonic
	heartYOrigin     = 13.0 // y-axis shift for normalisation (= first harmonic coefficient)
	heartNormXRange  = 32.0 // normalises x from [−16, 16] to [0, 1]  (= 2 × heartXScale)
	heartNormYRange  = 25.0 // normalises y from [−12, 13] to [0, 1]
)

func getShapeBounds(s shapes.Shape) (float64, float64, float64, float64) {
	return emuToPt(int64(s.X)), emuToPt(int64(s.Y)), emuToPt(int64(s.CX)), emuToPt(int64(s.CY))
}

func setPDFShapeFill(pdf *gopdf.GoPdf, s shapes.Shape, gradientRendered bool) {
	if s.Fill != nil && s.Fill.Color != "" {
		pdf.SetFillColor(hexToRGB(s.Fill.Color))
	} else if !gradientRendered && s.GradientFill != nil && len(s.GradientFill.Stops) > 0 {
		pdf.SetFillColor(hexToRGB(s.GradientFill.Stops[0].Color))
	}
}

func setPDFShapeStroke(pdf *gopdf.GoPdf, s shapes.Shape) bool {
	if s.Line == nil || s.Line.Width <= 0 {
		return false
	}
	strokeWidth := emuToPt(int64(s.Line.Width))
	if strokeWidth < minStrokeWidth {
		strokeWidth = minStrokeWidth
	}
	pdf.SetLineWidth(strokeWidth)
	pdf.SetStrokeColor(hexToRGB(s.Line.Color))
	applyPDFLineDash(pdf, shapeLineDash(s), strokeWidth)
	return true
}

// flipState mirrors a shape's outline about its own centre. gopdf exposes no
// coordinate transform — only Rotate — so a flip has to be applied to the points
// themselves rather than to the graphics state.
//
// Only the asymmetric presets change under it: a rectangle or an ellipse is its
// own mirror, but a chevron, an arrow or a flowchart card is not, and those used
// to draw the wrong way round when a diagram's cached layout flipped them.
type flipState struct {
	H, V           bool
	cx, cy         float64
	unflippedShape bool
}

func flipFor(s shapes.Shape, x, y, w, h float64) flipState {
	return flipState{
		H:              s.FlipH,
		V:              s.FlipV,
		cx:             x + w/2,
		cy:             y + h/2,
		unflippedShape: !s.FlipH && !s.FlipV,
	}
}

// polygon draws points with the shape's flip applied.
func (f flipState) polygon(pdf *gopdf.GoPdf, points []gopdf.Point, style string) {
	pdf.Polygon(f.points(points), style)
}

func (f flipState) points(points []gopdf.Point) []gopdf.Point {
	if f.unflippedShape {
		return points
	}
	out := make([]gopdf.Point, len(points))
	for i, p := range points {
		out[i] = p
		if f.H {
			out[i].X = 2*f.cx - p.X
		}
		if f.V {
			out[i].Y = 2*f.cy - p.Y
		}
	}
	return out
}

func drawPDFGeometry( //nolint:funlen // Shape dispatch requires one branch per supported shape type.
	pdf *gopdf.GoPdf, s shapes.Shape, x, y, w, h float64, style string,
) {
	fl := flipFor(s, x, y, w, h)
	switch s.Type {
	case shapes.ShapeTypeRectangle:
		pdf.RectFromUpperLeftWithStyle(x, y, w, h, style)
	case shapes.ShapeTypeRoundedRectangle:
		drawRoundedRect(pdf, x, y, w, h, style)
	case shapes.ShapeTypePie, shapes.ShapeTypePieWedge, shapes.ShapeTypeChord:
		drawPieShape(pdf, s, x, y, w, h, style)
	case shapes.ShapeTypeEllipse:
		// gopdf Oval only strokes; use a polygon approximation so fill works.
		fl.polygon(pdf, ellipsePoints(x+w/2, y+h/2, w/2, h/2, calloutPolyRes), style)
	case shapes.ShapeTypeTriangle:
		fl.polygon(pdf,
			[]gopdf.Point{{X: x + w/2, Y: y}, {X: x, Y: y + h}, {X: x + w, Y: y + h}},
			style,
		)
	case shapes.ShapeTypeRightTriangle:
		fl.polygon(pdf,
			[]gopdf.Point{{X: x, Y: y}, {X: x + w, Y: y + h}, {X: x, Y: y + h}},
			style,
		)
	case shapes.ShapeTypeDiamond:
		fl.polygon(pdf,
			[]gopdf.Point{{X: x + w/2, Y: y}, {X: x + w, Y: y + h/2}, {X: x + w/2, Y: y + h}, {X: x, Y: y + h/2}},
			style,
		)
	case shapes.ShapeTypeHexagon:
		fl.polygon(pdf, hexagonPoints(x, y, w, h), style)
	case shapes.ShapeTypePentagon:
		fl.polygon(pdf, regularPolygonPoints(x+w/2, y+h/2, w/2, h/2, pentagonSides, -math.Pi/2), style)
	case shapes.ShapeTypeOctagon:
		fl.polygon(pdf, regularPolygonPoints(x+w/2, y+h/2, w/2, h/2, octagonSides, -math.Pi/octagonSides), style)
	case shapes.ShapeTypeParallelogram:
		off := w * shapeOffsetRatio
		fl.polygon(pdf, []gopdf.Point{
			{X: x + off, Y: y}, {X: x + w, Y: y},
			{X: x + w - off, Y: y + h}, {X: x, Y: y + h},
		}, style)
	case shapes.ShapeTypeTrapezoid:
		off := w * shapeOffsetRatio
		fl.polygon(pdf, []gopdf.Point{
			{X: x + off, Y: y}, {X: x + w - off, Y: y},
			{X: x + w, Y: y + h}, {X: x, Y: y + h},
		}, style)
	case shapes.ShapeTypeRightArrow:
		fl.polygon(pdf, rightArrowPoints(x, y, w, h, arrowGeometryFor(s.Adjustments, w, h)), style)
	case shapes.ShapeTypeLeftArrow:
		fl.polygon(pdf, leftArrowPoints(x, y, w, h, arrowGeometryFor(s.Adjustments, w, h)), style)
	case shapes.ShapeTypeUpArrow:
		fl.polygon(pdf, upArrowPoints(x, y, w, h, arrowGeometryFor(s.Adjustments, h, w)), style)
	case shapes.ShapeTypeDownArrow:
		fl.polygon(pdf, downArrowPoints(x, y, w, h, arrowGeometryFor(s.Adjustments, h, w)), style)
	case shapes.ShapeTypeLeftRightArrow:
		fl.polygon(pdf, leftRightArrowPoints(x, y, w, h, arrowGeometryFor(s.Adjustments, w, h)), style)
	case shapes.ShapeTypeUpDownArrow:
		fl.polygon(pdf, upDownArrowPoints(x, y, w, h, arrowGeometryFor(s.Adjustments, h, w)), style)
	case shapes.ShapeTypeChevronArrow:
		fl.polygon(pdf, chevronPoints(x, y, w, h), style)
	case shapes.ShapeTypeStar4, shapes.ShapeTypeStar5, shapes.ShapeTypeStar6,
		shapes.ShapeTypeStar7, shapes.ShapeTypeStar8, shapes.ShapeTypeStar10,
		shapes.ShapeTypeStar12, shapes.ShapeTypeStar16, shapes.ShapeTypeStar24, shapes.ShapeTypeStar32:
		drawPDFStarShape(pdf, fl, s.Type, x, y, w, h, style)
	case shapes.ShapeTypeHeart:
		fl.polygon(pdf, heartPoints(x, y, w, h, calloutPolyRes), style)
	case shapes.ShapeTypeWedgeRectCallout, shapes.ShapeTypeWedgeRRectCallout:
		fl.polygon(pdf, wedgeRectCalloutPoints(x, y, w, h), style)
	case shapes.ShapeTypeWedgeEllipseCallout:
		fl.polygon(pdf, wedgeEllipseCalloutPoints(x, y, w, h, calloutPolyRes), style)
	case shapes.ShapeTypeCloudCallout, shapes.ShapeTypeCloud:
		fl.polygon(pdf, ellipsePoints(x+w/2, y+h/2, w/2, h/2, calloutPolyRes), style)
	default:
		// The presets that draw as a single outline live in their own tables; a
		// rectangle is only the last resort.
		if drawPDFExtendedGeometry(pdf, fl, s.Type, x, y, w, h, style) ||
			drawPDFCornerGeometry(pdf, fl, s.Type, x, y, w, h, style) ||
			drawPDFSolidGeometry(pdf, fl, s.Type, x, y, w, h, style) {
			return
		}
		pdf.RectFromUpperLeftWithStyle(x, y, w, h, style)
	}
}

// drawPDFStarShape renders n-pointed star shapes; extracted to reduce drawPDFGeometry cyclomatic complexity.
func drawPDFStarShape(pdf *gopdf.GoPdf, fl flipState, shapeType string, x, y, w, h float64, style string) {
	cx, cy, rx, ry := x+w/2, y+h/2, w/2, h/2
	switch shapeType {
	case shapes.ShapeTypeStar5:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints5, starInnerRatio5, -math.Pi/2), style)
	case shapes.ShapeTypeStar4:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints4, starInnerRatio4, -math.Pi/starPoints4), style)
	case shapes.ShapeTypeStar6:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints6, starInnerRatio6, -math.Pi/starPoints6), style)
	case shapes.ShapeTypeStar7:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints7, starInnerRatio7, -math.Pi/2), style)
	case shapes.ShapeTypeStar8:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints8, starInnerRatio8, -math.Pi/starPoints8), style)
	case shapes.ShapeTypeStar10:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints10, starInnerRatio10, -math.Pi/starPoints10), style)
	case shapes.ShapeTypeStar12:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints12, starInnerRatio12, -math.Pi/starPoints12), style)
	case shapes.ShapeTypeStar16:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints16, starInnerRatio16, -math.Pi/starPoints16), style)
	case shapes.ShapeTypeStar24:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints24, starInnerRatio24, -math.Pi/starPoints24), style)
	case shapes.ShapeTypeStar32:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints32, starInnerRatio32, -math.Pi/starPoints32), style)
	}
}

// ellipsePoints returns n evenly-spaced points on the perimeter of an ellipse.
func ellipsePoints(cx, cy, rx, ry float64, n int) []gopdf.Point {
	pts := make([]gopdf.Point, n)
	for i := range n {
		a := 2 * math.Pi * float64(i) / float64(n)
		pts[i] = gopdf.Point{X: cx + rx*math.Cos(a), Y: cy + ry*math.Sin(a)}
	}
	return pts
}

// hexagonPoints returns the vertices of the OOXML "hexagon" preset: a box with
// its left and right ends cut back by adj x the shorter side.
//
// Drawing it as a regular hexagon inscribed in the box was wrong for any shape
// that is not square — a wide flowchart node came out with long diagonal ends
// where PowerPoint draws short ones.
func hexagonPoints(x, y, w, h float64) []gopdf.Point {
	inset := math.Min(w, h) * hexagonAdjustment
	if inset > w/2 {
		inset = w / 2
	}
	midY := y + h/2
	return []gopdf.Point{
		{X: x, Y: midY},
		{X: x + inset, Y: y},
		{X: x + w - inset, Y: y},
		{X: x + w, Y: midY},
		{X: x + w - inset, Y: y + h},
		{X: x + inset, Y: y + h},
	}
}

// regularPolygonPoints returns n evenly-spaced vertices of a polygon inscribed
// in an ellipse with semi-axes rx, ry, starting at angle startAngle.
func regularPolygonPoints(cx, cy, rx, ry float64, n int, startAngle float64) []gopdf.Point {
	pts := make([]gopdf.Point, n)
	for i := range n {
		a := startAngle + 2*math.Pi*float64(i)/float64(n)
		pts[i] = gopdf.Point{X: cx + rx*math.Cos(a), Y: cy + ry*math.Sin(a)}
	}
	return pts
}

// starPoints returns the vertices of an n-pointed star.
// innerRatio is inner-radius / outer-radius (e.g. 0.382 for a classic 5-star).
func starPoints(cx, cy, rx, ry float64, n int, innerRatio, startAngle float64) []gopdf.Point {
	pts := make([]gopdf.Point, 2*n)
	step := math.Pi / float64(n)
	for i := range n {
		outerA := startAngle + 2*math.Pi*float64(i)/float64(n)
		innerA := outerA + step
		pts[2*i] = gopdf.Point{X: cx + rx*math.Cos(outerA), Y: cy + ry*math.Sin(outerA)}
		pts[2*i+1] = gopdf.Point{X: cx + rx*innerRatio*math.Cos(innerA), Y: cy + ry*innerRatio*math.Sin(innerA)}
	}
	return pts
}

// heartPoints returns a polygon approximation of a heart shape.
func heartPoints(x, y, w, h float64, n int) []gopdf.Point {
	pts := make([]gopdf.Point, n)
	for i := range n {
		t := 2 * math.Pi * float64(i) / float64(n)
		// Parametric heart: x=16sin³(t), y=13cos(t)-5cos(2t)-2cos(3t)-cos(4t)
		hx := heartXScale * math.Pow(math.Sin(t), heartSinExponent)
		hy := 13*math.Cos(t) - 5*math.Cos(2*t) - 2*math.Cos(heartCos3Freq*t) - math.Cos(heartCos4Freq*t)
		// Normalise from [-16,16]×[-12,13] to [x,x+w]×[y,y+h]
		nx := (hx + heartXScale) / heartNormXRange
		ny := (-hy + heartYOrigin) / heartNormYRange // flip Y (PDF Y grows downward)
		pts[i] = gopdf.Point{X: x + nx*w, Y: y + ny*h}
	}
	return pts
}

// wedgeRectCalloutPoints returns a rectangle with a small triangular pointer
// protruding from the bottom-left corner.
func wedgeRectCalloutPoints(x, y, w, h float64) []gopdf.Point {
	tailH := h * tailHeightRatio
	return []gopdf.Point{
		{X: x, Y: y},
		{X: x + w, Y: y},
		{X: x + w, Y: y + h},
		{X: x + w*0.4, Y: y + h},
		{X: x + w*0.2, Y: y + h + tailH},
		{X: x + w*0.25, Y: y + h},
		{X: x, Y: y + h},
	}
}

// wedgeEllipseCalloutPoints returns an ellipse polygon with a pointer tail.
func wedgeEllipseCalloutPoints(x, y, w, h float64, n int) []gopdf.Point {
	cx, cy := x+w/2, y+h/2
	rx, ry := w/2, h/2
	pts := make([]gopdf.Point, 0, n+2)
	for i := range n {
		a := 2 * math.Pi * float64(i) / float64(n)
		pts = append(pts, gopdf.Point{X: cx + rx*math.Cos(a), Y: cy + ry*math.Sin(a)})
	}
	// Append tail tip below the ellipse.
	pts = append(pts, gopdf.Point{X: cx - rx*0.3, Y: cy + ry + ry*0.5})
	return pts
}
