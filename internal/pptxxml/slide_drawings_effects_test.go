package pptxxml

import (
	"strings"
	"testing"
)

func TestCustomShapeXML_EmitsEffectList(t *testing.T) {
	xml := customShapeXML(ShapeSpec{
		Type: "rect",
		X:    1,
		Y:    1,
		CX:   100,
		CY:   80,
		Text: "fx",
		Effects: &ShapeEffectsSpec{
			Shadow:     true,
			Glow:       true,
			SoftEdges:  true,
			Reflection: true,
		},
	}, 11)

	for _, want := range []string{"<a:effectLst>", "<a:outerShdw", "<a:glow", "<a:softEdge", "<a:ref"} {
		if !strings.Contains(xml, want) {
			t.Fatalf("expected %q in shape XML", want)
		}
	}
}
