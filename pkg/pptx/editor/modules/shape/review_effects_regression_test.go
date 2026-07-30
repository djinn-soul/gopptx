package shape

import "testing"

func TestParseGlowWithoutShadowDisablesShadowInheritance(t *testing.T) {
	shapeXML := []byte(
		`<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" ` +
			`xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
			`<p:nvSpPr><p:cNvPr id="2" name="Glow"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>` +
			`<p:spPr><a:prstGeom prst="rect"><a:avLst/></a:prstGeom>` +
			`<a:effectLst><a:glow rad="1000"><a:srgbClr val="FF0000"/></a:glow></a:effectLst>` +
			`</p:spPr></p:sp>`,
	)

	parsed, err := ParseShapeProperties(shapeXML)
	if err != nil {
		t.Fatalf("parse shape: %v", err)
	}
	if parsed.Shadow == nil || parsed.Shadow.Inherit == nil || *parsed.Shadow.Inherit {
		t.Fatalf("expected explicit inherit=false, got %+v", parsed.Shadow)
	}
	if parsed.Glow == nil {
		t.Fatal("expected glow to remain parsed")
	}
}
