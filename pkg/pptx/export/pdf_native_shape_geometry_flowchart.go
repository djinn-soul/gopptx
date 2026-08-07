//nolint:mnd // Flowchart presets are defined by fixed proportions in the DrawingML preset geometry.
package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// Every preset the geometry switch does not name is drawn as a plain rectangle,
// silently. That is a poor answer for the flowchart family, which is most of
// what a process deck is made of: a decision came out as a box, and so did the
// terminator next to it, leaving a diagram of identical rectangles.
//
// These are the presets whose outline is a single polygon or a simple pair of
// them, which covers the flowchart set, the two math signs and the cross.

// drawPDFExtendedGeometry draws one of the presets below and reports whether it
// recognised the type.
func drawPDFExtendedGeometry(
	pdf *gopdf.GoPdf,
	shapeType string,
	x, y, w, h float64,
	style string,
) bool {
	switch shapeType {
	case shapes.ShapeTypeFlowChartProcess, shapes.ShapeTypeFlowChartPredefinedProcess,
		shapes.ShapeTypeFlowChartInternalStorage, shapes.ShapeTypeFrame,
		shapes.ShapeTypeBevel:
		// The inner rules of a predefined-process or internal-storage box are
		// decoration on a rectangle; the frames and bevel are a rectangle with
		// an edge treatment this renderer does not draw.
		pdf.RectFromUpperLeftWithStyle(x, y, w, h, style)
	case shapes.ShapeTypeFlowChartAlternateProcess, shapes.ShapeTypeFlowChartTerminator:
		drawRoundedRect(pdf, x, y, w, h, style)
	case shapes.ShapeTypeFlowChartDecision, shapes.ShapeTypeFlowChartSort:
		pdf.Polygon(diamondPoints(x, y, w, h), style)
	case shapes.ShapeTypeFlowChartInputOutput:
		pdf.Polygon(flowParallelogramPoints(x, y, w, h), style)
	case shapes.ShapeTypeFlowChartPreparation:
		pdf.Polygon(hexagonPoints(x, y, w, h), style)
	case shapes.ShapeTypeFlowChartConnector, shapes.ShapeTypeFlowChartOr,
		shapes.ShapeTypeFlowChartSummingJunction, shapes.ShapeTypeDonut,
		shapes.ShapeTypeNoSmoking, shapes.ShapeTypeBlockArc:
		pdf.Polygon(ellipsePoints(x+w/2, y+h/2, w/2, h/2, calloutPolyRes), style)
	case shapes.ShapeTypeFlowChartExtract:
		pdf.Polygon([]gopdf.Point{
			{X: x + w/2, Y: y}, {X: x + w, Y: y + h}, {X: x, Y: y + h},
		}, style)
	case shapes.ShapeTypeFlowChartMerge:
		pdf.Polygon([]gopdf.Point{
			{X: x, Y: y}, {X: x + w, Y: y}, {X: x + w/2, Y: y + h},
		}, style)
	case shapes.ShapeTypeFlowChartCollate:
		pdf.Polygon([]gopdf.Point{
			{X: x, Y: y}, {X: x + w, Y: y}, {X: x, Y: y + h}, {X: x + w, Y: y + h},
		}, style)
	case shapes.ShapeTypeFlowChartOffpageConnector:
		pdf.Polygon(offpageConnectorPoints(x, y, w, h), style)
	case shapes.ShapeTypeFlowChartCard, shapes.ShapeTypeFlowChartPunchedCard:
		pdf.Polygon(clippedCornerPoints(x, y, w, h), style)
	case shapes.ShapeTypeFlowChartManualInput:
		pdf.Polygon([]gopdf.Point{
			{X: x, Y: y + h*flowManualInputSlope}, {X: x + w, Y: y},
			{X: x + w, Y: y + h}, {X: x, Y: y + h},
		}, style)
	case shapes.ShapeTypeFlowChartManualOperation:
		pdf.Polygon([]gopdf.Point{
			{X: x, Y: y}, {X: x + w, Y: y},
			{X: x + w - w*flowTrapezoidInset, Y: y + h}, {X: x + w*flowTrapezoidInset, Y: y + h},
		}, style)
	case shapes.ShapeTypeFlowChartDocument, shapes.ShapeTypeFlowChartMultidocument:
		pdf.Polygon(flowDocumentPoints(x, y, w, h), style)
	case shapes.ShapeTypeFlowChartDelay, shapes.ShapeTypeFlowChartDisplay:
		pdf.Polygon(flowDelayPoints(x, y, w, h), style)
	case shapes.ShapeTypeFlowChartMagneticDisk, shapes.ShapeTypeFlowChartDirectAccessStorage,
		shapes.ShapeTypeFlowChartMagneticDrum, shapes.ShapeTypeCan:
		pdf.Polygon(cylinderPoints(x, y, w, h), style)
	case shapes.ShapeTypeFlowChartStoredData, shapes.ShapeTypeFlowChartOnlineStorage,
		shapes.ShapeTypeFlowChartSequentialAccessStorage, shapes.ShapeTypeFlowChartMagneticTape:
		pdf.Polygon(flowStoredDataPoints(x, y, w, h), style)
	case shapes.ShapeTypePlus, shapes.ShapeTypeMathPlus:
		pdf.Polygon(crossPoints(x, y, w, h), style)
	case shapes.ShapeTypeMathMinus:
		bar := h * mathBarThickness
		pdf.RectFromUpperLeftWithStyle(x, y+(h-bar)/2, w, bar, style)
	default:
		return false
	}
	return true
}

const (
	// flowManualInputSlope is how far down the left edge the sloped top of a
	// manual-input box starts, as a fraction of its height.
	flowManualInputSlope = 0.2
	// flowTrapezoidInset is how far in each bottom corner of a manual-operation
	// box sits, as a fraction of its width.
	flowTrapezoidInset = 0.2
	// flowWaveDepth is the height of the wave along the bottom of a document
	// box, as a fraction of the box height.
	flowWaveDepth = 0.12
	// flowClipRatio is the corner a card box has cut off, as a fraction of the
	// shorter side.
	flowClipRatio = 0.2
	// flowOffpageNotch is how much of the height the point of an off-page
	// connector takes.
	flowOffpageNotch = 0.25
	// cylinderCapRatio is the height of a cylinder's elliptical cap, as a
	// fraction of the whole.
	cylinderCapRatio = 0.16
	// crossArmRatio is the thickness of a cross's arms, as a fraction of the
	// side they cross.
	crossArmRatio = 0.33
	// mathBarThickness is the thickness of a minus sign's bar, as a fraction of
	// the shape's height.
	mathBarThickness = 0.23
	// flowWaveResolution is how many segments the document wave is drawn with.
	flowWaveResolution = 24
)

func diamondPoints(x, y, w, h float64) []gopdf.Point {
	return []gopdf.Point{
		{X: x + w/2, Y: y}, {X: x + w, Y: y + h/2},
		{X: x + w/2, Y: y + h}, {X: x, Y: y + h/2},
	}
}

func flowParallelogramPoints(x, y, w, h float64) []gopdf.Point {
	off := w * shapeOffsetRatio
	return []gopdf.Point{
		{X: x + off, Y: y}, {X: x + w, Y: y},
		{X: x + w - off, Y: y + h}, {X: x, Y: y + h},
	}
}

// offpageConnectorPoints is a rectangle with a point at the bottom.
func offpageConnectorPoints(x, y, w, h float64) []gopdf.Point {
	notch := h * flowOffpageNotch
	return []gopdf.Point{
		{X: x, Y: y}, {X: x + w, Y: y},
		{X: x + w, Y: y + h - notch}, {X: x + w/2, Y: y + h},
		{X: x, Y: y + h - notch},
	}
}

// clippedCornerPoints is a rectangle with its top-left corner cut off, which is
// the punched card of a flowchart.
func clippedCornerPoints(x, y, w, h float64) []gopdf.Point {
	clip := math.Min(w, h) * flowClipRatio
	return []gopdf.Point{
		{X: x + clip, Y: y}, {X: x + w, Y: y},
		{X: x + w, Y: y + h}, {X: x, Y: y + h},
		{X: x, Y: y + clip},
	}
}

// flowDocumentPoints is a rectangle whose bottom edge is a wave.
func flowDocumentPoints(x, y, w, h float64) []gopdf.Point {
	depth := h * flowWaveDepth
	points := make([]gopdf.Point, 0, flowWaveResolution+3)
	points = append(points,
		gopdf.Point{X: x, Y: y},
		gopdf.Point{X: x + w, Y: y},
		gopdf.Point{X: x + w, Y: y + h - depth},
	)
	// One full sine period across the width: down on the right, up on the left,
	// which is the direction PowerPoint draws it.
	for i := range flowWaveResolution + 1 {
		t := float64(i) / float64(flowWaveResolution)
		points = append(points, gopdf.Point{
			X: x + w*(1-t),
			Y: y + h - depth + depth*math.Sin(2*math.Pi*t),
		})
	}
	return points
}

// flowDelayPoints is a rectangle with a semicircular right end.
func flowDelayPoints(x, y, w, h float64) []gopdf.Point {
	radius := h / 2
	flat := math.Max(0, w-radius)
	points := []gopdf.Point{{X: x, Y: y}, {X: x + flat, Y: y}}
	for i := range flowWaveResolution + 1 {
		angle := -math.Pi/2 + math.Pi*float64(i)/float64(flowWaveResolution)
		points = append(points, gopdf.Point{
			X: x + flat + radius*math.Cos(angle),
			Y: y + radius + radius*math.Sin(angle),
		})
	}
	return append(points, gopdf.Point{X: x, Y: y + h})
}

// cylinderPoints is the silhouette of a can standing on its end: the outline of
// the body with the bottom cap curving away from the viewer.
func cylinderPoints(x, y, w, h float64) []gopdf.Point {
	capH := h * cylinderCapRatio
	points := make([]gopdf.Point, 0, 2*flowWaveResolution+4)
	// Top cap, left to right over the top.
	for i := range flowWaveResolution + 1 {
		angle := math.Pi - math.Pi*float64(i)/float64(flowWaveResolution)
		points = append(points, gopdf.Point{
			X: x + w/2 + (w/2)*math.Cos(angle),
			Y: y + capH - capH*math.Sin(angle),
		})
	}
	// Bottom cap, right to left under the bottom.
	for i := range flowWaveResolution + 1 {
		angle := math.Pi * float64(i) / float64(flowWaveResolution)
		points = append(points, gopdf.Point{
			X: x + w/2 + (w/2)*math.Cos(angle),
			Y: y + h - capH + capH*math.Sin(angle),
		})
	}
	return points
}

// flowStoredDataPoints is a rectangle whose left edge curves in, which is the
// stored-data and online-storage symbol.
func flowStoredDataPoints(x, y, w, h float64) []gopdf.Point {
	inset := w * flowClipRatio
	points := []gopdf.Point{{X: x + inset, Y: y}, {X: x + w, Y: y}}
	// Right edge curves the same way, so the shape reads as a rolled sheet.
	for i := range flowWaveResolution + 1 {
		t := float64(i) / float64(flowWaveResolution)
		points = append(points, gopdf.Point{
			X: x + w - inset + inset*math.Sin(math.Pi*t),
			Y: y + h*t,
		})
	}
	points = append(points, gopdf.Point{X: x + inset, Y: y + h})
	for i := range flowWaveResolution + 1 {
		t := 1 - float64(i)/float64(flowWaveResolution)
		points = append(points, gopdf.Point{
			X: x + inset*math.Sin(math.Pi*t),
			Y: y + h*t,
		})
	}
	return points
}

// crossPoints is the plus sign: a rectangle with a bite out of each corner.
func crossPoints(x, y, w, h float64) []gopdf.Point {
	armW := w * crossArmRatio
	armH := h * crossArmRatio
	left, right := x+(w-armW)/2, x+(w+armW)/2
	top, bottom := y+(h-armH)/2, y+(h+armH)/2
	return []gopdf.Point{
		{X: left, Y: y}, {X: right, Y: y},
		{X: right, Y: top}, {X: x + w, Y: top},
		{X: x + w, Y: bottom}, {X: right, Y: bottom},
		{X: right, Y: y + h}, {X: left, Y: y + h},
		{X: left, Y: bottom}, {X: x, Y: bottom},
		{X: x, Y: top}, {X: left, Y: top},
	}
}
