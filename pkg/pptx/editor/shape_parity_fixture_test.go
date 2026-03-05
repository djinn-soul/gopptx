package editor

import (
	"fmt"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestShapeParityFixture_DeterministicShapeIDAllocation(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-id-parity-fixture.pptx", []elements.SlideContent{
		elements.NewSlide("Shape ID Parity"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	const shapeCount = 40
	ids := make([]int, 0, shapeCount)
	for i := range shapeCount {
		shapeType := "rect"
		if i%2 == 1 {
			shapeType = "ellipse"
		}
		shapeID, addErr := editor.AddShape(
			0,
			shapeType,
			float64(100+i*20),
			float64(100+i*10),
			240,
			140,
		)
		if addErr != nil {
			t.Fatalf("add shape %d: %v", i, addErr)
		}
		ids = append(ids, shapeID)
	}

	for i := 1; i < len(ids); i++ {
		if ids[i] != ids[i-1]+1 {
			t.Fatalf("shape IDs must be deterministic and sequential, got %d then %d", ids[i-1], ids[i])
		}
	}

	// Ensure IDs continue monotonically across other shape-tree mutations.
	groupID, err := editor.GroupShapes(0, []int{ids[0], ids[2]})
	if err != nil {
		t.Fatalf("group shapes: %v", err)
	}
	freeformID, err := editor.AddFreeformShape(0, []freeformPoint{
		{X: 1000, Y: 1000},
		{X: 1300, Y: 1100},
		{X: 1400, Y: 1500},
	}, true)
	if err != nil {
		t.Fatalf("add freeform shape: %v", err)
	}
	if freeformID <= groupID {
		t.Fatalf("expected monotonic ID allocation across shape operations: group=%d freeform=%d", groupID, freeformID)
	}
}

func TestShapeParityFixture_GroupingPreservesTreeOrderingSemantics(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-ordering-grouping-parity.pptx", []elements.SlideContent{
		elements.NewSlide("Grouping Order"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	shapeA, err := editor.AddShape(0, "rect", 100, 100, 200, 120)
	if err != nil {
		t.Fatalf("add shape A: %v", err)
	}
	shapeB, err := editor.AddShape(0, "ellipse", 400, 100, 200, 120)
	if err != nil {
		t.Fatalf("add shape B: %v", err)
	}
	shapeC, err := editor.AddShape(0, "triangle", 700, 100, 200, 120)
	if err != nil {
		t.Fatalf("add shape C: %v", err)
	}

	groupID, err := editor.GroupShapes(0, []int{shapeC, shapeA})
	if err != nil {
		t.Fatalf("group non-contiguous shapes: %v", err)
	}
	if groupID <= 0 {
		t.Fatalf("expected valid group id, got %d", groupID)
	}

	partPath := editor.slides[0].Part
	slideXML, ok := editor.parts.Get(partPath)
	if !ok {
		t.Fatalf("missing slide xml for %s", partPath)
	}
	xmlText := string(slideXML)

	// Group insertion should occur at earliest selected shape position.
	groupPos := strings.Index(xmlText, fmt.Sprintf(`<p:cNvPr id="%d" name="Group %d"`, groupID, groupID))
	shapeBPos := strings.Index(xmlText, fmt.Sprintf(`<p:cNvPr id="%d" name="ellipse %d"`, shapeB, shapeB))
	if groupPos == -1 || shapeBPos == -1 {
		t.Fatalf("expected group and untouched shape markers in XML")
	}
	if groupPos > shapeBPos {
		t.Fatalf("group node should preserve tree ordering at first selected shape position")
	}

	// Child order should follow original tree order, not caller-provided list order.
	groupStart := strings.Index(xmlText, "<p:grpSp")
	groupEnd := strings.Index(xmlText[groupStart:], "</p:grpSp>")
	if groupStart == -1 || groupEnd == -1 {
		t.Fatalf("expected group wrapper in XML")
	}
	groupXML := xmlText[groupStart : groupStart+groupEnd]
	shapeAPos := strings.Index(groupXML, fmt.Sprintf(`<p:cNvPr id="%d" name="rect %d"`, shapeA, shapeA))
	shapeCPos := strings.Index(groupXML, fmt.Sprintf(`<p:cNvPr id="%d" name="triangle %d"`, shapeC, shapeC))
	if shapeAPos == -1 || shapeCPos == -1 {
		t.Fatalf("expected grouped child shapes in XML")
	}
	if shapeAPos > shapeCPos {
		t.Fatalf("expected grouped children in original tree order")
	}
}
