//nolint:mnd // Preset outlines use fixed proportions from the DrawingML geometry.
package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// The presets whose outline is a plain polygon: the remaining regular n-gons,
// the trapezoid and wedge family, the maths glyphs, the seals, the tabbed
// rectangles and the simple solids. Each fell through to a rectangle before.

// drawPDFPolygonGeometry draws one of the presets below and reports whether it
// recognised the type.
func drawPDFPolygonGeometry( //nolint:funlen // One branch per preset is the clearest form for a geometry table.
	pdf *gopdf.GoPdf,
	fl flipState,
	shapeType string,
	x, y, w, h float64,
	style string,
) bool {
	if drawPDFSealGeometry(pdf, fl, shapeType, x, y, w, h, style) {
		return true
	}
	switch shapeType {
	case shapes.ShapeTypeHeptagon:
		fl.polygon(pdf, regularPolygonPoints(x+w/2, y+h/2, w/2, h/2, 7, -math.Pi/2), style)
	case shapes.ShapeTypeDecagon:
		fl.polygon(pdf, regularPolygonPoints(x+w/2, y+h/2, w/2, h/2, 10, -math.Pi/10), style)
	case shapes.ShapeTypeDodecagon:
		fl.polygon(pdf, regularPolygonPoints(x+w/2, y+h/2, w/2, h/2, 12, -math.Pi/12), style)
	case shapes.ShapeTypeIsoscelesTrapezoid:
		fl.polygon(pdf, trapezoidPoints(x, y, w, h, shapeOffsetRatio, shapeOffsetRatio), style)
	case shapes.ShapeTypeNonIsoscelesTrapezoid:
		fl.polygon(pdf, trapezoidPoints(x, y, w, h, shapeOffsetRatio, shapeOffsetRatio*2), style)
	case shapes.ShapeTypeHomePlate, shapes.ShapeTypePentagonArrow:
		fl.polygon(pdf, homePlatePoints(x, y, w, h), style)
	case shapes.ShapeTypeLShape:
		fl.polygon(pdf, cornerPoints(x, y, w, h), style)
	case shapes.ShapeTypeCone:
		fl.polygon(pdf, conePoints(x, y, w, h), style)
	case shapes.ShapeTypeCylinder:
		fl.polygon(pdf, cylinderPoints(x, y, w, h), style)
	case shapes.ShapeTypeFunnel:
		fl.polygon(pdf, funnelPoints(x, y, w, h), style)
	case shapes.ShapeTypeMathMultiply, shapes.ShapeTypeChartX:
		fl.polygon(pdf, saltirePoints(x, y, w, h), style)
	case shapes.ShapeTypeChartPlus:
		fl.polygon(pdf, crossPoints(x, y, w, h), style)
	case shapes.ShapeTypeMathEqual:
		fl.polygon(pdf, equalsBarPoints(x, y, w, h, true), style)
		fl.polygon(pdf, equalsBarPoints(x, y, w, h, false), style)
	case shapes.ShapeTypeMathNotEqual:
		fl.polygon(pdf, equalsBarPoints(x, y, w, h, true), style)
		fl.polygon(pdf, equalsBarPoints(x, y, w, h, false), style)
		fl.polygon(pdf, notEqualSlashPoints(x, y, w, h), style)
	case shapes.ShapeTypeMathDivide:
		drawDivideGlyph(pdf, fl, x, y, w, h, style)
	case shapes.ShapeTypeSquareTabs:
		drawTabbedRect(pdf, fl, x, y, w, h, style, tabsSquare)
	case shapes.ShapeTypeCornerTabs:
		drawTabbedRect(pdf, fl, x, y, w, h, style, tabsCorner)
	case shapes.ShapeTypePlaqueTabs:
		drawTabbedRect(pdf, fl, x, y, w, h, style, tabsPlaque)
	case shapes.ShapeTypeFlowChartData:
		fl.polygon(pdf, flowParallelogramPoints(x, y, w, h), style)
	case shapes.ShapeTypeFlowChartPunchedTape:
		fl.polygon(pdf, flowDocumentPoints(x, y, w, h), style)
	case shapes.ShapeTypePictureFrame:
		pdf.RectFromUpperLeftWithStyle(x, y, w, h, style)
	case shapes.ShapeTypeLightningBolt:
		fl.polygon(pdf, lightningBoltPoints(x, y, w, h), style)
	case shapes.ShapeTypeMusicNote:
		fl.polygon(pdf, musicNotePoints(x, y, w, h), style)
	default:
		return false
	}
	return true
}

const (
	// mathBarRatio is the thickness of an equals bar, as a fraction of height.
	mathBarRatio = 0.15
	// mathBarGapRatio is the gap between the two equals bars.
	mathBarGapRatio = 0.2
	// saltireArmRatio is the thickness of a multiplication cross's arms.
	saltireArmRatio = 0.25
	// divideDotRatio is the size of a division sign's dots.
	divideDotRatio = 0.12
	// coneTopRatio is how wide a cone's top ellipse is relative to its base.
	coneTopRatio = 0.25
	// funnelNeckRatio is how wide a funnel's neck is relative to its mouth.
	funnelNeckRatio = 0.25
	// sunRayRatio is the inner radius of the sun's rays.
	sunRayRatio = 0.62
	// sunRays is how many rays the sun preset draws.
	sunRays = 8
	// tabRatio is the size of a tab, as a fraction of the shorter side.
	tabRatio = 0.2
)

// trapezoidPoints is a trapezoid with independent left and right insets, which
// is what separates the isosceles preset from the non-isosceles one.
func trapezoidPoints(x, y, w, h, leftInset, rightInset float64) []gopdf.Point {
	return []gopdf.Point{
		{X: x + w*leftInset, Y: y}, {X: x + w - w*rightInset, Y: y},
		{X: x + w, Y: y + h}, {X: x, Y: y + h},
	}
}

// homePlatePoints is the pentagon arrow: a rectangle with a point on its right.
func homePlatePoints(x, y, w, h float64) []gopdf.Point {
	point := w * shapeOffsetRatio
	return []gopdf.Point{
		{X: x, Y: y}, {X: x + w - point, Y: y}, {X: x + w, Y: y + h/2},
		{X: x + w - point, Y: y + h}, {X: x, Y: y + h},
	}
}

// conePoints is the silhouette of a cone: a narrow top over a round base.
func conePoints(x, y, w, h float64) []gopdf.Point {
	topRX := (w / 2) * coneTopRatio
	capH := h * cylinderCapRatio
	points := []gopdf.Point{{X: x + w/2 - topRX, Y: y + capH}, {X: x + w/2 + topRX, Y: y + capH}}
	// The base curves away from the viewer, as the cylinder's does.
	for i := range arcSegments + 1 {
		angle := math.Pi * float64(i) / float64(arcSegments)
		points = append(points, gopdf.Point{
			X: x + w/2 + (w/2)*math.Cos(-angle),
			Y: y + h - capH + capH*math.Sin(-angle),
		})
	}
	return points
}

// funnelPoints is a wide mouth tapering to a neck.
func funnelPoints(x, y, w, h float64) []gopdf.Point {
	neck := (w / 2) * funnelNeckRatio
	capH := h * cylinderCapRatio
	points := make([]gopdf.Point, 0, arcSegments+4)
	for i := range arcSegments + 1 {
		angle := math.Pi - math.Pi*float64(i)/float64(arcSegments)
		points = append(points, gopdf.Point{
			X: x + w/2 + (w/2)*math.Cos(angle),
			Y: y + capH - capH*math.Sin(angle),
		})
	}
	return append(points,
		gopdf.Point{X: x + w/2 + neck, Y: y + h},
		gopdf.Point{X: x + w/2 - neck, Y: y + h},
	)
}

// saltirePoints is the diagonal cross of the multiplication sign.
func saltirePoints(x, y, w, h float64) []gopdf.Point {
	armX := w * saltireArmRatio / 2
	armY := h * saltireArmRatio / 2
	cx, cy := x+w/2, y+h/2
	return []gopdf.Point{
		{X: x, Y: y + armY}, {X: x + armX, Y: y},
		{X: cx, Y: cy - armY}, {X: x + w - armX, Y: y},
		{X: x + w, Y: y + armY}, {X: cx + armX, Y: cy},
		{X: x + w, Y: y + h - armY}, {X: x + w - armX, Y: y + h},
		{X: cx, Y: cy + armY}, {X: x + armX, Y: y + h},
		{X: x, Y: y + h - armY}, {X: cx - armX, Y: cy},
	}
}

// equalsBarPoints is one of the two bars of an equals sign.
func equalsBarPoints(x, y, w, h float64, upper bool) []gopdf.Point {
	bar := h * mathBarRatio
	gap := h * mathBarGapRatio
	top := y + h/2 - gap/2 - bar
	if !upper {
		top = y + h/2 + gap/2
	}
	return []gopdf.Point{
		{X: x, Y: top}, {X: x + w, Y: top},
		{X: x + w, Y: top + bar}, {X: x, Y: top + bar},
	}
}

// notEqualSlashPoints is the stroke drawn through an equals sign.
func notEqualSlashPoints(x, y, w, h float64) []gopdf.Point {
	half := w * mathBarRatio / 2
	return []gopdf.Point{
		{X: x + w*0.6 - half, Y: y}, {X: x + w*0.6 + half, Y: y},
		{X: x + w*0.4 + half, Y: y + h}, {X: x + w*0.4 - half, Y: y + h},
	}
}

// drawDivideGlyph draws the division sign: a bar between two dots.
func drawDivideGlyph(pdf *gopdf.GoPdf, fl flipState, x, y, w, h float64, style string) {
	bar := h * mathBarRatio
	fl.polygon(pdf, []gopdf.Point{
		{X: x, Y: y + h/2 - bar/2}, {X: x + w, Y: y + h/2 - bar/2},
		{X: x + w, Y: y + h/2 + bar/2}, {X: x, Y: y + h/2 + bar/2},
	}, style)
	radius := math.Min(w, h) * divideDotRatio
	for _, dotY := range []float64{y + h*0.22, y + h*0.78} {
		fl.polygon(pdf, ellipsePoints(x+w/2, dotY, radius, radius), style)
	}
}

// irregularSealPoints is the ragged burst of the two irregular seals. The radii
// repeat a fixed pattern rather than being random, so the same deck always
// exports the same shape.
func irregularSealPoints(cx, cy, rx, ry float64, seed []float64) []gopdf.Point {
	points := make([]gopdf.Point, 0, len(seed))
	for i, factor := range seed {
		angle := 2 * math.Pi * float64(i) / float64(len(seed))
		points = append(points, gopdf.Point{
			X: cx + rx*factor*math.Cos(angle),
			Y: cy + ry*factor*math.Sin(angle),
		})
	}
	return points
}

// The seal radii, alternating long and short spikes with the irregularity the
// preset is named for.
//
//nolint:gochecknoglobals // Immutable shape data shared by every seal drawn.
var (
	irregularSeal1Seed = []float64{
		1.00, 0.62, 0.92, 0.55, 0.98, 0.60, 0.86, 0.52,
		0.96, 0.58, 0.90, 0.50, 0.94, 0.64, 0.88, 0.56,
	}
	irregularSeal2Seed = []float64{
		0.95, 0.52, 1.00, 0.60, 0.88, 0.50, 0.97, 0.58,
		0.90, 0.54, 0.99, 0.56, 0.85, 0.62, 0.93, 0.51,
	}
)

// sunPoints is a disc with triangular rays around it.
func sunPoints(cx, cy, rx, ry float64) []gopdf.Point {
	return starPoints(cx, cy, rx, ry, sunRays, sunRayRatio, -math.Pi/2)
}

// tabStyle is which corners of a tabbed rectangle carry a tab.
type tabStyle uint8

const (
	tabsSquare tabStyle = iota
	tabsCorner
	tabsPlaque
)

// drawTabbedRect draws the tab family: a rectangle with small squares at its
// corners. The preset differs only in which corners are tabbed.
func drawTabbedRect(pdf *gopdf.GoPdf, fl flipState, x, y, w, h float64, style string, tabs tabStyle) {
	pdf.RectFromUpperLeftWithStyle(x, y, w, h, style)
	tab := math.Min(w, h) * tabRatio
	corners := [][2]float64{{x, y}, {x + w - tab, y}, {x + w - tab, y + h - tab}, {x, y + h - tab}}
	if tabs == tabsCorner {
		// Corner tabs sit outside the box's own corners.
		corners = [][2]float64{{x - tab, y - tab}, {x + w, y - tab}, {x + w, y + h}, {x - tab, y + h}}
	}
	if tabs == tabsPlaque {
		// The plaque's tabs are only on the diagonal pair.
		corners = corners[:1:1]
		corners = append(corners, [2]float64{x + w - tab, y + h - tab})
	}
	for _, corner := range corners {
		fl.polygon(pdf, []gopdf.Point{
			{X: corner[0], Y: corner[1]}, {X: corner[0] + tab, Y: corner[1]},
			{X: corner[0] + tab, Y: corner[1] + tab}, {X: corner[0], Y: corner[1] + tab},
		}, style)
	}
}

// lightningBoltPoints is the zig-zag bolt.
func lightningBoltPoints(x, y, w, h float64) []gopdf.Point {
	return []gopdf.Point{
		{X: x + w*0.45, Y: y}, {X: x + w*0.85, Y: y + h*0.38},
		{X: x + w*0.58, Y: y + h*0.42}, {X: x + w, Y: y + h},
		{X: x + w*0.38, Y: y + h*0.60}, {X: x + w*0.62, Y: y + h*0.55},
		{X: x, Y: y + h*0.22},
	}
}

// musicNotePoints is a quaver: the head with its stem and flag.
func musicNotePoints(x, y, w, h float64) []gopdf.Point {
	headRX, headRY := w*0.28, h*0.16
	headCX, headCY := x+headRX, y+h-headRY
	points := ellipsePoints(headCX, headCY, headRX, headRY)
	return append(points,
		gopdf.Point{X: headCX + headRX, Y: y + h*0.2},
		gopdf.Point{X: x + w, Y: y},
		gopdf.Point{X: x + w, Y: y + h*0.18},
		gopdf.Point{X: headCX + headRX*0.6, Y: y + h*0.34},
	)
}

// drawPDFSealGeometry draws the star-shaped presets — the seals, the two
// irregular bursts and the sun — and reports whether it recognised the type.
// They are split out of the polygon table so neither dispatch grows past the
// complexity the linter allows.
func drawPDFSealGeometry(
	pdf *gopdf.GoPdf,
	fl flipState,
	shapeType string,
	x, y, w, h float64,
	style string,
) bool {
	cx, cy, rx, ry := x+w/2, y+h/2, w/2, h/2
	switch shapeType {
	case shapes.ShapeTypeSeal, shapes.ShapeTypeChartStar:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints5, starInnerRatio5, -math.Pi/2), style)
	case shapes.ShapeTypeSeal4:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints4, starInnerRatio4, -math.Pi/starPoints4), style)
	case shapes.ShapeTypeSeal8:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints8, starInnerRatio8, -math.Pi/starPoints8), style)
	case shapes.ShapeTypeSeal16:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints16, starInnerRatio16, -math.Pi/starPoints16), style)
	case shapes.ShapeTypeSeal32:
		fl.polygon(pdf, starPoints(cx, cy, rx, ry, starPoints32, starInnerRatio32, -math.Pi/starPoints32), style)
	case shapes.ShapeTypeIrregularSeal1:
		fl.polygon(pdf, irregularSealPoints(cx, cy, rx, ry, irregularSeal1Seed), style)
	case shapes.ShapeTypeIrregularSeal2:
		fl.polygon(pdf, irregularSealPoints(cx, cy, rx, ry, irregularSeal2Seed), style)
	case shapes.ShapeTypeSun:
		fl.polygon(pdf, sunPoints(cx, cy, rx, ry), style)
	default:
		return false
	}
	return true
}
