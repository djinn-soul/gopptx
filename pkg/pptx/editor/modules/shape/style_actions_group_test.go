package shape

import (
	"regexp"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func strPtr2(v string) *string { return &v }
func intPtr2(v int) *int       { return &v }
func boolPtr2(v bool) *bool    { return &v }
func floatPtr2(v float64) *float64 {
	return &v
}

func TestRenderFillLineAndEffectsXML(t *testing.T) {
	fillXML, err := RenderFillXML(&common.ShapeFill{Solid: strPtr2("#00ff00")})
	if err != nil || !strings.Contains(fillXML, `val="00FF00"`) {
		t.Fatalf("RenderFillXML(solid) failed: xml=%q err=%v", fillXML, err)
	}
	_, err = RenderFillXML(&common.ShapeFill{Solid: strPtr2("FF0000"), Background: boolPtr2(true)})
	if err == nil {
		t.Fatal("expected mutually exclusive fill mode error")
	}

	grad := &common.GradientFill{
		AngleDeg: floatPtr2(45),
		Stops: []common.GradientStop{
			{Color: "FF0000", PositionPct: floatPtr2(0)},
			{Color: "0000FF", PositionPct: floatPtr2(100)},
		},
	}
	gradXML, err := RenderFillXML(&common.ShapeFill{Gradient: grad})
	if err != nil || !strings.Contains(gradXML, "<a:gradFill>") || !strings.Contains(gradXML, "<a:lin") {
		t.Fatalf("RenderFillXML(gradient) failed: xml=%q err=%v", gradXML, err)
	}

	patternXML, err := RenderFillXML(&common.ShapeFill{
		Pattern: &common.PatternedFill{
			Preset:  strPtr2("diagCross"),
			FgColor: strPtr2("#112233"),
			BgColor: strPtr2("#445566"),
		},
	})
	if err != nil || !strings.Contains(patternXML, `prst="diagCross"`) || !strings.Contains(patternXML, `val="112233"`) {
		t.Fatalf("RenderFillXML(pattern) failed: xml=%q err=%v", patternXML, err)
	}

	lineXML, err := RenderLineXML(&common.ShapeLine{
		WidthEmu:  intPtr2(12700),
		Color:     strPtr2("ABCDEF"),
		DashStyle: strPtr2("long_dash"),
	})
	if err != nil || !strings.Contains(lineXML, `w="12700"`) || !strings.Contains(lineXML, `prstDash val="lgDash"`) {
		t.Fatalf("RenderLineXML failed: xml=%q err=%v", lineXML, err)
	}
	_, err = RenderLineXML(&common.ShapeLine{WidthEmu: intPtr2(0)})
	if err == nil {
		t.Fatal("expected invalid line width error")
	}

	effectsXML, err := RenderEffectsXML(
		&common.ShapeShadow{Color: strPtr2("#000000"), BlurEmu: intPtr2(100), DistanceEmu: intPtr2(200), AngleDeg: floatPtr2(10)},
		&common.ShapeGlow{Color: strPtr2("FF00FF"), RadiusEmu: intPtr2(300)},
		&common.ShapeBlur{RadiusEmu: intPtr2(400)},
		&common.ShapeSoftEdge{RadiusEmu: intPtr2(500)},
		&common.ShapeReflection{BlurEmu: intPtr2(600), DistanceEmu: intPtr2(700)},
	)
	if err != nil || !strings.Contains(effectsXML, "<a:effectLst>") || !strings.Contains(effectsXML, "<a:outerShdw") {
		t.Fatalf("RenderEffectsXML(explicit) failed: xml=%q err=%v", effectsXML, err)
	}
	inheritXML, err := RenderEffectsXML(&common.ShapeShadow{Inherit: boolPtr2(false)}, nil, nil, nil, nil)
	if err != nil || inheritXML != "<a:effectLst/>" {
		t.Fatalf("RenderEffectsXML(inherit false) failed: xml=%q err=%v", inheritXML, err)
	}
	inheritXML, err = RenderEffectsXML(&common.ShapeShadow{Inherit: boolPtr2(true)}, nil, nil, nil, nil)
	if err != nil || inheritXML != "" {
		t.Fatalf("RenderEffectsXML(inherit true) failed: xml=%q err=%v", inheritXML, err)
	}
	_, err = RenderEffectsXML(&common.ShapeShadow{Inherit: boolPtr2(true), Color: strPtr2("FF0000")}, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected inherit + explicit shadow validation error")
	}
}

func TestActionMutationHelpers(t *testing.T) {
	xmlOpenClose := `<p:sp><p:nvSpPr><p:cNvPr id="1"><a:hlinkClick r:id="rId1"/></p:cNvPr></p:nvSpPr></p:sp>`
	updated, err := ApplyCNvPrActions(
		[]byte(xmlOpenClose),
		true,
		false,
		`<a:hlinkClick r:id="rId2"/>`,
		"",
	)
	if err != nil {
		t.Fatalf("ApplyCNvPrActions(open-close) failed: %v", err)
	}
	updatedXML := string(updated)
	if strings.Contains(updatedXML, `r:id="rId1"`) || !strings.Contains(updatedXML, `r:id="rId2"`) {
		t.Fatalf("ApplyCNvPrActions did not replace click action: %s", updatedXML)
	}

	xmlSelfClosing := `<p:sp><p:nvSpPr><p:cNvPr id="1"/></p:nvSpPr></p:sp>`
	updated, err = ApplyCNvPrActions(
		[]byte(xmlSelfClosing),
		true,
		true,
		`<a:hlinkClick r:id="rId3"/>`,
		`<a:hlinkMouseOver r:id="rId4"/>`,
	)
	if err != nil {
		t.Fatalf("ApplyCNvPrActions(self-closing) failed: %v", err)
	}
	if !strings.Contains(string(updated), "<p:cNvPr id=\"1\">") || !strings.Contains(string(updated), `r:id="rId4"`) {
		t.Fatalf("ApplyCNvPrActions did not expand/append actions: %s", string(updated))
	}

	_, err = ApplyCNvPrActions([]byte(`<p:sp/>`), true, false, `<a:hlinkClick/>`, "")
	if err == nil {
		t.Fatal("expected missing cNvPr error when actions requested")
	}

	cleaned := RemoveMatchedTags(`<a>x</a><b>y</b><a>z</a>`, regexp.MustCompile(`(?s)<a>.*?</a>`))
	if cleaned != "<b>y</b>" {
		t.Fatalf("RemoveMatchedTags unexpected result: %q", cleaned)
	}
}

func TestGroupExtractionHelpers(t *testing.T) {
	groupXML := []byte(`
<p:grpSp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:nvGrpSpPr/>
  <p:grpSpPr/>
  <p:sp><p:nvSpPr><p:cNvPr id="12" name="Text"/></p:nvSpPr></p:sp>
  <p:pic><p:nvPicPr><p:cNvPr id="17" name="Picture"/></p:nvPicPr></p:pic>
</p:grpSp>`)

	children, err := ExtractGroupChildShapeNodes(groupXML, "grpSp", "pic")
	if err != nil {
		t.Fatalf("ExtractGroupChildShapeNodes failed: %v", err)
	}
	if len(children) != 2 {
		t.Fatalf("expected 2 child shape nodes, got %d", len(children))
	}

	pattern := regexp.MustCompile(`<p:cNvPr[^>]*id="(\d+)"`)
	if got := FirstShapeIDInXML(children, pattern, 2); got != 12 {
		t.Fatalf("FirstShapeIDInXML=%d, want 12", got)
	}

	invalidXML := []byte(`<p:grpSp><p:sp>`)
	if _, err = ExtractGroupChildShapeNodes(invalidXML, "grpSp", "pic"); err == nil {
		t.Fatal("expected invalid XML parse error")
	}
}
