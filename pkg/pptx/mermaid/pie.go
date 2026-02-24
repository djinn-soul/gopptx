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

	layout := pieLayout{
		centerX: styling.Inches(4),
		centerY: styling.Inches(3.5),
		radius:  styling.Inches(2.0),
	}
	palette := piePalette(theme)

	if pie.Title != "" {
		shapesList = append(shapesList, pieTitleShape(pie.Title, layout))
	}

	shapesList = append(shapesList, pieSliceShapes(pie, total, layout, theme, palette)...)
	shapesList = append(shapesList, pieLegendShapes(pie, total, layout, theme, palette)...)
	minX, minY, maxX, maxY := pieBounds(pie, layout)

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

type pieLayout struct {
	centerX styling.Length
	centerY styling.Length
	radius  styling.Length
}

func piePalette(theme Theme) []string {
	return []string{
		theme.PrimaryFill,
		theme.SecondaryFill,
		"4285F4", "EA4335", "FBBC05", "34A853", "FF6D01", "46BDC6", "7BAAF7", "F07B72",
	}
}

func pieTitleShape(title string, layout pieLayout) shapes.Shape {
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		layout.centerX-styling.Inches(3),
		layout.centerY-layout.radius-styling.Inches(1.0),
		styling.Inches(6),
		styling.Inches(0.6),
	).WithText(title).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
}

func pieSliceShapes(pie *PieDiagram, total float64, layout pieLayout, theme Theme, palette []string) []shapes.Shape {
	out := make([]shapes.Shape, 0, len(pie.Data))
	currentAngle := 0.0

	for i, d := range pie.Data {
		endAngle := currentAngle + (d.Value/total)*360.0
		color := palette[i%len(palette)]
		out = append(out, shapes.NewPieSlice(
			layout.centerX-layout.radius,
			layout.centerY-layout.radius,
			layout.radius*2,
			layout.radius*2,
			currentAngle,
			endAngle,
		).WithFill(shapes.NewShapeFill(color)).
			WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)))
		currentAngle = endAngle
	}
	return out
}

func pieLegendShapes(pie *PieDiagram, total float64, layout pieLayout, theme Theme, palette []string) []shapes.Shape {
	legendX := layout.centerX + layout.radius + styling.Inches(0.5)
	legendY := layout.centerY - layout.radius
	itemHeight := styling.Inches(0.35)
	boxSize := styling.Inches(0.15)

	out := make([]shapes.Shape, 0, len(pie.Data)*2)
	for i, d := range pie.Data {
		color := palette[i%len(palette)]
		percentage := (d.Value / total) * 100
		out = append(out, pieLegendColorBox(i, legendX, legendY, itemHeight, boxSize, color, theme))
		out = append(out, pieLegendLabel(i, legendX, legendY, itemHeight, boxSize, d.Label, percentage))
	}
	return out
}

func pieLegendColorBox(
	index int,
	legendX styling.Length,
	legendY styling.Length,
	itemHeight styling.Length,
	boxSize styling.Length,
	color string,
	theme Theme,
) shapes.Shape {
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		legendX,
		legendY+styling.Length(index)*itemHeight+(itemHeight-boxSize)/2,
		boxSize,
		boxSize,
	).WithFill(shapes.NewShapeFill(color)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, styling.Emu(12700)))
}

func pieLegendLabel(
	index int,
	legendX styling.Length,
	legendY styling.Length,
	itemHeight styling.Length,
	boxSize styling.Length,
	label string,
	percentage float64,
) shapes.Shape {
	text := fmt.Sprintf("%s: %.1f%%", label, percentage)
	labelShape := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		legendX+boxSize+styling.Inches(0.1),
		legendY+styling.Length(index)*itemHeight,
		styling.Inches(3.0),
		itemHeight,
	).WithText(text).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.05), styling.Inches(0.02), styling.Inches(0.05), styling.Inches(0.02))
	labelShape.Fill = nil
	labelShape.Line = nil
	return labelShape
}

func pieBounds(pie *PieDiagram, layout pieLayout) (styling.Length, styling.Length, styling.Length, styling.Length) {
	minX := layout.centerX - layout.radius
	minY := layout.centerY - layout.radius
	if pie.Title != "" {
		minY = layout.centerY - layout.radius - styling.Inches(0.8)
	}

	legendX := layout.centerX + layout.radius + styling.Inches(0.5)
	itemHeight := styling.Inches(0.35)
	legendY := layout.centerY - layout.radius
	maxX := legendX + styling.Inches(3)
	maxY := layout.centerY + layout.radius
	if styling.Length(len(pie.Data))*itemHeight > layout.radius*2 {
		maxY = legendY + styling.Length(len(pie.Data))*itemHeight
	}
	return minX, minY, maxX, maxY
}
