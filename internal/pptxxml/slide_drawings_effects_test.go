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

	for _, want := range []string{"<a:effectLst>", "<a:outerShdw", "<a:glow", "<a:softEdge", "<a:reflection"} {
		if !strings.Contains(xml, want) {
			t.Fatalf("expected %q in shape XML", want)
		}
	}
}

func TestCustomShapeXML_EmitsDetailedEffectValues(t *testing.T) {
	xml := customShapeXML(ShapeSpec{
		Type: "rect",
		Effects: &ShapeEffectsSpec{
			Glow:       true,
			SoftEdges:  true,
			Reflection: true,
			GlowSpec: &ShapeGlowSpec{
				Color:     "AABBCC",
				RadiusEmu: 1234,
			},
			BlurSpec: &ShapeBlurSpec{
				RadiusEmu: 2345,
			},
			SoftEdgeSpec: &ShapeSoftEdgeSpec{
				RadiusEmu: 3456,
			},
			ReflectionSpec: &ShapeReflectionSpec{
				BlurEmu:     4567,
				DistanceEmu: 5678,
			},
		},
	}, 12)

	for _, want := range []string{
		`<a:glow rad="1234"><a:srgbClr val="AABBCC"/></a:glow>`,
		`<a:blur rad="2345"/>`,
		`<a:softEdge rad="3456"/>`,
		`<a:reflection blurRad="4567" dist="5678"/>`,
	} {
		if !strings.Contains(xml, want) {
			t.Fatalf("expected %q in shape XML", want)
		}
	}
}
