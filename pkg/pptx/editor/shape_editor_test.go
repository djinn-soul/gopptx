package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

func TestParseSlideShapes(t *testing.T) {
	xmlContent := `
<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <p:cSld>
	<p:spTree>
	  <p:nvGrpSpPr/>
	  <p:grpSpPr/>
	  <p:sp>
		<p:nvSpPr><p:cNvPr id="2" name="Title 1"/></p:nvSpPr>
		<p:spPr>
			<a:xfrm><a:off x="100" y="100"/><a:ext cx="500" cy="50"/></a:xfrm>
		</p:spPr>
		<p:txBody><p:p><p:r><a:t>Hello World</a:t></p:r></p:p></p:txBody>
	  </p:sp>
	  <p:pic>
		<p:nvPicPr><p:cNvPr id="3" name="Picture 2"/></p:nvPicPr>
		<p:spPr/>
	  </p:pic>
	</p:spTree>
  </p:cSld>
</p:sld>
`
	shapes, err := parseSlideShapes([]byte(xmlContent))
	if err != nil {
		t.Fatalf("parseSlideShapes failed: %v", err)
	}

	if len(shapes) != 2 {
		t.Errorf("expected 2 shapes, got %d", len(shapes))
	}

	// Check Shape 1
	if shapes[0].Type != "sp" {
		t.Errorf("expected type sp, got %s", shapes[0].Type)
	}
	if shapes[0].ID != 2 {
		t.Errorf("expected ID 2, got %d", shapes[0].ID)
	}
	if shapes[0].Name != "Title 1" {
		t.Errorf("expected Name 'Title 1', got '%s'", shapes[0].Name)
	}
	if strings.TrimSpace(shapes[0].Text) != "Hello World" {
		t.Errorf("expected Text 'Hello World', got '%s'", shapes[0].Text)
	}
	if shapes[0].X != 100 {
		t.Errorf("expected X 100, got %d", shapes[0].X)
	}

	// Check Shape 2 (Picture)
	if shapes[1].Type != "pic" {
		t.Errorf("expected type pic, got %s", shapes[1].Type)
	}
	if shapes[1].ID != 3 {
		t.Errorf("expected ID 3, got %d", shapes[1].ID)
	}
}

func TestRenderShapeXML(t *testing.T) {
	s := &parsedShape{
		ID:   10,
		Name: "Test Shape",
		Type: "sp",
		Text: "Updated Text",
		X:    50,
		Y:    60,
		W:    200,
		H:    100,
	}

	// Create a minimal editor instance to call the method
	e := &PresentationEditor{
		nextRelIDNum: 1,
	}
	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)

	if !strings.Contains(xmlStr, `<p:cNvPr id="10" name="Test Shape">`) {
		t.Errorf("rendered XML missing ID/Name: %s", xmlStr)
	}
	if !strings.Contains(xmlStr, `<a:t>Updated Text</a:t>`) {
		t.Errorf("rendered XML missing text: %s", xmlStr)
	}
	if !strings.Contains(xmlStr, `x="50" y="60"`) {
		t.Errorf("rendered XML missing pos: %s", xmlStr)
	}
}

func TestRenderShapeXMLWithTextFrameOrientationAndColumns(t *testing.T) {
	orientation := "vertical270"
	columns := 2
	s := &parsedShape{
		ID:   11,
		Name: "Text Frame Shape",
		Type: "sp",
		Text: "Two columns",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		TextFrame: &common.TextFrame{
			Orientation: &orientation,
			Columns:     &columns,
		},
	}

	e := &PresentationEditor{nextRelIDNum: 1}
	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)

	if !strings.Contains(xmlStr, `vert="vert270"`) {
		t.Fatalf("expected text-frame orientation in bodyPr, got: %s", xmlStr)
	}
	if !strings.Contains(xmlStr, `numCol="2"`) {
		t.Fatalf("expected text-frame columns in bodyPr, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLWithTextFrameInvalidOrientation(t *testing.T) {
	orientation := "diagonal"
	s := &parsedShape{
		ID:   12,
		Name: "Invalid Orientation",
		Type: "sp",
		Text: "Invalid",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		TextFrame: &common.TextFrame{
			Orientation: &orientation,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}

	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected invalid orientation to fail")
	}
	if !strings.Contains(err.Error(), "unsupported text_frame.orientation") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLWithTextFrameInvalidColumns(t *testing.T) {
	columns := 0
	s := &parsedShape{
		ID:   13,
		Name: "Invalid Columns",
		Type: "sp",
		Text: "Invalid",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		TextFrame: &common.TextFrame{
			Columns: &columns,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}

	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected invalid columns to fail")
	}
	if !strings.Contains(err.Error(), "text_frame.columns must be >=") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLWithTextFrameRotation(t *testing.T) {
	rotation := 45.0
	s := &parsedShape{
		ID:   14,
		Name: "Rotated Text",
		Type: "sp",
		Text: "Rotated",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		TextFrame: &common.TextFrame{
			Rotation: &rotation,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}

	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `rot="2700000"`) {
		t.Fatalf("expected text-frame rotation in bodyPr, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLWithTextFrameInvalidRotation(t *testing.T) {
	rotation := 720.0
	s := &parsedShape{
		ID:   15,
		Name: "Invalid Rotation",
		Type: "sp",
		Text: "Invalid",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		TextFrame: &common.TextFrame{
			Rotation: &rotation,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}

	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected invalid rotation to fail")
	}
	if !strings.Contains(err.Error(), "text_frame.rotation must be between -360 and 360") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLWithParagraphIndentAndHanging(t *testing.T) {
	indent := 228600
	hanging := 228600
	s := &parsedShape{
		ID:   16,
		Name: "Paragraph Shape",
		Type: "sp",
		Text: "Indented",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Paragraph: &common.Paragraph{
			Indent:  &indent,
			Hanging: &hanging,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}

	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `marL="228600"`) {
		t.Fatalf("expected paragraph marL, got: %s", xmlStr)
	}
	if !strings.Contains(xmlStr, `indent="-228600"`) {
		t.Fatalf("expected paragraph hanging indent, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLWithParagraphInvalidHanging(t *testing.T) {
	hanging := -1
	s := &parsedShape{
		ID:   17,
		Name: "Invalid Paragraph",
		Type: "sp",
		Text: "Invalid",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Paragraph: &common.Paragraph{
			Hanging: &hanging,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}

	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected invalid paragraph hanging to fail")
	}
	if !strings.Contains(err.Error(), "paragraph.hanging must be >=") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLWithHoverAction(t *testing.T) {
	macro := "HoverMacro"
	s := &parsedShape{
		ID:   18,
		Name: "Hover Shape",
		Type: "sp",
		Text: "Hover",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		HoverAction: &common.Hyperlink{
			Macro: &macro,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}

	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `hlinkMouseOver`) {
		t.Fatalf("expected hover action in cNvPr, got: %s", xmlStr)
	}
	if !strings.Contains(xmlStr, `ppaction://macro?name=HoverMacro`) {
		t.Fatalf("expected macro action URL, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLWithFillAndLine(t *testing.T) {
	fill := "FF0000"
	lineColor := "00FF00"
	lineWidth := 25400
	lineDash := "dashDot"
	s := &parsedShape{
		ID:   19,
		Name: "Styled Shape",
		Type: "sp",
		Text: "Styled",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Fill: &common.ShapeFill{Solid: &fill},
		Line: &common.ShapeLine{Color: &lineColor, WidthEmu: &lineWidth, DashStyle: &lineDash},
	}
	e := &PresentationEditor{nextRelIDNum: 1}

	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `<a:solidFill><a:srgbClr val="FF0000"/></a:solidFill>`) {
		t.Fatalf("expected solid fill in shape XML, got: %s", xmlStr)
	}
	if !strings.Contains(
		xmlStr,
		`<a:ln w="25400"><a:prstDash val="dashDot"/><a:solidFill><a:srgbClr val="00FF00"/></a:solidFill></a:ln>`,
	) {
		t.Fatalf("expected line style in shape XML, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLStyleOrderingBeforePresetGeometry(t *testing.T) {
	fill := "112233"
	lineColor := "445566"
	lineDash := "dash"
	width := 25400
	shadowColor := "778899"
	shadowDist := 20000
	glowColor := "AABBCC"
	glowRadius := 40000
	s := &parsedShape{
		ID:   190,
		Name: "Ordered Style",
		Type: "sp",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Fill: &common.ShapeFill{Solid: &fill},
		Line: &common.ShapeLine{Color: &lineColor, WidthEmu: &width, DashStyle: &lineDash},
		Shadow: &common.ShapeShadow{
			Color:       &shadowColor,
			DistanceEmu: &shadowDist,
		},
		Glow: &common.ShapeGlow{
			Color:     &glowColor,
			RadiusEmu: &glowRadius,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	idxFill := strings.Index(xmlStr, `<a:solidFill><a:srgbClr val="112233"/></a:solidFill>`)
	idxLine := strings.Index(
		xmlStr,
		`<a:ln w="25400"><a:prstDash val="dash"/><a:solidFill><a:srgbClr val="445566"/></a:solidFill></a:ln>`,
	)
	idxEffect := strings.Index(xmlStr, `<a:effectLst><a:outerShdw`)
	idxGeom := strings.Index(xmlStr, `<a:prstGeom`)
	if idxFill == -1 || idxLine == -1 || idxEffect == -1 || idxGeom == -1 {
		t.Fatalf("missing expected style/geom tokens: %s", xmlStr)
	}
	if idxFill >= idxLine || idxLine >= idxEffect || idxEffect >= idxGeom {
		t.Fatalf("unexpected style ordering fill/line/effect/geom: %s", xmlStr)
	}
}

func TestRenderShapeXMLRejectsUnsupportedLineDashStyle(t *testing.T) {
	lineDash := "zigzag"
	s := &parsedShape{
		ID:   119,
		Name: "Invalid Dash",
		Type: "sp",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Line: &common.ShapeLine{DashStyle: &lineDash},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected invalid line dash style to fail")
	}
	if !strings.Contains(err.Error(), "line.dash_style") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReplaceShapeClickActionPreservesExistingCNvPrChildren(t *testing.T) {
	xmlIn := []byte(
		`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Shape 1"><a:extLst><a:ext uri="{A}"/></a:extLst><a:hlinkClick action="ppaction://old"/></p:cNvPr></p:nvSpPr></p:sp>`,
	)
	action := "ppaction://hlinksldjump"
	clickAction := &common.Hyperlink{Action: &action}

	got, err := replaceShapeClickAction(&PresentationEditor{}, "ppt/slides/slide1.xml", xmlIn, clickAction)
	if err != nil {
		t.Fatalf("replaceShapeClickAction failed: %v", err)
	}
	out := string(got)

	if !strings.Contains(out, `<a:extLst><a:ext uri="{A}"/></a:extLst>`) {
		t.Fatalf("existing cNvPr child extension list was not preserved: %s", out)
	}
	if strings.Contains(out, `action="ppaction://old"`) {
		t.Fatalf("stale hyperlink click action was not removed: %s", out)
	}
	if !strings.Contains(out, `action="ppaction://hlinksldjump"`) {
		t.Fatalf("new hyperlink click action missing: %s", out)
	}
}

func TestParseShapePropertiesExtractsParagraphIndentAndHanging(t *testing.T) {
	shapeXML := []byte(
		`<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
			`<p:nvSpPr><p:cNvPr id="20" name="Paragraph Parse"/></p:nvSpPr>` +
			`<p:spPr><a:xfrm><a:off x="10" y="20"/><a:ext cx="300" cy="200"/></a:xfrm></p:spPr>` +
			`<p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:pPr marL="228600" indent="-114300"/>` +
			`<a:r><a:rPr lang="en-US"/><a:t>Text</a:t></a:r></a:p></p:txBody></p:sp>`,
	)

	shape, err := parseShapeProperties(shapeXML)
	if err != nil {
		t.Fatalf("parseShapeProperties failed: %v", err)
	}
	if shape.Paragraph == nil {
		t.Fatal("expected paragraph properties to be parsed")
	}
	if shape.Paragraph.Indent == nil || *shape.Paragraph.Indent != 228600 {
		t.Fatalf("expected indent 228600, got %#v", shape.Paragraph.Indent)
	}
	if shape.Paragraph.Hanging == nil || *shape.Paragraph.Hanging != 114300 {
		t.Fatalf("expected hanging 114300, got %#v", shape.Paragraph.Hanging)
	}
}

func TestBuildClickActionXMLCreatesSlideRelsWhenMissing(t *testing.T) {
	parts := NewPartStore()
	e := &PresentationEditor{parts: parts}
	addr := "https://example.com"

	xml, err := e.buildClickActionXML("ppt/slides/slide1.xml", &common.Hyperlink{Address: &addr})
	if err != nil {
		t.Fatalf("buildClickActionXML failed: %v", err)
	}
	if !strings.Contains(xml, `r:id="rId1"`) {
		t.Fatalf("expected hyperlink xml to reference rId1, got: %s", xml)
	}
	rels, ok := parts.Get(common.SlideRelsPartName("ppt/slides/slide1.xml"))
	if !ok {
		t.Fatal("expected slide rels part to be created")
	}
	if !strings.Contains(string(rels), "/relationships/hyperlink") {
		t.Fatalf("expected hyperlink relationship in rels part, got: %s", string(rels))
	}
}

func TestBuildClickActionXMLTargetSlideCreatesSlideRelationship(t *testing.T) {
	parts := NewPartStore()
	parts.Set("ppt/slides/_rels/slide1.xml.rels", []byte(
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`,
	))
	e := &PresentationEditor{
		parts: parts,
		slides: []common.EditorSlideRef{
			{Part: "ppt/slides/slide1.xml"},
			{Part: "ppt/slides/slide2.xml"},
		},
	}
	target := 1
	xml, err := e.buildClickActionXML("ppt/slides/slide1.xml", &common.Hyperlink{TargetSlide: &target})
	if err != nil {
		t.Fatalf("buildClickActionXML failed: %v", err)
	}
	if !strings.Contains(xml, `hlinkClick`) || !strings.Contains(xml, `ppaction://hlinksldjump`) {
		t.Fatalf("unexpected click action XML: %s", xml)
	}
	rels, ok := parts.Get("ppt/slides/_rels/slide1.xml.rels")
	if !ok {
		t.Fatal("expected slide rels part")
	}
	relsStr := string(rels)
	if !strings.Contains(relsStr, `/relationships/slide`) {
		t.Fatalf("expected slide relationship in rels part, got: %s", relsStr)
	}
	if !strings.Contains(relsStr, `Target="slide2.xml"`) {
		t.Fatalf("expected target slide relationship to slide2.xml, got: %s", relsStr)
	}
}

func TestBuildClickActionXMLRejectsMutuallyExclusiveSelectors(t *testing.T) {
	e := &PresentationEditor{nextRelIDNum: 1}
	addr := "https://example.com"
	jump := "nextslide"
	_, err := e.buildClickActionXML("ppt/slides/slide1.xml", &common.Hyperlink{
		Address:    &addr,
		TargetJump: &jump,
	})
	if err == nil {
		t.Fatal("expected mutually exclusive hyperlink selector validation error")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildClickActionXMLRejectsUnsupportedJump(t *testing.T) {
	e := &PresentationEditor{nextRelIDNum: 1}
	jump := "homeslide"
	_, err := e.buildClickActionXML("ppt/slides/slide1.xml", &common.Hyperlink{
		TargetJump: &jump,
	})
	if err == nil {
		t.Fatal("expected unsupported jump validation error")
	}
	if !strings.Contains(err.Error(), "unsupported jump target") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReplaceShapeClickActionErrorsWithoutCNvPr(t *testing.T) {
	xmlIn := []byte(`<p:sp><p:spPr/></p:sp>`)
	action := "ppaction://hlinksldjump"
	clickAction := &common.Hyperlink{Action: &action}

	_, err := replaceShapeClickAction(&PresentationEditor{}, "ppt/slides/slide1.xml", xmlIn, clickAction)
	if err == nil {
		t.Fatal("expected error when cNvPr is missing for click action update")
	}
	if !strings.Contains(err.Error(), "no cNvPr") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReplaceShapeStyleUpdatesSpPrFillAndLine(t *testing.T) {
	xmlIn := []byte(
		`<p:sp><p:spPr><a:xfrm/><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr></p:sp>`,
	)
	fill := "112233"
	lineColor := "334455"
	width := 38100
	lineDash := "long_dash_dot"
	got, err := replaceShapeStyle(
		xmlIn,
		&common.ShapeFill{Solid: &fill},
		&common.ShapeLine{Color: &lineColor, WidthEmu: &width, DashStyle: &lineDash},
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("replaceShapeStyle failed: %v", err)
	}
	out := string(got)
	if !strings.Contains(out, `val="112233"`) {
		t.Fatalf("expected fill color in updated style: %s", out)
	}
	if !strings.Contains(out, `a:ln w="38100"`) {
		t.Fatalf("expected line width in updated style: %s", out)
	}
	if !strings.Contains(out, `a:prstDash val="lgDashDot"`) {
		t.Fatalf("expected line dash style in updated style: %s", out)
	}
	if !strings.Contains(out, `val="334455"`) {
		t.Fatalf("expected line color in updated style: %s", out)
	}
}

func TestReplaceShapeStyleRemovesOldStyleNodesAndKeepsOrdering(t *testing.T) {
	xmlIn := []byte(
		`<p:sp><p:spPr><a:xfrm/><a:pattFill prst="pct5"/><a:ln w="12700"><a:prstDash val="sysDot"/></a:ln>` +
			`<a:effectLst><a:glow rad="1000"/></a:effectLst><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr></p:sp>`,
	)
	fill := "ABCDEF"
	dash := "round_dot"
	width := 38100
	got, err := replaceShapeStyle(
		xmlIn,
		&common.ShapeFill{Solid: &fill},
		&common.ShapeLine{WidthEmu: &width, DashStyle: &dash},
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("replaceShapeStyle failed: %v", err)
	}
	out := string(got)
	if strings.Contains(out, `<a:pattFill`) || strings.Contains(out, `<a:effectLst`) {
		t.Fatalf("expected old style nodes removed, got: %s", out)
	}
	idxFill := strings.Index(out, `<a:solidFill><a:srgbClr val="ABCDEF"/></a:solidFill>`)
	idxLine := strings.Index(out, `<a:ln w="38100"><a:prstDash val="sysDot"/></a:ln>`)
	idxGeom := strings.Index(out, `<a:prstGeom`)
	if idxFill == -1 || idxLine == -1 || idxGeom == -1 {
		t.Fatalf("expected fill+line+geom after replace, got: %s", out)
	}
	if idxFill >= idxLine || idxLine >= idxGeom {
		t.Fatalf("unexpected ordering after replace: %s", out)
	}
}

func TestReplaceShapeStyleSelectivePreservesSchemeFillWhenFillNotUpdated(t *testing.T) {
	xmlIn := []byte(
		`<p:sp><p:spPr><a:xfrm/><a:solidFill><a:schemeClr val="accent1"/></a:solidFill>` +
			`<a:ln w="12700"><a:solidFill><a:srgbClr val="222222"/></a:solidFill></a:ln>` +
			`<a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr></p:sp>`,
	)
	width := 25400
	lineColor := "445566"
	got, err := replaceShapeStyleSelective(
		xmlIn,
		nil,
		&common.ShapeLine{Color: &lineColor, WidthEmu: &width},
		nil,
		nil,
		nil,
		nil,
		nil,
		false,
		true,
		false,
	)
	if err != nil {
		t.Fatalf("replaceShapeStyleSelective failed: %v", err)
	}
	out := string(got)
	if !strings.Contains(out, `<a:schemeClr val="accent1"/>`) {
		t.Fatalf("expected scheme fill to be preserved, got: %s", out)
	}
	if !strings.Contains(out, `<a:ln w="25400"><a:solidFill><a:srgbClr val="445566"/></a:solidFill></a:ln>`) {
		t.Fatalf("expected updated line to be present, got: %s", out)
	}
}

func TestParseShapePropertiesExtractsLineDashStyle(t *testing.T) {
	shapeXML := []byte(
		`<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
			`<p:nvSpPr><p:cNvPr id="120" name="Line Parse"/></p:nvSpPr>` +
			`<p:spPr><a:xfrm><a:off x="10" y="20"/><a:ext cx="300" cy="200"/></a:xfrm><a:ln w="12700"><a:prstDash val="sysDashDot"/></a:ln></p:spPr>` +
			`</p:sp>`,
	)
	shape, err := parseShapeProperties(shapeXML)
	if err != nil {
		t.Fatalf("parseShapeProperties failed: %v", err)
	}
	if shape.Line == nil || shape.Line.DashStyle == nil {
		t.Fatalf("expected parsed line dash style, got %#v", shape.Line)
	}
	if *shape.Line.DashStyle != "sysDashDot" {
		t.Fatalf("expected dash style sysDashDot, got %q", *shape.Line.DashStyle)
	}
}

func TestParseShapePropertiesExtractsBackgroundFill(t *testing.T) {
	shapeXML := []byte(
		`<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
			`<p:nvSpPr><p:cNvPr id="121" name="Fill Parse"/></p:nvSpPr>` +
			`<p:spPr><a:xfrm><a:off x="10" y="20"/><a:ext cx="300" cy="200"/></a:xfrm><a:noFill/></p:spPr>` +
			`</p:sp>`,
	)
	shape, err := parseShapeProperties(shapeXML)
	if err != nil {
		t.Fatalf("parseShapeProperties failed: %v", err)
	}
	if shape.Fill == nil || shape.Fill.Background == nil || !*shape.Fill.Background {
		t.Fatalf("expected parsed background fill, got %#v", shape.Fill)
	}
}

func TestParseShapePropertiesExtractsBlurEffect(t *testing.T) {
	shapeXML := []byte(
		`<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
			`<p:nvSpPr><p:cNvPr id="122" name="Blur Parse"/></p:nvSpPr>` +
			`<p:spPr><a:xfrm><a:off x="10" y="20"/><a:ext cx="300" cy="200"/></a:xfrm><a:effectLst><a:blur rad="61000"/></a:effectLst></p:spPr>` +
			`</p:sp>`,
	)
	shape, err := parseShapeProperties(shapeXML)
	if err != nil {
		t.Fatalf("parseShapeProperties failed: %v", err)
	}
	if shape.Blur == nil || shape.Blur.RadiusEmu == nil || *shape.Blur.RadiusEmu != 61000 {
		t.Fatalf("expected parsed blur effect, got %#v", shape.Blur)
	}
}

func TestParseShapePropertiesExtractsSoftEdgeEffect(t *testing.T) {
	shapeXML := []byte(
		`<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
			`<p:nvSpPr><p:cNvPr id="123" name="SoftEdge Parse"/></p:nvSpPr>` +
			`<p:spPr><a:xfrm><a:off x="10" y="20"/><a:ext cx="300" cy="200"/></a:xfrm><a:effectLst><a:softEdge rad="62000"/></a:effectLst></p:spPr>` +
			`</p:sp>`,
	)
	shape, err := parseShapeProperties(shapeXML)
	if err != nil {
		t.Fatalf("parseShapeProperties failed: %v", err)
	}
	if shape.SoftEdge == nil || shape.SoftEdge.RadiusEmu == nil || *shape.SoftEdge.RadiusEmu != 62000 {
		t.Fatalf("expected parsed soft-edge effect, got %#v", shape.SoftEdge)
	}
}

func TestParseShapePropertiesExtractsReflectionEffect(t *testing.T) {
	shapeXML := []byte(
		`<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
			`<p:nvSpPr><p:cNvPr id="124" name="Reflection Parse"/></p:nvSpPr>` +
			`<p:spPr><a:xfrm><a:off x="10" y="20"/><a:ext cx="300" cy="200"/></a:xfrm><a:effectLst><a:reflection blurRad="14000" dist="9000"/></a:effectLst></p:spPr>` +
			`</p:sp>`,
	)
	shape, err := parseShapeProperties(shapeXML)
	if err != nil {
		t.Fatalf("parseShapeProperties failed: %v", err)
	}
	if shape.Reflection == nil ||
		shape.Reflection.BlurEmu == nil ||
		shape.Reflection.DistanceEmu == nil ||
		*shape.Reflection.BlurEmu != 14000 ||
		*shape.Reflection.DistanceEmu != 9000 {
		t.Fatalf("expected parsed reflection effect, got %#v", shape.Reflection)
	}
}

func TestRenderShapeXMLWithShadow(t *testing.T) {
	color := "123456"
	blur := 60000
	dist := 40000
	angle := 45.0
	s := &parsedShape{
		ID:   20,
		Name: "Shadowed Shape",
		Type: "sp",
		Text: "Shadow",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Shadow: &common.ShapeShadow{
			Color:       &color,
			BlurEmu:     &blur,
			DistanceEmu: &dist,
			AngleDeg:    &angle,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}

	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `<a:effectLst><a:outerShdw`) {
		t.Fatalf("expected effectLst outer shadow in shape XML, got: %s", xmlStr)
	}
	if !strings.Contains(xmlStr, `blurRad="60000"`) || !strings.Contains(xmlStr, `dist="40000"`) {
		t.Fatalf("expected blur/dist in shadow XML, got: %s", xmlStr)
	}
	if !strings.Contains(xmlStr, `dir="2700000"`) || !strings.Contains(xmlStr, `val="123456"`) {
		t.Fatalf("expected direction/color in shadow XML, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLWithInvalidShadowCombination(t *testing.T) {
	inherit := true
	color := "123456"
	s := &parsedShape{
		ID:   21,
		Name: "Invalid Shadow",
		Type: "sp",
		Text: "Invalid",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Shadow: &common.ShapeShadow{
			Inherit: &inherit,
			Color:   &color,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected invalid shadow combination to fail")
	}
	if !strings.Contains(err.Error(), "shadow.inherit") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLWithShadowInheritFalseEmitsEmptyEffectList(t *testing.T) {
	inherit := false
	s := &parsedShape{
		ID:   25,
		Name: "Inherit False",
		Type: "sp",
		Text: "No effect",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Shadow: &common.ShapeShadow{
			Inherit: &inherit,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `<a:effectLst/>`) {
		t.Fatalf("expected empty effect list, got: %s", xmlStr)
	}
	if strings.Contains(xmlStr, `outerShdw`) {
		t.Fatalf("expected no explicit shadow when inherit is false, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLWithGradientFill(t *testing.T) {
	angle := 90.0
	s := &parsedShape{
		ID:   22,
		Name: "Gradient Fill",
		Type: "sp",
		Text: "Gradient",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Fill: &common.ShapeFill{
			Gradient: &common.GradientFill{
				AngleDeg: &angle,
				Stops: []common.GradientStop{
					{Color: "FF0000"},
					{Color: "0000FF"},
				},
			},
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `<a:gradFill>`) || !strings.Contains(xmlStr, `ang="5400000"`) {
		t.Fatalf("expected gradient fill with angle, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLWithPatternFill(t *testing.T) {
	preset := "diagCross"
	fg := "112233"
	bg := "AABBCC"
	s := &parsedShape{
		ID:   23,
		Name: "Pattern Fill",
		Type: "sp",
		Text: "Pattern",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Fill: &common.ShapeFill{
			Pattern: &common.PatternedFill{
				Preset:  &preset,
				FgColor: &fg,
				BgColor: &bg,
			},
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `a:pattFill prst="diagCross"`) {
		t.Fatalf("expected pattern fill preset, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLWithBackgroundFill(t *testing.T) {
	background := true
	s := &parsedShape{
		ID:   25,
		Name: "Background Fill",
		Type: "sp",
		Text: "Background",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Fill: &common.ShapeFill{
			Background: &background,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `<a:noFill/>`) {
		t.Fatalf("expected background/noFill token, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLRejectsMutuallyExclusiveFillModes(t *testing.T) {
	solid := "FF0000"
	s := &parsedShape{
		ID:   24,
		Name: "Invalid Fill",
		Type: "sp",
		Text: "Invalid",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Fill: &common.ShapeFill{
			Solid:    &solid,
			Gradient: &common.GradientFill{Stops: []common.GradientStop{{Color: "000000"}}},
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected mutually exclusive fill mode error")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLRejectsBackgroundWithOtherFillModes(t *testing.T) {
	solid := "FF0000"
	background := true
	s := &parsedShape{
		ID:   241,
		Name: "Invalid Background Mode",
		Type: "sp",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Fill: &common.ShapeFill{
			Solid:      &solid,
			Background: &background,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected mutually exclusive fill mode error")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLRejectsInvalidBackgroundFillFlag(t *testing.T) {
	background := false
	s := &parsedShape{
		ID:   240,
		Name: "Invalid Background Fill",
		Type: "sp",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Fill: &common.ShapeFill{
			Background: &background,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected invalid background fill flag to fail")
	}
	if !strings.Contains(err.Error(), "fill.background") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLWithGlow(t *testing.T) {
	color := "ABCDEF"
	radius := 50000
	s := &parsedShape{
		ID:   26,
		Name: "Glow Shape",
		Type: "sp",
		Text: "Glow",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Glow: &common.ShapeGlow{
			Color:     &color,
			RadiusEmu: &radius,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `<a:effectLst><a:glow rad="50000"><a:srgbClr val="ABCDEF"/></a:glow></a:effectLst>`) {
		t.Fatalf("expected glow effect XML, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLWithBlur(t *testing.T) {
	radius := 61000
	s := &parsedShape{
		ID:   28,
		Name: "Blur Shape",
		Type: "sp",
		Text: "Blur",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Blur: &common.ShapeBlur{
			RadiusEmu: &radius,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `<a:effectLst><a:blur rad="61000"/></a:effectLst>`) {
		t.Fatalf("expected blur effect XML, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLWithSoftEdge(t *testing.T) {
	radius := 62000
	s := &parsedShape{
		ID:   31,
		Name: "SoftEdge Shape",
		Type: "sp",
		Text: "SoftEdge",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		SoftEdge: &common.ShapeSoftEdge{
			RadiusEmu: &radius,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `<a:effectLst><a:softEdge rad="62000"/></a:effectLst>`) {
		t.Fatalf("expected soft-edge effect XML, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLRejectsInvalidSoftEdgeRadius(t *testing.T) {
	radius := -1
	s := &parsedShape{
		ID:   32,
		Name: "Invalid SoftEdge",
		Type: "sp",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		SoftEdge: &common.ShapeSoftEdge{
			RadiusEmu: &radius,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected invalid soft-edge radius to fail")
	}
	if !strings.Contains(err.Error(), "soft_edge.radius_emu") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLRejectsShadowInheritWithSoftEdge(t *testing.T) {
	inherit := true
	radius := 62000
	s := &parsedShape{
		ID:   33,
		Name: "Invalid Effects SoftEdge",
		Type: "sp",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Shadow: &common.ShapeShadow{
			Inherit: &inherit,
		},
		SoftEdge: &common.ShapeSoftEdge{
			RadiusEmu: &radius,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected shadow.inherit + soft-edge combination to fail")
	}
	if !strings.Contains(err.Error(), "shadow.inherit") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLWithReflection(t *testing.T) {
	blur := 14000
	dist := 9000
	s := &parsedShape{
		ID:   34,
		Name: "Reflection Shape",
		Type: "sp",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Reflection: &common.ShapeReflection{
			BlurEmu:     &blur,
			DistanceEmu: &dist,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	xmlBytes, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err != nil {
		t.Fatalf("renderShapeXML failed: %v", err)
	}
	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, `<a:effectLst><a:reflection blurRad="14000" dist="9000"/></a:effectLst>`) {
		t.Fatalf("expected reflection effect XML, got: %s", xmlStr)
	}
}

func TestRenderShapeXMLRejectsInvalidReflectionValues(t *testing.T) {
	blur := -1
	s := &parsedShape{
		ID:   35,
		Name: "Invalid Reflection",
		Type: "sp",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Reflection: &common.ShapeReflection{
			BlurEmu: &blur,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected invalid reflection values to fail")
	}
	if !strings.Contains(err.Error(), "reflection.blur_emu") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLRejectsShadowInheritWithReflection(t *testing.T) {
	inherit := true
	dist := 9000
	s := &parsedShape{
		ID:   36,
		Name: "Invalid Effects Reflection",
		Type: "sp",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Shadow: &common.ShapeShadow{
			Inherit: &inherit,
		},
		Reflection: &common.ShapeReflection{
			DistanceEmu: &dist,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected shadow.inherit + reflection combination to fail")
	}
	if !strings.Contains(err.Error(), "shadow.inherit") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLRejectsInvalidBlurRadius(t *testing.T) {
	radius := -1
	s := &parsedShape{
		ID:   29,
		Name: "Invalid Blur",
		Type: "sp",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Blur: &common.ShapeBlur{
			RadiusEmu: &radius,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected invalid blur radius to fail")
	}
	if !strings.Contains(err.Error(), "blur.radius_emu") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLRejectsShadowInheritWithBlur(t *testing.T) {
	inherit := true
	radius := 61000
	s := &parsedShape{
		ID:   30,
		Name: "Invalid Effects Blur",
		Type: "sp",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Shadow: &common.ShapeShadow{
			Inherit: &inherit,
		},
		Blur: &common.ShapeBlur{
			RadiusEmu: &radius,
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected shadow.inherit + blur combination to fail")
	}
	if !strings.Contains(err.Error(), "shadow.inherit") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderShapeXMLRejectsShadowInheritWithGlow(t *testing.T) {
	inherit := true
	s := &parsedShape{
		ID:   27,
		Name: "Invalid Effects",
		Type: "sp",
		Text: "Invalid",
		X:    10,
		Y:    20,
		W:    300,
		H:    200,
		Shadow: &common.ShapeShadow{
			Inherit: &inherit,
		},
		Glow: &common.ShapeGlow{
			RadiusEmu: ptrInt(50000),
		},
	}
	e := &PresentationEditor{nextRelIDNum: 1}
	_, err := e.renderShapeXML("ppt/slides/slide1.xml", s)
	if err == nil {
		t.Fatal("expected shadow.inherit+glow validation error")
	}
	if !strings.Contains(err.Error(), "shadow.inherit") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReplaceShapeNodes(t *testing.T) {
	// Use valid XML that matches what our parser expects (namespaces) to ensure correct parsing
	original := []byte(
		`<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><p:sp>Shape1</p:sp> MIDDLE <p:sp>Shape2</p:sp></p:sld>`,
	)

	// Parse to get real offsets
	shapes, err := parseSlideShapes(original)
	if err != nil {
		t.Fatalf("setup parse failed: %v", err)
	}
	if len(shapes) != 2 {
		t.Fatalf("expected 2 shapes, got %d", len(shapes))
	}

	// Modify second shape only
	modified := replaceShapeNodes(original, shapes, func(i int, _ *parsedShape) ([]byte, bool) {
		if i == 1 {
			return []byte(`<p:sp>REPLACED</p:sp>`), true
		}
		return nil, false
	})

	expected := `<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><p:sp>Shape1</p:sp> MIDDLE <p:sp>REPLACED</p:sp></p:sld>`
	if string(modified) != expected {
		t.Errorf("replace mismatch.\nExpected: %s\nGot:      %s", expected, string(modified))
	}
}

func TestMaxObjectIDIncludesGraphicFrame(t *testing.T) {
	xmlContent := []byte(`
<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld>
    <p:spTree>
      <p:nvGrpSpPr><p:cNvPr id="1" name=""/></p:nvGrpSpPr>
      <p:sp><p:nvSpPr><p:cNvPr id="3" name="Title"/></p:nvSpPr></p:sp>
      <p:graphicFrame><p:nvGraphicFramePr><p:cNvPr id="9" name="Chart 1"/></p:nvGraphicFramePr></p:graphicFrame>
    </p:spTree>
  </p:cSld>
</p:sld>`)

	got := editorshape.MaxObjectID(xmlContent, cNvPrIDPattern, cNvPrSubmatchSize)
	if got != 9 {
		t.Fatalf("maxObjectID() = %d, want 9", got)
	}
}

func ptrInt(value int) *int {
	return &value
}
