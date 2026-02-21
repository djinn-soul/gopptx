package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// EREntity represents an entity in an ER diagram.
type EREntity struct {
	Name       string
	Attributes []string
}

// ERRelationship represents a relationship between two entities.
type ERRelationship struct {
	From  string
	To    string
	Type  string // ||--o{, etc.
	Label string
}

// ERDiagram represents the parsed structure of a Mermaid ER diagram.
type ERDiagram struct {
	Entities      []EREntity
	Relationships []ERRelationship
}

// renderER parses and renders a Mermaid ER diagram into PowerPoint elements.
func renderER(code string, theme Theme) DiagramElements {
	diagram := parseER(code)
	return generateERElements(diagram, theme)
}

func parseER(code string) *ERDiagram {
	lines := ParseLines(code)
	entities := make(map[string]*EREntity)
	var relationships []ERRelationship

	var currentEntity *EREntity

	for _, line := range lines {
		if strings.HasPrefix(line, "erDiagram") {
			continue
		}

		// Handle entity definition with braces
		if strings.HasSuffix(line, "{") {
			name := strings.TrimSpace(strings.TrimSuffix(line, "{"))
			if _, ok := entities[name]; !ok {
				entities[name] = &EREntity{Name: name}
			}
			currentEntity = entities[name]
			continue
		}

		if line == "}" {
			currentEntity = nil
			continue
		}

		if currentEntity != nil {
			currentEntity.Attributes = append(currentEntity.Attributes, strings.TrimSpace(line))
			continue
		}

		// Handle relationships
		if from, relType, to, label, found := splitERRelationship(line); found {
			relationships = append(relationships, ERRelationship{
				From:  from,
				To:    to,
				Type:  relType,
				Label: label,
			})
			if _, ok := entities[from]; !ok {
				entities[from] = &EREntity{Name: from}
			}
			if _, ok := entities[to]; !ok {
				entities[to] = &EREntity{Name: to}
			}
			continue
		}

		// Simple entity definition
		if !strings.Contains(line, " ") && !strings.Contains(line, "-") {
			name := strings.TrimSpace(line)
			if _, ok := entities[name]; !ok {
				entities[name] = &EREntity{Name: name}
			}
		}
	}

	entityList := make([]EREntity, 0, len(entities))
	for _, e := range entities {
		entityList = append(entityList, *e)
	}

	return &ERDiagram{
		Entities:      entityList,
		Relationships: relationships,
	}
}

func splitERRelationship(line string) (string, string, string, string, bool) {
	relTypes := []string{"||--o{", "||--|{", "}|--|{", "}|--o{", "|o--o{", "|o--|{", "o{--}o", "o{--|{", "o{--o{", "||--||", "||--|o", "|o--|o", "|o--||"}
	// Also simpler ones
	relTypes = append(relTypes, "||--", "}|--", "o{--", "--o{", "--|{", "--}o", "--")

	for _, rt := range relTypes {
		if idx := strings.Index(line, rt); idx != -1 {
			from := strings.TrimSpace(line[:idx])
			rest := strings.TrimSpace(line[idx+len(rt):])
			to := rest
			label := ""
			if labelIdx := strings.Index(rest, ":"); labelIdx != -1 {
				to = strings.TrimSpace(rest[:labelIdx])
				label = strings.TrimSpace(rest[labelIdx+1:])
			}
			return from, rt, to, label, true
		}
	}
	return "", "", "", "", false
}

func generateERElements(diagram *ERDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape
	var connectors []shapes.Connector

	if len(diagram.Entities) == 0 {
		return DiagramElements{Grouped: true}
	}

	// Layout parameters
	entityWidth := styling.Inches(2.2)
	headerHeight := styling.Inches(0.5)
	itemHeight := styling.Inches(0.35)
	hSpacing := styling.Inches(3.0)
	vSpacing := styling.Inches(2.5)

	entityPositions := make(map[string]struct{ x, y styling.Length })
	entityShapeIndices := make(map[string]int)

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

	// Simple grid layout
	startX := styling.Inches(1.0)
	startY := styling.Inches(1.0)
	cols := 3

	for i, entity := range diagram.Entities {
		col := i % cols
		row := i / cols

		x := startX + styling.Length(col)*hSpacing
		y := startY + styling.Length(row)*vSpacing

		entityPositions[entity.Name] = struct{ x, y styling.Length }{x, y}

		// Calculate total height
		attrCount := len(entity.Attributes)
		totalHeight := headerHeight + styling.Length(attrCount)*itemHeight

		// Header
		header := shapes.NewShape(shapes.ShapeTypeRectangle, x, y, entityWidth, headerHeight).
			WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
			WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
			WithText(entity.Name).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
		shapesList = append(shapesList, header)
		entityShapeIndices[entity.Name] = len(shapesList)
		updateBounds(x, y, entityWidth, totalHeight)

		// Attributes box
		if attrCount > 0 {
			attrY := y + headerHeight
			attrHeight := styling.Length(attrCount) * itemHeight
			attrBox := shapes.NewShape(shapes.ShapeTypeRectangle, x, attrY, entityWidth, attrHeight).
				WithFill(shapes.NewShapeFill(theme.Background)).
				WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
				WithText(strings.Join(entity.Attributes, "\n")).
				WithAutoFit(shapes.TextAutoFitNormal).
				WithVerticalAnchor(shapes.TextAnchorTop).
				WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
			shapesList = append(shapesList, attrBox)
		}
	}

	// Create connectors
	for _, rel := range diagram.Relationships {
		fromPos, fromExists := entityPositions[rel.From]
		toPos, toExists := entityPositions[rel.To]

		if fromExists && toExists {
			var startX, startY, endX, endY styling.Length
			var startSite, endSite string

			if fromPos.y < toPos.y {
				startX = fromPos.x + entityWidth/2
				startY = fromPos.y + headerHeight
				startSite = shapes.ConnectionSiteBottom
				endSite = shapes.ConnectionSiteTop
				endX = toPos.x + entityWidth/2
				endY = toPos.y
			} else {
				startX = fromPos.x + entityWidth/2
				startY = fromPos.y
				startSite = shapes.ConnectionSiteTop
				endSite = shapes.ConnectionSiteBottom
				endX = toPos.x + entityWidth/2
				endY = toPos.y + headerHeight
			}

			connector := shapes.NewConnector(shapes.ConnectorTypeElbow, startX, startY, endX, endY).
				WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight))

			// Set arrow types based on relationship type
			// ER notation is complex, we'll use simple arrows for now
			startArrow := shapes.ArrowTypeNone
			endArrow := shapes.ArrowTypeTriangle

			if strings.Contains(rel.Type, "o") {
				// Optional
			}
			if strings.Contains(rel.Type, "{") || strings.Contains(rel.Type, "}") {
				// Many
				endArrow = shapes.ArrowTypeStealth
			}

			connector = connector.WithArrows(startArrow, endArrow)

			if idx, ok := entityShapeIndices[rel.From]; ok {
				connector = connector.ConnectStart(idx, startSite)
			}
			if idx, ok := entityShapeIndices[rel.To]; ok {
				connector = connector.ConnectEnd(idx, endSite)
			}

			connectors = append(connectors, connector)

			if rel.Label != "" {
				labelWidth := styling.Inches(1.0)
				labelHeight := styling.Inches(0.4)
				midX := (startX + endX) / 2
				midY := (startY + endY) / 2

				labelShape := shapes.NewShape(shapes.ShapeTypeRectangle, midX-labelWidth/2, midY-labelHeight/2, labelWidth, labelHeight).
					WithFill(shapes.NewShapeFill(theme.Background)).
					WithText(rel.Label).
					WithAutoFit(shapes.TextAutoFitNormal).
					WithTextMargins(styling.Inches(0.05), styling.Inches(0.02), styling.Inches(0.05), styling.Inches(0.02))
				shapesList = append(shapesList, labelShape)
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
