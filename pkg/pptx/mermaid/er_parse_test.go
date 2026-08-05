package mermaid

import "testing"

// Non-identifying `..` relations were absent from the hand-written token list,
// so the relation and any entity named only there were dropped without a word.
func TestERParsesBothNotations(t *testing.T) {
	tests := []struct {
		line    string
		relType string
	}{
		{"CUSTOMER ||--o{ ORDER : places", "||--o{"},
		{"CUSTOMER ||..o{ ORDER : places", "||..o{"},
		{"ORDER }|..|{ LINE_ITEM : contains", "}|..|{"},
		{"ORDER }|--|{ LINE_ITEM : contains", "}|--|{"},
		{"A |o..o| B : maybe", "|o..o|"},
		{"A ||--|| B : one", "||--||"},
		{"A -- B : plain", "--"},
	}

	for _, tt := range tests {
		from, relType, to, label, ok := splitERRelationship(tt.line)
		if !ok {
			t.Errorf("%q was not recognised as a relation", tt.line)
			continue
		}
		if relType != tt.relType {
			t.Errorf("%q gave type %q want %q", tt.line, relType, tt.relType)
		}
		if from == "" || to == "" {
			t.Errorf("%q gave from=%q to=%q", tt.line, from, to)
		}
		if label == "" {
			t.Errorf("%q lost its label", tt.line)
		}
	}
}

func TestERKeepsEntitiesFromDottedRelations(t *testing.T) {
	diagram := parseER(`erDiagram
    CUSTOMER ||..o{ ORDER : places`)

	if len(diagram.Relationships) != 1 {
		t.Fatalf("got %d relationships, want 1", len(diagram.Relationships))
	}
	names := make(map[string]bool, len(diagram.Entities))
	for _, e := range diagram.Entities {
		names[e.Name] = true
	}
	for _, want := range []string{"CUSTOMER", "ORDER"} {
		if !names[want] {
			t.Errorf("entity %q was dropped, got %v", want, names)
		}
	}
}

// A label may contain dots; they must not be mistaken for the connector.
func TestERLabelDotsAreNotAConnector(t *testing.T) {
	from, relType, to, label, ok := splitERRelationship("A ||--o{ B : version 1.2.3")
	if !ok {
		t.Fatal("relation not recognised")
	}
	if from != "A" || to != "B" || relType != "||--o{" {
		t.Errorf("parsed as from=%q type=%q to=%q", from, relType, to)
	}
	if label != "version 1.2.3" {
		t.Errorf("label=%q want %q", label, "version 1.2.3")
	}
}

func TestERRejectsNonRelationLines(t *testing.T) {
	for _, line := range []string{"erDiagram", "CUSTOMER {", "string name", ""} {
		if _, _, _, _, ok := splitERRelationship(line); ok {
			t.Errorf("%q was wrongly read as a relation", line)
		}
	}
}
