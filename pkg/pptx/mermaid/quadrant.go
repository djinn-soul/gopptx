package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// QuadrantPoint represents a point in a quadrant chart.
type QuadrantPoint struct {
	Label string
	X     float64
	Y     float64
}

// QuadrantDiagram represents the parsed structure of a Mermaid quadrant chart.
type QuadrantDiagram struct {
	Title     string
	XAxis     string
	YAxis     string
	Quadrants [4]string
	Points    []QuadrantPoint
}

// renderQuadrant parses and renders a Mermaid quadrant chart into PowerPoint elements.
func renderQuadrant(code string, theme Theme) DiagramElements {
	quadrant := parseQuadrant(code)
	return generateQuadrantElements(quadrant, theme)
}

func generateQuadrantElements(quadrant *QuadrantDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape

	layout := quadrantLayout{
		startX:    styling.Inches(1),
		startY:    styling.Inches(1.5),
		chartSize: styling.Inches(5),
	}
	quadSize := layout.chartSize / 2

	if quadrant.Title != "" {
		shapesList = append(shapesList, quadrantTitleShape(quadrant.Title, layout, theme))
	}

	shapesList = append(shapesList, quadrantAreaShapes(quadrant, layout, quadSize, theme)...)
	shapesList = append(shapesList, quadrantAxisShapes(quadrant, layout, quadSize, theme)...)
	shapesList = append(shapesList, quadrantPointShapes(quadrant, layout, theme)...)

	return DiagramElements{
		Shapes:  shapesList,
		Grouped: true,
		Bounds: &DiagramBounds{
			X:  layout.startX - styling.Inches(0.5),
			Y:  layout.startY - styling.Inches(0.7),
			CX: layout.chartSize + styling.Inches(2),
			CY: layout.chartSize + styling.Inches(1.5),
		},
	}
}

type quadrantLayout struct {
	startX    styling.Length
	startY    styling.Length
	chartSize styling.Length
}

func quadrantTitleShape(title string, layout quadrantLayout, theme Theme) shapes.Shape {
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		layout.startX,
		layout.startY-styling.Inches(0.7),
		layout.chartSize,
		styling.Inches(0.5),
	).WithText(title).
		WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight))
}

func quadrantAreaShapes(
	quadrant *QuadrantDiagram,
	layout quadrantLayout,
	quadSize styling.Length,
	theme Theme,
) []shapes.Shape {
	out := make([]shapes.Shape, 0, 8)
	quadColors := []string{theme.PrimaryFill, theme.SecondaryFill, theme.Background, theme.PrimaryFill}
	quadPositions := []struct{ x, y styling.Length }{
		{layout.startX + quadSize, layout.startY},
		{layout.startX, layout.startY},
		{layout.startX, layout.startY + quadSize},
		{layout.startX + quadSize, layout.startY + quadSize},
	}

	for i, pos := range quadPositions {
		out = append(out, shapes.NewShape(
			shapes.ShapeTypeRectangle,
			pos.x,
			pos.y,
			quadSize,
			quadSize,
		).WithFill(shapes.NewShapeFill(quadColors[i])).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)))
		if quadrant.Quadrants[i] != "" {
			out = append(out, quadrantLabelShape(pos.x, pos.y, quadSize, quadrant.Quadrants[i], theme))
		}
	}
	return out
}

func quadrantLabelShape(
	x styling.Length,
	y styling.Length,
	quadSize styling.Length,
	label string,
	theme Theme,
) shapes.Shape {
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		x+styling.Inches(0.08),
		y+styling.Inches(0.08),
		quadSize-styling.Inches(0.16),
		styling.Inches(0.30),
	).WithText(label).
		WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.Background, styling.Emu(0))).
		WithAutoFit(shapes.TextAutoFitNormal)
}

func quadrantAxisShapes(
	quadrant *QuadrantDiagram,
	layout quadrantLayout,
	quadSize styling.Length,
	theme Theme,
) []shapes.Shape {
	out := make([]shapes.Shape, 0, 4)
	out = append(out, shapes.NewShape(
		shapes.ShapeTypeRectangle,
		layout.startX,
		layout.startY+quadSize,
		layout.chartSize,
		styling.Emu(25400),
	).WithFill(shapes.NewShapeFill(theme.PrimaryStroke)))

	out = append(out, shapes.NewShape(
		shapes.ShapeTypeRectangle,
		layout.startX+quadSize,
		layout.startY,
		styling.Emu(25400),
		layout.chartSize,
	).WithFill(shapes.NewShapeFill(theme.PrimaryStroke)))

	if quadrant.XAxis != "" {
		out = append(out, quadrantXAxisLabelShape(quadrant.XAxis, layout, theme))
	}
	if quadrant.YAxis != "" {
		out = append(out, quadrantYAxisLabelShape(quadrant.YAxis, layout, theme))
	}
	return out
}

func quadrantXAxisLabelShape(label string, layout quadrantLayout, theme Theme) shapes.Shape {
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		layout.startX,
		layout.startY+layout.chartSize+styling.Inches(0.1),
		layout.chartSize,
		styling.Inches(0.3),
	).WithText(label).
		WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight))
}

func quadrantYAxisLabelShape(label string, layout quadrantLayout, theme Theme) shapes.Shape {
	label = strings.TrimSpace(strings.ReplaceAll(label, "-->", " -> "))
	centerY := layout.startY + (layout.chartSize / 2)
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		layout.startX-styling.Inches(0.8),
		centerY-styling.Inches(1.6),
		styling.Inches(0.32),
		styling.Inches(3.2),
	).WithText(label).
		WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.Background, styling.Emu(0))).
		WithRotation(-90).
		WithTextMargins(styling.Inches(0.02), styling.Inches(0.04), styling.Inches(0.02), styling.Inches(0.04)).
		WithAutoFit(shapes.TextAutoFitNormal)
}

func quadrantPointShapes(quadrant *QuadrantDiagram, layout quadrantLayout, theme Theme) []shapes.Shape {
	out := make([]shapes.Shape, 0, len(quadrant.Points)*2)
	for _, p := range quadrant.Points {
		px := layout.startX + styling.Length(p.X*float64(layout.chartSize))
		py := layout.startY + styling.Length((1.0-p.Y)*float64(layout.chartSize))
		out = append(out, quadrantPointShape(px, py, theme))
		out = append(out, quadrantPointLabelShape(px, py, p.Label, theme))
	}
	return out
}

func quadrantPointShape(px styling.Length, py styling.Length, theme Theme) shapes.Shape {
	pointSize := styling.Inches(0.1)
	return shapes.NewShape(
		shapes.ShapeTypeEllipse,
		px-pointSize/2,
		py-pointSize/2,
		pointSize,
		pointSize,
	).WithFill(shapes.NewShapeFill(theme.PrimaryStroke)).
		WithLine(shapes.NewShapeLine(theme.Background, theme.LineWeight))
}

func quadrantPointLabelShape(px styling.Length, py styling.Length, label string, theme Theme) shapes.Shape {
	labelX := px + styling.Inches(0.14)
	if px > styling.Inches(4.7) {
		labelX = px - styling.Inches(1.35)
	}
	labelY := py - styling.Inches(0.34)
	if py < styling.Inches(2.1) {
		labelY = py + styling.Inches(0.10)
	}
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		labelX,
		labelY,
		styling.Inches(1.2),
		styling.Inches(0.28),
	).WithText(label).
		WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.Background, styling.Emu(0))).
		WithAutoFit(shapes.TextAutoFitNormal)
}
