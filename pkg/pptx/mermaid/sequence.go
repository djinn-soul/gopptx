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
		line := lines[i]

		if after, ok := strings.CutPrefix(line, "participant"); ok {
			rest := strings.TrimSpace(after)
			if before, after, ok := strings.Cut(rest, " as "); ok {
				id := strings.TrimSpace(before)
				displayName := strings.TrimSpace(after)
				addParticipant(id, displayName)
			} else {
				id := strings.Fields(rest)[0]
				addParticipant(id, id)
			}
		} else if strings.Contains(line, "->>") || strings.Contains(line, "-->>") {
			arrow := "->>"
			if strings.Contains(line, "-->>") {
				arrow = "-->>"
			}

			parts := strings.SplitN(line, arrow, 2)
			if len(parts) == 2 {
				from := strings.TrimSpace(parts[0])
				rest := parts[1]
				msgParts := strings.SplitN(rest, ":", 2)
				if len(msgParts) == 2 {
					to := strings.TrimSpace(msgParts[0])
					text := strings.TrimSpace(msgParts[1])

					addParticipant(from, from)
					addParticipant(to, to)

					messages = append(messages, Message{
						From:  from,
						To:    to,
						Text:  text,
						Arrow: arrow,
					})
				}
			}
		}
	}

	return &SequenceDiagram{
		Participants: participants,
		Messages:     messages,
	}
}

func generateSequenceElements(diagram *SequenceDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape
	participantCount := len(diagram.Participants)

	if participantCount == 0 {
		return DiagramElements{Grouped: true}
	}

	// Layout parameters
	startX := styling.Inches(0.5)
	startY := styling.Inches(1.5)
	participantWidth := styling.Inches(1.6)
	participantHeight := styling.Inches(0.6)
	hSpacing := styling.Inches(2.2)
	lifelineHeight := styling.Inches(4.0)
	messageSpacing := styling.Inches(0.6)

	participantX := make(map[string]styling.Length)

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

	for i, p := range diagram.Participants {
		x := startX + styling.Length(i)*hSpacing
		participantX[p.ID] = x

		// Participant box at top
		boxTop := shapes.NewShape(shapes.ShapeTypeRectangle, x, startY, participantWidth, participantHeight).
			WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
			WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
			WithText(p.DisplayName).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
		shapesList = append(shapesList, boxTop)
		updateBounds(x, startY, participantWidth, participantHeight)

		// Lifeline
		lifelineX := x + participantWidth/2 - styling.Emu(10000)
		lifelineY := startY + participantHeight
		lifeline := shapes.NewShape(shapes.ShapeTypeRectangle, lifelineX, lifelineY, styling.Emu(20000), lifelineHeight).
			WithFill(shapes.NewShapeFill(theme.SecondaryStroke))
		shapesList = append(shapesList, lifeline)
		updateBounds(lifelineX, lifelineY, styling.Emu(20000), lifelineHeight)

		// Participant box at bottom
		boxBottom := shapes.NewShape(shapes.ShapeTypeRectangle, x, startY+participantHeight+lifelineHeight, participantWidth, participantHeight).
			WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
			WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
			WithText(p.DisplayName).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
		shapesList = append(shapesList, boxBottom)
		updateBounds(x, startY+participantHeight+lifelineHeight, participantWidth, participantHeight)
	}

	// Message arrows
	messageYStart := startY + participantHeight + styling.Inches(0.3)

	for i, msg := range diagram.Messages {
		fromX, fromExists := participantX[msg.From]
		toX, toExists := participantX[msg.To]

		if fromExists && toExists {
			y := messageYStart + styling.Length(i)*messageSpacing
			fromCenter := fromX + participantWidth/2
			toCenter := toX + participantWidth/2

			var arrowX, arrowWidth styling.Length
			var arrowType string
			if fromCenter < toCenter {
				arrowX = fromCenter
				arrowWidth = toCenter - fromCenter
				arrowType = shapes.ShapeTypeRightArrow
			} else {
				arrowX = toCenter
				arrowWidth = fromCenter - toCenter
				arrowType = shapes.ShapeTypeLeftArrow
			}

			arrow := shapes.NewShape(arrowType, arrowX, y, arrowWidth, styling.Inches(0.15)).
				WithFill(shapes.NewShapeFill(theme.PrimaryStroke))
			shapesList = append(shapesList, arrow)
			updateBounds(arrowX, y, arrowWidth, styling.Inches(0.15))

			textShape := shapes.NewShape(shapes.ShapeTypeRectangle, arrowX, y-styling.Inches(0.25), arrowWidth, styling.Inches(0.2)).
				WithText(msg.Text).
				WithAutoFit(shapes.TextAutoFitNormal).
				WithTextMargins(styling.Inches(0.05), styling.Inches(0.02), styling.Inches(0.05), styling.Inches(0.02))
			// Make text background transparent or no line
			textShape.Line = nil
			textShape.Fill = nil
			shapesList = append(shapesList, textShape)
		}
	}

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
