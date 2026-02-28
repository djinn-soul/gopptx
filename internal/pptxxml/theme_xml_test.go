package pptxxml

import (
	"strings"
	"testing"
)

func TestThemeXML(t *testing.T) {
	spec := &ThemeSpec{
		Name: "Test Theme",
		Colors: ColorSchemeSpec{
			Dk1: "000000",
			Lt1: "FFFFFF",
		},
	}
	xml := Theme(spec)
	if !strings.Contains(xml, `name="Test Theme marketing"`) {
		t.Error("Theme name missing or wrong")
	}
	if !strings.Contains(xml, `lastClr="000000"`) {
		t.Error("Dk1 color missing")
	}
}

func TestSlideLayoutXML(t *testing.T) {
	t.Run("TitleAndContent", func(t *testing.T) {
		xml := SlideLayoutTitleAndContent()
		if !strings.Contains(xml, "Title and Content") {
			t.Error("Title missing")
		}
	})
	
	t.Run("Blank", func(t *testing.T) {
		xml := SlideLayoutBlank()
		if !strings.Contains(xml, "Blank") {
			t.Error("Name missing")
		}
	})
	
	t.Run("Relationships", func(t *testing.T) {
		xml := SlideLayoutRelationships(1)
		if !strings.Contains(xml, "slideMaster1.xml") {
			t.Error("Master relationship missing")
		}
	})
}

func TestSlideMasterRelationships(t *testing.T) {
	xml := SlideMasterRelationships([]string{"media/image1.png"}, 1, 1)
	if !strings.Contains(xml, "theme1.xml") {
		t.Error("Theme relationship missing")
	}
	if !strings.Contains(xml, "image1.png") {
		t.Error("Image relationship missing")
	}
}
