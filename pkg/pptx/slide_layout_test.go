package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesUsesTitleOnlyLayout(t *testing.T) {
	slide := NewSlide("Title Only").WithTitleOnlyLayout()

	data, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	if !strings.Contains(slideXML, `name="Title"`) {
		t.Fatalf("expected title shape in title-only layout")
	}
	if strings.Contains(slideXML, `name="Content"`) {
		t.Fatalf("did not expect content shape in title-only layout")
	}

	relsXML := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(relsXML, `Target="../slideLayouts/slideLayout2.xml"`) {
		t.Fatalf("expected title-only slide layout target")
	}
}

func TestCreateWithSlidesUsesBlankLayout(t *testing.T) {
	slide := NewSlide("").WithBlankLayout()

	data, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	if strings.Contains(slideXML, `name="Title"`) {
		t.Fatalf("did not expect title shape in blank layout")
	}
	if strings.Contains(slideXML, `name="Content"`) {
		t.Fatalf("did not expect content shape in blank layout")
	}

	relsXML := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(relsXML, `Target="../slideLayouts/slideLayout3.xml"`) {
		t.Fatalf("expected blank slide layout target")
	}
}

func TestCreateWithSlidesEmbedsAllLayoutFiles(t *testing.T) {
	slides := []SlideContent{
		NewSlide("A"),
		NewSlide("B").WithTitleOnlyLayout(),
		NewSlide("").WithBlankLayout(),
	}

	data, err := CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	required := []string{
		"ppt/slideLayouts/slideLayout1.xml",
		"ppt/slideLayouts/slideLayout2.xml",
		"ppt/slideLayouts/slideLayout3.xml",
		"ppt/slideLayouts/_rels/slideLayout1.xml.rels",
		"ppt/slideLayouts/_rels/slideLayout2.xml.rels",
		"ppt/slideLayouts/_rels/slideLayout3.xml.rels",
	}
	for _, name := range required {
		if !zipHasFile(zr, name) {
			t.Fatalf("missing %s", name)
		}
	}
}

func TestCreateWithSlidesRejectsTitleOnlyWithBullets(t *testing.T) {
	_, err := CreateWithSlides("Demo", []SlideContent{
		NewSlide("No Bullets").WithTitleOnlyLayout().AddBullet("x"),
	})
	if err == nil {
		t.Fatalf("expected title-only bullets validation error")
	}
	if !strings.Contains(err.Error(), "title_only layout does not support bullets") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsBlankLayoutWithTitle(t *testing.T) {
	_, err := CreateWithSlides("Demo", []SlideContent{
		NewSlide("Should be Empty").WithBlankLayout(),
	})
	if err == nil {
		t.Fatalf("expected blank-layout title validation error")
	}
	if !strings.Contains(err.Error(), "blank layout requires empty title") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsInvalidSlideLayout(t *testing.T) {
	_, err := CreateWithSlides("Demo", []SlideContent{
		NewSlide("Bad").WithLayout("hero"),
	})
	if err == nil {
		t.Fatalf("expected invalid layout validation error")
	}
	if !strings.Contains(err.Error(), "layout must be one of title_and_content|title_only|blank") {
		t.Fatalf("unexpected error: %v", err)
	}
}
