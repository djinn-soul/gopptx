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

func TestParseShapeProperties_ParsesTextAndEffects(t *testing.T) {
	xml := []byte(`
<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"
      xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <p:nvSpPr>
    <p:cNvPr id="9" name="Shape 9"/>
    <p:cNvSpPr/>
    <p:nvPr><p:ph idx="3" type="body"/></p:nvPr>
  </p:nvSpPr>
  <p:spPr>
    <a:xfrm>
      <a:off x="100" y="200"/>
      <a:ext cx="300" cy="400"/>
    </a:xfrm>
    <a:effectLst>
      <a:outerShdw blurRad="1000" dist="2000" dir="5400000">
        <a:srgbClr val="112233"/>
      </a:outerShdw>
      <a:glow rad="300"><a:srgbClr val="AABBCC"/></a:glow>
      <a:blur rad="250"/>
      <a:softEdge rad="120"/>
      <a:reflection blurRad="80" dist="160"/>
    </a:effectLst>
  </p:spPr>
  <p:txBody>
    <a:bodyPr lIns="10" rIns="20" tIns="30" bIns="40" anchor="ctr"/>
    <a:lstStyle/>
    <a:p>
      <a:pPr marL="91440" indent="-45720" lvl="1" algn="ctr">
        <a:lnSp><a:spcPct val="120000"/></a:lnSp>
        <a:spcBef><a:spcPts val="600"/></a:spcBef>
        <a:spcAft><a:spcPts val="400"/></a:spcAft>
      </a:pPr>
      <a:r>
        <a:rPr b="1" i="1" u="sng" baseline="-25000" cap="small" sz="1800">
          <a:solidFill><a:srgbClr val="FF00FF"/></a:solidFill>
          <a:latin typeface="Calibri"/>
        </a:rPr>
        <a:t>Hello</a:t>
      </a:r>
    </a:p>
  </p:txBody>
</p:sp>`)

	props, err := ParseShapeProperties(xml)
	if err != nil {
		t.Fatalf("ParseShapeProperties error: %v", err)
	}
	if props.Name != "Shape 9" || props.X != 100 || props.Y != 200 || props.W != 300 ||
		props.H != 400 {
		t.Fatalf("unexpected identity/transform parse: %+v", props)
	}
	if props.PhType != "body" || props.PhIndex != 3 {
		t.Fatalf("unexpected placeholder parse: type=%q idx=%d", props.PhType, props.PhIndex)
	}
	if props.Shadow == nil || props.Shadow.Color == nil || *props.Shadow.Color != "112233" {
		t.Fatalf("unexpected shadow parse: %+v", props.Shadow)
	}
	if props.Glow == nil || props.Blur == nil || props.SoftEdge == nil || props.Reflection == nil {
		t.Fatalf(
			"expected all effects parsed: glow=%+v blur=%+v soft=%+v refl=%+v",
			props.Glow,
			props.Blur,
			props.SoftEdge,
			props.Reflection,
		)
	}
	if len(props.Runs) != 1 || props.Runs[0].Text != "Hello" {
		t.Fatalf("unexpected run parse: %+v", props.Runs)
	}
	if props.Runs[0].Color == nil || *props.Runs[0].Color != "FF00FF" {
		t.Fatalf("unexpected run style parse: %+v", props.Runs[0])
	}
	if props.Paragraph == nil || props.Paragraph.Alignment == nil ||
		*props.Paragraph.Alignment != "ctr" {
		t.Fatalf("unexpected paragraph parse: %+v", props.Paragraph)
	}
}
