package mermaid

import (
	"fmt"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestCreateDiagramFitsToSafeArea(t *testing.T) {
	code := `classDiagram
class Slide {
    +String title
    +List~Element~ elements
    +addElement()
}
class Element {
    +render()
}
class Presentation {
    +String title
    +List~Slide~ slides
    +addSlide()
    +save()
}
class Table {
    +List~Row~ rows
    +render()
}
class TextBox {
    +String text
    +render()
}
Slide "1" --> "*" Element
Presentation "1" --> "*" Slide
Element <|-- Table
Element <|-- TextBox`
	elements, err := CreateDiagram(code)
	if err != nil {
		t.Fatalf("CreateDiagram failed: %v", err)
	}

	left := styling.Inches(0.5)
	top := styling.Inches(1.2)
	right := styling.Inches(9.5)
	bottom := styling.Inches(7.1)
	for i, shape := range elements.Shapes {
		if shape.X < left || shape.Y < top || shape.X+shape.CX > right || shape.Y+shape.CY > bottom {
			t.Fatalf("shape %d out of safe area: x=%d y=%d cx=%d cy=%d", i, shape.X, shape.Y, shape.CX, shape.CY)
		}
	}
}

func TestPieSliceAnglesAreNormalized(t *testing.T) {
	pie := parsePie(`pie title Market Share 2024
"Our Product" : 35
"Competitor A" : 25
"Competitor B" : 20
"Competitor C" : 12
"Others" : 8`)
	elements := generatePieElements(pie, GetTheme("default"))
	limit := 21600000 // 360 * 60000
	pieCount := 0
	for _, shape := range elements.Shapes {
		if shape.Type != "pie" {
			continue
		}
		pieCount++
		for _, adj := range shape.Adjustments {
			if !strings.HasPrefix(adj.Formula, "val ") {
				t.Fatalf("invalid pie adjustment formula %q", adj.Formula)
			}
			var value int
			if _, err := fmt.Sscanf(adj.Formula, "val %d", &value); err != nil {
				t.Fatalf("invalid pie adjustment formula %q: %v", adj.Formula, err)
			}
			if value < 0 || value > limit {
				t.Fatalf("pie angle adjustment out of range: %d", value)
			}
		}
	}
	if pieCount != len(pie.Data) {
		t.Fatalf("expected %d pie slices, got %d", len(pie.Data), pieCount)
	}
}
