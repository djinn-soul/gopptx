package editor

import (
	"strings"
	"testing"
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

	xmlBytes := renderShapeXML(s)
	xmlStr := string(xmlBytes)

	if !strings.Contains(xmlStr, `<p:cNvPr id="10" name="Test Shape"/>`) {
		t.Errorf("rendered XML missing ID/Name: %s", xmlStr)
	}
	if !strings.Contains(xmlStr, `<a:t>Updated Text</a:t>`) {
		t.Errorf("rendered XML missing text: %s", xmlStr)
	}
	if !strings.Contains(xmlStr, `x="50" y="60"`) {
		t.Errorf("rendered XML missing pos: %s", xmlStr)
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
	modified := replaceShapeNodes(original, shapes, func(i int, p *parsedShape) ([]byte, bool) {
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

	got := maxObjectID(xmlContent)
	if got != 9 {
		t.Fatalf("maxObjectID() = %d, want 9", got)
	}
}
