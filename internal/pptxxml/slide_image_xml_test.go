package pptxxml

import (
	"strings"
	"testing"
)

func TestImageShapeXML(t *testing.T) {
	spec := ImageRef{
		RelID:   "rId1",
		X:       100,
		Y:       200,
		CX:      300,
		CY:      400,
		AltText: "Alt Text",
		Crop: &ImageCropRef{
			Left: 10000,
			Top:  20000,
		},
		Shadow:     true,
		Reflection: true,
		Rotation:   5400000,
		FlipH:      true,
		FlipV:      true,
	}
	xml := imageShape(spec, 5)
	if !strings.Contains(xml, `r:embed="rId1"`) {
		t.Error("RelID missing")
	}
	if !strings.Contains(xml, `<a:srcRect l="10000" t="20000"/>`) {
		t.Error("Crop missing")
	}
	if !strings.Contains(xml, `<a:outerShdw`) {
		t.Error("Shadow missing")
	}
	if !strings.Contains(xml, `<a:ref`) {
		t.Error("Reflection missing")
	}
	if !strings.Contains(xml, `rot="5400000"`) || !strings.Contains(xml, `flipH="1"`) {
		t.Error("Transforms missing")
	}
}
