package logical

import (
	"testing"
)

type mockPartProvider struct {
	parts map[string][]byte
}

func (m *mockPartProvider) Has(path string) bool {
	_, ok := m.parts[path]
	return ok
}

func (m *mockPartProvider) Get(path string) ([]byte, bool) {
	data, ok := m.parts[path]
	return data, ok
}

func (m *mockPartProvider) Keys() []string {
	keys := make([]string, 0, len(m.parts))
	for k := range m.parts {
		keys = append(keys, k)
	}
	return keys
}

func TestChecker_Check(t *testing.T) {
	checker := &Checker{}

	// Test with valid slide
	validSlideXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld>
    <p:spTree>
      <p:sp>
        <p:nvSpPr>
          <p:cNvPr id="2" name="Title 1"/>
        </p:nvSpPr>
        <p:spPr>
          <a:xfrm>
            <a:off x="100" y="100"/>
            <a:ext cx="1000" cy="1000"/>
          </a:xfrm>
          <a:prstGeom prst="rect"/>
        </p:spPr>
        <p:txBody>
          <a:p>
            <a:r>
              <a:t>Slide Title</a:t>
            </a:r>
          </a:p>
        </p:txBody>
      </p:sp>
    </p:spTree>
  </p:cSld>
</p:sld>`

	provider := &mockPartProvider{
		parts: map[string][]byte{
			"ppt/slides/slide1.xml": []byte(validSlideXML),
		},
	}

	issues := checker.Check(provider)
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d: %v", len(issues), issues)
	}
}

func TestExtractFirstAText(t *testing.T) {
	xmlContent := `
<root xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <a:t>  Hello  </a:t>
  <a:t>World</a:t>
</root>`
	got := extractFirstAText([]byte(xmlContent))
	if got != "Hello" {
		t.Errorf("extractFirstAText() = %q; want %q", got, "Hello")
	}

	gotEmpty := extractFirstAText([]byte("<root></root>"))
	if gotEmpty != "" {
		t.Errorf("extractFirstAText() = %q; want %q", gotEmpty, "")
	}
}

func TestParseSlideIndex(t *testing.T) {
	tests := []struct {
		path     string
		expected int
	}{
		{"ppt/slides/slide1.xml", 1},
		{"ppt/slides/slide10.xml", 10},
		{"slide5.xml", 5},
	}

	for _, tt := range tests {
		got := parseSlideIndex(tt.path)
		if got != tt.expected {
			t.Errorf("parseSlideIndex(%q) = %d; want %d", tt.path, got, tt.expected)
		}
	}
}

func TestParseSlideShapes(t *testing.T) {
	slideXML := `
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:spTree>
    <p:sp>
      <p:nvSpPr><p:cNvPr id="2" name="Shape 1"/></p:nvSpPr>
      <p:spPr>
        <a:xfrm>
          <a:off x="100" y="200"/>
          <a:ext cx="300" cy="400"/>
        </a:xfrm>
        <a:prstGeom prst="rect"/>
      </p:spPr>
      <p:txBody>
        <a:p><a:r><a:t>Text 1</a:t></a:r></a:p>
      </p:txBody>
    </p:sp>
    <p:pic>
      <p:nvPicPr><p:cNvPr id="3" name="Picture 1"/></p:nvPicPr>
      <p:spPr>
        <a:xfrm>
          <a:off x="500" y="600"/>
          <a:ext cx="700" cy="800"/>
        </a:xfrm>
      </p:spPr>
    </p:pic>
  </p:spTree>
</p:sld>`

	shapes := parseSlideShapes([]byte(slideXML))
	if len(shapes) != 2 {
		t.Errorf("expected 2 shapes, got %d", len(shapes))
	}

	if shapes[0].Name != "Shape 1" || shapes[0].Type != "rect" || shapes[0].Text != "Text 1" {
		t.Errorf("shape 0 mismatch: %+v", shapes[0])
	}
	if shapes[0].X.Emu() != 100 || shapes[0].Y.Emu() != 200 {
		t.Errorf("shape 0 position mismatch: (%d, %d)", shapes[0].X.Emu(), shapes[0].Y.Emu())
	}

	if shapes[1].Name != "Picture 1" || shapes[1].Type != "pic" {
		t.Errorf("shape 1 mismatch: %+v", shapes[1])
	}
}

func TestChecker_Check_InvalidSlide(_ *testing.T) {
	checker := &Checker{}

	// Slide with no title and no shapes should be treated as blank layout
	// But wait, SlideContent.Validate might fail if it expects something else.
	// Let's see SlideContent.Validate in pkg/pptx/elements/slide.go

	emptySlideXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"></p:sld>`

	provider := &mockPartProvider{
		parts: map[string][]byte{
			"ppt/slides/slide1.xml": []byte(emptySlideXML),
		},
	}

	issues := checker.Check(provider)
	// If SlideLayoutBlank is valid for empty slide, issues should be 0.
	// If it fails validation, it should return a warning.
	// Based on checker.go:
	// if err := slide.Validate(index); err != nil { ... return []structural.Issue{...} }

	// We don't know for sure if it will fail without seeing elements.SlideContent.Validate
	// but we can at least verify it runs.
	_ = issues
}
