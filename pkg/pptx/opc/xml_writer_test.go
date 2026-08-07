package opc_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/opc"
)

func TestXMLWriterBuildsADocument(t *testing.T) {
	doc, err := opc.NewXMLWriter().
		StartElement("p:custom", "xmlns:p", "urn:example").
		Element("p:name", "Deck & Co").
		EmptyElement("p:flag", "val", "1").
		Raw("<p:passthrough/>").
		EndElement().
		Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	for _, want := range []string{
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`,
		`<p:custom xmlns:p="urn:example">`,
		"<p:name>Deck &amp; Co</p:name>",
		`<p:flag val="1"/>`,
		"<p:passthrough/>",
		"</p:custom>",
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("document is missing %s:\n%s", want, doc)
		}
	}
}

func TestXMLWriterReportsUnbalancedMarkup(t *testing.T) {
	if _, err := opc.NewXMLWriter().StartElement("a").Finish(); err == nil {
		t.Error("an element left open should be an error")
	}
	if _, err := opc.NewXMLFragment().EndElement().Finish(); err == nil {
		t.Error("closing with nothing open should be an error")
	}
	if _, err := opc.NewXMLFragment().StartElement("a", "only-a-name").Finish(); err == nil {
		t.Error("an odd attribute list should be an error")
	}
}

// A fragment carries no declaration, so it can be spliced into an existing part.
func TestXMLFragmentHasNoDeclaration(t *testing.T) {
	doc, err := opc.NewXMLFragment().EmptyElement("a:ext", "uri", "{GUID}").Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}
	if strings.Contains(doc, "<?xml") {
		t.Fatalf("fragment should carry no declaration: %s", doc)
	}
}
