package mermaid

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// PieData represents a single slice in a pie chart.
type PieData struct {
	Label string
	Value float64
}

// PieDiagram represents the parsed structure of a Mermaid pie chart.
type PieDiagram struct {
	Title string
	Data  []PieData
}

// renderPie parses and renders a Mermaid pie chart into PowerPoint elements.
func renderPie(code string, theme Theme) DiagramElements {
	pie := parsePie(code)
	return generatePieElements(pie, theme)
}

func parsePie(code string) *PieDiagram {
	lines := ParseLines(code)
	pie := &PieDiagram{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(trimmed), "pie") {
			if strings.Contains(strings.ToLower(trimmed), "title") {
				pie.Title = strings.TrimSpace(trimmed[strings.Index(strings.ToLower(trimmed), "title")+5:])
			}
			continue
		}

		if strings.Contains(trimmed, ":") {
			parts := strings.Split(trimmed, ":")
			if len(parts) == 2 {
				label := strings.TrimSpace(parts[0])
				// Remove quotes if present
				label = strings.Trim(label, "\"")
				valueStr := strings.TrimSpace(parts[1])
				value, err := strconv.ParseFloat(valueStr, 64)
				if err == nil {
					pie.Data = append(pie.Data, PieData{Label: label, Value: value})
				}
			}
		} else if strings.HasPrefix(strings.ToLower(trimmed), "title") {
			pie.Title = strings.TrimSpace(trimmed[5:])
		}
	}

	return pie
}

func generatePieElements(pie *PieDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape

	if len(pie.Data) == 0 {
		return createPlaceholder("pie (no data)", theme)
	}

	total := 0.0
	for _, d := range pie.Data {
		total += d.Value
	}
	if total <= 0 {
		return createPlaceholder("pie (invalid totals)", theme)
	}

	// Layout parameters
	centerX := styling.Inches(4)
	centerY := styling.Inches(3.5)
	radius := styling.Inches(2.0)

	// Title
	if pie.Title != "" {
		titleShape := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			centerX-styling.Inches(3),
			centerY-radius-styling.Inches(1.0),
			styling.Inches(6),
			styling.Inches(0.6),
		).WithText(pie.Title).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
		shapesList = append(shapesList, titleShape)
	}

	// Pie slices
	currentAngle := 0.0
	palette := []string{
		theme.PrimaryFill,
		theme.SecondaryFill,
		"4285F4", "EA4335", "FBBC05", "34A853", "FF6D01", "46BDC6", "7BAAF7", "F07B72",
	}

	for i, d := range pie.Data {
		sliceAngle := (d.Value / total) * 360.0
		endAngle := currentAngle + sliceAngle

		color := palette[i%len(palette)]

		slice := shapes.NewPieSlice(
			centerX-radius,
			centerY-radius,
			radius*2,
			radius*2,
			currentAngle,
			endAngle,
		).WithFill(shapes.NewShapeFill(color)).
			WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight))

		shapesList = append(shapesList, slice)
		currentAngle = endAngle
	}

	// Legend
	legendX := centerX + radius + styling.Inches(0.5)
	legendY := centerY - radius
	itemHeight := styling.Inches(0.35)

	for i, d := range pie.Data {
		percentage := (d.Value / total) * 100
		color := palette[i%len(palette)]

		// Color box
		boxSize := styling.Inches(0.15)
		colorBox := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			legendX,
			legendY+styling.Length(i)*itemHeight+(itemHeight-boxSize)/2,
			boxSize,
			boxSize,
		).WithFill(shapes.NewShapeFill(color)).
			WithLine(shapes.NewShapeLine(theme.PrimaryStroke, styling.Emu(12700)))
		shapesList = append(shapesList, colorBox)

		// Label
		label := fmt.Sprintf("%s: %.1f%%", d.Label, percentage)
		labelShape := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			legendX+boxSize+styling.Inches(0.1),
			legendY+styling.Length(i)*itemHeight,
			styling.Inches(3.0),
			itemHeight,
		).WithText(label).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.05), styling.Inches(0.02), styling.Inches(0.05), styling.Inches(0.02))
		labelShape.Fill = nil
		labelShape.Line = nil

		shapesList = append(shapesList, labelShape)
	}

	// Calculate bounds
	minX := centerX - radius
	minY := centerY - radius
	if pie.Title != "" {
		minY = centerY - radius - styling.Inches(0.8)
	}
	maxX := legendX + styling.Inches(3)
	maxY := centerY + radius
	if styling.Length(len(pie.Data))*itemHeight > radius*2 {
		maxY = legendY + styling.Length(len(pie.Data))*itemHeight
	}

	return DiagramElements{
		Shapes:  shapesList,
		Grouped: true,
		Bounds: &DiagramBounds{
			X:  minX,
			Y:  minY,
			CX: maxX - minX,
			CY: maxY - minY,
		},
	}
}
