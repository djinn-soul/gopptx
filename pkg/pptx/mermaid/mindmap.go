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
		_, cleanLine, indent, ok := parseMindmapLine(line)
		if !ok {
			continue
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

		root, stack = appendMindmapNode(root, stack, node)
	}

	return root
}

func parseMindmapLine(line string) (string, string, int, bool) {
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == "" || strings.HasPrefix(strings.TrimSpace(trimmed), "%%") {
		return "", "", 0, false
	}

	cleanLine := strings.TrimSpace(trimmed)
	if strings.EqualFold(cleanLine, "mindmap") {
		return "", "", 0, false
	}
	return trimmed, cleanLine, leadingIndent(line), true
}

func leadingIndent(line string) int {
	indent := 0
	for _, char := range line {
		switch char {
		case ' ':
			indent++
		case '\t':
			indent += 4
		default:
			return indent
		}
	}
	return indent
}

func appendMindmapNode(root *MindmapNode, stack []*MindmapNode, node *MindmapNode) (*MindmapNode, []*MindmapNode) {
	if root == nil {
		return node, []*MindmapNode{node}
	}

	stack = popMindmapParents(stack, node.Level)
	if len(stack) > 0 {
		parent := stack[len(stack)-1]
		parent.Children = append(parent.Children, node)
	}
	stack = append(stack, node)
	return root, stack
}

func popMindmapParents(stack []*MindmapNode, indent int) []*MindmapNode {
	for len(stack) > 0 && stack[len(stack)-1].Level >= indent {
		stack = stack[:len(stack)-1]
	}
	return stack
}

func generateMindmapElements(root *MindmapNode, theme Theme) DiagramElements {
	renderer := newMindmapRenderer(theme)
	renderer.layoutNode(root, 0, 0, 0, 2*math.Pi, renderer.rootRadius)
	if len(renderer.shapes) == 0 {
		return DiagramElements{}
	}
	renderer.applyOffset()
	return renderer.diagramElements()
}

type mindmapRenderer struct {
	theme         Theme
	shapes        []shapes.Shape
	connectors    []shapes.Connector
	bounds        mindmapBounds
	nodeWidth     styling.Length
	nodeHeight    styling.Length
	radiusStep    styling.Length
	rootRadius    styling.Length
	textMarginX   styling.Length
	textMarginY   styling.Length
	layoutMargin  styling.Length
	initialBounds bool
}

type mindmapBounds struct {
	minX styling.Length
	minY styling.Length
	maxX styling.Length
	maxY styling.Length
}

func newMindmapRenderer(theme Theme) *mindmapRenderer {
	return &mindmapRenderer{
		theme:        theme,
		nodeWidth:    styling.Inches(1.4),
		nodeHeight:   styling.Inches(0.6),
		radiusStep:   styling.Inches(2.0),
		rootRadius:   styling.Inches(2.5),
		textMarginX:  styling.Inches(0.1),
		textMarginY:  styling.Inches(0.05),
		layoutMargin: styling.Inches(1),
	}
}

func (r *mindmapRenderer) layoutNode(
	node *MindmapNode,
	x, y styling.Length,
	angleStart, angleEnd float64,
	radius styling.Length,
) {
	nodeShape := r.nodeShape(node, x, y)
	r.shapes = append(r.shapes, nodeShape)
	r.includeShape(nodeShape)

	if len(node.Children) == 0 {
		return
	}
	angleStep := (angleEnd - angleStart) / float64(len(node.Children))
	currentAngle := angleStart + angleStep/2
	nextRadius := radius + r.radiusStep
	for _, child := range node.Children {
		childX := x + styling.Length(float64(radius)*math.Cos(currentAngle))
		childY := y + styling.Length(float64(radius)*math.Sin(currentAngle))
		connector := shapes.NewConnector(
			shapes.ConnectorTypeStraight,
			x,
			y,
			childX,
			childY,
		).WithLine(shapes.NewShapeLine(r.theme.SecondaryStroke, r.theme.LineWeight))
		r.connectors = append(r.connectors, connector)
		r.includeConnector(connector)
		childAngleStart := currentAngle - angleStep/2
		childAngleEnd := currentAngle + angleStep/2
		r.layoutNode(child, childX, childY, childAngleStart, childAngleEnd, nextRadius)
		currentAngle += angleStep
	}
}

func (r *mindmapRenderer) nodeShape(node *MindmapNode, x, y styling.Length) shapes.Shape {
	fillColor, strokeColor := r.nodeColors(node)
	return shapes.NewShape(
		mindmapShapeType(node.Shape),
		x-r.nodeWidth/2,
		y-r.nodeHeight/2,
		r.nodeWidth,
		r.nodeHeight,
	).WithFill(shapes.NewShapeFill(fillColor)).
		WithLine(shapes.NewShapeLine(strokeColor, r.theme.LineWeight)).
		WithText(node.Label).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(r.textMarginX, r.textMarginY, r.textMarginX, r.textMarginY)
}

func (r *mindmapRenderer) nodeColors(node *MindmapNode) (string, string) {
	if node.Level == 0 {
		return r.theme.PrimaryStroke, r.theme.PrimaryFill
	}
	if len(node.Children) > 0 {
		return r.theme.SecondaryFill, r.theme.SecondaryStroke
	}
	return r.theme.PrimaryFill, r.theme.PrimaryStroke
}

func mindmapShapeType(nodeShape NodeShape) string {
	switch nodeShape {
	case NodeShapeCircle:
		return shapes.ShapeTypeEllipse
	case NodeShapeRoundedRect:
		return shapes.ShapeTypeRoundedRectangle
	case NodeShapeStadium:
		return shapes.ShapeTypeFlowChartConnector
	case NodeShapeDiamond:
		return shapes.ShapeTypeDiamond
	case NodeShapeHexagon:
		return shapes.ShapeTypeHexagon
	default:
		return shapes.ShapeTypeRectangle
	}
}

func (r *mindmapRenderer) includeShape(shape shapes.Shape) {
	r.include(shape.X, shape.Y, shape.CX, shape.CY)
}

func (r *mindmapRenderer) includeConnector(connector shapes.Connector) {
	r.includePoint(connector.StartX, connector.StartY)
	r.includePoint(connector.EndX, connector.EndY)
}

func (r *mindmapRenderer) includePoint(x, y styling.Length) {
	r.include(x, y, 0, 0)
}

func (r *mindmapRenderer) include(x, y, cx, cy styling.Length) {
	if !r.initialBounds {
		r.bounds = mindmapBounds{
			minX: x,
			minY: y,
			maxX: x + cx,
			maxY: y + cy,
		}
		r.initialBounds = true
		return
	}
	if x < r.bounds.minX {
		r.bounds.minX = x
	}
	if y < r.bounds.minY {
		r.bounds.minY = y
	}
	if x+cx > r.bounds.maxX {
		r.bounds.maxX = x + cx
	}
	if y+cy > r.bounds.maxY {
		r.bounds.maxY = y + cy
	}
}

func (r *mindmapRenderer) applyOffset() {
	offsetX := r.layoutMargin - r.bounds.minX
	offsetY := r.layoutMargin - r.bounds.minY
	for i := range r.shapes {
		r.shapes[i].X += offsetX
		r.shapes[i].Y += offsetY
	}
	for i := range r.connectors {
		r.connectors[i].StartX += offsetX
		r.connectors[i].StartY += offsetY
		r.connectors[i].EndX += offsetX
		r.connectors[i].EndY += offsetY
	}
	r.bounds.minX += offsetX
	r.bounds.minY += offsetY
	r.bounds.maxX += offsetX
	r.bounds.maxY += offsetY
}

func (r *mindmapRenderer) diagramElements() DiagramElements {
	return DiagramElements{
		Shapes:     r.shapes,
		Connectors: r.connectors,
		Grouped:    true,
		Bounds: &DiagramBounds{
			X:  r.bounds.minX,
			Y:  r.bounds.minY,
			CX: r.bounds.maxX - r.bounds.minX,
			CY: r.bounds.maxY - r.bounds.minY,
		},
	}
}
