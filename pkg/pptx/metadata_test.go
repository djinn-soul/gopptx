package pptx_test

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

func TestCreateWithMetadata(t *testing.T) {
	meta := pptx.Metadata{
		Metadata: common.Metadata{
			Title:       "Test Title",
			Subject:     "Test Subject",
			Creator:     "Test Creator",
			Description: "Test Description",
		},
	}
	slides := []pptx.SlideContent{pptx.NewSlide("Slide 1")}

	data, err := pptx.CreateWithMetadata(meta, slides)
	if err != nil {
		t.Fatalf("CreateWithMetadata failed: %v", err)
	}

	inspectCoreXML(t, data, meta)
}

func inspectCoreXML(t *testing.T, data []byte, meta pptx.Metadata) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip.NewReader failed: %v", err)
	}

	found := false
	for _, f := range zr.File {
		if f.Name == "docProps/core.xml" {
			found = true
			validateCoreXMLContent(t, f, meta)
		}
	}

	if !found {
		t.Fatal("docProps/core.xml not found in package")
	}
}

func validateCoreXMLContent(t *testing.T, f *zip.File, meta pptx.Metadata) {
	rc, openErr := f.Open()
	if openErr != nil {
		t.Fatalf("failed to open core.xml: %v", openErr)
	}
	content, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("failed to read core.xml: %v", err)
	}
	if closeErr := rc.Close(); closeErr != nil {
		t.Fatalf("failed to close core.xml: %v", closeErr)
	}

	xml := string(content)
	checks := []struct {
		pattern string
		msg     string
	}{
		{"<dc:title>" + meta.Title + "</dc:title>", "missing title"},
		{"<dc:subject>" + meta.Subject + "</dc:subject>", "missing subject"},
		{"<dc:creator>" + meta.Creator + "</dc:creator>", "missing creator"},
		{"<cp:lastModifiedBy>" + meta.Creator + "</cp:lastModifiedBy>", "missing lastModifiedBy"},
		{"<dc:description>" + meta.Description + "</dc:description>", "missing description"},
	}

	for _, check := range checks {
		if !strings.Contains(xml, check.pattern) {
			t.Errorf("%s in core.xml: %s", check.msg, xml)
		}
	}
}

func TestSlideSize(t *testing.T) {
	meta := pptx.Metadata{
		Metadata: common.Metadata{
			Title:     "16:9 Test",
			SlideSize: pptx.SlideSize16x9(),
		},
	}
	slides := []pptx.SlideContent{pptx.NewSlide("Slide 1")}

	data, err := pptx.CreateWithMetadata(meta, slides)
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
			if closeErr := rc.Close(); closeErr != nil {
				t.Errorf("failed to close presentation.xml: %v", closeErr)
			}
			xml := string(content)
			if !strings.Contains(xml, `cx="12192000" cy="6858000" type="screen16x9"`) {
				t.Errorf("incorrect slide size in presentation.xml: %s", xml)
			}
		}
		if f.Name == "docProps/app.xml" {
			rc, _ := f.Open()
			content, _ := io.ReadAll(rc)
			if closeErr := rc.Close(); closeErr != nil {
				t.Errorf("failed to close app.xml: %v", closeErr)
			}
			xml := string(content)
			if !strings.Contains(xml, "<PresentationFormat>Widescreen</PresentationFormat>") {
				t.Errorf("incorrect presentation format in app.xml: %s", xml)
			}
		}
	}
}
