package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Participant represents a participant in a sequence diagram.
type Participant struct {
	ID          string
	DisplayName string
}

// Message represents a message between participants in a sequence diagram.
type Message struct {
	From  string
	To    string
	Text  string
	Arrow string // ->> or -->>
}

// SequenceDiagram represents the parsed structure of a Mermaid sequence diagram.
type SequenceDiagram struct {
	Participants []Participant
	Messages     []Message
}

// renderSequence parses and renders a Mermaid sequence diagram into PowerPoint elements.
func renderSequence(code string, theme Theme) DiagramElements {
	diagram := parseSequence(code)
	return generateSequenceElements(diagram, theme)
}

func parseSequence(code string) *SequenceDiagram {
	lines := ParseLines(code)
	var participants []Participant
	var messages []Message
	participantMap := make(map[string]bool)

	addParticipant := func(id, displayName string) {
		if !participantMap[id] {
			participants = append(participants, Participant{ID: id, DisplayName: displayName})
			participantMap[id] = true
		}
	}

	// Skip header
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		if id, displayName, ok := parseParticipantLine(line); ok {
			addParticipant(id, displayName)
			continue
		}

		if msg, ok := parseMessageLine(line); ok {
			addParticipant(msg.From, msg.From)
			addParticipant(msg.To, msg.To)
			messages = append(messages, msg)
		}
	}

	return &SequenceDiagram{
		Participants: participants,
		Messages:     messages,
	}
}

func parseParticipantLine(line string) (string, string, bool) {
	after, ok := strings.CutPrefix(line, "participant")
	if !ok {
		return "", "", false
	}
	rest := strings.TrimSpace(after)
	if before, after, ok := strings.Cut(rest, " as "); ok {
		id := strings.TrimSpace(before)
		displayName := strings.TrimSpace(after)
		return id, displayName, true
	}
	fields := strings.Fields(rest)
	if len(fields) == 0 {
		return "", "", false
	}
	return fields[0], fields[0], true
}

func parseMessageLine(line string) (Message, bool) {
	arrow, ok := detectSequenceArrow(line)
	if !ok {
		return Message{}, false
	}
	fromPart, rest, ok := strings.Cut(line, arrow)
	if !ok {
		return Message{}, false
	}
	toPart, textPart, ok := strings.Cut(rest, ":")
	if !ok {
		return Message{}, false
	}
	return Message{
		From:  strings.TrimSpace(fromPart),
		To:    strings.TrimSpace(toPart),
		Text:  strings.TrimSpace(textPart),
		Arrow: arrow,
	}, true
}

func detectSequenceArrow(line string) (string, bool) {
	if strings.Contains(line, "-->>") {
		return "-->>", true
	}
	if strings.Contains(line, "->>") {
		return "->>", true
	}
	return "", false
}

func generateSequenceElements(diagram *SequenceDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape
	participantCount := len(diagram.Participants)

	if participantCount == 0 {
		return DiagramElements{Grouped: true}
	}

	layout := sequenceLayout{
		startX:            styling.Inches(0.5),
		startY:            styling.Inches(1.5),
		participantWidth:  styling.Inches(1.6),
		participantHeight: styling.Inches(0.6),
		hSpacing:          styling.Inches(2.2),
		lifelineHeight:    styling.Inches(4.0),
		messageSpacing:    styling.Inches(0.6),
	}

	participantX := make(map[string]styling.Length)
	bounds := newSequenceBounds()

	for i, p := range diagram.Participants {
		x := layout.startX + styling.Length(i)*layout.hSpacing
		participantX[p.ID] = x

		for _, s := range sequenceParticipantShapes(p.DisplayName, x, layout, theme) {
			shapesList = append(shapesList, s)
			bounds.includeShape(s)
		}
	}

	messageYStart := layout.startY + layout.participantHeight + styling.Inches(0.3)

	for i, msg := range diagram.Messages {
		y := messageYStart + styling.Length(i)*layout.messageSpacing
		rendered, ok := sequenceMessageShapes(msg, participantX, y, layout, theme)
		if !ok {
			continue
		}
		shapesList = append(shapesList, rendered.arrow, rendered.text)
		bounds.includeShape(rendered.arrow)
	}

	return DiagramElements{
		Shapes:  shapesList,
		Grouped: true,
		Bounds: &DiagramBounds{
			X:  bounds.minX,
			Y:  bounds.minY,
			CX: bounds.maxX - bounds.minX,
			CY: bounds.maxY - bounds.minY,
		},
	}
}

type sequenceLayout struct {
	startX            styling.Length
	startY            styling.Length
	participantWidth  styling.Length
	participantHeight styling.Length
	hSpacing          styling.Length
	lifelineHeight    styling.Length
	messageSpacing    styling.Length
}

type sequenceBounds struct {
	minX  styling.Length
	minY  styling.Length
	maxX  styling.Length
	maxY  styling.Length
	first bool
}

func newSequenceBounds() *sequenceBounds {
	return &sequenceBounds{first: true}
}

func (b *sequenceBounds) includeShape(s shapes.Shape) {
	b.include(s.X, s.Y, s.CX, s.CY)
}

func (b *sequenceBounds) include(x, y, cx, cy styling.Length) {
	if b.first {
		b.minX, b.minY = x, y
		b.maxX, b.maxY = x+cx, y+cy
		b.first = false
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

func sequenceParticipantShapes(name string, x styling.Length, layout sequenceLayout, theme Theme) []shapes.Shape {
	top := sequenceParticipantBox(name, x, layout.startY, layout, theme)
	lifeline := sequenceLifeline(x, layout, theme)
	bottomY := layout.startY + layout.participantHeight + layout.lifelineHeight
	bottom := sequenceParticipantBox(name, x, bottomY, layout, theme)
	return []shapes.Shape{top, lifeline, bottom}
}

func sequenceParticipantBox(
	name string,
	x styling.Length,
	y styling.Length,
	layout sequenceLayout,
	theme Theme,
) shapes.Shape {
	return shapes.NewShape(shapes.ShapeTypeRectangle, x, y, layout.participantWidth, layout.participantHeight).
		WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithText(name).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
}

func sequenceLifeline(x styling.Length, layout sequenceLayout, theme Theme) shapes.Shape {
	lifelineX := x + layout.participantWidth/2 - styling.Emu(10000)
	lifelineY := layout.startY + layout.participantHeight
	return shapes.NewShape(shapes.ShapeTypeRectangle, lifelineX, lifelineY, styling.Emu(20000), layout.lifelineHeight).
		WithFill(shapes.NewShapeFill(theme.SecondaryStroke))
}

type sequenceRenderedMessage struct {
	arrow shapes.Shape
	text  shapes.Shape
}

func sequenceMessageShapes(
	msg Message,
	participantX map[string]styling.Length,
	y styling.Length,
	layout sequenceLayout,
	theme Theme,
) (sequenceRenderedMessage, bool) {
	fromX, fromExists := participantX[msg.From]
	toX, toExists := participantX[msg.To]
	if !fromExists || !toExists {
		return sequenceRenderedMessage{}, false
	}

	arrowX, arrowWidth, arrowType := sequenceArrowGeometry(fromX, toX, layout.participantWidth)
	arrow := shapes.NewShape(arrowType, arrowX, y, arrowWidth, styling.Inches(0.15)).
		WithFill(shapes.NewShapeFill(theme.PrimaryStroke))
	textShape := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		arrowX,
		y-styling.Inches(0.25),
		arrowWidth,
		styling.Inches(0.2),
	).WithText(msg.Text).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.05), styling.Inches(0.02), styling.Inches(0.05), styling.Inches(0.02))
	textShape.Line = nil
	textShape.Fill = nil
	return sequenceRenderedMessage{arrow: arrow, text: textShape}, true
}

func sequenceArrowGeometry(
	fromX styling.Length,
	toX styling.Length,
	participantWidth styling.Length,
) (styling.Length, styling.Length, string) {
	fromCenter := fromX + participantWidth/2
	toCenter := toX + participantWidth/2
	if fromCenter < toCenter {
		return fromCenter, toCenter - fromCenter, shapes.ShapeTypeRightArrow
	}
	return toCenter, fromCenter - toCenter, shapes.ShapeTypeLeftArrow
}
