package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// FlowNode represents a node in a flowchart.
type FlowNode struct {
	ID    string
	Label string
	Shape NodeShape
}

// FlowConnection represents a connection between two nodes in a flowchart.
type FlowConnection struct {
	From       string
	To         string
	Label      string
	ArrowStyle ArrowStyle
}

// Subgraph represents a grouped set of nodes in a flowchart.
type Subgraph struct {
	Name  string
	Nodes []string
}

// FlowchartDiagram represents the parsed structure of a Mermaid flowchart.
type FlowchartDiagram struct {
	Direction   FlowDirection
	Nodes       []FlowNode
	Connections []FlowConnection
	Subgraphs   []Subgraph
}

// renderFlowchart parses and renders a Mermaid flowchart into PowerPoint elements.
func renderFlowchart(code string, theme Theme) DiagramElements {
	flowchart := parseFlowchart(code)
	return generateFlowchartElements(flowchart, theme)
}

func parseFlowchart(code string) *FlowchartDiagram {
	lines := ParseLines(code)
	if len(lines) == 0 {
		return &FlowchartDiagram{Direction: FlowDirectionTB}
	}

	direction := ExtractDirection(lines[0])
	nodes := make(map[string]FlowNode)
	var connections []FlowConnection
	var subgraphs []Subgraph
	var currentSubgraph *Subgraph

	// Skip the header line (graph/flowchart ...)
	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if after, ok := strings.CutPrefix(line, "subgraph"); ok {
			name := strings.TrimSpace(after)
			currentSubgraph = &Subgraph{Name: name, Nodes: []string{}}
			continue
		}

		if line == "end" {
			if currentSubgraph != nil {
				subgraphs = append(subgraphs, *currentSubgraph)
				currentSubgraph = nil
			}
			continue
		}

		if fromPart, arrow, rest, found := SplitConnection(line); found {
			arrowStyle := parseArrowStyle(arrow)

			// Parse from node
			fromID, fromLabel, fromShape := ParseNodeDef(fromPart)
			if _, exists := nodes[fromID]; !exists {
				nodes[fromID] = FlowNode{ID: fromID, Label: fromLabel, Shape: fromShape}
				if currentSubgraph != nil {
					currentSubgraph.Nodes = append(currentSubgraph.Nodes, fromID)
				}
			}

			// Parse to node (may have label on arrow)
			arrowLabel, toPart := ExtractArrowLabel(rest)
			toID, toLabel, toShape := ParseNodeDef(toPart)
			if _, exists := nodes[toID]; !exists {
				nodes[toID] = FlowNode{ID: toID, Label: toLabel, Shape: toShape}
				if currentSubgraph != nil {
					currentSubgraph.Nodes = append(currentSubgraph.Nodes, toID)
				}
			}

			connections = append(connections, FlowConnection{
				From:       fromID,
				To:         toID,
				Label:      arrowLabel,
				ArrowStyle: arrowStyle,
			})
		} else {
			// Standalone node definition
			id, label, shape := ParseNodeDef(line)
			if _, exists := nodes[id]; !exists {
				nodes[id] = FlowNode{ID: id, Label: label, Shape: shape}
				if currentSubgraph != nil {
					currentSubgraph.Nodes = append(currentSubgraph.Nodes, id)
				}
			}
		}
	}

	nodeList := make([]FlowNode, 0, len(nodes))
	for _, node := range nodes {
		nodeList = append(nodeList, node)
	}

	return &FlowchartDiagram{
		Direction:   direction,
		Nodes:       nodeList,
		Connections: connections,
		Subgraphs:   subgraphs,
	}
}

func parseArrowStyle(arrow string) ArrowStyle {
	switch arrow {
	case "==>":
		return ArrowStyleThick
	case "-.->":
		return ArrowStyleDotted
	case "---":
		return ArrowStyleOpen
	default:
		return ArrowStyleArrow
	}
}

//nolint:gocyclo,cyclop // layout/render branches are intentionally explicit for diagram readability
func generateFlowchartElements(flowchart *FlowchartDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape
	var connectors []shapes.Connector
	nodeCount := len(flowchart.Nodes)

	if nodeCount == 0 {
		return DiagramElements{Grouped: true}
	}

	// Layout parameters
	baseNodeWidth := styling.Inches(1.6)
	nodeHeight := styling.Inches(0.6)
	hSpacing := styling.Inches(2.5)
	vSpacing := styling.Inches(1.2)

	isHorizontal := flowchart.Direction == FlowDirectionLR || flowchart.Direction == FlowDirectionRL

	nodePositions := make(map[string]struct{ x, y styling.Length })
	nodeActualWidths := make(map[string]styling.Length)
	nodeShapeIndices := make(map[string]int)

	var minX, minY, maxX, maxY styling.Length
	firstElement := true

	updateBounds := func(x, y, cx, cy styling.Length) {
		if firstElement {
			minX, minY = x, y
			maxX, maxY = x+cx, y+cy
			firstElement = false
		} else {
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
			if x+cx > maxX {
				maxX = x + cx
			}
			if y+cy > maxY {
				maxY = y + cy
			}
		}
	}

	calculateWidth := func(label string) styling.Length {
		w := baseNodeWidth
		if len(label) > 15 {
			w += styling.Length(len(label)-15) * styling.Inches(0.08)
		}
		return w
	}

	if len(flowchart.Subgraphs) > 0 {
		subgraphX := styling.Inches(0.5)
		subgraphStartY := styling.Inches(1.5)

		for _, sg := range flowchart.Subgraphs {
			maxSgNodeWidth := baseNodeWidth
			for _, nodeID := range sg.Nodes {
				for _, n := range flowchart.Nodes {
					if n.ID == nodeID {
						nw := calculateWidth(n.Label)
						if nw > maxSgNodeWidth {
							maxSgNodeWidth = nw
						}
					}
				}
			}

			sgWidth := maxSgNodeWidth + styling.Inches(0.5)
			titleHeight := styling.Inches(0.3)
			sgHeight := styling.Length(len(sg.Nodes))*vSpacing + styling.Inches(0.5) + titleHeight
			sgX := subgraphX
			sgY := subgraphStartY

			// Subgraph background
			sgShape := shapes.NewShape(shapes.ShapeTypeRoundedRectangle, sgX, sgY, sgWidth, sgHeight).
				WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
				WithLine(shapes.NewShapeLine(theme.SecondaryStroke, styling.Emu(12700)))
			shapesList = append(shapesList, sgShape)
			updateBounds(sgX, sgY, sgWidth, sgHeight)

			// Subgraph title
			titleShape := shapes.NewShape(shapes.ShapeTypeRectangle, sgX+styling.Inches(0.1), sgY+styling.Inches(0.1), sgWidth-styling.Inches(0.2), titleHeight).
				WithText(sg.Name).
				WithAutoFit(shapes.TextAutoFitNormal)
			titleShape.Fill = nil
			titleShape.Line = nil
			shapesList = append(shapesList, titleShape)

			// Layout nodes within subgraph
			for j, nodeID := range sg.Nodes {
				var node *FlowNode
				for k := range flowchart.Nodes {
					if flowchart.Nodes[k].ID == nodeID {
						node = &flowchart.Nodes[k]
						break
					}
				}

				if node != nil {
					nw := calculateWidth(node.Label)
					x := sgX + (sgWidth-nw)/2
					y := sgY + titleHeight + styling.Inches(0.3) + styling.Length(j)*vSpacing

					nodePositions[node.ID] = struct{ x, y styling.Length }{x, y}
					nodeActualWidths[node.ID] = nw
					shape := createNodeShape(node, x, y, nw, nodeHeight, theme)
					shapesList = append(shapesList, shape)
					nodeShapeIndices[node.ID] = len(shapesList)
					updateBounds(x, y, nw, nodeHeight)
				}
			}

			subgraphX += sgWidth + styling.Inches(0.8)
		}

		// Layout orphan nodes
		orphanY := subgraphStartY
		for _, node := range flowchart.Nodes {
			if _, exists := nodePositions[node.ID]; !exists {
				nw := calculateWidth(node.Label)
				x := subgraphX
				y := orphanY

				nodePositions[node.ID] = struct{ x, y styling.Length }{x, y}
				nodeActualWidths[node.ID] = nw
				shape := createNodeShape(&node, x, y, nw, nodeHeight, theme)
				shapesList = append(shapesList, shape)
				nodeShapeIndices[node.ID] = len(shapesList)
				updateBounds(x, y, nw, nodeHeight)

				orphanY += vSpacing
			}
		}
	} else {
		// Simple grid layout
		startX := styling.Inches(1.0)
		startY := styling.Inches(1.5)
		cols := 1
		if isHorizontal {
			cols = min(nodeCount, 5)
		}

		for i, node := range flowchart.Nodes {
			col := i % cols
			row := i / cols

			nw := calculateWidth(node.Label)
			var x, y styling.Length
			if isHorizontal {
				x = startX + styling.Length(col)*hSpacing
				y = startY + styling.Length(row)*vSpacing
			} else {
				x = startX + styling.Length(col)*hSpacing
				y = startY + styling.Length(i)*vSpacing
			}

			nodePositions[node.ID] = struct{ x, y styling.Length }{x, y}
			nodeActualWidths[node.ID] = nw
			shape := createNodeShape(&node, x, y, nw, nodeHeight, theme)
			shapesList = append(shapesList, shape)
			nodeShapeIndices[node.ID] = len(shapesList)
			updateBounds(x, y, nw, nodeHeight)
		}
	}

	// Create connectors
	for _, conn := range flowchart.Connections {
		fromPos, fromExists := nodePositions[conn.From]
		toPos, toExists := nodePositions[conn.To]

		if fromExists && toExists {
			fromWidth := nodeActualWidths[conn.From]
			toWidth := nodeActualWidths[conn.To]
			var startSite, endSite string
			var startX, startY, endX, endY styling.Length

			if isHorizontal {
				startSite = shapes.ConnectionSiteRight
				endSite = shapes.ConnectionSiteLeft
				startX = fromPos.x + fromWidth
				startY = fromPos.y + nodeHeight/2
				endX = toPos.x
				endY = toPos.y + nodeHeight/2
			} else {
				startSite = shapes.ConnectionSiteBottom
				endSite = shapes.ConnectionSiteTop
				startX = fromPos.x + fromWidth/2
				startY = fromPos.y + nodeHeight
				endX = toPos.x + toWidth/2
				endY = toPos.y
			}

			connectorType := shapes.ConnectorTypeElbow
			diffX := startX - endX
			if diffX < 0 {
				diffX = -diffX
			}
			diffY := startY - endY
			if diffY < 0 {
				diffY = -diffY
			}

			if diffX < styling.Inches(0.1) || diffY < styling.Inches(0.1) {
				connectorType = shapes.ConnectorTypeStraight
			}

			lineColor := theme.PrimaryStroke
			lineWidth := theme.LineWeight
			lineDash := shapes.LineDashSolid
			switch conn.ArrowStyle {
			case ArrowStyleArrow:
				// Default style; nothing to override.
			case ArrowStyleThick:
				lineWidth = theme.LineWeight * 2
			case ArrowStyleDotted:
				lineDash = shapes.LineDashDash
			case ArrowStyleOpen:
				// Handled by WithArrows
			}

			connector := shapes.NewConnector(connectorType, startX, startY, endX, endY).
				WithLine(shapes.NewShapeLine(lineColor, lineWidth).WithDash(lineDash)).
				WithArrows(shapes.ArrowTypeNone, shapes.ArrowTypeTriangle)

			if conn.ArrowStyle == ArrowStyleOpen {
				connector = connector.WithArrows(shapes.ArrowTypeNone, shapes.ArrowTypeNone)
			}

			if idx, ok := nodeShapeIndices[conn.From]; ok {
				connector = connector.ConnectStart(idx, startSite)
			}
			if idx, ok := nodeShapeIndices[conn.To]; ok {
				connector = connector.ConnectEnd(idx, endSite)
			}

			connectors = append(connectors, connector)

			if conn.Label != "" {
				labelWidth := styling.Inches(0.8)
				labelHeight := styling.Inches(0.3)
				midX := (startX + endX) / 2
				midY := (startY + endY) / 2

				labelShape := shapes.NewShape(
					shapes.ShapeTypeRectangle,
					midX-labelWidth/2,
					midY-labelHeight/2,
					labelWidth,
					labelHeight,
				).WithFill(shapes.NewShapeFill(theme.Background)).
					WithLine(shapes.NewShapeLine(theme.SecondaryStroke, styling.Emu(12700))).
					WithText(conn.Label).
					WithAutoFit(shapes.TextAutoFitNormal)

				shapesList = append(shapesList, labelShape)
				updateBounds(midX-labelWidth/2, midY-labelHeight/2, labelWidth, labelHeight)
			}
		}
	}

	return DiagramElements{
		Shapes:     shapesList,
		Connectors: connectors,
		Grouped:    true,
		Bounds: &DiagramBounds{
			X:  minX,
			Y:  minY,
			CX: maxX - minX,
			CY: maxY - minY,
		},
	}
}

func getSubgraphColor(index int) string {
	colors := []string{"E3F2FD", "F3E5F5", "E8F5E9", "FFF3E0", "E0F7FA", "FCE4EC"}
	return colors[index%len(colors)]
}

func createNodeShape(node *FlowNode, x, y, width, height styling.Length, theme Theme) shapes.Shape {
	shapeType := shapes.ShapeTypeRectangle
	switch node.Shape {
	case NodeShapeRectangle:
		shapeType = shapes.ShapeTypeRectangle
	case NodeShapeRoundedRect:
		shapeType = shapes.ShapeTypeRoundedRectangle
	case NodeShapeStadium:
		shapeType = shapes.ShapeTypeRoundedRectangle // Stadium is often represented as rounded rect in PPT
	case NodeShapeDiamond:
		shapeType = shapes.ShapeTypeDiamond
	case NodeShapeCircle:
		shapeType = shapes.ShapeTypeEllipse
	case NodeShapeHexagon:
		shapeType = shapes.ShapeTypeHexagon
	}

	fillColor := theme.PrimaryFill
	if node.Shape == NodeShapeDiamond {
		fillColor = theme.SecondaryFill
	}

	return shapes.NewShape(shapeType, x, y, width, height).
		WithFill(shapes.NewShapeFill(fillColor)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithText(node.Label).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
}
