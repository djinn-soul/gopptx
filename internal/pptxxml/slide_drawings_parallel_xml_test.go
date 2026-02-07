package pptxxml

import (
	"strings"
	"testing"
)

func TestRenderCustomShapeXMLConcurrentlyPreservesOrder(t *testing.T) {
	shapes := []ShapeSpec{
		{Type: "rect", X: 1, Y: 1, CX: 1, CY: 1, Text: "A"},
		{Type: "rect", X: 2, Y: 2, CX: 2, CY: 2, Text: "B"},
		{Type: "rect", X: 3, Y: 3, CX: 3, CY: 3, Text: "C"},
	}

	shapeIDs, xmlParts := renderCustomShapeXMLConcurrently(shapes, 7)
	if len(shapeIDs) != 3 || len(xmlParts) != 3 {
		t.Fatalf("expected 3 shape ids/xml parts, got ids=%d xml=%d", len(shapeIDs), len(xmlParts))
	}

	expectedIDs := []int{7, 8, 9}
	for i, want := range expectedIDs {
		if shapeIDs[i] != want {
			t.Fatalf("unexpected shape id at index %d: got %d want %d", i, shapeIDs[i], want)
		}
	}

	expectations := []string{
		`<p:cNvPr id="7" name="Shape 7"/>`,
		`<a:t>A</a:t>`,
		`<p:cNvPr id="8" name="Shape 8"/>`,
		`<a:t>B</a:t>`,
		`<p:cNvPr id="9" name="Shape 9"/>`,
		`<a:t>C</a:t>`,
	}
	for idx, needle := range expectations {
		part := xmlParts[idx/2]
		if !strings.Contains(part, needle) {
			t.Fatalf("expected %q in xml part %d", needle, idx/2)
		}
	}
}
