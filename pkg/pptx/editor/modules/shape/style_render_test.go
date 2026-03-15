package shape

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func intPtrSR(v int) *int           { return &v }
func boolPtrSR(v bool) *bool        { return &v }
func strPtrSR(v string) *string     { return &v }
func floatPtrSR(v float64) *float64 { return &v }

func TestRenderLineAndFillXML(t *testing.T) {
	lineXML, err := RenderLineXML(&common.ShapeLine{
		Color:     strPtrSR("#00FF00"),
		WidthEmu:  intPtrSR(12700),
		DashStyle: strPtrSR("long-dash"),
	})
	if err != nil {
		t.Fatalf("RenderLineXML failed: %v", err)
	}
	if !strings.Contains(lineXML, `w="12700"`) || !strings.Contains(lineXML, `prstDash val="lgDash"`) {
		t.Fatalf("unexpected line xml: %s", lineXML)
	}
	if _, err = RenderLineXML(&common.ShapeLine{WidthEmu: intPtrSR(0)}); err == nil {
		t.Fatal("expected line width validation error")
	}
	if _, err = RenderLineXML(&common.ShapeLine{Color: strPtrSR("GGGGGG")}); err == nil {
		t.Fatal("expected line color validation error")
	}
	if _, err = RenderLineXML(&common.ShapeLine{DashStyle: strPtrSR("bogus")}); err == nil {
		t.Fatal("expected dash style validation error")
	}

	fillXML, err := RenderFillXML(&common.ShapeFill{Solid: strPtrSR("AABBCC")})
	if err != nil || !strings.Contains(fillXML, `srgbClr val="AABBCC"`) {
		t.Fatalf("solid fill render failed: xml=%q err=%v", fillXML, err)
	}
	noFillXML, err := RenderFillXML(&common.ShapeFill{Background: boolPtrSR(true)})
	if err != nil || noFillXML != `<a:noFill/>` {
		t.Fatalf("background fill render failed: xml=%q err=%v", noFillXML, err)
	}
	if _, err = RenderFillXML(&common.ShapeFill{Background: boolPtrSR(false)}); err == nil {
		t.Fatal("expected background=false validation error")
	}
	if _, err = RenderFillXML(&common.ShapeFill{
		Solid:      strPtrSR("112233"),
		Background: boolPtrSR(true),
	}); err == nil {
		t.Fatal("expected mutually exclusive fill mode error")
	}

	gradientXML, err := RenderFillXML(&common.ShapeFill{
		Gradient: &common.GradientFill{
			AngleDeg: floatPtrSR(45),
			Stops: []common.GradientStop{
				{Color: "#FF0000", PositionPct: floatPtrSR(0)},
				{Color: "#00FF00", PositionPct: floatPtrSR(100)},
			},
		},
	})
	if err != nil || !strings.Contains(gradientXML, "<a:gradFill>") || !strings.Contains(gradientXML, `ang="2700000"`) {
		t.Fatalf("gradient fill render failed: xml=%q err=%v", gradientXML, err)
	}
	if _, err = RenderFillXML(&common.ShapeFill{Gradient: &common.GradientFill{Stops: nil}}); err == nil {
		t.Fatal("expected empty gradient stops error")
	}
	if _, err = RenderFillXML(&common.ShapeFill{
		Gradient: &common.GradientFill{
			Stops: []common.GradientStop{{Color: "#FF0000", PositionPct: floatPtrSR(101)}},
		},
	}); err == nil {
		t.Fatal("expected out-of-range gradient position error")
	}

	patternXML, err := RenderFillXML(&common.ShapeFill{
		Pattern: &common.PatternedFill{
			Preset:  strPtrSR("diagCross"),
			FgColor: strPtrSR("111111"),
			BgColor: strPtrSR("FFFFFF"),
		},
	})
	if err != nil || !strings.Contains(patternXML, `prst="diagCross"`) {
		t.Fatalf("pattern fill render failed: xml=%q err=%v", patternXML, err)
	}
	if _, err = RenderFillXML(&common.ShapeFill{Pattern: &common.PatternedFill{FgColor: strPtrSR("BADHEX")}}); err == nil {
		t.Fatal("expected pattern fg color validation error")
	}
}

func TestRenderEffectsXMLAndRotationValidation(t *testing.T) {
	effectsXML, err := RenderEffectsXML(
		&common.ShapeShadow{
			Color:       strPtrSR("222222"),
			BlurEmu:     intPtrSR(50800),
			DistanceEmu: intPtrSR(38100),
			AngleDeg:    floatPtrSR(15),
		},
		&common.ShapeGlow{Color: strPtrSR("00FF00"), RadiusEmu: intPtrSR(1000)},
		&common.ShapeBlur{RadiusEmu: intPtrSR(900)},
		&common.ShapeSoftEdge{RadiusEmu: intPtrSR(800)},
		&common.ShapeReflection{BlurEmu: intPtrSR(700), DistanceEmu: intPtrSR(600)},
	)
	if err != nil {
		t.Fatalf("RenderEffectsXML failed: %v", err)
	}
	if !strings.Contains(effectsXML, "<a:effectLst>") || !strings.Contains(effectsXML, "<a:outerShdw") {
		t.Fatalf("unexpected effects XML: %s", effectsXML)
	}

	if xml, err := RenderEffectsXML(&common.ShapeShadow{Inherit: boolPtrSR(true)}, nil, nil, nil, nil); err != nil || xml != "" {
		t.Fatalf("shadow inherit=true should emit empty xml without error: xml=%q err=%v", xml, err)
	}
	if xml, err := RenderEffectsXML(&common.ShapeShadow{Inherit: boolPtrSR(false)}, nil, nil, nil, nil); err != nil || xml != `<a:effectLst/>` {
		t.Fatalf("shadow inherit=false should emit empty effect list: xml=%q err=%v", xml, err)
	}

	if _, err := RenderEffectsXML(&common.ShapeShadow{Inherit: boolPtrSR(true), Color: strPtrSR("FFFFFF")}, nil, nil, nil, nil); err == nil {
		t.Fatal("expected inherit + explicit shadow attribute error")
	}
	if _, err := RenderEffectsXML(&common.ShapeShadow{Inherit: boolPtrSR(true)}, &common.ShapeGlow{}, nil, nil, nil); err == nil {
		t.Fatal("expected inherit + other effects error")
	}
	if _, err := RenderEffectsXML(&common.ShapeShadow{BlurEmu: intPtrSR(-1)}, nil, nil, nil, nil); err == nil {
		t.Fatal("expected negative shadow blur error")
	}
	if _, err := RenderEffectsXML(nil, &common.ShapeGlow{RadiusEmu: intPtrSR(-1)}, nil, nil, nil); err == nil {
		t.Fatal("expected negative glow radius error")
	}
	if _, err := RenderEffectsXML(nil, nil, &common.ShapeBlur{RadiusEmu: intPtrSR(-1)}, nil, nil); err == nil {
		t.Fatal("expected negative blur radius error")
	}
	if _, err := RenderEffectsXML(nil, nil, nil, &common.ShapeSoftEdge{RadiusEmu: intPtrSR(-1)}, nil); err == nil {
		t.Fatal("expected negative soft edge radius error")
	}
	if _, err := RenderEffectsXML(nil, nil, nil, nil, &common.ShapeReflection{DistanceEmu: intPtrSR(-1)}); err == nil {
		t.Fatal("expected negative reflection distance error")
	}
}

func TestRenderShapeStyleXMLComposition(t *testing.T) {
	xml, err := RenderShapeStyleXML(
		&common.ShapeFill{Solid: strPtrSR("123456")},
		&common.ShapeLine{Color: strPtrSR("654321"), WidthEmu: intPtrSR(1000)},
		&common.ShapeShadow{Color: strPtrSR("222222")},
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("RenderShapeStyleXML failed: %v", err)
	}
	if !strings.Contains(xml, "<a:solidFill>") || !strings.Contains(xml, "<a:ln") || !strings.Contains(xml, "<a:effectLst>") {
		t.Fatalf("expected fill+line+effects composition, got: %s", xml)
	}
}
