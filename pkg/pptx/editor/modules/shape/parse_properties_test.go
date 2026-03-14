package shape

import "testing"

func TestParseShapeProperties_ParsesPresetGeometryAndAdjustments(t *testing.T) {
	xml := []byte(`
<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"
      xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <p:nvSpPr>
    <p:cNvPr id="4" name="Shape 4"/>
    <p:cNvSpPr/>
    <p:nvPr/>
  </p:nvSpPr>
  <p:spPr>
    <a:xfrm>
      <a:off x="1485900" y="2423160"/>
      <a:ext cx="3657600" cy="3657600"/>
    </a:xfrm>
    <a:prstGeom prst="pie">
      <a:avLst>
        <a:gd name="adj1" fmla="val 0"/>
        <a:gd name="adj2" fmla="val 17155555"/>
      </a:avLst>
    </a:prstGeom>
    <a:solidFill><a:srgbClr val="ECECFF"/></a:solidFill>
    <a:ln w="12700"><a:solidFill><a:srgbClr val="9370DB"/></a:solidFill></a:ln>
  </p:spPr>
  <p:txBody><a:bodyPr/><a:lstStyle/><a:p/></p:txBody>
</p:sp>`)

	props, err := ParseShapeProperties(xml)
	if err != nil {
		t.Fatalf("ParseShapeProperties error: %v", err)
	}
	if props.Type != "pie" {
		t.Fatalf("expected type=pie, got %q", props.Type)
	}
	if len(props.Adjustments) != 2 {
		t.Fatalf("expected 2 adjustments, got %d", len(props.Adjustments))
	}
	if props.Adjustments[0].Name != "adj1" || props.Adjustments[0].Formula != "val 0" {
		t.Fatalf("unexpected adjustment[0]: %#v", props.Adjustments[0])
	}
	if props.Fill == nil || props.Fill.Solid == nil || *props.Fill.Solid != "ECECFF" {
		t.Fatalf("unexpected fill: %#v", props.Fill)
	}
	if props.Line == nil || props.Line.Color == nil || *props.Line.Color != "9370DB" {
		t.Fatalf("unexpected line: %#v", props.Line)
	}
}
