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
    <a:solidFill><a:srgbClr val="ECECFF"><a:alpha val="65000"/></a:srgbClr></a:solidFill>
    <a:ln w="12700">
      <a:solidFill><a:srgbClr val="9370DB"/></a:solidFill>
      <a:headEnd type="triangle" w="lg" len="sm"/>
      <a:tailEnd type="arrow"/>
    </a:ln>
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
	if props.Fill.Transparency == nil || *props.Fill.Transparency != 0.35 {
		t.Fatalf("expected transparency 0.35, got %#v", props.Fill.Transparency)
	}
	if props.Line == nil || props.Line.Color == nil || *props.Line.Color != "9370DB" {
		t.Fatalf("unexpected line: %#v", props.Line)
	}
	if props.Line.StartArrow == nil || *props.Line.StartArrow != "triangle" {
		t.Fatalf("expected start arrow triangle, got %#v", props.Line)
	}
	if props.Line.StartArrowWidth == nil || *props.Line.StartArrowWidth != "lg" {
		t.Fatalf("expected start arrow width lg, got %#v", props.Line)
	}
	if props.Line.StartArrowLength == nil || *props.Line.StartArrowLength != "sm" {
		t.Fatalf("expected start arrow length sm, got %#v", props.Line)
	}
	if props.Line.EndArrow == nil || *props.Line.EndArrow != "arrow" {
		t.Fatalf("expected end arrow arrow, got %#v", props.Line)
	}
}

//nolint:gocyclo,cyclop // This regression test intentionally validates many parsed shape fields in one fixture.
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
        <a:buAutoNum type="romanLcPeriod"/>
        <a:buClr><a:srgbClr val="ABCDEF"/></a:buClr>
        <a:buSzPct val="85000"/>
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
	if props.Paragraph.BulletStyle == nil || *props.Paragraph.BulletStyle != "roman_lower" {
		t.Fatalf("expected roman_lower bullet style, got %+v", props.Paragraph)
	}
	if props.Paragraph.BulletColor == nil || *props.Paragraph.BulletColor != "ABCDEF" {
		t.Fatalf("expected bullet color ABCDEF, got %+v", props.Paragraph)
	}
	if props.Paragraph.BulletSizePct == nil || *props.Paragraph.BulletSizePct != 85 {
		t.Fatalf("expected bullet_size_pct 85, got %+v", props.Paragraph)
	}
	if props.TextFrame == nil {
		t.Fatalf("expected text frame parse, got nil")
	}
	if props.TextFrame.MarginLeft == nil || *props.TextFrame.MarginLeft != 10 {
		t.Fatalf("expected left margin 10, got %+v", props.TextFrame)
	}
	if props.TextFrame.MarginRight == nil || *props.TextFrame.MarginRight != 20 {
		t.Fatalf("expected right margin 20, got %+v", props.TextFrame)
	}
	if props.TextFrame.MarginTop == nil || *props.TextFrame.MarginTop != 30 {
		t.Fatalf("expected top margin 30, got %+v", props.TextFrame)
	}
	if props.TextFrame.MarginBottom == nil || *props.TextFrame.MarginBottom != 40 {
		t.Fatalf("expected bottom margin 40, got %+v", props.TextFrame)
	}
	if props.TextFrame.VerticalAlign == nil || *props.TextFrame.VerticalAlign != "ctr" {
		t.Fatalf("expected vertical align ctr, got %+v", props.TextFrame)
	}
	if props.Runs[0].Font == nil || *props.Runs[0].Font != "Calibri" {
		t.Fatalf("expected parsed run font Calibri, got %+v", props.Runs[0])
	}
	if props.Runs[0].SizePt == nil || *props.Runs[0].SizePt != 18 {
		t.Fatalf("expected parsed run size 18pt, got %+v", props.Runs[0])
	}
}

func TestParseShapeProperties_ParsesTextFrameAutoFitAndShapeRotation(t *testing.T) {
	xml := []byte(`
<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"
      xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <p:nvSpPr>
    <p:cNvPr id="11" name="Rotated Shape"/>
    <p:cNvSpPr/>
    <p:nvPr/>
  </p:nvSpPr>
  <p:spPr>
    <a:xfrm rot="5400000">
      <a:off x="10" y="20"/>
      <a:ext cx="300" cy="400"/>
    </a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
  </p:spPr>
  <p:txBody>
    <a:bodyPr wrap="none" vert="vert270" numCol="2" rot="-2700000"><a:normAutoFit/></a:bodyPr>
    <a:lstStyle/>
    <a:p><a:r><a:t>Rotated</a:t></a:r></a:p>
  </p:txBody>
</p:sp>`)

	props, err := ParseShapeProperties(xml)
	if err != nil {
		t.Fatalf("ParseShapeProperties error: %v", err)
	}
	if props.Rotation == nil || *props.Rotation != 90 {
		t.Fatalf("expected shape rotation 90deg, got %+v", props.Rotation)
	}
	if props.TextFrame == nil {
		t.Fatalf("expected text frame parse, got nil")
	}
	if props.TextFrame.WordWrap == nil || *props.TextFrame.WordWrap {
		t.Fatalf("expected word wrap false, got %+v", props.TextFrame)
	}
	if props.TextFrame.Orientation == nil || *props.TextFrame.Orientation != "vert270" {
		t.Fatalf("expected orientation vert270, got %+v", props.TextFrame)
	}
	if props.TextFrame.Columns == nil || *props.TextFrame.Columns != 2 {
		t.Fatalf("expected columns 2, got %+v", props.TextFrame)
	}
	if props.TextFrame.Rotation == nil || *props.TextFrame.Rotation != -45 {
		t.Fatalf("expected text rotation -45deg, got %+v", props.TextFrame)
	}
	if props.TextFrame.AutoFitType == nil || *props.TextFrame.AutoFitType != "normal" {
		t.Fatalf("expected autofit normal, got %+v", props.TextFrame)
	}
	if props.TextFrame.AutoFit == nil || !*props.TextFrame.AutoFit {
		t.Fatalf("expected legacy autofit true, got %+v", props.TextFrame)
	}
}

func TestParseShapeProperties_ParsesRunOutline(t *testing.T) {
	xml := []byte(`
<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"
      xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <p:nvSpPr>
    <p:cNvPr id="14" name="Outlined Text"/>
    <p:cNvSpPr/>
    <p:nvPr/>
  </p:nvSpPr>
  <p:spPr>
    <a:xfrm><a:off x="0" y="0"/><a:ext cx="300" cy="200"/></a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
  </p:spPr>
  <p:txBody>
    <a:bodyPr/>
    <a:lstStyle/>
    <a:p>
      <a:r>
        <a:rPr>
          <a:ln w="38100"><a:solidFill><a:srgbClr val="112233"/></a:solidFill></a:ln>
        </a:rPr>
        <a:t>Outlined</a:t>
      </a:r>
    </a:p>
  </p:txBody>
</p:sp>`)

	props, err := ParseShapeProperties(xml)
	if err != nil {
		t.Fatalf("ParseShapeProperties error: %v", err)
	}
	if len(props.Runs) != 1 {
		t.Fatalf("expected one run, got %+v", props.Runs)
	}
	if props.Runs[0].OutlineColor == nil || *props.Runs[0].OutlineColor != "112233" {
		t.Fatalf("expected outline color 112233, got %+v", props.Runs[0].OutlineColor)
	}
	if props.Runs[0].OutlineWidthPt == nil || *props.Runs[0].OutlineWidthPt != 3 {
		t.Fatalf("expected outline width 3pt, got %+v", props.Runs[0].OutlineWidthPt)
	}
}

func TestParseShapeProperties_ParsesMultipleTextParagraphs(t *testing.T) {
	xml := []byte(`
<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"
      xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <p:nvSpPr>
    <p:cNvPr id="12" name="Multi Paragraph Shape"/>
    <p:cNvSpPr/>
    <p:nvPr/>
  </p:nvSpPr>
  <p:spPr>
    <a:xfrm>
      <a:off x="10" y="20"/>
      <a:ext cx="300" cy="400"/>
    </a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
  </p:spPr>
  <p:txBody>
    <a:bodyPr/>
    <a:lstStyle/>
    <a:p>
      <a:pPr lvl="1"><a:buAutoNum type="arabicPeriod"/></a:pPr>
      <a:r><a:rPr b="1"/><a:t>Alpha</a:t></a:r>
    </a:p>
    <a:p>
      <a:r><a:t>Beta</a:t></a:r>
    </a:p>
  </p:txBody>
</p:sp>`)

	props, err := ParseShapeProperties(xml)
	if err != nil {
		t.Fatalf("ParseShapeProperties error: %v", err)
	}
	if len(props.Paragraphs) != 2 {
		t.Fatalf("expected 2 parsed paragraphs, got %d", len(props.Paragraphs))
	}
	if len(props.Paragraphs[0].Runs) != 1 || props.Paragraphs[0].Runs[0].Text != "Alpha" {
		t.Fatalf("unexpected first paragraph runs: %+v", props.Paragraphs[0])
	}
	if props.Paragraphs[0].Paragraph == nil || props.Paragraphs[0].Paragraph.BulletStyle == nil ||
		*props.Paragraphs[0].Paragraph.BulletStyle != "number" {
		t.Fatalf("unexpected first paragraph style: %+v", props.Paragraphs[0].Paragraph)
	}
	if len(props.Paragraphs[1].Runs) != 1 || props.Paragraphs[1].Runs[0].Text != "Beta" {
		t.Fatalf("unexpected second paragraph runs: %+v", props.Paragraphs[1])
	}
	if props.Text != "Alpha\nBeta" {
		t.Fatalf("expected combined text with newline, got %q", props.Text)
	}
}

func TestParseShapeReaderMetadata_ParsesAltTextAndActions(t *testing.T) {
	xml := []byte(`
<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"
      xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
      xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <p:nvSpPr>
    <p:cNvPr id="13" name="Action Shape" descr="Accessible rectangle">
      <a:hlinkClick r:id="rId1" tooltip="Jump" action="ppaction://hlinksldjump"/>
      <a:hlinkHover action="ppaction://macro?name=HoverMacro"/>
    </p:cNvPr>
    <p:cNvSpPr/>
    <p:nvPr/>
  </p:nvSpPr>
  <p:txBody>
    <a:bodyPr/>
    <a:lstStyle/>
    <a:p>
      <a:r>
        <a:rPr>
          <a:hlinkClick r:id="rId2" tooltip="Run link"/>
          <a:hlinkMouseOver action="ppaction://hlinkshowjump?jump=nextslide"/>
        </a:rPr>
        <a:t>Hello</a:t>
      </a:r>
    </a:p>
  </p:txBody>
</p:sp>`)

	metadata, err := ParseShapeReaderMetadata(xml)
	if err != nil {
		t.Fatalf("ParseShapeReaderMetadata error: %v", err)
	}
	if !metadata.HasAltText || metadata.AltText != "Accessible rectangle" {
		t.Fatalf("expected alt text to be preserved, got %+v", metadata)
	}
	if metadata.ClickAction == nil || metadata.ClickAction.RelID != "rId1" {
		t.Fatalf("expected shape click action rel id rId1, got %+v", metadata.ClickAction)
	}
	if metadata.HoverAction == nil || metadata.HoverAction.Action == nil ||
		*metadata.HoverAction.Action != "ppaction://macro?name=HoverMacro" {
		t.Fatalf("expected macro hover action, got %+v", metadata.HoverAction)
	}
	if len(metadata.RunActions) != 1 || len(metadata.RunActions[0]) != 1 {
		t.Fatalf("expected one run action row, got %+v", metadata.RunActions)
	}
	if metadata.RunActions[0][0].ClickAction == nil || metadata.RunActions[0][0].ClickAction.RelID != "rId2" {
		t.Fatalf("expected run click action rel id rId2, got %+v", metadata.RunActions[0][0].ClickAction)
	}
	if metadata.RunActions[0][0].HoverAction == nil || metadata.RunActions[0][0].HoverAction.Action == nil ||
		*metadata.RunActions[0][0].HoverAction.Action != "ppaction://hlinkshowjump?jump=nextslide" {
		t.Fatalf("expected run hover action, got %+v", metadata.RunActions[0][0].HoverAction)
	}
}
