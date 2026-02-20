package mermaid

import (
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// MindmapNode represents a node in a mindmap.
type MindmapNode struct {
	Label    string
	Level    int
	Children []*MindmapNode
	Shape    NodeShape
}

// renderMindmap parses and renders a Mermaid mindmap into PowerPoint elements.
func renderMindmap(code string, theme Theme) DiagramElements {
	root := parseMindmap(code)
	if root == nil {
		return createPlaceholder("mindmap (no data)", theme)
	}
	return generateMindmapElements(root, theme)
}

func parseMindmap(code string) *MindmapNode {
	lines := strings.Split(code, "\n")
	var root *MindmapNode
	var stack []*MindmapNode

	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " 	")
		if trimmed == "" || strings.HasPrefix(strings.TrimSpace(trimmed), "%%") {
			continue
		}

		cleanLine := strings.TrimSpace(trimmed)
		if strings.ToLower(cleanLine) == "mindmap" {
			continue
		}

		// Determine level by counting leading spaces
		indent := 0
		for _, char := range line {
			if char == ' ' {
				indent++
			} else if char == '	' {
				indent += 4
			} else {
				break
			}
		}

		// Parse node label and shape
		_, label, shape := ParseNodeDef(cleanLine)
		
		// Handle icons or other mindmap specific syntax (simplified)
		if strings.Contains(label, "::icon") {
			label = strings.Split(label, "::icon")[0]
		}

		node := &MindmapNode{
			Label: strings.TrimSpace(label),
			Level: indent,
			Shape: shape,
		}

		if root == nil {
			root = node
			stack = []*MindmapNode{node}
			continue
		}

		// Find parent based on indent
		for len(stack) > 0 && stack[len(stack)-1].Level >= indent {
			stack = stack[:len(stack)-1]
		}

		if len(stack) > 0 {
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, node)
		}
		stack = append(stack, node)
	}

	return root
}

func generateMindmapElements(root *MindmapNode, theme Theme) DiagramElements {
	var shapesList []shapes.Shape
	var connectors []shapes.Connector

	// Start at origin for initial layout
	initialX := styling.Length(0)
	initialY := styling.Length(0)

	// Recursive function to layout nodes
	var layoutNode func(node *MindmapNode, x, y styling.Length, angleStart, angleEnd float64, radius styling.Length)
	layoutNode = func(node *MindmapNode, x, y styling.Length, angleStart, angleEnd float64, radius styling.Length) {
		// Create shape for current node
		nodeWidth := styling.Inches(1.4)
		nodeHeight := styling.Inches(0.6)

		shapeType := shapes.ShapeTypeRectangle
		switch node.Shape {
		case NodeShapeCircle:
			shapeType = shapes.ShapeTypeEllipse
		case NodeShapeRoundedRect:
			shapeType = shapes.ShapeTypeRoundedRectangle
		case NodeShapeStadium:
			shapeType = shapes.ShapeTypeFlowChartConnector
		}

		fillColor := theme.PrimaryFill
		strokeColor := theme.PrimaryStroke
		if node.Level == 0 {
			fillColor = theme.PrimaryStroke
			strokeColor = theme.PrimaryFill
		} else if len(node.Children) > 0 {
			fillColor = theme.SecondaryFill
			strokeColor = theme.SecondaryStroke
		}

		nodeShape := shapes.NewShape(
			shapeType,
			x-nodeWidth/2,
			y-nodeHeight/2,
			nodeWidth,
			nodeHeight,
		).WithFill(shapes.NewShapeFill(fillColor)).
			WithLine(shapes.NewShapeLine(strokeColor, theme.LineWeight)).
			WithText(node.Label).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))

		shapesList = append(shapesList, nodeShape)

		if len(node.Children) == 0 {
			return
		}

		// Distribute children in the given angle arc
		angleStep := (angleEnd - angleStart) / float64(len(node.Children))
		currentAngle := angleStart + angleStep/2

		nextRadius := radius + styling.Inches(2.0)

		for _, child := range node.Children {
			childX := x + styling.Length(float64(radius)*math.Cos(currentAngle))
			childY := y + styling.Length(float64(radius)*math.Sin(currentAngle))

			// Add connector
			connector := shapes.NewConnector(
				shapes.ConnectorTypeStraight,
				x, y, childX, childY,
			).WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight))
			connectors = append(connectors, connector)

			// Recurse
			childAngleStart := currentAngle - angleStep/2
			childAngleEnd := currentAngle + angleStep/2
			layoutNode(child, childX, childY, childAngleStart, childAngleEnd, nextRadius)

			currentAngle += angleStep
		}
	}

	layoutNode(root, initialX, initialY, 0, 2*math.Pi, styling.Inches(2.5))

	// Calculate bounds
	if len(shapesList) == 0 {
		return DiagramElements{}
	}

	minX := shapesList[0].X
	minY := shapesList[0].Y
	maxX := shapesList[0].X + shapesList[0].CX
	maxY := shapesList[0].Y + shapesList[0].CY

	for _, s := range shapesList {
		if s.X < minX {
			minX = s.X
		}
		if s.Y < minY {
			minY = s.Y
		}
		if s.X+s.CX > maxX {
			maxX = s.X + s.CX
		}
		if s.Y+s.CY > maxY {
			maxY = s.Y + s.CY
		}
	}

	// Also consider connectors for bounds
	for _, c := range connectors {
		if c.StartX < minX {
			minX = c.StartX
		}
		if c.EndX < minX {
			minX = c.EndX
		}
		if c.StartY < minY {
			minY = c.StartY
		}
		if c.EndY < minY {
			minY = c.EndY
		}
		if c.StartX > maxX {
			maxX = c.StartX
		}
		if c.EndX > maxX {
			maxX = c.EndX
		}
		if c.StartY > maxY {
			maxY = c.StartY
		}
		if c.EndY > maxY {
			maxY = c.EndY
		}
	}

	// Apply offset to ensure all coordinates are positive with a margin
	margin := styling.Inches(1)
	offsetX := margin - minX
	offsetY := margin - minY

	for i := range shapesList {
		shapesList[i].X += offsetX
		shapesList[i].Y += offsetY
	}

	for i := range connectors {
		connectors[i].StartX += offsetX
		connectors[i].StartY += offsetY
		connectors[i].EndX += offsetX
		connectors[i].EndY += offsetY
	}

	return DiagramElements{
		Shapes:     shapesList,
		Connectors: connectors,
		Grouped:    true,
		Bounds: &DiagramBounds{
			X:  minX + offsetX,
			Y:  minY + offsetY,
			CX: maxX - minX,
			CY: maxY - minY,
		},
	}
}
