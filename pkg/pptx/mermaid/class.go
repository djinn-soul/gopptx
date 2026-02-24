package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// ClassNode represents a class in a class diagram.
type ClassNode struct {
	ID         string
	Name       string
	Attributes []string
	Methods    []string
}

// ClassRelationship represents a relationship between two classes.
type ClassRelationship struct {
	From  string
	To    string
	Type  string // <|--, --, -->, etc.
	Label string
}

// ClassDiagram represents the parsed structure of a Mermaid class diagram.
type ClassDiagram struct {
	Classes       []ClassNode
	Relationships []ClassRelationship
}

// renderClass parses and renders a Mermaid class diagram into PowerPoint elements.
func renderClass(code string, theme Theme) DiagramElements {
	diagram := parseClass(code)
	return generateClassElements(diagram, theme)
}

func parseClass(code string) *ClassDiagram {
	lines := ParseLines(code)
	classes := make(map[string]*ClassNode)
	var relationships []ClassRelationship

	var currentClass *ClassNode

	for i := range lines {
		line := lines[i]
		if strings.HasPrefix(line, "classDiagram") ||
			strings.HasPrefix(line, "class ") && !strings.Contains(line, "{") && !strings.Contains(line, ":") &&
				i == 0 {
			continue
		}

		// Handle class definition with braces
		if strings.HasPrefix(line, "class ") && strings.HasSuffix(line, "{") {
			name := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(line, "{"), "class"))
			if _, ok := classes[name]; !ok {
				classes[name] = &ClassNode{ID: name, Name: name}
			}
			currentClass = classes[name]
			continue
		}

		if line == "}" {
			currentClass = nil
			continue
		}

		if currentClass != nil {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(trimmed, "(") {
				currentClass.Methods = append(currentClass.Methods, trimmed)
			} else {
				currentClass.Attributes = append(currentClass.Attributes, trimmed)
			}
			continue
		}

		// Handle class members with colon: ClassName : member
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			className := strings.TrimSpace(parts[0])
			member := strings.TrimSpace(parts[1])

			if _, ok := classes[className]; !ok {
				classes[className] = &ClassNode{ID: className, Name: className}
			}
			c := classes[className]
			if strings.Contains(member, "(") {
				c.Methods = append(c.Methods, member)
			} else {
				c.Attributes = append(c.Attributes, member)
			}
			continue
		}

		// Handle relationships
		if from, relType, to, found := splitClassRelationship(line); found {
			relationships = append(relationships, ClassRelationship{
				From: from,
				To:   to,
				Type: relType,
			})
			if _, ok := classes[from]; !ok {
				classes[from] = &ClassNode{ID: from, Name: from}
			}
			if _, ok := classes[to]; !ok {
				classes[to] = &ClassNode{ID: to, Name: to}
			}
			continue
		}

		// Handle simple class definition: class ClassName
		if strings.HasPrefix(line, "class ") {
			name := strings.TrimSpace(strings.TrimPrefix(line, "class"))
			if _, ok := classes[name]; !ok {
				classes[name] = &ClassNode{ID: name, Name: name}
			}
			continue
		}
	}

	classList := make([]ClassNode, 0, len(classes))
	for _, c := range classes {
		classList = append(classList, *c)
	}

	return &ClassDiagram{
		Classes:       classList,
		Relationships: relationships,
	}
}

func splitClassRelationship(line string) (string, string, string, bool) {
	relTypes := []string{"<|--", "*--", "o--", "-->", "--", "..>", "..", "<|..", "*..", "o.."}
	for _, rt := range relTypes {
		if before, after, ok := strings.Cut(line, rt); ok {
			from := strings.TrimSpace(before)
			rest := strings.TrimSpace(after)
			// Rest might contain a label: "To : label"
			to := rest
			if before, _, ok := strings.Cut(rest, ":"); ok {
				to = strings.TrimSpace(before)
			}
			return from, rt, to, true
		}
	}
	return "", "", "", false
}

func generateClassElements(diagram *ClassDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape
	var connectors []shapes.Connector

	if len(diagram.Classes) == 0 {
		return DiagramElements{Grouped: true}
	}

	// Layout parameters
	classWidth := styling.Inches(2.2)
	headerHeight := styling.Inches(0.5)
	itemHeight := styling.Inches(0.35)
	hSpacing := styling.Inches(3.0)
	vSpacing := styling.Inches(2.5)

	classPositions := make(map[string]struct{ x, y styling.Length })
	classShapeIndices := make(map[string]int)

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

	for i, class := range diagram.Classes {
		col := i % cols
		row := i / cols

		x := startX + styling.Length(col)*hSpacing
		y := startY + styling.Length(row)*vSpacing

		classPositions[class.ID] = struct{ x, y styling.Length }{x, y}

		// Calculate total height
		attrCount := len(class.Attributes)
		methodCount := len(class.Methods)
		if attrCount == 0 {
			attrCount = 1 // Empty slot
		}
		if methodCount == 0 {
			methodCount = 1 // Empty slot
		}
		totalHeight := headerHeight + styling.Length(attrCount+methodCount)*itemHeight + styling.Inches(0.1)

		// Header
		header := shapes.NewShape(shapes.ShapeTypeRectangle, x, y, classWidth, headerHeight).
			WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
			WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
			WithText(class.Name).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
		shapesList = append(shapesList, header)
		classShapeIndices[class.ID] = len(shapesList)
		updateBounds(x, y, classWidth, totalHeight)

		// Attributes box
		attrY := y + headerHeight
		attrHeight := styling.Length(attrCount) * itemHeight
		attrText := strings.Join(class.Attributes, "\n")
		if attrText == "" {
			attrText = " "
		}
		attrBox := shapes.NewShape(shapes.ShapeTypeRectangle, x, attrY, classWidth, attrHeight).
			WithFill(shapes.NewShapeFill(theme.Background)).
			WithLine(shapes.NewShapeLine(theme.PrimaryStroke, styling.Emu(12700))).
			WithText(attrText).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithVerticalAnchor(shapes.TextAnchorTop).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
		shapesList = append(shapesList, attrBox)

		// Methods box
		methodY := attrY + attrHeight
		methodHeight := styling.Length(methodCount) * itemHeight
		methodText := strings.Join(class.Methods, "\n")
		if methodText == "" {
			methodText = " "
		}
		methodBox := shapes.NewShape(shapes.ShapeTypeRectangle, x, methodY, classWidth, methodHeight).
			WithFill(shapes.NewShapeFill(theme.Background)).
			WithLine(shapes.NewShapeLine(theme.PrimaryStroke, styling.Emu(12700))).
			WithText(methodText).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithVerticalAnchor(shapes.TextAnchorTop).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
		shapesList = append(shapesList, methodBox)
	}

	// Create connectors
	for _, rel := range diagram.Relationships {
		fromPos, fromExists := classPositions[rel.From]
		toPos, toExists := classPositions[rel.To]

		if fromExists && toExists {
			// Determine best connection points (simple top/bottom for now)
			var startX, startY, endX, endY styling.Length
			var startSite, endSite string

			if fromPos.y < toPos.y {
				startX = fromPos.x + classWidth/2
				startY = fromPos.y + headerHeight // Anchor to header bottom or just use the whole box?
				// Actually, it's better to anchor to the whole class box, but we only have the header index.
				// Let's just use the header for now as the anchor point.
				startSite = shapes.ConnectionSiteBottom
				endSite = shapes.ConnectionSiteTop
				endX = toPos.x + classWidth/2
				endY = toPos.y
			} else {
				startX = fromPos.x + classWidth/2
				startY = fromPos.y
				startSite = shapes.ConnectionSiteTop
				endSite = shapes.ConnectionSiteBottom
				endX = toPos.x + classWidth/2
				endY = toPos.y + headerHeight
			}

			connector := shapes.NewConnector(shapes.ConnectorTypeElbow, startX, startY, endX, endY).
				WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight))

			// Set arrow types based on relationship type
			startArrow := shapes.ArrowTypeNone
			endArrow := shapes.ArrowTypeTriangle

			switch rel.Type {
			case "<|--", "<|..":
				endArrow = shapes.ArrowTypeStealth // Inheritance
			case "*--", "*..":
				startArrow = shapes.ArrowTypeDiamond // Composition
			case "o--", "o..":
				startArrow = shapes.ArrowTypeOval // Aggregation (should be hollow diamond, but Oval is closest if not available)
			case "-->", "..>":
				endArrow = shapes.ArrowTypeTriangle // Association
			case "--", "..":
				endArrow = shapes.ArrowTypeNone // Link
			}

			if strings.Contains(rel.Type, "..") {
				connector = connector.WithLine(
					shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight).WithDash(shapes.LineDashDash),
				)
			}

			connector = connector.WithArrows(startArrow, endArrow)

			if idx, ok := classShapeIndices[rel.From]; ok {
				connector = connector.ConnectStart(idx, startSite)
			}
			if idx, ok := classShapeIndices[rel.To]; ok {
				connector = connector.ConnectEnd(idx, endSite)
			}

			connectors = append(connectors, connector)
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
