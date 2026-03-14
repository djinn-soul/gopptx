package mermaid

import (
	"testing"
)

func TestDetectType(t *testing.T) {
	tests := []struct {
		code     string
		expected Type
	}{
		{"flowchart LR", Flowchart},
		{"graph TD", Flowchart},
		{"sequenceDiagram", Sequence},
		{"pie", Pie},
		{"gantt", Gantt},
		{"classDiagram", Class},
		{"stateDiagram", State},
		{"erDiagram", ER},
		{"mindmap", Mindmap},
		{"timeline", Timeline},
		{"journey", Journey},
		{"quadrantChart", Quadrant},
		{"gitGraph", GitGraph},
		{"unknown", Unknown},
		{"%% comment\nflowchart LR", Flowchart},
	}

	for _, tt := range tests {
		if got := DetectType(tt.code); got != tt.expected {
			t.Errorf("DetectType(%q) = %v, want %v", tt.code, got, tt.expected)
		}
	}
}

func TestParseNodeDef(t *testing.T) {
	tests := []struct {
		input         string
		expectedID    string
		expectedLabel string
		expectedShape NodeShape
	}{
		{"A[Rectangle]", "A", "Rectangle", NodeShapeRectangle},
		{"B(Rounded)", "B", "Rounded", NodeShapeRoundedRect},
		{"C{Diamond}", "C", "Diamond", NodeShapeDiamond},
		{"D((Circle))", "D", "Circle", NodeShapeCircle},
		{"E([Stadium])", "E", "Stadium", NodeShapeStadium},
		{"F{{Hexagon}}", "F", "Hexagon", NodeShapeHexagon},
		{"Plain", "Plain", "Plain", NodeShapeRectangle},
	}

	for _, tt := range tests {
		id, label, shape := ParseNodeDef(tt.input)
		if id != tt.expectedID || label != tt.expectedLabel || shape != tt.expectedShape {
			t.Errorf("ParseNodeDef(%q) = (%q, %q, %v), want (%q, %q, %v)",
				tt.input, id, label, shape, tt.expectedID, tt.expectedLabel, tt.expectedShape)
		}
	}
}

func TestCreateDiagram(t *testing.T) {
	tests := []diagramCase{
		{
			name: "Flowchart",
			code: `flowchart LR
A[Start] --> B{Decision}
B -- Yes --> C[End]
B -- No --> D[Wait]`,
			minShapes:       4,
			exactConnectors: 3,
			checkConnectors: true,
			requireBounds:   true,
		},
		{
			name: "Sequence",
			code: `sequenceDiagram
Alice->>Bob: Hello Bob, how are you?
Bob-->>Alice: Jolly good!`,
			minShapes:     6,
			requireBounds: true,
		},
		{
			name: "Pie",
			code: `pie title Pets adopted by volunteers
"Dogs" : 386
"Cats" : 85
"Rats" : 15`,
			minShapes:     5,
			requireBounds: true,
		},
		{
			name: "Gantt",
			code: `gantt
title A Gantt Diagram
section Section
A task :a1, 2014-01-01, 30d
Another task :after a1, 20d`,
			minShapes:     4,
			requireBounds: true,
		},
		{
			name: "Timeline",
			code: `timeline
title History of Social Media Platform
2002 : LinkedIn
2004 : Facebook
: Google`,
			minShapes:     5,
			requireBounds: true,
		},
		{
			name: "Quadrant",
			code: `quadrantChart
title Reach and engagement of campaigns
x-axis Low Reach --> High Reach
y-axis Low Engagement --> High Engagement
quadrant-1 We should expand
Campaign A: [0.3, 0.6]`,
			minShapes:     8,
			requireBounds: true,
		},
		{
			name: "Class",
			code: `classDiagram
class Animal {
    +String name
    +isMammal()
}
class Dog {
    +bark()
}
Animal <|-- Dog`,
			minShapes:       6,
			exactConnectors: 1,
			checkConnectors: true,
		},
		{
			name: "State",
			code: `stateDiagram-v2
[*] --> First
First --> Second
Second --> [*]`,
			minShapes:       3,
			exactConnectors: 3,
			checkConnectors: true,
		},
		{
			name: "ER",
			code: `erDiagram
CUSTOMER ||--o{ ORDER : places
CUSTOMER {
    string name
    string email
}
ORDER {
    int orderNumber
}`,
			minShapes:       4,
			exactConnectors: 1,
			checkConnectors: true,
		},
		{
			name: "Mindmap",
			code: `mindmap
root((mindmap))
    Origins
        Long history
    Research
        On effectiveness`,
			exactShapes:     5,
			exactConnectors: 4,
			checkConnectors: true,
		},
		{
			name: "Journey",
			code: `journey
title My working day
section Go to work
    Make tea: 5: Me
    Do work: 1: Me, Cat`,
			exactShapes: 4,
		},
		{
			name: "GitGraph",
			code: `gitGraph
commit
commit
branch develop
checkout develop
commit
checkout main
merge develop`,
			exactShapes:     8,
			exactConnectors: 3,
			checkConnectors: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertDiagramCase(t, tc)
		})
	}
}

type diagramCase struct {
	name            string
	code            string
	minShapes       int
	exactShapes     int
	exactConnectors int
	checkConnectors bool
	requireBounds   bool
}

func assertDiagramCase(t *testing.T, tc diagramCase) {
	t.Helper()

	elements, err := CreateDiagram(tc.code)
	if err != nil {
		t.Fatalf("CreateDiagram failed: %v", err)
	}

	if tc.minShapes > 0 && len(elements.Shapes) < tc.minShapes {
		t.Errorf("Expected at least %d shapes, got %d", tc.minShapes, len(elements.Shapes))
	}
	if tc.exactShapes > 0 && len(elements.Shapes) != tc.exactShapes {
		t.Errorf("Expected %d shapes, got %d", tc.exactShapes, len(elements.Shapes))
	}
	if tc.checkConnectors && len(elements.Connectors) != tc.exactConnectors {
		t.Errorf("Expected %d connectors, got %d", tc.exactConnectors, len(elements.Connectors))
	}
	if tc.requireBounds && elements.Bounds == nil {
		t.Error("Expected bounds to be set")
	}
}

func TestFlowchartParsesSpacedEdgeLabelsWithoutGhostNodes(t *testing.T) {
	diagram := parseFlowchart(`flowchart LR
A[Start] --> B{Decision}
B -- Yes --> C[Ship]
B -- No --> D[Revise]`)

	if len(diagram.Connections) != 3 {
		t.Fatalf("expected 3 connections, got %d", len(diagram.Connections))
	}

	if diagram.Connections[1].Label != "Yes" {
		t.Fatalf("expected Yes label, got %q", diagram.Connections[1].Label)
	}
	if diagram.Connections[2].Label != "No" {
		t.Fatalf("expected No label, got %q", diagram.Connections[2].Label)
	}

	for _, n := range diagram.Nodes {
		if n.ID == "B -- Yes" || n.ID == "B -- No" {
			t.Fatalf("unexpected ghost node id parsed from edge label: %q", n.ID)
		}
	}
}

func TestMindmapTwoChildrenUsesHorizontalBranches(t *testing.T) {
	diagram, err := CreateDiagram(`mindmap
root((mindmap))
    Left
    Right`)
	if err != nil {
		t.Fatalf("CreateDiagram failed: %v", err)
	}

	var rootX, leftX, rightX int64
	for _, s := range diagram.Shapes {
		switch s.Text {
		case "mindmap":
			rootX = int64(s.X)
		case "Left":
			leftX = int64(s.X)
		case "Right":
			rightX = int64(s.X)
		}
	}

	if leftX == 0 || rightX == 0 || rootX == 0 {
		t.Fatal("expected root and child nodes to be present")
	}
	if (leftX > rootX && rightX > rootX) || (leftX < rootX && rightX < rootX) {
		t.Fatalf("expected mindmap children on opposite sides of root; root=%d left=%d right=%d", rootX, leftX, rightX)
	}
}
