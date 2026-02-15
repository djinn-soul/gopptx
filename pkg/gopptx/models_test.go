package gopptx

import (
	"encoding/xml"
	"testing"
)

func TestPresentationXML(t *testing.T) {
	pres := &Presentation{}
	output, err := xml.Marshal(pres)
	if err != nil {
		t.Fatalf("Failed to marshal Presentation: %v", err)
	}

	expected := `<presentation xmlns="http://schemas.openxmlformats.org/presentationml/2006/main"></presentation>`
	if string(output) != expected {
		t.Errorf("Expected %s, got %s", expected, string(output))
	}
}

func TestSlideXML(t *testing.T) {
	slide := &Slide{}
	output, err := xml.Marshal(slide)
	if err != nil {
		t.Fatalf("Failed to marshal Slide: %v", err)
	}

	expected := `<sld xmlns="http://schemas.openxmlformats.org/presentationml/2006/main"></sld>`
	if string(output) != expected {
		t.Errorf("Expected %s, got %s", expected, string(output))
	}
}
