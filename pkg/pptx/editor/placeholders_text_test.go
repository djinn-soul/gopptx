package editor

import "testing"

// list_placeholders reports the text already present in a placeholder so that
// callers can read placeholder content without a separate shape lookup.
func TestParsePlaceholdersFromSlideXMLCapturesText(t *testing.T) {
	slideXML := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"` +
		` xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
		`<p:cSld><p:spTree>` +
		`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Title 1"/><p:nvPr>` +
		`<p:ph type="title" idx="0"/></p:nvPr></p:nvSpPr><p:spPr/>` +
		`<p:txBody><a:bodyPr/><a:p><a:r><a:t>Deck Title</a:t></a:r></a:p></p:txBody>` +
		`</p:sp>` +
		`</p:spTree></p:cSld></p:sld>`)

	placeholders := parsePlaceholdersFromSlideXML(slideXML)

	if len(placeholders) != 1 {
		t.Fatalf("expected 1 placeholder, got %d", len(placeholders))
	}
	if placeholders[0].Text != "Deck Title" {
		t.Fatalf("expected placeholder text %q, got %q", "Deck Title", placeholders[0].Text)
	}
	if placeholders[0].Type != "title" {
		t.Fatalf("expected placeholder type %q, got %q", "title", placeholders[0].Type)
	}
}
