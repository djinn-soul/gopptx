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
		currentClass = parseClassLine(line, i, classes, currentClass, &relationships)
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

func parseClassLine(
	line string,
	lineIndex int,
	classes map[string]*ClassNode,
	currentClass *ClassNode,
	relationships *[]ClassRelationship,
) *ClassNode {
	if shouldSkipClassLine(line, lineIndex) {
		return currentClass
	}
	if className, ok := classBlockStart(line); ok {
		return ensureClassNode(classes, className)
	}
	if line == "}" {
		return nil
	}
	if currentClass != nil {
		appendClassMember(currentClass, line)
		return currentClass
	}
	if className, member, ok := parseClassInlineMember(line); ok {
		appendClassMember(ensureClassNode(classes, className), member)
		return currentClass
	}
	if rel, ok := parseClassRelationshipLine(line); ok {
		*relationships = append(*relationships, rel)
		ensureClassNode(classes, rel.From)
		ensureClassNode(classes, rel.To)
		return currentClass
	}
	if className, ok := parseSimpleClassDefinition(line); ok {
		ensureClassNode(classes, className)
	}
	return currentClass
}

func shouldSkipClassLine(line string, lineIndex int) bool {
	if strings.HasPrefix(line, "classDiagram") {
		return true
	}
	return strings.HasPrefix(line, "class ") &&
		!strings.Contains(line, "{") &&
		!strings.Contains(line, ":") &&
		lineIndex == 0
}

func classBlockStart(line string) (string, bool) {
	if !strings.HasPrefix(line, "class ") || !strings.HasSuffix(line, "{") {
		return "", false
	}
	name := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(line, "{"), "class"))
	return name, name != ""
}

func ensureClassNode(classes map[string]*ClassNode, name string) *ClassNode {
	if _, ok := classes[name]; !ok {
		classes[name] = &ClassNode{ID: name, Name: name}
	}
	return classes[name]
}

func appendClassMember(class *ClassNode, member string) {
	trimmed := strings.TrimSpace(member)
	if strings.Contains(trimmed, "(") {
		class.Methods = append(class.Methods, trimmed)
		return
	}
	class.Attributes = append(class.Attributes, trimmed)
}

func parseClassInlineMember(line string) (string, string, bool) {
	if !strings.Contains(line, ":") {
		return "", "", false
	}
	parts := strings.SplitN(line, ":", 2)
	className := strings.TrimSpace(parts[0])
	member := strings.TrimSpace(parts[1])
	return className, member, className != "" && member != ""
}

func parseClassRelationshipLine(line string) (ClassRelationship, bool) {
	from, relType, to, found := splitClassRelationship(line)
	if !found {
		return ClassRelationship{}, false
	}
	return ClassRelationship{From: from, To: to, Type: relType}, true
}

func parseSimpleClassDefinition(line string) (string, bool) {
	if !strings.HasPrefix(line, "class ") {
		return "", false
	}
	name := strings.TrimSpace(strings.TrimPrefix(line, "class"))
	return name, name != ""
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
	if len(diagram.Classes) == 0 {
		return DiagramElements{Grouped: true}
	}

	layout := defaultClassLayout()
	state := newClassRenderState(layout)
	for i, class := range diagram.Classes {
		classShapes := classShapesForNode(class, i, layout, theme)
		state.addClass(class.ID, classShapes)
	}
	for _, rel := range diagram.Relationships {
		connector, marker, ok := classRelationshipConnector(rel, state, layout, theme)
		if !ok {
			continue
		}
		state.connectors = append(state.connectors, connector)
		if marker != nil {
			state.shapes = append(state.shapes, *marker)
		}
	}
	return state.diagramElements()
}

type classLayout struct {
	classWidth   styling.Length
	headerHeight styling.Length
	itemHeight   styling.Length
	hSpacing     styling.Length
	vSpacing     styling.Length
	startX       styling.Length
	startY       styling.Length
	cols         int
}

type classRenderState struct {
	layout       classLayout
	shapes       []shapes.Shape
	connectors   []shapes.Connector
	positions    map[string]classPosition
	shapeIndices map[string]int
	bounds       classBounds
}

type classPosition struct {
	x styling.Length
	y styling.Length
}

type classBounds struct {
	minX  styling.Length
	minY  styling.Length
	maxX  styling.Length
	maxY  styling.Length
	empty bool
}

type classNodeShapes struct {
	position    classPosition
	totalHeight styling.Length
	header      shapes.Shape
	attrBox     shapes.Shape
	methodBox   shapes.Shape
}

type classConnectorGeometry struct {
	startX    styling.Length
	startY    styling.Length
	endX      styling.Length
	endY      styling.Length
	startSite string
	endSite   string
}

func defaultClassLayout() classLayout {
	return classLayout{
		classWidth:   styling.Inches(2.2),
		headerHeight: styling.Inches(0.5),
		itemHeight:   styling.Inches(0.35),
		hSpacing:     styling.Inches(3.0),
		vSpacing:     styling.Inches(2.5),
		startX:       styling.Inches(1.0),
		startY:       styling.Inches(1.0),
		cols:         3,
	}
}

func newClassRenderState(layout classLayout) *classRenderState {
	return &classRenderState{
		layout:       layout,
		positions:    make(map[string]classPosition),
		shapeIndices: make(map[string]int),
		bounds:       classBounds{empty: true},
	}
}

func classPositionForIndex(index int, layout classLayout) classPosition {
	col := index % layout.cols
	row := index / layout.cols
	return classPosition{
		x: layout.startX + styling.Length(col)*layout.hSpacing,
		y: layout.startY + styling.Length(row)*layout.vSpacing,
	}
}

func classHeights(class ClassNode, layout classLayout) (styling.Length, styling.Length, styling.Length) {
	attrCount := len(class.Attributes)
	methodCount := len(class.Methods)
	attrHeight := styling.Length(attrCount) * layout.itemHeight
	methodHeight := styling.Length(methodCount) * layout.itemHeight
	total := layout.headerHeight + attrHeight + methodHeight
	return attrHeight, methodHeight, total
}

func classShapesForNode(class ClassNode, index int, layout classLayout, theme Theme) classNodeShapes {
	position := classPositionForIndex(index, layout)
	attrHeight, methodHeight, totalHeight := classHeights(class, layout)
	header := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		position.x,
		position.y,
		layout.classWidth,
		layout.headerHeight,
	).WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithText(class.Name).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))

	attrText := strings.Join(class.Attributes, "\n")
	attrY := position.y + layout.headerHeight
	attrBox := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		position.x,
		attrY,
		layout.classWidth,
		attrHeight,
	).WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, styling.Emu(12700))).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithVerticalAnchor(shapes.TextAnchorTop).
		WithText(attrText).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))

	methodText := strings.Join(class.Methods, "\n")
	methodY := attrY + attrHeight
	methodBox := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		position.x,
		methodY,
		layout.classWidth,
		methodHeight,
	).WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, styling.Emu(12700))).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithVerticalAnchor(shapes.TextAnchorTop).
		WithText(methodText).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))

	return classNodeShapes{
		position:    position,
		totalHeight: totalHeight,
		header:      header,
		attrBox:     attrBox,
		methodBox:   methodBox,
	}
}

func (s *classRenderState) addClass(classID string, nodeShapes classNodeShapes) {
	s.positions[classID] = nodeShapes.position
	s.shapes = append(s.shapes, nodeShapes.header)
	s.shapeIndices[classID] = len(s.shapes)
	s.bounds.include(nodeShapes.position.x, nodeShapes.position.y, s.layout.classWidth, nodeShapes.totalHeight)
	if nodeShapes.attrBox.CY > 0 {
		s.shapes = append(s.shapes, nodeShapes.attrBox)
	}
	if nodeShapes.methodBox.CY > 0 {
		s.shapes = append(s.shapes, nodeShapes.methodBox)
	}
}

func classRelationshipConnector(
	rel ClassRelationship,
	state *classRenderState,
	layout classLayout,
	theme Theme,
) (shapes.Connector, *shapes.Shape, bool) {
	fromID := rel.From
	toID := rel.To
	// Mermaid inheritance arrows point toward the parent class.
	if strings.HasPrefix(rel.Type, "<|") {
		fromID, toID = rel.To, rel.From
	}

	fromPos, fromExists := state.positions[fromID]
	toPos, toExists := state.positions[toID]
	if !fromExists || !toExists {
		return shapes.Connector{}, nil, false
	}
	geometry := classConnectorPoints(fromPos, toPos, layout)
	line := shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)
	if strings.Contains(rel.Type, "..") {
		line = line.WithDash(shapes.LineDashDash)
	}
	startArrow, endArrow := classArrowTypes(rel.Type)
	var marker *shapes.Shape
	if classNeedsHollowInheritanceMarker(rel.Type) {
		endArrow = shapes.ArrowTypeNone
		marker = classInheritanceMarker(geometry, theme)
	}
	connector := shapes.NewConnector(
		shapes.ConnectorTypeElbow,
		geometry.startX,
		geometry.startY,
		geometry.endX,
		geometry.endY,
	).WithLine(line).WithArrows(startArrow, endArrow)
	if idx, ok := state.shapeIndices[fromID]; ok {
		connector = connector.ConnectStart(idx, geometry.startSite)
	}
	if idx, ok := state.shapeIndices[toID]; ok {
		connector = connector.ConnectEnd(idx, geometry.endSite)
	}
	return connector, marker, true
}
