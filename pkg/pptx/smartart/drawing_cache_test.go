package smartart

import "testing"

const sampleDrawingXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<dsp:drawing xmlns:dsp="http://schemas.microsoft.com/office/drawing/2008/diagram"
 xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><dsp:spTree>
<dsp:nvGrpSpPr><dsp:cNvPr id="0" name=""/><dsp:cNvGrpSpPr/></dsp:nvGrpSpPr><dsp:grpSpPr/>
<dsp:sp modelId="{NODE-1}"><dsp:spPr><a:xfrm rot="5400000"><a:off x="4042" y="829182"/>
<a:ext cx="1838086" cy="1252799"/></a:xfrm><a:prstGeom prst="roundRect"><a:avLst>
<a:gd name="adj" fmla="val 10000"/></a:avLst></a:prstGeom><a:solidFill><a:schemeClr val="accent1">
<a:hueOff val="0"/><a:satOff val="0"/><a:lumOff val="20000"/></a:schemeClr></a:solidFill>
<a:ln w="19050"><a:solidFill><a:schemeClr val="lt1"/></a:solidFill></a:ln></dsp:spPr>
<dsp:txBody><a:bodyPr anchor="ctr"/><a:p><a:pPr algn="ctr"/><a:r><a:rPr sz="1800" b="1">
<a:solidFill><a:srgbClr val="FFFFFF"/></a:solidFill><a:latin typeface="Calibri"/></a:rPr>
<a:t>Plan</a:t></a:r></a:p></dsp:txBody></dsp:sp>
<dsp:sp modelId="{SPACER}"><dsp:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="100" cy="100"/></a:xfrm>
<a:noFill/></dsp:spPr></dsp:sp>
</dsp:spTree></dsp:drawing>`

func TestParseDrawingShapesReadsGeometryFillAndText(t *testing.T) {
	shapes := ParseDrawingShapes([]byte(sampleDrawingXML))
	if len(shapes) != 2 {
		t.Fatalf("got %d shapes, want 2", len(shapes))
	}

	node := shapes[0]
	if node.ModelID != "{NODE-1}" {
		t.Errorf("ModelID = %q, want {NODE-1}", node.ModelID)
	}
	if node.X != 4042 || node.Y != 829182 || node.CX != 1838086 || node.CY != 1252799 {
		t.Errorf("geometry = %d,%d %dx%d", node.X, node.Y, node.CX, node.CY)
	}
	if node.RotationDeg != 90 {
		t.Errorf("RotationDeg = %v, want 90", node.RotationDeg)
	}
	if node.PresetGeom != "roundRect" {
		t.Errorf("PresetGeom = %q, want roundRect", node.PresetGeom)
	}
	if node.Adjustments["adj"] != 10000 {
		t.Errorf("adj = %d, want 10000", node.Adjustments["adj"])
	}
	if node.Fill.Scheme != "accent1" || node.Fill.LumOff != 20000 {
		t.Errorf("fill = %+v, want accent1 with lumOff 20000", node.Fill)
	}
	if node.Line.Scheme != "lt1" || node.LineWidthEMU != 19050 {
		t.Errorf("line = %+v w=%d", node.Line, node.LineWidthEMU)
	}
	if node.Anchor != "ctr" {
		t.Errorf("Anchor = %q, want ctr", node.Anchor)
	}
}

func TestParseDrawingShapesReadsRunFormatting(t *testing.T) {
	shapes := ParseDrawingShapes([]byte(sampleDrawingXML))
	if len(shapes) == 0 {
		t.Fatal("no shapes parsed")
	}
	paragraphs := shapes[0].Paragraphs
	if len(paragraphs) != 1 || len(paragraphs[0].Runs) != 1 {
		t.Fatalf("got %d paragraphs, want 1 with 1 run", len(paragraphs))
	}
	run := paragraphs[0].Runs[0]
	if run.Text != "Plan" {
		t.Errorf("Text = %q, want Plan", run.Text)
	}
	if run.SizePt != 18 {
		t.Errorf("SizePt = %v, want 18", run.SizePt)
	}
	if !run.Bold || run.Italic {
		t.Errorf("bold/italic = %v/%v, want true/false", run.Bold, run.Italic)
	}
	if run.Color.SRGB != "FFFFFF" {
		t.Errorf("Color = %+v, want FFFFFF", run.Color)
	}
	if run.Font != "Calibri" {
		t.Errorf("Font = %q, want Calibri", run.Font)
	}
	if paragraphs[0].Align != "ctr" {
		t.Errorf("Align = %q, want ctr", paragraphs[0].Align)
	}
}

func TestParseDrawingShapesNoFillLeavesColorUnset(t *testing.T) {
	shapes := ParseDrawingShapes([]byte(sampleDrawingXML))
	if len(shapes) < 2 {
		t.Fatalf("got %d shapes, want 2", len(shapes))
	}
	if shapes[1].Fill.IsSet() {
		t.Errorf("noFill shape carries a fill: %+v", shapes[1].Fill)
	}
	if shapes[1].HasText() {
		t.Error("spacer shape reports text")
	}
}

func TestParseDrawingShapesRejectsGarbage(t *testing.T) {
	if got := ParseDrawingShapes([]byte("not xml at all")); got != nil {
		t.Errorf("got %d shapes from garbage, want none", len(got))
	}
	empty := `<dsp:drawing xmlns:dsp="x"><dsp:spTree/></dsp:drawing>`
	if got := ParseDrawingShapes([]byte(empty)); got != nil {
		t.Errorf("got %d shapes from an empty drawing, want none", len(got))
	}
}
