package pptxxml

import (
	"strings"
	"testing"
)

func TestSignatureOriginXML(t *testing.T) {
	xml := SignatureOrigin()

	if !strings.Contains(
		xml,
		`<SignatureOrigin xmlns="http://schemas.openxmlformats.org/package/2006/digital-signature"/>`,
	) {
		t.Fatalf("unexpected signature origin xml: %s", xml)
	}
}

func TestContentTypes_Full(t *testing.T) {
	xml := ContentTypes(
		2, // 2 slides
		[]string{"png", "jpg", "mp3"},
		1,        // 1 chart
		1,        // 1 smartart
		[]int{1}, // 1 notes slide
		true,     // has notes master
		1,        // 1 custom xml
		1,        // 1 master
		1,        // notes theme
		true,     // has sections
		[]int{1}, // 1 slide with comments
		true,     // has custom props
		true,     // has signatures
		true,     // has vba
		true,     // has handout master
		true,     // has embedded fonts
	)

	checks := []string{
		`PartName="/ppt/slides/slide1.xml"`,
		`PartName="/ppt/slides/slide2.xml"`,
		`Extension="png"`,
		`Extension="jpg"`,
		`Extension="mp3"`,
		`PartName="/ppt/slideMasters/slideMaster1.xml"`,
		`PartName="/ppt/slideLayouts/slideLayout1.xml"`,
		`PartName="/ppt/notesSlides/notesSlide1.xml"`,
		`PartName="/ppt/notesMasters/notesMaster1.xml"`,
		`PartName="/customXml/item1.xml"`,
		`PartName="/ppt/theme/theme1.xml"`,
		`PartName="/ppt/charts/chart1.xml"`,
		`PartName="/ppt/commentAuthors.xml"`,
		`PartName="/ppt/comments/comment1.xml"`,
		`PartName="/docProps/custom.xml"`,
		`PartName="/ppt/vbaProject.bin"`,
		`PartName="/ppt/handoutMasters/handoutMaster1.xml"`,
		`PartName="/ppt/sectionList.xml"`,
		`PartName="/docProps/core.xml"`,
		`PartName="/docProps/app.xml"`,
	}

	for _, c := range checks {
		if !strings.Contains(xml, c) {
			t.Errorf("missing %s in ContentTypes", c)
		}
	}
}

func TestWriteRID(t *testing.T) {
	var b strings.Builder
	WriteRID(&b, "rId5")
	if b.String() != "rId5" {
		t.Error("WriteRID failed")
	}
}

func TestFastEscapeRID(t *testing.T) {
	if FastEscapeRID("rId1") != "rId1" {
		t.Error("FastEscapeRID failed")
	}
}

func TestCustomPropertiesXML(t *testing.T) {
	xml := CustomProperties(true)
	if !strings.Contains(xml, "_MarkAsFinal") {
		t.Error("CustomProperties XML failed")
	}
	if CustomProperties(false) != "" {
		t.Error("CustomProperties should be empty when false")
	}
}
