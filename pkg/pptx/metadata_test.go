package pptx

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestCreateWithMetadata(t *testing.T) {
	meta := PresentationMetadata{
		Title:       "Test Title",
		Subject:     "Test Subject",
		Creator:     "Test Creator",
		Description: "Test Description",
	}
	slides := []SlideContent{NewSlide("Slide 1")}

	data, err := CreateWithMetadata(meta, slides)
	if err != nil {
		t.Fatalf("CreateWithMetadata failed: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip.NewReader failed: %v", err)
	}

	found := false
	for _, f := range zr.File {
		if f.Name == "docProps/core.xml" {
			found = true
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("failed to open core.xml: %v", err)
			}
			content, err := io.ReadAll(rc)
			if err != nil {
				t.Fatalf("failed to read core.xml: %v", err)
			}
			if err := rc.Close(); err != nil {
				t.Fatalf("failed to close core.xml: %v", err)
			}

			xml := string(content)
			if !strings.Contains(xml, "<dc:title>Test Title</dc:title>") {
				t.Errorf("missing title in core.xml: %s", xml)
			}
			if !strings.Contains(xml, "<dc:subject>Test Subject</dc:subject>") {
				t.Errorf("missing subject in core.xml: %s", xml)
			}
			if !strings.Contains(xml, "<dc:creator>Test Creator</dc:creator>") {
				t.Errorf("missing creator in core.xml: %s", xml)
			}
			if !strings.Contains(xml, "<cp:lastModifiedBy>Test Creator</cp:lastModifiedBy>") {
				t.Errorf("missing lastModifiedBy in core.xml: %s", xml)
			}
			if !strings.Contains(xml, "<dc:description>Test Description</dc:description>") {
				t.Errorf("missing description in core.xml: %s", xml)
			}
		}
	}

	if !found {
		t.Fatal("docProps/core.xml not found in package")
	}
}

func TestSlideSize(t *testing.T) {
	meta := PresentationMetadata{
		Title:     "16:9 Test",
		SlideSize: SlideSize16x9,
	}
	slides := []SlideContent{NewSlide("Slide 1")}

	data, err := CreateWithMetadata(meta, slides)
	if err != nil {
		t.Fatalf("CreateWithMetadata failed: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip.NewReader failed: %v", err)
	}

	for _, f := range zr.File {
		if f.Name == "ppt/presentation.xml" {
			rc, _ := f.Open()
			content, _ := io.ReadAll(rc)
			if err := rc.Close(); err != nil {
				t.Errorf("failed to close presentation.xml: %v", err)
			}
			xml := string(content)
			if !strings.Contains(xml, `cx="12192000" cy="6858000" type="screen16x9"`) {
				t.Errorf("incorrect slide size in presentation.xml: %s", xml)
			}
		}
		if f.Name == "docProps/app.xml" {
			rc, _ := f.Open()
			content, _ := io.ReadAll(rc)
			if err := rc.Close(); err != nil {
				t.Errorf("failed to close app.xml: %v", err)
			}
			xml := string(content)
			if !strings.Contains(xml, "<PresentationFormat>Widescreen</PresentationFormat>") {
				t.Errorf("incorrect presentation format in app.xml: %s", xml)
			}
		}
	}
}
