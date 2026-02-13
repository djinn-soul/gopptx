package editor

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestSlide_Placeholders_DiscoversPHElements(t *testing.T) {
	path := writeDeckFixture(t, "placeholder_discovery.pptx", []elements.SlideContent{
		elements.NewSlide("Discovery Test"),
	})
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("OpenPresentationEditor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	// Inject slide XML with placeholder elements
	slidePart := editor.slides[0].Part
	editor.parts.Set(slidePart, []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:cSld>
<p:spTree>
<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="2" name="Title 1"/>
    <p:cNvSpPr/>
    <p:nvPr>
      <p:ph type="title"/>
    </p:nvPr>
  </p:nvSpPr>
</p:sp>
<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="3" name="Body 2"/>
    <p:cNvSpPr/>
    <p:nvPr>
      <p:ph idx="1" type="body"/>
    </p:nvPr>
  </p:nvSpPr>
</p:sp>
<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="4" name="Footer 3"/>
    <p:cNvSpPr/>
    <p:nvPr>
      <p:ph idx="10" type="ftr"/>
    </p:nvPr>
  </p:nvSpPr>
</p:sp>
</p:spTree>
</p:cSld>
</p:sld>`))

	slide, err := editor.GetSlide(0)
	if err != nil {
		t.Fatalf("GetSlide: %v", err)
	}

	placeholders, err := slide.Placeholders()
	if err != nil {
		t.Fatalf("Placeholders: %v", err)
	}

	if len(placeholders) != 3 {
		t.Fatalf("expected 3 placeholders, got %d", len(placeholders))
	}

	// First placeholder: title (no idx attr means idx=0)
	if placeholders[0].Type != "title" {
		t.Errorf("expected type 'title', got %q", placeholders[0].Type)
	}
	if placeholders[0].Index != 0 {
		t.Errorf("expected index 0 for title, got %d", placeholders[0].Index)
	}

	// Second placeholder: body idx=1
	if placeholders[1].Type != "body" {
		t.Errorf("expected type 'body', got %q", placeholders[1].Type)
	}
	if placeholders[1].Index != 1 {
		t.Errorf("expected index 1 for body, got %d", placeholders[1].Index)
	}

	// Third placeholder: footer idx=10
	if placeholders[2].Type != "ftr" {
		t.Errorf("expected type 'ftr', got %q", placeholders[2].Type)
	}
	if placeholders[2].Index != 10 {
		t.Errorf("expected index 10 for footer, got %d", placeholders[2].Index)
	}
}
