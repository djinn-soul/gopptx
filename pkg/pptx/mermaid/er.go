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
		currentEntity = parseERLine(line, entities, currentEntity, &relationships)
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

func parseERLine(
	line string,
	entities map[string]*EREntity,
	currentEntity *EREntity,
	relationships *[]ERRelationship,
) *EREntity {
	if strings.HasPrefix(line, "erDiagram") {
		return currentEntity
	}
	if entityName, ok := parseEREntityBlockStart(line); ok {
		return ensureEREntity(entities, entityName)
	}
	if line == "}" {
		return nil
	}
	if currentEntity != nil {
		currentEntity.Attributes = append(currentEntity.Attributes, strings.TrimSpace(line))
		return currentEntity
	}
	if rel, ok := parseERRelationshipLine(line); ok {
		*relationships = append(*relationships, rel)
		ensureEREntity(entities, rel.From)
		ensureEREntity(entities, rel.To)
		return currentEntity
	}
	if entityName, ok := parseERSimpleEntity(line); ok {
		ensureEREntity(entities, entityName)
	}
	return currentEntity
}

func parseEREntityBlockStart(line string) (string, bool) {
	before, ok := strings.CutSuffix(line, "{")
	if !ok {
		return "", false
	}
	name := strings.TrimSpace(before)
	return name, name != ""
}

func ensureEREntity(entities map[string]*EREntity, name string) *EREntity {
	if _, ok := entities[name]; !ok {
		entities[name] = &EREntity{Name: name}
	}
	return entities[name]
}

func parseERRelationshipLine(line string) (ERRelationship, bool) {
	from, relType, to, label, found := splitERRelationship(line)
	if !found {
		return ERRelationship{}, false
	}
	return ERRelationship{
		From:  from,
		To:    to,
		Type:  relType,
		Label: label,
	}, true
}

func parseERSimpleEntity(line string) (string, bool) {
	if strings.Contains(line, " ") || strings.Contains(line, "-") {
		return "", false
	}
	name := strings.TrimSpace(line)
	return name, name != ""
}

func splitERRelationship(line string) (string, string, string, string, bool) {
	relTypes := []string{
		"||--o{",
		"||--|{",
		"}|--|{",
		"}|--o{",
		"|o--o{",
		"|o--|{",
		"o{--}o",
		"o{--|{",
		"o{--o{",
		"||--||",
		"||--|o",
		"|o--|o",
		"|o--||",
	}
	// Also simpler ones
	relTypes = append(relTypes, "||--", "}|--", "o{--", "--o{", "--|{", "--}o", "--")

	for _, rt := range relTypes {
		if before, after, ok := strings.Cut(line, rt); ok {
			from := strings.TrimSpace(before)
			rest := strings.TrimSpace(after)
			to := rest
			label := ""
			if before, after, ok := strings.Cut(rest, ":"); ok {
				to = strings.TrimSpace(before)
				label = strings.TrimSpace(after)
			}
			return from, rt, to, label, true
		}
	}
	return "", "", "", "", false
}

func generateERElements(diagram *ERDiagram, theme Theme) DiagramElements {
	if len(diagram.Entities) == 0 {
		return DiagramElements{Grouped: true}
	}

	layout := defaultERLayout()
	state := newERRenderState(layout)
	for i, entity := range diagram.Entities {
		entityShapes := erShapesForEntity(entity, i, layout, theme)
		state.addEntity(entity.Name, entityShapes)
	}
	for _, rel := range diagram.Relationships {
		connector, labelShape, hasLabel, ok := erRelationshipConnector(rel, state, layout, theme)
		if !ok {
			continue
		}
		state.connectors = append(state.connectors, connector)
		if hasLabel {
			state.shapes = append(state.shapes, labelShape)
			state.bounds.includeShape(labelShape)
		}
	}
	return state.diagramElements()
}

type erLayout struct {
	entityWidth   styling.Length
	headerHeight  styling.Length
	itemHeight    styling.Length
	hSpacing      styling.Length
	vSpacing      styling.Length
	startX        styling.Length
	startY        styling.Length
	cols          int
	labelWidth    styling.Length
	labelHeight   styling.Length
	labelMarginX  styling.Length
	labelMarginY  styling.Length
	entityLineWgt styling.Length
}

type erPosition struct {
	x styling.Length
	y styling.Length
}

type erBounds struct {
	minX  styling.Length
	minY  styling.Length
	maxX  styling.Length
	maxY  styling.Length
	empty bool
}

type erRenderState struct {
	layout       erLayout
	shapes       []shapes.Shape
	connectors   []shapes.Connector
	positions    map[string]erPosition
	shapeIndices map[string]int
	bounds       erBounds
}

type erEntityShapes struct {
	position    erPosition
	totalHeight styling.Length
	header      shapes.Shape
	attrBox     shapes.Shape
	hasAttrBox  bool
}

type erConnectorGeometry struct {
	startX    styling.Length
	startY    styling.Length
	endX      styling.Length
	endY      styling.Length
	startSite string
	endSite   string
}

func defaultERLayout() erLayout {
	return erLayout{
		entityWidth:   styling.Inches(2.2),
		headerHeight:  styling.Inches(0.5),
		itemHeight:    styling.Inches(0.35),
		hSpacing:      styling.Inches(3.0),
		vSpacing:      styling.Inches(2.5),
		startX:        styling.Inches(1.0),
		startY:        styling.Inches(1.0),
		cols:          3,
		labelWidth:    styling.Inches(1.0),
		labelHeight:   styling.Inches(0.4),
		labelMarginX:  styling.Inches(0.05),
		labelMarginY:  styling.Inches(0.02),
		entityLineWgt: styling.Emu(12700),
	}
}

func newERRenderState(layout erLayout) *erRenderState {
	return &erRenderState{
		layout:       layout,
		positions:    make(map[string]erPosition),
		shapeIndices: make(map[string]int),
		bounds:       erBounds{empty: true},
	}
}

func erPositionForIndex(index int, layout erLayout) erPosition {
	col := index % layout.cols
	row := index / layout.cols
	return erPosition{
		x: layout.startX + styling.Length(col)*layout.hSpacing,
		y: layout.startY + styling.Length(row)*layout.vSpacing,
	}
}

func erShapesForEntity(entity EREntity, index int, layout erLayout, theme Theme) erEntityShapes {
	position := erPositionForIndex(index, layout)
	attrCount := len(entity.Attributes)
	totalHeight := layout.headerHeight + styling.Length(attrCount)*layout.itemHeight
	header := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		position.x,
		position.y,
		layout.entityWidth,
		layout.headerHeight,
	).WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithText(entity.Name).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
	if attrCount == 0 {
		return erEntityShapes{
			position:    position,
			totalHeight: totalHeight,
			header:      header,
			hasAttrBox:  false,
		}
	}
	attrY := position.y + layout.headerHeight
	attrHeight := styling.Length(attrCount) * layout.itemHeight
	attrBox := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		position.x,
		attrY,
		layout.entityWidth,
		attrHeight,
	).WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, layout.entityLineWgt)).
		WithText(strings.Join(entity.Attributes, "\n")).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithVerticalAnchor(shapes.TextAnchorTop).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
	return erEntityShapes{
		position:    position,
		totalHeight: totalHeight,
		header:      header,
		attrBox:     attrBox,
		hasAttrBox:  true,
	}
}

func (s *erRenderState) addEntity(entityName string, entityShapes erEntityShapes) {
	s.positions[entityName] = entityShapes.position
	s.shapes = append(s.shapes, entityShapes.header)
	s.shapeIndices[entityName] = len(s.shapes)
	s.bounds.include(entityShapes.position.x, entityShapes.position.y, s.layout.entityWidth, entityShapes.totalHeight)
	if entityShapes.hasAttrBox {
		s.shapes = append(s.shapes, entityShapes.attrBox)
	}
}

func erRelationshipConnector(
	rel ERRelationship,
	state *erRenderState,
	layout erLayout,
	theme Theme,
) (shapes.Connector, shapes.Shape, bool, bool) {
	fromPos, fromExists := state.positions[rel.From]
	toPos, toExists := state.positions[rel.To]
	if !fromExists || !toExists {
		return shapes.Connector{}, shapes.Shape{}, false, false
	}
	geometry := erConnectorPoints(fromPos, toPos, layout)
	endArrow := shapes.ArrowTypeTriangle
	if strings.Contains(rel.Type, "{") || strings.Contains(rel.Type, "}") {
		endArrow = shapes.ArrowTypeStealth
	}
	connector := shapes.NewConnector(
		shapes.ConnectorTypeElbow,
		geometry.startX,
		geometry.startY,
		geometry.endX,
		geometry.endY,
	).WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
		WithArrows(shapes.ArrowTypeNone, endArrow)
	if idx, ok := state.shapeIndices[rel.From]; ok {
		connector = connector.ConnectStart(idx, geometry.startSite)
	}
	if idx, ok := state.shapeIndices[rel.To]; ok {
		connector = connector.ConnectEnd(idx, geometry.endSite)
	}
	if rel.Label == "" {
		return connector, shapes.Shape{}, false, true
	}
	midX := (geometry.startX + geometry.endX) / 2
	midY := (geometry.startY + geometry.endY) / 2
	labelShape := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		midX-layout.labelWidth/2,
		midY-layout.labelHeight/2,
		layout.labelWidth,
		layout.labelHeight,
	).WithFill(shapes.NewShapeFill(theme.Background)).
		WithText(rel.Label).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(layout.labelMarginX, layout.labelMarginY, layout.labelMarginX, layout.labelMarginY)
	return connector, labelShape, true, true
}

func erConnectorPoints(fromPos, toPos erPosition, layout erLayout) erConnectorGeometry {
	if fromPos.y < toPos.y {
		return erConnectorGeometry{
			startX:    fromPos.x + layout.entityWidth/2,
			startY:    fromPos.y + layout.headerHeight,
			endX:      toPos.x + layout.entityWidth/2,
			endY:      toPos.y,
			startSite: shapes.ConnectionSiteBottom,
			endSite:   shapes.ConnectionSiteTop,
		}
	}
	return erConnectorGeometry{
		startX:    fromPos.x + layout.entityWidth/2,
		startY:    fromPos.y,
		endX:      toPos.x + layout.entityWidth/2,
		endY:      toPos.y + layout.headerHeight,
		startSite: shapes.ConnectionSiteTop,
		endSite:   shapes.ConnectionSiteBottom,
	}
}

func (b *erBounds) includeShape(s shapes.Shape) {
	b.include(s.X, s.Y, s.CX, s.CY)
}

func (b *erBounds) include(x, y, cx, cy styling.Length) {
	if b.empty {
		b.minX, b.minY = x, y
		b.maxX, b.maxY = x+cx, y+cy
		b.empty = false
		return
	}
	if x < b.minX {
		b.minX = x
	}
	if y < b.minY {
		b.minY = y
	}
	if x+cx > b.maxX {
		b.maxX = x + cx
	}
	if y+cy > b.maxY {
		b.maxY = y + cy
	}
}

func (s *erRenderState) diagramElements() DiagramElements {
	return DiagramElements{
		Shapes:     s.shapes,
		Connectors: s.connectors,
		Grouped:    true,
		Bounds: &DiagramBounds{
			X:  s.bounds.minX,
			Y:  s.bounds.minY,
			CX: s.bounds.maxX - s.bounds.minX,
			CY: s.bounds.maxY - s.bounds.minY,
		},
	}
}
