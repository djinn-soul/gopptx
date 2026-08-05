package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Geometry of the framed blocks, notes and activation bars.
const (
	seqFramePadX      styling.Length = 274320 // 0.3" of air either side of the lifelines
	seqFrameHeaderY   styling.Length = 365760 // 0.4" of room under a frame's caption tab
	seqFrameTabWidth  styling.Length = 640080 // 0.7"
	seqFrameTabHeight styling.Length = 219456 // 0.24"
	seqBranchGap      styling.Length = 228600 // 0.25" for an else/and divider
	seqNoteHeight     styling.Length = 365760 // 0.4"
	seqNoteGap        styling.Length = 91440  // 0.1"
	seqActivationHalf styling.Length = 45720  // half the 0.1" activation bar
	seqFrameEndPad    styling.Length = 137160 // 0.15" below the last message

	// Notes keep Mermaid's pale yellow, matching the state diagram's notes.
	seqNoteFill   = "FFF5AD"
	seqNoteStroke = "D6B656"
)

// sequenceOpenBlock is a block whose "end" has not been read yet.
type sequenceOpenBlock struct {
	kind     string
	label    string
	startY   styling.Length
	depth    int
	dividers []sequenceDivider
}

// sequenceDivider is an "else"/"and" branch line inside an open block.
type sequenceDivider struct {
	y     styling.Length
	label string
}

// sequenceScript turns the ordered statements of a diagram into shapes,
// carrying the state a statement needs from the ones before it: which blocks
// are open, and which participants are currently activated.
type sequenceScript struct {
	participants []Participant
	participantX map[string]styling.Length
	layout       sequenceLayout
	theme        Theme

	open        []sequenceOpenBlock
	activations map[string]styling.Length

	backdrops  []shapes.Shape
	messages   []shapes.Shape
	foreground []shapes.Shape
}

func newSequenceScript(
	participants []Participant,
	participantX map[string]styling.Length,
	layout sequenceLayout,
	theme Theme,
) *sequenceScript {
	return &sequenceScript{
		participants: participants,
		participantX: participantX,
		layout:       layout,
		theme:        theme,
		activations:  make(map[string]styling.Length),
	}
}

// consume draws one statement and returns the y the next one starts at.
func (s *sequenceScript) consume(event SequenceEvent, y styling.Length) styling.Length {
	switch event.Kind {
	case seqEventMessage:
		return s.consumeMessage(event.Message, y)
	case seqEventBlockStart:
		s.open = append(s.open, sequenceOpenBlock{
			kind:   event.BlockKind,
			label:  event.Label,
			startY: y,
			depth:  len(s.open),
		})
		return y + seqFrameHeaderY
	case seqEventBlockElse:
		s.addDivider(y, event.Label)
		return y + seqBranchGap
	case seqEventBlockEnd:
		return s.closeBlock(y)
	case seqEventNote:
		return s.consumeNote(event.Note, y)
	case seqEventActivate:
		s.activations[s.resolveParticipant(event.Participant)] = y
		return y
	case seqEventDeactivate:
		s.closeActivation(s.resolveParticipant(event.Participant), y)
		return y
	default:
		return y
	}
}

// finish closes anything the diagram left open, so a missing "end" or
// "deactivate" still draws rather than dropping the frame entirely.
func (s *sequenceScript) finish(y styling.Length) {
	for len(s.open) > 0 {
		y = s.closeBlock(y)
	}
	for participant := range s.activations {
		s.closeActivation(participant, y)
	}
}

func (s *sequenceScript) consumeMessage(msg Message, y styling.Length) styling.Length {
	rendered, ok := sequenceMessageShapes(msg, s.participantX, y, s.layout, s.theme)
	if !ok {
		return y
	}
	s.messages = append(s.messages, rendered.shapes...)
	return y + rendered.height
}

func (s *sequenceScript) addDivider(y styling.Length, label string) {
	if len(s.open) == 0 {
		return
	}
	current := &s.open[len(s.open)-1]
	current.dividers = append(current.dividers, sequenceDivider{y: y, label: label})
}

// closeBlock emits the frame for the innermost open block.
func (s *sequenceScript) closeBlock(y styling.Length) styling.Length {
	if len(s.open) == 0 {
		return y
	}
	block := s.open[len(s.open)-1]
	s.open = s.open[:len(s.open)-1]

	left, right := s.frameSpan(block.depth)
	bottom := y + seqFrameEndPad

	frame := shapes.NewShape(shapes.ShapeTypeRectangle, left, block.startY, right-left, bottom-block.startY).
		WithLine(shapes.NewShapeLine(s.theme.SecondaryStroke, s.theme.LineWeight))
	frame.Fill = nil
	s.foreground = append(s.foreground, frame)

	s.foreground = append(s.foreground, s.frameTab(block, left))
	for _, divider := range block.dividers {
		s.foreground = append(s.foreground, s.dividerShapes(divider, left, right)...)
	}
	return bottom + seqNoteGap
}

// frameTab is the caption in the frame's top-left corner: "loop", "alt" and so
// on, followed by the condition the diagram gave.
func (s *sequenceScript) frameTab(block sequenceOpenBlock, left styling.Length) shapes.Shape {
	caption := block.kind
	if block.label != "" {
		caption += " [" + block.label + "]"
	}
	width := seqFrameTabWidth + styling.Length(len(caption))*styling.Emu(45720)
	return shapes.NewShape(shapes.ShapeTypeRectangle, left, block.startY, width, seqFrameTabHeight).
		WithFill(shapes.NewShapeFill(s.theme.SecondaryFill)).
		WithLine(shapes.NewShapeLine(s.theme.SecondaryStroke, s.theme.LineWeight)).
		WithText(caption).
		WithVerticalAnchor(shapes.TextAnchorMiddle).
		WithTextMargins(styling.Inches(0.05), styling.Inches(0.01), styling.Inches(0.05), styling.Inches(0.01))
}

// dividerShapes draw the dashed branch line of an "else"/"and", with its label.
func (s *sequenceScript) dividerShapes(
	divider sequenceDivider,
	left, right styling.Length,
) []shapes.Shape {
	line := shapes.NewShape(shapes.ShapeTypeRectangle, left, divider.y, right-left, styling.Emu(12700)).
		WithFill(shapes.NewShapeFill(s.theme.SecondaryStroke))
	out := []shapes.Shape{line}
	if divider.label == "" {
		return out
	}
	label := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		left+styling.Inches(0.05),
		divider.y,
		seqFrameTabWidth+styling.Length(len(divider.label))*styling.Emu(45720),
		seqFrameTabHeight,
	).WithText("["+divider.label+"]").
		WithVerticalAnchor(shapes.TextAnchorMiddle).
		WithTextMargins(styling.Inches(0.05), styling.Inches(0.01), styling.Inches(0.05), styling.Inches(0.01))
	label.Fill = nil
	label.Line = nil
	return append(out, label)
}

// frameSpan is the horizontal extent of a frame. Nested blocks are inset so
// their outlines do not sit on top of each other.
func (s *sequenceScript) frameSpan(depth int) (styling.Length, styling.Length) {
	inset := styling.Length(depth) * styling.Inches(0.12)
	left := s.layout.startX
	right := s.layout.startX + s.layout.participantWidth
	for _, participant := range s.participants {
		x, ok := s.participantX[participant.ID]
		if !ok {
			continue
		}
		left = min(left, x)
		right = max(right, x+s.layout.participantWidth)
	}
	return left - seqFramePadX + inset, right + seqFramePadX - inset
}

func (s *sequenceScript) consumeNote(note SequenceNote, y styling.Length) styling.Length {
	left, right, ok := s.noteSpan(note)
	if !ok {
		return y
	}
	shape := shapes.NewShape(shapes.ShapeTypeRectangle, left, y, right-left, seqNoteHeight).
		WithFill(shapes.NewShapeFill(seqNoteFill)).
		WithLine(shapes.NewShapeLine(seqNoteStroke, s.theme.LineWeight)).
		WithText(note.Text).
		WithVerticalAnchor(shapes.TextAnchorMiddle).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.08), styling.Inches(0.03), styling.Inches(0.08), styling.Inches(0.03))
	s.messages = append(s.messages, shape)
	return y + seqNoteHeight + seqNoteGap
}

// noteSpan places a note: over its participants' lifelines, or beside one.
func (s *sequenceScript) noteSpan(note SequenceNote) (styling.Length, styling.Length, bool) {
	var left, right styling.Length
	found := false
	for _, name := range note.Participants {
		x, ok := s.participantX[s.resolveParticipant(name)]
		if !ok {
			continue
		}
		if !found {
			left, right, found = x, x+s.layout.participantWidth, true
			continue
		}
		left = min(left, x)
		right = max(right, x+s.layout.participantWidth)
	}
	if !found {
		return 0, 0, false
	}

	switch note.Placement {
	case notePlacementLeft:
		return left - s.layout.participantWidth - seqNoteGap, left - seqNoteGap, true
	case notePlacementRight:
		return right + seqNoteGap, right + s.layout.participantWidth + seqNoteGap, true
	default:
		return left, right, true
	}
}

// closeActivation draws the bar that has been running on a participant's
// lifeline since its "activate".
func (s *sequenceScript) closeActivation(participant string, y styling.Length) {
	startY, active := s.activations[participant]
	if !active {
		return
	}
	delete(s.activations, participant)

	x, ok := s.participantX[participant]
	if !ok {
		return
	}
	height := y - startY
	if height <= 0 {
		height = seqNoteGap
	}
	bar := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		x+s.layout.participantWidth/2-seqActivationHalf,
		startY,
		seqActivationHalf*2,
		height,
	).WithFill(shapes.NewShapeFill(s.theme.PrimaryFill)).
		WithLine(shapes.NewShapeLine(s.theme.PrimaryStroke, s.theme.LineWeight))
	s.backdrops = append(s.backdrops, bar)
}

// resolveParticipant maps a name onto a participant id. A note or activation
// may name a participant by its display name where the diagram declared an
// alias with "participant A as Alice".
func (s *sequenceScript) resolveParticipant(name string) string {
	for _, participant := range s.participants {
		if strings.EqualFold(participant.ID, name) || strings.EqualFold(participant.DisplayName, name) {
			return participant.ID
		}
	}
	return name
}
