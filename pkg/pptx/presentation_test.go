package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesProducesPackage(t *testing.T) {
	slides := []SlideContent{
		NewSlide("Intro").AddBullet("First").AddBullet("Second"),
		NewSlide("Details"),
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
		"[Content_Types].xml",
		"_rels/.rels",
		"ppt/presentation.xml",
		"ppt/_rels/presentation.xml.rels",
		"ppt/slides/slide1.xml",
		"ppt/slides/slide2.xml",
		"ppt/slides/_rels/slide1.xml.rels",
		"ppt/slides/_rels/slide2.xml.rels",
		"ppt/slideLayouts/slideLayout1.xml",
		"ppt/slideMasters/slideMaster1.xml",
		"ppt/theme/theme1.xml",
		"docProps/core.xml",
		"docProps/app.xml",
	}

	present := map[string]bool{}
	for _, f := range zr.File {
		present[f.Name] = true
	}
	for _, name := range required {
		if !present[name] {
			t.Fatalf("missing %s", name)
		}
	}
}

func TestCreateWithSlidesEscapesText(t *testing.T) {
	slides := []SlideContent{
		NewSlide("Title & More").AddBullet("Use <tag>").AddBullet("5 > 3"),
	}

	data, err := CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	var slideXML string
	for _, f := range zr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			r, err := f.Open()
			if err != nil {
				t.Fatalf("open slide: %v", err)
			}
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(r); err != nil {
				_ = r.Close()
				t.Fatalf("read slide: %v", err)
			}
			_ = r.Close()
			slideXML = buf.String()
			break
		}
	}

	if slideXML == "" {
		t.Fatalf("slide XML not found")
	}

	checks := []string{
		"Title &amp; More",
		"Use &lt;tag&gt;",
		"5 &gt; 3",
	}

	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesValidation(t *testing.T) {
	_, err := CreateWithSlides("Demo", []SlideContent{{}})
	if err == nil {
		t.Fatalf("expected error for empty title")
	}
}
