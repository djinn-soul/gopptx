package pptxxml

import (
	"strings"
	"testing"
)

func TestCorePropertiesXML(t *testing.T) {
	info := CorePropertiesInfo{
		Title:       "Test Title",
		Subject:     "Test Subject",
		Creator:     "Test Creator",
		Description: "Test Description",
	}
	xml := CoreProperties(info)
	if !strings.Contains(xml, "<dc:title>Test Title</dc:title>") {
		t.Error("Title missing")
	}
	if !strings.Contains(xml, "<dc:subject>Test Subject</dc:subject>") {
		t.Error("Subject missing")
	}
}

func TestAppPropertiesXML(t *testing.T) {
	xml := AppProperties(5, 2, 100, 200)
	if !strings.Contains(xml, "<Slides>5</Slides>") {
		t.Error("Slide count missing")
	}
	if !strings.Contains(xml, "<Notes>2</Notes>") {
		t.Error("Notes count missing")
	}
}
