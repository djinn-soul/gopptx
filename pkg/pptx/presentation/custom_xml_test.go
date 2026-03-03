package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

func TestGenerateCustomXMLItemEscapesStructuredValues(t *testing.T) {
	xmlContent, err := generateCustomXMLItem(common.CustomXMLPart{
		RootElement: "meta",
		Namespace:   `https://example.com/ns?a=1&b=2`,
		Properties: []common.CustomXMLKV{
			{Key: "title", Value: `R&D <draft> "v1"`},
		},
	})
	if err != nil {
		t.Fatalf("generateCustomXMLItem: %v", err)
	}

	if !strings.Contains(xmlContent, `xmlns="https://example.com/ns?a=1&amp;b=2"`) {
		t.Fatalf("namespace was not escaped in structured custom xml: %s", xmlContent)
	}
	if !strings.Contains(xmlContent, `<title>R&amp;D &lt;draft&gt; &quot;v1&quot;</title>`) {
		t.Fatalf("property value was not escaped in structured custom xml: %s", xmlContent)
	}
}

func TestGenerateCustomXMLItemRejectsInvalidNames(t *testing.T) {
	_, err := generateCustomXMLItem(common.CustomXMLPart{
		RootElement: "bad root",
	})
	if err == nil {
		t.Fatal("expected invalid root element name to fail")
	}

	_, err = generateCustomXMLItem(common.CustomXMLPart{
		RootElement: "meta",
		Properties: []common.CustomXMLKV{
			{Key: "bad key", Value: "x"},
		},
	})
	if err == nil {
		t.Fatal("expected invalid property name to fail")
	}
}

func TestGenerateCustomXMLItemPreservesRawInnerXML(t *testing.T) {
	xmlContent, err := generateCustomXMLItem(common.CustomXMLPart{
		RootElement: "meta",
		Content:     `<child>v</child>`,
	})
	if err != nil {
		t.Fatalf("generateCustomXMLItem: %v", err)
	}

	if !strings.Contains(xmlContent, `<child>v</child>`) {
		t.Fatalf("expected raw inner XML to be preserved, got: %s", xmlContent)
	}
	if strings.Contains(xmlContent, `&lt;child&gt;`) {
		t.Fatalf("expected inner XML not to be escaped, got: %s", xmlContent)
	}
}
