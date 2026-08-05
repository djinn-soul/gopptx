package mermaid

import "testing"

// Only an unlabelled <|-- became a connector. Every other relation was turned
// into a class box holding the raw relation text, because the inline-member
// check ran first and a labelled relation also contains a colon.
func TestClassParsesEveryRelationForm(t *testing.T) {
	tests := []struct {
		line    string
		from    string
		relType string
		to      string
		label   string
	}{
		{"Animal <|-- Duck", "Animal", "<|--", "Duck", ""},
		{"Animal <|-- Duck : extends", "Animal", "<|--", "Duck", "extends"},
		{"Animal <|.. Duck", "Animal", "<|..", "Duck", ""},
		{"Order *-- LineItem", "Order", "*--", "LineItem", ""},
		{"Order o-- Customer", "Order", "o--", "Customer", ""},
		{"Shape ..> Drawable", "Shape", "..>", "Drawable", ""},
		{"Shape --> Canvas : draws", "Shape", "-->", "Canvas", "draws"},
		{`Customer "1" --> "*" Order`, "Customer", "-->", "Order", ""},
	}

	for _, tt := range tests {
		parts, ok := splitClassRelationship(tt.line)
		if !ok {
			t.Errorf("%q was not recognised as a relation", tt.line)
			continue
		}
		if parts.from != tt.from || parts.to != tt.to || parts.relType != tt.relType || parts.label != tt.label {
			t.Errorf("%q parsed as from=%q type=%q to=%q label=%q, want %q %q %q %q",
				tt.line, parts.from, parts.relType, parts.to, parts.label,
				tt.from, tt.relType, tt.to, tt.label)
		}
	}
}

// A labelled relation must become a relationship, not a class named after the
// whole line.
func TestClassLabelledRelationIsNotAMember(t *testing.T) {
	diagram := parseClass(`classDiagram
    Animal <|-- Duck : extends`)

	if len(diagram.Relationships) != 1 {
		t.Fatalf("got %d relationships, want 1", len(diagram.Relationships))
	}
	for _, c := range diagram.Classes {
		if c.Name != "Animal" && c.Name != "Duck" {
			t.Errorf("unexpected class %q built from the relation line", c.Name)
		}
		if len(c.Attributes) > 0 || len(c.Methods) > 0 {
			t.Errorf("class %q picked up members %v %v", c.Name, c.Attributes, c.Methods)
		}
	}
}

// An `o` ending a class name must not be mistaken for the aggregation marker.
func TestClassNameEndingInOSurvives(t *testing.T) {
	parts, ok := splitClassRelationship("Foo--Bar")
	if !ok {
		t.Fatal("relation not recognised")
	}
	if parts.from != "Foo" || parts.to != "Bar" || parts.relType != "--" {
		t.Errorf("parsed as from=%q type=%q to=%q", parts.from, parts.relType, parts.to)
	}
}

// Real members must still parse as members.
func TestClassMembersStillParse(t *testing.T) {
	diagram := parseClass(`classDiagram
    class Duck {
        +String beakColor
        +swim()
    }`)

	if len(diagram.Classes) != 1 {
		t.Fatalf("got %d classes, want 1", len(diagram.Classes))
	}
	duck := diagram.Classes[0]
	if len(duck.Attributes) != 1 || len(duck.Methods) != 1 {
		t.Errorf("duck has attributes %v and methods %v", duck.Attributes, duck.Methods)
	}
}
