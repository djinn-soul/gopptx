package mermaid

import (
	"testing"
)

func TestDetectType(t *testing.T) {
	tests := []struct {
		code     string
		expected MermaidType
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
	t.Run("Flowchart", func(t *testing.T) {
		code := `flowchart LR
A[Start] --> B{Decision}
B -- Yes --> C[End]
B -- No --> D[Wait]`
		elements, err := CreateDiagram(code)
		if err != nil {
			t.Fatalf("CreateDiagram failed: %v", err)
		}
		// 4 nodes + 3 labels = 7 shapes
		if len(elements.Shapes) < 4 {
			t.Errorf("Expected at least 4 shapes, got %d", len(elements.Shapes))
		}
		if len(elements.Connectors) != 3 {
			t.Errorf("Expected 3 connectors, got %d", len(elements.Connectors))
		}
		if elements.Bounds == nil {
			t.Error("Expected bounds to be set")
		}
	})

	t.Run("Sequence", func(t *testing.T) {
		code := `sequenceDiagram
Alice->>Bob: Hello Bob, how are you?
Bob-->>Alice: Jolly good!`
		elements, err := CreateDiagram(code)
		if err != nil {
			t.Fatalf("CreateDiagram failed: %v", err)
		}
		// 2 participants * (2 boxes + 1 lifeline) + 2 arrows + 2 text = 10 shapes
		if len(elements.Shapes) < 6 {
			t.Errorf("Expected at least 6 shapes, got %d", len(elements.Shapes))
		}
		if elements.Bounds == nil {
			t.Error("Expected bounds to be set")
		}
	})

	t.Run("Pie", func(t *testing.T) {
		code := `pie title Pets adopted by volunteers
"Dogs" : 386
"Cats" : 85
"Rats" : 15`
		elements, err := CreateDiagram(code)
		if err != nil {
			t.Fatalf("CreateDiagram failed: %v", err)
		}
		// 1 title + 1 circle + 3 legend items (box + text) = 8 shapes
		if len(elements.Shapes) < 5 {
			t.Errorf("Expected at least 5 shapes, got %d", len(elements.Shapes))
		}
		if elements.Bounds == nil {
			t.Error("Expected bounds to be set")
		}
	})

	t.Run("Gantt", func(t *testing.T) {
		code := `gantt
title A Gantt Diagram
section Section
A task :a1, 2014-01-01, 30d
Another task :after a1, 20d`
		elements, err := CreateDiagram(code)
		if err != nil {
			t.Fatalf("CreateDiagram failed: %v", err)
		}
		// 1 title + 1 section + 2 tasks (label + bar) = 6 shapes
		if len(elements.Shapes) < 4 {
			t.Errorf("Expected at least 4 shapes, got %d", len(elements.Shapes))
		}
		if elements.Bounds == nil {
			t.Error("Expected bounds to be set")
		}
	})

	t.Run("Timeline", func(t *testing.T) {
		code := `timeline
title History of Social Media Platform
2002 : LinkedIn
2004 : Facebook
: Google`
		elements, err := CreateDiagram(code)
		if err != nil {
			t.Fatalf("CreateDiagram failed: %v", err)
		}
		// 1 title + 1 line + 2 dates (marker + label) + 3 events = 9 shapes
		if len(elements.Shapes) < 5 {
			t.Errorf("Expected at least 5 shapes, got %d", len(elements.Shapes))
		}
		if elements.Bounds == nil {
			t.Error("Expected bounds to be set")
		}
	})

	t.Run("Quadrant", func(t *testing.T) {
		code := `quadrantChart
title Reach and engagement of campaigns
x-axis Low Reach --> High Reach
y-axis Low Engagement --> High Engagement
quadrant-1 We should expand
Campaign A: [0.3, 0.6]`
		elements, err := CreateDiagram(code)
		if err != nil {
			t.Fatalf("CreateDiagram failed: %v", err)
		}
		// 1 title + 4 quadrants + 2 axes + 2 axis labels + 1 point (dot + label) = 12 shapes
		if len(elements.Shapes) < 8 {
			t.Errorf("Expected at least 8 shapes, got %d", len(elements.Shapes))
		}
		if elements.Bounds == nil {
			t.Error("Expected bounds to be set")
		}
	})

	t.Run("Class", func(t *testing.T) {
		code := `classDiagram
class Animal {
    +String name
    +isMammal()
}
class Dog {
    +bark()
}
Animal <|-- Dog`
		elements, err := CreateDiagram(code)
		if err != nil {
			t.Fatalf("CreateDiagram failed: %v", err)
		}
		// 2 classes * 3 shapes (header, attr, method) = 6 shapes
		if len(elements.Shapes) < 6 {
			t.Errorf("Expected at least 6 shapes, got %d", len(elements.Shapes))
		}
		if len(elements.Connectors) != 1 {
			t.Errorf("Expected 1 connector, got %d", len(elements.Connectors))
		}
	})

	t.Run("State", func(t *testing.T) {
		code := `stateDiagram-v2
[*] --> First
First --> Second
Second --> [*]`
		elements, err := CreateDiagram(code)
		if err != nil {
			t.Fatalf("CreateDiagram failed: %v", err)
		}
		// 3 states ([*], First, Second) = 3 shapes
		if len(elements.Shapes) < 3 {
			t.Errorf("Expected at least 3 shapes, got %d", len(elements.Shapes))
		}
		if len(elements.Connectors) != 3 {
			t.Errorf("Expected 3 connectors, got %d", len(elements.Connectors))
		}
	})

	t.Run("ER", func(t *testing.T) {
		code := `erDiagram
CUSTOMER ||--o{ ORDER : places
CUSTOMER {
    string name
    string email
}
ORDER {
    int orderNumber
}`
		elements, err := CreateDiagram(code)
		if err != nil {
			t.Fatalf("CreateDiagram failed: %v", err)
		}
		// 2 entities * 2 shapes (header, attr) = 4 shapes
		if len(elements.Shapes) < 4 {
			t.Errorf("Expected at least 4 shapes, got %d", len(elements.Shapes))
		}
		if len(elements.Connectors) != 1 {
			t.Errorf("Expected 1 connector, got %d", len(elements.Connectors))
		}
	})

	t.Run("Mindmap", func(t *testing.T) {
		code := `mindmap
root((mindmap))
    Origins
        Long history
    Research
        On effectivness`
		elements, err := CreateDiagram(code)
		if err != nil {
			t.Fatalf("CreateDiagram failed: %v", err)
		}
		// 5 nodes = 5 shapes
		if len(elements.Shapes) != 5 {
			t.Errorf("Expected 5 shapes, got %d", len(elements.Shapes))
		}
		if len(elements.Connectors) != 4 {
			t.Errorf("Expected 4 connectors, got %d", len(elements.Connectors))
		}
	})

	t.Run("Journey", func(t *testing.T) {
		code := `journey
title My working day
section Go to work
    Make tea: 5: Me
    Do work: 1: Me, Cat`
		elements, err := CreateDiagram(code)
		if err != nil {
			t.Fatalf("CreateDiagram failed: %v", err)
		}
		// 1 title + 1 section + 2 tasks = 4 shapes
		if len(elements.Shapes) != 4 {
			t.Errorf("Expected 4 shapes, got %d", len(elements.Shapes))
		}
	})

	t.Run("GitGraph", func(t *testing.T) {
		code := `gitGraph
commit
commit
branch develop
checkout develop
commit
checkout main
merge develop`
		elements, err := CreateDiagram(code)
		if err != nil {
			t.Fatalf("CreateDiagram failed: %v", err)
		}
		// 2 branch lines + 2 branch labels + 4 commits = 8 shapes
		if len(elements.Shapes) != 8 {
			t.Errorf("Expected 8 shapes, got %d", len(elements.Shapes))
		}
		// 2 horizontal + 1 merge = 3 connectors
		if len(elements.Connectors) != 3 {
			t.Errorf("Expected 3 connectors, got %d", len(elements.Connectors))
		}
	})
}
