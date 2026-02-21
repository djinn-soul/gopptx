package mermaid

import (
	"strconv"
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

func parseQuadrant(code string) *QuadrantDiagram {
	lines := ParseLines(code)
	quadrant := &QuadrantDiagram{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		if strings.HasPrefix(lower, "quadrantchart") {
			continue
		}

		if strings.HasPrefix(lower, "title") {
			quadrant.Title = strings.TrimSpace(trimmed[5:])
			continue
		}

		if strings.HasPrefix(lower, "x-axis") {
			quadrant.XAxis = strings.TrimSpace(trimmed[6:])
			continue
		}

		if strings.HasPrefix(lower, "y-axis") {
			quadrant.YAxis = strings.TrimSpace(trimmed[6:])
			continue
		}

		if strings.HasPrefix(lower, "quadrant-") {
			idxStr := trimmed[9:10]
			idx, err := strconv.Atoi(idxStr)
			if err == nil && idx >= 1 && idx <= 4 {
				quadrant.Quadrants[idx-1] = strings.TrimSpace(trimmed[10:])
			}
			continue
		}

		if strings.Contains(trimmed, ":") && strings.Contains(trimmed, "[") && strings.Contains(trimmed, "]") {
			parts := strings.Split(trimmed, ":")
			label := strings.TrimSpace(parts[0])
			coords := strings.TrimSpace(parts[1])
			coords = strings.Trim(coords, "[]")
			coordParts := strings.Split(coords, ",")
			if len(coordParts) == 2 {
				x, errX := strconv.ParseFloat(strings.TrimSpace(coordParts[0]), 64)
				y, errY := strconv.ParseFloat(strings.TrimSpace(coordParts[1]), 64)
				if errX == nil && errY == nil {
					quadrant.Points = append(quadrant.Points, QuadrantPoint{
						Label: label,
						X:     x,
						Y:     y,
					})
				}
			}
		}
	}

	return quadrant
}

func generateQuadrantElements(quadrant *QuadrantDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape

	// Layout parameters
	startX := styling.Inches(1)
	startY := styling.Inches(1.5)
	chartSize := styling.Inches(5)
	quadSize := chartSize / 2

	// Title
	if quadrant.Title != "" {
		titleShape := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			startX,
			startY-styling.Inches(0.7),
			chartSize,
			styling.Inches(0.5),
		).WithText(quadrant.Title).
			WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight))
		shapesList = append(shapesList, titleShape)
	}

	// Quadrants
	// Use theme colors for quadrants
	quadColors := []string{theme.PrimaryFill, theme.SecondaryFill, theme.Background, theme.PrimaryFill}

	quadPositions := []struct{ x, y styling.Length }{
		{startX + quadSize, startY},            // Q1: Top-Right
		{startX, startY},                       // Q2: Top-Left
		{startX, startY + quadSize},            // Q3: Bottom-Left
		{startX + quadSize, startY + quadSize}, // Q4: Bottom-Right
	}

	for i, pos := range quadPositions {
		quadShape := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			pos.x,
			pos.y,
			quadSize,
			quadSize,
		).WithFill(shapes.NewShapeFill(quadColors[i])).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight))
		shapesList = append(shapesList, quadShape)

		if quadrant.Quadrants[i] != "" {
			labelShape := shapes.NewShape(
				shapes.ShapeTypeRectangle,
				pos.x,
				pos.y+quadSize-styling.Inches(0.4),
				quadSize,
				styling.Inches(0.3),
			).WithText(quadrant.Quadrants[i]).
				WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
				WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight))
			shapesList = append(shapesList, labelShape)
		}
	}

	// Axes
	xAxisLine := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		startX,
		startY+quadSize,
		chartSize,
		styling.Emu(25400),
	).WithFill(shapes.NewShapeFill(theme.PrimaryStroke))
	shapesList = append(shapesList, xAxisLine)

	yAxisLine := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		startX+quadSize,
		startY,
		styling.Emu(25400),
		chartSize,
	).WithFill(shapes.NewShapeFill(theme.PrimaryStroke))
	shapesList = append(shapesList, yAxisLine)

	// Axis labels
	if quadrant.XAxis != "" {
		xLabel := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			startX,
			startY+chartSize+styling.Inches(0.1),
			chartSize,
			styling.Inches(0.3),
		).WithText(quadrant.XAxis).
			WithFill(shapes.NewShapeFill(theme.Background)).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight))
		shapesList = append(shapesList, xLabel)
	}

	if quadrant.YAxis != "" {
		yLabel := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			startX-styling.Inches(1.4),
			startY-styling.Inches(0.4),
			styling.Inches(1.3),
			styling.Inches(0.3),
		).WithText(quadrant.YAxis).
			WithFill(shapes.NewShapeFill(theme.Background)).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
			WithAutoFit(shapes.TextAutoFitNormal)
		shapesList = append(shapesList, yLabel)
	}

	// Points
	for _, p := range quadrant.Points {
		// Map 0-1 coordinates to chart area
		// Mermaid coordinates: 0,0 is bottom-left, 1,1 is top-right
		px := startX + styling.Length(p.X*float64(chartSize))
		py := startY + styling.Length((1.0-p.Y)*float64(chartSize))

		pointSize := styling.Inches(0.1)
		point := shapes.NewShape(
			shapes.ShapeTypeEllipse,
			px-pointSize/2,
			py-pointSize/2,
			pointSize,
			pointSize,
		).WithFill(shapes.NewShapeFill(theme.PrimaryStroke)).
			WithLine(shapes.NewShapeLine(theme.Background, theme.LineWeight))
		shapesList = append(shapesList, point)

		label := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			px+pointSize,
			py-styling.Inches(0.15),
			styling.Inches(1.5),
			styling.Inches(0.3),
		).WithText(p.Label).
			WithFill(shapes.NewShapeFill(theme.Background)).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight))
		shapesList = append(shapesList, label)
	}

	return DiagramElements{
		Shapes:  shapesList,
		Grouped: true,
		Bounds: &DiagramBounds{
			X:  startX - styling.Inches(0.5),
			Y:  startY - styling.Inches(0.7),
			CX: chartSize + styling.Inches(2),
			CY: chartSize + styling.Inches(1.5),
		},
	}
}
