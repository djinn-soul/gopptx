package shape

import "testing"

func TestParseShapePropertiesExtractsParagraphLineSpacingFromLnSpc(t *testing.T) {
	shapeXML := []byte(
		`<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
			`<p:nvSpPr><p:cNvPr id="22" name="Paragraph Advanced Parse"/></p:nvSpPr>` +
			`<p:spPr><a:xfrm><a:off x="10" y="20"/><a:ext cx="300" cy="200"/></a:xfrm></p:spPr>` +
			`<p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:pPr algn="ctr" lvl="1"><a:lnSpc><a:spcPct val="120000"/></a:lnSpc><a:spcBef><a:spcPts val="200"/></a:spcBef><a:spcAft><a:spcPts val="100"/></a:spcAft></a:pPr>` +
			`<a:r><a:rPr lang="en-US"/><a:t>Text</a:t></a:r></a:p></p:txBody></p:sp>`,
	)

	shape, err := ParseShapeProperties(shapeXML)
	if err != nil {
		t.Fatalf("parseShapeProperties failed: %v", err)
	}
	if shape.Paragraph == nil {
		t.Fatal("expected paragraph properties to be parsed")
	}
	if shape.Paragraph.LineSpacingPct == nil || *shape.Paragraph.LineSpacingPct != 120000 {
		t.Fatalf("expected line spacing pct 120000, got %#v", shape.Paragraph.LineSpacingPct)
	}
	if shape.Paragraph.SpaceBeforePts == nil || *shape.Paragraph.SpaceBeforePts != 200 {
		t.Fatalf("expected space before 200, got %#v", shape.Paragraph.SpaceBeforePts)
	}
	if shape.Paragraph.SpaceAfterPts == nil || *shape.Paragraph.SpaceAfterPts != 100 {
		t.Fatalf("expected space after 100, got %#v", shape.Paragraph.SpaceAfterPts)
	}
}
