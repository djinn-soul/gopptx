package mermaid

import (
	"math"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	fitAreaX = 0.5
	fitAreaY = 1.2
	fitAreaW = 9.0
	fitAreaH = 5.9
)

type rawBounds struct {
	minX styling.Length
	minY styling.Length
	maxX styling.Length
	maxY styling.Length
	ok   bool
}

func fitDiagramToSlide(diagram DiagramElements) DiagramElements {
	bounds := collectBounds(diagram)
	if !bounds.ok {
		return diagram
	}

	srcW := bounds.maxX - bounds.minX
	srcH := bounds.maxY - bounds.minY
	if srcW <= 0 || srcH <= 0 {
		return diagram
	}

	targetX := styling.Inches(fitAreaX)
	targetY := styling.Inches(fitAreaY)
	targetW := styling.Inches(fitAreaW)
	targetH := styling.Inches(fitAreaH)

	scaleX := float64(targetW) / float64(srcW)
	scaleY := float64(targetH) / float64(srcH)
	scale := math.Min(scaleX, scaleY)
	if scale > 1 {
		scale = 1
	}
	if scale <= 0 || math.IsNaN(scale) || math.IsInf(scale, 0) {
		return diagram
	}

	scaledW := styling.Emu(int64(math.Round(float64(srcW) * scale)))
	scaledH := styling.Emu(int64(math.Round(float64(srcH) * scale)))
	offsetX := targetX + (targetW-scaledW)/2
	offsetY := targetY + (targetH-scaledH)/2

	transformPoint := func(v, minValue, offset styling.Length) styling.Length {
		return scaleLength(v-minValue, scale) + offset
	}

	for i := range diagram.Shapes {
		diagram.Shapes[i] = transformShape(diagram.Shapes[i], bounds.minX, bounds.minY, offsetX, offsetY, scale)
	}
	for i := range diagram.Connectors {
		diagram.Connectors[i] = transformConnector(
			diagram.Connectors[i],
			bounds.minX,
			bounds.minY,
			offsetX,
			offsetY,
			scale,
		)
	}

	diagram.Bounds = &DiagramBounds{
		X:  offsetX,
		Y:  offsetY,
		CX: transformPoint(bounds.maxX, bounds.minX, offsetX) - offsetX,
		CY: transformPoint(bounds.maxY, bounds.minY, offsetY) - offsetY,
	}
	return diagram
}

func collectBounds(diagram DiagramElements) rawBounds {
	bounds := rawBounds{}
	addRect := func(x, y, cx, cy styling.Length) {
		if cx <= 0 || cy <= 0 {
			return
		}
		minX := x
		minY := y
		maxX := x + cx
		maxY := y + cy
		if !bounds.ok {
			bounds.minX = minX
			bounds.minY = minY
			bounds.maxX = maxX
			bounds.maxY = maxY
			bounds.ok = true
			return
		}
		if minX < bounds.minX {
			bounds.minX = minX
		}
		if minY < bounds.minY {
			bounds.minY = minY
		}
		if maxX > bounds.maxX {
			bounds.maxX = maxX
		}
		if maxY > bounds.maxY {
			bounds.maxY = maxY
		}
	}

	for _, s := range diagram.Shapes {
		addRect(s.X, s.Y, s.CX, s.CY)
	}
	for _, c := range diagram.Connectors {
		minX, maxX := ordered(c.StartX, c.EndX)
		minY, maxY := ordered(c.StartY, c.EndY)
		addRect(minX, minY, maxX-minX, maxY-minY)
	}

	if !bounds.ok && diagram.Bounds != nil {
		addRect(diagram.Bounds.X, diagram.Bounds.Y, diagram.Bounds.CX, diagram.Bounds.CY)
	}
	return bounds
}

func ordered(a, b styling.Length) (styling.Length, styling.Length) {
	if a <= b {
		return a, b
	}
	return b, a
}

func transformShape(
	shape shapes.Shape,
	minX, minY, offsetX, offsetY styling.Length,
	scale float64,
) shapes.Shape {
	shape.X = scaleLength(shape.X-minX, scale) + offsetX
	shape.Y = scaleLength(shape.Y-minY, scale) + offsetY
	shape.CX = scaleLength(shape.CX, scale)
	shape.CY = scaleLength(shape.CY, scale)
	if shape.Line != nil {
		line := *shape.Line
		line.Width = minLineWidth(scaleLength(line.Width, scale))
		shape.Line = &line
	}
	if shape.TextFrame != nil {
		tf := *shape.TextFrame
		tf.MarginLeft = scaleLength(tf.MarginLeft, scale)
		tf.MarginRight = scaleLength(tf.MarginRight, scale)
		tf.MarginTop = scaleLength(tf.MarginTop, scale)
		tf.MarginBottom = scaleLength(tf.MarginBottom, scale)
		shape.TextFrame = &tf
	}
	return shape
}

func transformConnector(
	connector shapes.Connector,
	minX, minY, offsetX, offsetY styling.Length,
	scale float64,
) shapes.Connector {
	connector.StartX = scaleLength(connector.StartX-minX, scale) + offsetX
	connector.StartY = scaleLength(connector.StartY-minY, scale) + offsetY
	connector.EndX = scaleLength(connector.EndX-minX, scale) + offsetX
	connector.EndY = scaleLength(connector.EndY-minY, scale) + offsetY
	connector.Line.Width = minLineWidth(scaleLength(connector.Line.Width, scale))
	return connector
}

func scaleLength(value styling.Length, scale float64) styling.Length {
	return styling.Emu(int64(math.Round(float64(value) * scale)))
}

func minLineWidth(width styling.Length) styling.Length {
	halfPoint := styling.Points(0.5)
	if width < halfPoint {
		return halfPoint
	}
	return width
}
