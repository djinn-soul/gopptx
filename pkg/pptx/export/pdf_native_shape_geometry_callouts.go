//nolint:mnd // Callout, button and connector outlines use fixed proportions from the DrawingML geometry.
package export

import (
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// The last three families: the line callouts, which are a box with a leader
// drawn out to a point; the action buttons, which are a rounded box with a
// glyph on it; and the connector presets, which are lines rather than areas.

// drawPDFCalloutGeometry draws one of the presets below and reports whether it
// recognised the type.
func drawPDFCalloutGeometry(
	pdf *gopdf.GoPdf,
	fl flipState,
	shapeType string,
	x, y, w, h float64,
	style string,
) bool {
	switch {
	case isLineCalloutPreset(shapeType):
		drawLineCallout(pdf, fl, shapeType, x, y, w, h, style)
	case strings.HasPrefix(shapeType, actionButtonPrefix):
		drawActionButton(pdf, fl, shapeType, x, y, w, h, style)
	case isConnectorPreset(shapeType):
		fl.polyline(pdf, connectorPresetPoints(shapeType, x, y, w, h))
	default:
		return false
	}
	return true
}

const (
	// actionButtonPrefix is the shared prefix of every action-button preset.
	actionButtonPrefix = "actionButton"
	// calloutBodyRatio is how much of the box the callout's text body takes,
	// leaving the rest of the height for the leader to run through.
	calloutBodyRatio = 0.65
	// calloutLeaderInset is where along the body's edge the leader starts.
	calloutLeaderInset = 0.25
	// buttonGlyphRatio is how much of a button the glyph in its middle takes.
	buttonGlyphRatio = 0.4
)

// isLineCalloutPreset reports the callout family drawn as a box with a leader:
// callout1..3 and their bordered and accented variants.
func isLineCalloutPreset(shapeType string) bool {
	switch shapeType {
	case shapes.ShapeTypeCallout1, shapes.ShapeTypeCallout2, shapes.ShapeTypeCallout3,
		shapes.ShapeTypeBorderCallout1, shapes.ShapeTypeBorderCallout2, shapes.ShapeTypeBorderCallout3,
		shapes.ShapeTypeAccentCallout1, shapes.ShapeTypeAccentCallout2, shapes.ShapeTypeAccentCallout3,
		shapes.ShapeTypeAccentBorderCallout1, shapes.ShapeTypeAccentBorderCallout2,
		shapes.ShapeTypeAccentBorderCallout3:
		return true
	default:
		return false
	}
}

// calloutLeaderSegments is how many bends the leader has: callout1 runs
// straight, callout2 turns once and callout3 twice.
func calloutLeaderSegments(shapeType string) int {
	switch {
	case strings.HasSuffix(shapeType, "3"):
		return 3
	case strings.HasSuffix(shapeType, "2"):
		return 2
	default:
		return 1
	}
}

// drawLineCallout draws the body, then the leader out of its lower-left corner,
// then the accent bar for the accented variants.
func drawLineCallout(
	pdf *gopdf.GoPdf,
	fl flipState,
	shapeType string,
	x, y, w, h float64,
	style string,
) {
	bodyH := h * calloutBodyRatio
	pdf.RectFromUpperLeftWithStyle(x, y, w, bodyH, style)

	if strings.Contains(shapeType, "accent") || strings.Contains(shapeType, "Accent") {
		// The accent is the bar down the body's leading edge.
		fl.polygon(pdf, []gopdf.Point{
			{X: x, Y: y}, {X: x + w*0.03, Y: y},
			{X: x + w*0.03, Y: y + bodyH}, {X: x, Y: y + bodyH},
		}, style)
	}
	fl.polyline(pdf, calloutLeaderPoints(x, y, w, h, bodyH, calloutLeaderSegments(shapeType)))
}

// calloutLeaderPoints is the leader running from the body down to the point it
// annotates, with one bend per segment beyond the first.
func calloutLeaderPoints(x, y, w, h, bodyH float64, segments int) []gopdf.Point {
	start := gopdf.Point{X: x + w*calloutLeaderInset, Y: y + bodyH}
	target := gopdf.Point{X: x, Y: y + h}
	switch segments {
	case 1:
		return []gopdf.Point{start, target}
	case 2:
		return []gopdf.Point{start, {X: start.X, Y: y + h}, target}
	default:
		mid := y + bodyH + (h-bodyH)/2
		return []gopdf.Point{start, {X: start.X, Y: mid}, {X: x + w*0.1, Y: mid}, target}
	}
}

// drawActionButton draws the button's rounded body and the glyph that names it.
// The glyphs are simple marks rather than the icons PowerPoint ships, which is
// still a great deal closer than the blank rectangle they used to be.
func drawActionButton(
	pdf *gopdf.GoPdf,
	fl flipState,
	shapeType string,
	x, y, w, h float64,
	style string,
) {
	drawRoundedRect(pdf, x, y, w, h, style)

	glyphW := w * buttonGlyphRatio
	glyphH := h * buttonGlyphRatio
	cx, cy := x+w/2, y+h/2
	gx, gy := cx-glyphW/2, cy-glyphH/2

	switch shapeType {
	case shapes.ShapeTypeActionButtonBlank:
		return
	case shapes.ShapeTypeActionButtonForwardNext:
		fl.polygon(pdf, trianglePointingRight(gx, gy, glyphW, glyphH), "F")
	case shapes.ShapeTypeActionButtonBackPrevious:
		fl.polygon(pdf, mirrorHorizontally(trianglePointingRight(gx, gy, glyphW, glyphH), gx, glyphW), "F")
	case shapes.ShapeTypeActionButtonBeginning:
		fl.polygon(pdf, mirrorHorizontally(trianglePointingRight(gx, gy, glyphW, glyphH), gx, glyphW), "F")
		fl.polygon(pdf, barPoints(gx, gy, glyphW*0.15, glyphH), "F")
	case shapes.ShapeTypeActionButtonEnd:
		fl.polygon(pdf, trianglePointingRight(gx, gy, glyphW, glyphH), "F")
		fl.polygon(pdf, barPoints(gx+glyphW*0.85, gy, glyphW*0.15, glyphH), "F")
	case shapes.ShapeTypeActionButtonHome:
		fl.polygon(pdf, homeGlyphPoints(gx, gy, glyphW, glyphH), "F")
	case shapes.ShapeTypeActionButtonReturn:
		fl.polyline(pdf, returnGlyphPoints(gx, gy, glyphW, glyphH))
	case shapes.ShapeTypeActionButtonDocument:
		fl.polygon(pdf, documentGlyphPoints(gx, gy, glyphW, glyphH), "F")
	case shapes.ShapeTypeActionButtonSound:
		fl.polygon(pdf, soundGlyphPoints(gx, gy, glyphW, glyphH), "F")
	case shapes.ShapeTypeActionButtonMovie:
		fl.polygon(pdf, documentGlyphPoints(gx, gy, glyphW, glyphH), "F")
		fl.polygon(pdf, trianglePointingRight(gx+glyphW*0.25, gy+glyphH*0.25, glyphW*0.5, glyphH*0.5), "F")
	default:
		// Help and Information are lettered marks; a centred dot stands in for
		// the glyph rather than mislabelling the button.
		fl.polygon(pdf, ellipsePoints(cx, cy, glyphW*0.18, glyphH*0.18), "F")
	}
}

func trianglePointingRight(x, y, w, h float64) []gopdf.Point {
	return []gopdf.Point{{X: x, Y: y}, {X: x + w, Y: y + h/2}, {X: x, Y: y + h}}
}

func barPoints(x, y, w, h float64) []gopdf.Point {
	return []gopdf.Point{{X: x, Y: y}, {X: x + w, Y: y}, {X: x + w, Y: y + h}, {X: x, Y: y + h}}
}

// homeGlyphPoints is a roof over a box.
func homeGlyphPoints(x, y, w, h float64) []gopdf.Point {
	return []gopdf.Point{
		{X: x + w/2, Y: y}, {X: x + w, Y: y + h*0.45}, {X: x + w*0.8, Y: y + h*0.45},
		{X: x + w*0.8, Y: y + h}, {X: x + w*0.2, Y: y + h}, {X: x + w*0.2, Y: y + h*0.45},
		{X: x, Y: y + h*0.45},
	}
}

// returnGlyphPoints is the hooked arrow of the return button.
func returnGlyphPoints(x, y, w, h float64) []gopdf.Point {
	return []gopdf.Point{
		{X: x + w, Y: y}, {X: x + w, Y: y + h*0.6}, {X: x + w*0.2, Y: y + h*0.6},
		{X: x + w*0.2, Y: y + h*0.35}, {X: x, Y: y + h*0.7},
		{X: x + w*0.2, Y: y + h}, {X: x + w*0.2, Y: y + h*0.8},
	}
}

// documentGlyphPoints is a page with its corner turned.
func documentGlyphPoints(x, y, w, h float64) []gopdf.Point {
	fold := w * 0.3
	return []gopdf.Point{
		{X: x, Y: y}, {X: x + w - fold, Y: y}, {X: x + w, Y: y + fold},
		{X: x + w, Y: y + h}, {X: x, Y: y + h},
	}
}

// soundGlyphPoints is a speaker cone.
func soundGlyphPoints(x, y, w, h float64) []gopdf.Point {
	return []gopdf.Point{
		{X: x, Y: y + h*0.35}, {X: x + w*0.35, Y: y + h*0.35}, {X: x + w*0.7, Y: y},
		{X: x + w*0.7, Y: y + h}, {X: x + w*0.35, Y: y + h*0.65}, {X: x, Y: y + h*0.65},
	}
}

// isConnectorPreset reports the presets that are a line rather than an area.
func isConnectorPreset(shapeType string) bool {
	switch shapeType {
	case shapes.ShapeTypeLine, shapes.ShapeTypeLineInv, shapes.ShapeTypeStraightConnector1,
		shapes.ShapeTypeBentConnector2, shapes.ShapeTypeBentConnector3,
		shapes.ShapeTypeBentConnector4, shapes.ShapeTypeBentConnector5,
		shapes.ShapeTypeCurvedConnector2, shapes.ShapeTypeCurvedConnector3,
		shapes.ShapeTypeCurvedConnector4, shapes.ShapeTypeCurvedConnector5:
		return true
	default:
		return false
	}
}

// connectorPresetPoints routes a connector across its own box: straight for the
// line presets, and with one bend per elbow for the bent and curved ones. A
// connector placed as a shape carries only its box, so the route is derived
// from that rather than from endpoints it does not have.
func connectorPresetPoints(shapeType string, x, y, w, h float64) []gopdf.Point {
	start := gopdf.Point{X: x, Y: y}
	end := gopdf.Point{X: x + w, Y: y + h}
	if shapeType == shapes.ShapeTypeLineInv {
		start, end = gopdf.Point{X: x, Y: y + h}, gopdf.Point{X: x + w, Y: y}
	}
	switch shapeType {
	case shapes.ShapeTypeBentConnector2, shapes.ShapeTypeCurvedConnector2:
		return []gopdf.Point{start, {X: end.X, Y: start.Y}, end}
	case shapes.ShapeTypeBentConnector3, shapes.ShapeTypeCurvedConnector3:
		mid := x + w/2
		return []gopdf.Point{start, {X: mid, Y: start.Y}, {X: mid, Y: end.Y}, end}
	case shapes.ShapeTypeBentConnector4, shapes.ShapeTypeCurvedConnector4:
		mid := y + h/2
		return []gopdf.Point{
			start, {X: x + w/2, Y: start.Y}, {X: x + w/2, Y: mid}, {X: end.X, Y: mid}, end,
		}
	case shapes.ShapeTypeBentConnector5, shapes.ShapeTypeCurvedConnector5:
		return []gopdf.Point{
			start, {X: x + w/3, Y: start.Y}, {X: x + w/3, Y: y + h/2},
			{X: x + 2*w/3, Y: y + h/2}, {X: x + 2*w/3, Y: end.Y}, end,
		}
	default:
		return []gopdf.Point{start, end}
	}
}
