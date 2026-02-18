package presentation_test

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/presentation"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestPlaceholderOverrideSmoke(t *testing.T) {
	// 1. Create a presentation with a placeholder override
	slide := elements.NewSlide("smoke test").
		WithPlaceholderText(0, "overridden title").
		WithPlaceholderOverride(shapes.PlaceholderTarget{Type: "body", Index: 1}, shapes.PlaceholderOverrideOptions{
			X:  ptrStyling(styling.Inches(1)),
			Y:  ptrStyling(styling.Inches(2)),
			CX: ptrStyling(styling.Inches(4)),
			CY: ptrStyling(styling.Inches(2)),
			TextStyle: &shapes.PlaceholderTextStyle{
				SizePt: ptrInt(36),
				Bold:   ptrBool(true),
			},
		})

	meta := presentation.Metadata{
		Master: elements.NewMaster(),
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	err := presentation.WritePresentationPackage(zw, meta, []elements.SlideContent{slide}, 1)
	if err != nil {
		t.Fatalf("failed to write package files: %v", err)
	}
	zw.Close()

	// 2. Verify the XML content
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}

	foundSlide := false
	for _, f := range zr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			foundSlide = true
			rc, _ := f.Open()
			content, _ := io.ReadAll(rc)
			rc.Close()

			xml := string(content)
			// Verify geometry override
			if !strings.Contains(xml, "off x=\"914400\" y=\"1828800\"") {
				t.Errorf("expected geometry override in XML, got %s", xml)
			}
			// Verify text style override
			if !strings.Contains(xml, "sz=\"3600\"") {
				t.Errorf("expected text style override (size) in XML, got %s", xml)
			}
			if !strings.Contains(xml, "b=\"1\"") {
				t.Errorf("expected text style override (bold) in XML, got %s", xml)
			}
		}
	}

	if !foundSlide {
		t.Fatal("slide 1 XML not found in package")
	}
}

func TestPlaceholderOverrideCreatePathRejectsNameOnlyTarget(t *testing.T) {
	slide := elements.NewSlide("name-only target").WithPlaceholderOverride(
		shapes.PlaceholderTarget{Name: "Body Placeholder"},
		shapes.PlaceholderOverrideOptions{},
	)
	meta := presentation.Metadata{Master: elements.NewMaster()}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	err := presentation.WritePresentationPackage(zw, meta, []elements.SlideContent{slide}, 1)
	if err == nil || !strings.Contains(err.Error(), "name-only target") {
		t.Fatalf("expected create-path name-only target error, got %v", err)
	}
}

func ptrStyling(l styling.Length) *styling.Length { return &l }
func ptrInt(i int) *int                           { return &i }
func ptrBool(b bool) *bool                        { return &b }
