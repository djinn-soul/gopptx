package mermaid

import "strings"

// Sequence event kinds. A sequence diagram is a script rather than a graph, so
// the parser keeps the statements in order instead of collecting messages only:
// a loop or a note means nothing without knowing which messages it wraps.
const (
	seqEventMessage    = "message"
	seqEventBlockStart = "block_start"
	seqEventBlockElse  = "block_else"
	seqEventBlockEnd   = "block_end"
	seqEventNote       = "note"
	seqEventActivate   = "activate"
	seqEventDeactivate = "deactivate"
)

// Note placements Mermaid understands.
const (
	notePlacementOver  = "over"
	notePlacementLeft  = "left"
	notePlacementRight = "right"
)

// SequenceNote is a "Note over A,B: text" annotation.
type SequenceNote struct {
	Placement    string
	Participants []string
	Text         string
}

// SequenceEvent is one statement of a sequence diagram, in source order.
type SequenceEvent struct {
	Kind string
	// Message is set for Kind == message.
	Message Message
	// BlockKind is "loop", "alt", "opt", "par", "critical", "break" or "rect".
	BlockKind string
	// Label is the block's caption, or the "else" branch caption.
	Label string
	// Note is set for Kind == note.
	Note SequenceNote
	// Participant is the target of an activate/deactivate.
	Participant string
}

// sequenceBlockKeywords are the statements that open a framed block.
func sequenceBlockKeywords() []string {
	return []string{"loop", "alt", "opt", "par", "critical", "break", "rect"}
}

// parseSequenceStatement classifies a non-message line. It returns false for
// lines that carry no event, such as "end note" or an unknown keyword.
func parseSequenceStatement(line string) (SequenceEvent, bool) {
	lower := strings.ToLower(line)

	switch {
	case lower == "end":
		return SequenceEvent{Kind: seqEventBlockEnd}, true
	case lower == "else" || strings.HasPrefix(lower, "else "):
		return SequenceEvent{Kind: seqEventBlockElse, Label: strings.TrimSpace(line[len("else"):])}, true
	case lower == "and" || strings.HasPrefix(lower, "and "):
		// A par branch divider reads the same way as an alt's "else".
		return SequenceEvent{Kind: seqEventBlockElse, Label: strings.TrimSpace(line[len("and"):])}, true
	}

	if after, ok := cutSequenceKeyword(line, lower, "activate"); ok {
		return SequenceEvent{Kind: seqEventActivate, Participant: after}, true
	}
	if after, ok := cutSequenceKeyword(line, lower, "deactivate"); ok {
		return SequenceEvent{Kind: seqEventDeactivate, Participant: after}, true
	}

	if note, ok := parseSequenceNote(line, lower); ok {
		return SequenceEvent{Kind: seqEventNote, Note: note}, true
	}

	for _, keyword := range sequenceBlockKeywords() {
		if after, ok := cutSequenceKeyword(line, lower, keyword); ok {
			return SequenceEvent{Kind: seqEventBlockStart, BlockKind: keyword, Label: after}, true
		}
	}
	return SequenceEvent{}, false
}

// cutSequenceKeyword matches a leading keyword, either alone or followed by a
// space, and returns the rest of the line with its original casing.
func cutSequenceKeyword(line, lower, keyword string) (string, bool) {
	if lower == keyword {
		return "", true
	}
	if !strings.HasPrefix(lower, keyword+" ") {
		return "", false
	}
	return strings.TrimSpace(line[len(keyword):]), true
}

// parseSequenceNote reads "Note over A,B: text" and its left/right of forms.
func parseSequenceNote(line, lower string) (SequenceNote, bool) {
	rest, ok := cutSequenceKeyword(line, lower, "note")
	if !ok || rest == "" {
		return SequenceNote{}, false
	}

	target, text, found := strings.Cut(rest, ":")
	if !found {
		return SequenceNote{}, false
	}

	placement, participantList := sequenceNotePlacement(strings.TrimSpace(target))
	if participantList == "" {
		return SequenceNote{}, false
	}

	participants := make([]string, 0, 2)
	for part := range strings.SplitSeq(participantList, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			participants = append(participants, trimmed)
		}
	}
	if len(participants) == 0 {
		return SequenceNote{}, false
	}

	return SequenceNote{
		Placement:    placement,
		Participants: participants,
		Text:         strings.TrimSpace(text),
	}, true
}

// sequenceNotePlacement splits "over A,B", "left of A" or "right of A" into the
// placement and the participant list.
func sequenceNotePlacement(target string) (string, string) {
	lower := strings.ToLower(target)
	switch {
	case strings.HasPrefix(lower, "over "):
		return notePlacementOver, strings.TrimSpace(target[len("over"):])
	case strings.HasPrefix(lower, "left of "):
		return notePlacementLeft, strings.TrimSpace(target[len("left of"):])
	case strings.HasPrefix(lower, "right of "):
		return notePlacementRight, strings.TrimSpace(target[len("right of"):])
	default:
		return notePlacementOver, target
	}
}
