package mermaid

import "testing"

// A self-message used to size its arrow from the distance between sender and
// receiver. For `Bob->>Bob` that distance is zero, so the arrow got cx=0 and
// PowerPoint refused the whole deck with "size must be > 0".
func TestSelfMessageProducesNoZeroSizedShape(t *testing.T) {
	code := `sequenceDiagram
    participant Bob
    Bob->>Bob: poll`

	elements := renderSequence(code, GetTheme("default"))
	if len(elements.Shapes) == 0 {
		t.Fatal("no shapes rendered")
	}
	for i, s := range elements.Shapes {
		if s.CX <= 0 || s.CY <= 0 {
			t.Errorf("shape %d (%s) has non-positive size cx=%d cy=%d", i, s.Type, s.CX, s.CY)
		}
	}
}

func TestSelfMessageKeepsItsLabel(t *testing.T) {
	code := `sequenceDiagram
    participant Bob
    Bob->>Bob: poll`

	elements := renderSequence(code, GetTheme("default"))
	found := false
	for _, s := range elements.Shapes {
		if s.Text == "poll" {
			found = true
		}
	}
	if !found {
		t.Fatal(`the self-message label "poll" was not drawn`)
	}
}

// Every arrow Mermaid accepts must produce a message. Only ->> and -->> were
// recognised, so -x, -) and a plain -> were dropped without a word.
func TestSequenceParsesEveryArrowToken(t *testing.T) {
	for _, arrow := range []string{"->", "->>", "-->", "-->>", "-x", "--x", "-)", "--)"} {
		code := "sequenceDiagram\n    Alice" + arrow + "Bob: hello"
		diagram := parseSequence(code)
		if len(diagram.Messages) != 1 {
			t.Errorf("arrow %q produced %d messages, want 1", arrow, len(diagram.Messages))
			continue
		}
		msg := diagram.Messages[0]
		if msg.From != "Alice" || msg.To != "Bob" || msg.Text != "hello" {
			t.Errorf("arrow %q parsed as %+v", arrow, msg)
		}
		if msg.Arrow != arrow {
			t.Errorf("arrow %q recorded as %q", arrow, msg.Arrow)
		}
	}
}

// Messages must not overlap: a self-message needs a loop's worth of room, so
// the message after one has to start below the loop, not at a fixed stride.
func TestSelfMessageAdvancesTheNextMessage(t *testing.T) {
	code := `sequenceDiagram
    participant Alice
    participant Bob
    Alice->>Bob: first
    Bob->>Bob: retry
    Bob->>Alice: last`

	elements := renderSequence(code, GetTheme("default"))
	var firstY, lastY int64
	for _, s := range elements.Shapes {
		switch s.Text {
		case "first":
			firstY = int64(s.Y)
		case "last":
			lastY = int64(s.Y)
		}
	}
	if firstY == 0 || lastY == 0 {
		t.Fatal("expected both plain messages to be drawn")
	}
	if lastY <= firstY {
		t.Errorf("last message at y=%d is not below the first at y=%d", lastY, firstY)
	}
}
