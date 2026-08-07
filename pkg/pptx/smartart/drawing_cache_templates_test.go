package smartart

import (
	"os"
	"path/filepath"
	"testing"
)

// The shipped layout templates carry the drawing PowerPoint cached for each of
// them, which is exactly the input the parser sees when it reads a real deck.
// Sweeping them catches a shape of the format the parser does not handle.
const templateLayoutsDir = "../../../internal/pptxxml/templates/smartart/layouts"

func TestParseDrawingShapesAcceptsEveryShippedTemplate(t *testing.T) {
	entries, err := os.ReadDir(templateLayoutsDir)
	if err != nil {
		t.Skipf("layout templates not readable from here: %v", err)
	}

	checked := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		drawingPath := filepath.Join(templateLayoutsDir, entry.Name(), "drawing.xml")
		data, err := os.ReadFile(drawingPath)
		if err != nil {
			continue
		}
		checked++
		shapes := ParseDrawingShapes(data)
		if len(shapes) == 0 {
			t.Errorf("%s: parsed no shapes", entry.Name())
			continue
		}
		assertTemplateShapesUsable(t, entry.Name(), shapes)
	}
	if checked == 0 {
		t.Skip("no template drawings found")
	}
	t.Logf("parsed %d template drawings", checked)
}

func assertTemplateShapesUsable(t *testing.T, layout string, shapes []DrawingShape) {
	t.Helper()
	sized := 0
	geometries := 0
	for _, shape := range shapes {
		if shape.CX > 0 && shape.CY > 0 {
			sized++
		}
		if shape.PresetGeom != "" {
			geometries++
		}
	}
	if sized == 0 {
		t.Errorf("%s: no shape has a size", layout)
	}
	if geometries == 0 {
		t.Errorf("%s: no shape has preset geometry", layout)
	}
}
