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
	// Events keeps every statement in source order, including the blocks,
	// notes and activations that Messages alone cannot express.
	Events []SequenceEvent
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
	var events []SequenceEvent
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
			events = append(events, SequenceEvent{Kind: seqEventMessage, Message: msg})
			continue
		}

		// A block, note or activation. Checked after messages so a message whose
		// text happens to start with a keyword is still read as a message.
		if event, ok := parseSequenceStatement(line); ok {
			events = append(events, event)
		}
	}

	return &SequenceDiagram{
		Participants: participants,
		Messages:     messages,
		Events:       events,
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

// detectSequenceArrow finds the message arrow in a line.
//
// The tokens are tried longest first so that "-->>" is never mistaken for the
// "-->" hiding inside it. Only ->> and -->> used to be recognised, so a diagram
// using -x, -) or a plain -> had those messages silently dropped.
func detectSequenceArrow(line string) (string, bool) {
	for _, token := range []string{"-->>", "--x", "--)", arrowSolid, "->>", "-x", "-)", "->"} {
		if strings.Contains(line, token) {
			return token, true
		}
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
		participantX[p.ID] = layout.startX + styling.Length(i)*layout.hSpacing
	}

	messagesTop := layout.startY + layout.participantHeight + styling.Inches(0.3)
	y := messagesTop

	script := newSequenceScript(diagram.Participants, participantX, layout, theme)
	for _, event := range diagram.Events {
		y = script.consume(event, y)
	}
	script.finish(y)

	// The lifelines can only be drawn once the script is laid out: a fixed
	// height left them stopping short on any diagram with blocks or notes, with
	// the closing participant boxes stranded in the middle of the messages.
	layout.lifelineHeight = max(
		y-(layout.startY+layout.participantHeight)+styling.Inches(0.3),
		styling.Inches(1.0),
	)
	for _, p := range diagram.Participants {
		for _, s := range sequenceParticipantShapes(p.DisplayName, participantX[p.ID], layout, theme) {
			shapesList = append(shapesList, s)
			bounds.includeShape(s)
		}
	}

	// Activation bars sit on the lifelines and so must paint before the arrows
	// crossing them; the block frames are outlines drawn over everything.
	for _, group := range [][]shapes.Shape{script.backdrops, script.messages, script.foreground} {
		for _, s := range group {
			shapesList = append(shapesList, s)
			bounds.includeShape(s)
		}
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

// sequenceRenderedMessage is everything one message draws, plus the vertical
// room it needs before the next message can start.
type sequenceRenderedMessage struct {
	shapes []shapes.Shape
	height styling.Length
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
	// A participant messaging itself has no horizontal span to draw an arrow
	// across. Sizing one from the distance gave it cx=0, which PowerPoint
	// rejects outright ("size must be > 0"), failing the whole deck rather than
	// just the diagram.
	if msg.From == msg.To {
		return sequenceSelfMessageShapes(msg, fromX, y, layout, theme), true
	}

	arrowX, arrowWidth, arrowType := sequenceArrowGeometry(fromX, toX, layout.participantWidth)
	arrow := shapes.NewShape(arrowType, arrowX, y, arrowWidth, sequenceArrowThickness).
		WithFill(shapes.NewShapeFill(theme.PrimaryStroke))
	textShape := sequenceMessageLabel(msg.Text, arrowX, y-styling.Inches(0.25), arrowWidth)
	return sequenceRenderedMessage{
		shapes: []shapes.Shape{arrow, textShape},
		height: layout.messageSpacing,
	}, true
}

// sequenceMessageLabel is the caption drawn with a message: text only, no box.
func sequenceMessageLabel(text string, x, y, width styling.Length) shapes.Shape {
	label := shapes.NewShape(shapes.ShapeTypeRectangle, x, y, width, styling.Inches(0.2)).
		WithText(text).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.05), styling.Inches(0.02), styling.Inches(0.05), styling.Inches(0.02))
	label.Line = nil
	label.Fill = nil
	return label
}
