package shape

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func inheritPtr(v bool) *bool     { return &v }
func inheritStr(v string) *string { return &v }
func inheritInt(v int) *int       { return &v }

// Shadow inheritance mirrors python-pptx's ShadowFormat.inherit (upstream #130):
// inherit=false suppresses the style effects, and explicit effects may replace
// them instead of being rejected.
func TestRenderEffectsXMLShadowInheritCombinations(t *testing.T) {
	t.Run("inherit false alone empties the effect list", func(t *testing.T) {
		xml, err := RenderEffectsXML(&common.ShapeShadow{Inherit: inheritPtr(false)}, nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if xml != effectListEmptyXML {
			t.Fatalf("xml = %q, want %q", xml, effectListEmptyXML)
		}
	})

	t.Run("inherit false with explicit shadow renders the shadow", func(t *testing.T) {
		xml, err := RenderEffectsXML(&common.ShapeShadow{
			Inherit:     inheritPtr(false),
			Color:       inheritStr("FF0000"),
			BlurEmu:     inheritInt(60000),
			DistanceEmu: inheritInt(40000),
		}, nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(xml, `<a:outerShdw blurRad="60000" dist="40000"`) ||
			!strings.Contains(xml, `val="FF0000"`) {
			t.Fatalf("unexpected xml: %s", xml)
		}
	})

	t.Run("inherit false with another effect renders that effect", func(t *testing.T) {
		xml, err := RenderEffectsXML(
			&common.ShapeShadow{Inherit: inheritPtr(false)},
			&common.ShapeGlow{RadiusEmu: inheritInt(1000)},
			nil, nil, nil,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(xml, `<a:glow rad="1000"`) {
			t.Fatalf("unexpected xml: %s", xml)
		}
		if strings.Contains(xml, "outerShdw") {
			t.Fatalf("inherit=false must not synthesize a shadow: %s", xml)
		}
	})

	t.Run("inherit true alone emits nothing", func(t *testing.T) {
		xml, err := RenderEffectsXML(&common.ShapeShadow{Inherit: inheritPtr(true)}, nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if xml != "" {
			t.Fatalf("xml = %q, want empty", xml)
		}
	})

	t.Run("inherit true with explicit shadow is rejected", func(t *testing.T) {
		if _, err := RenderEffectsXML(&common.ShapeShadow{
			Inherit: inheritPtr(true),
			Color:   inheritStr("FF0000"),
		}, nil, nil, nil, nil); err == nil {
			t.Fatal("expected contradiction error")
		}
	})

	t.Run("inherit true with another effect is rejected", func(t *testing.T) {
		if _, err := RenderEffectsXML(
			&common.ShapeShadow{Inherit: inheritPtr(true)},
			&common.ShapeGlow{RadiusEmu: inheritInt(1000)},
			nil, nil, nil,
		); err == nil {
			t.Fatal("expected contradiction error")
		}
	})
}
