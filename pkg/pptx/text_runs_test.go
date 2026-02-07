package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesEmbedsRichTextRunFormatting(t *testing.T) {
	slide := NewSlide("Rich").
		AddBulletRuns([]TextRun{
			NewTextRun("Heading ").
				WithBold(true).
				WithItalic(true).
				WithUnderline(true).
				WithColor("1A2B3C").
				WithFont("Calibri").
				WithSizePt(18),
			NewTextRun("Code").WithCode(true),
		})

	data, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")

	checks := []string{
		`<a:rPr lang="en-US" sz="1800" b="1" i="1" u="sng" dirty="0"><a:solidFill><a:srgbClr val="1A2B3C"/></a:solidFill><a:latin typeface="Calibri"/></a:rPr>`,
		`<a:rPr lang="en-US" sz="2800" b="0" i="0" u="none" dirty="0"><a:latin typeface="Consolas"/></a:rPr>`,
		`<a:t>Heading </a:t>`,
		`<a:t>Code</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidTextRunColor(t *testing.T) {
	slide := NewSlide("Rich").
		AddBulletRuns([]TextRun{
			NewTextRun("Bad").WithColor("GGHHII"),
		})

	_, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err == nil {
		t.Fatalf("expected color validation error")
	}
	if !strings.Contains(err.Error(), "color must be 6-digit RGB hex") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsInvalidTextRunSize(t *testing.T) {
	slide := NewSlide("Rich").
		AddBulletRuns([]TextRun{
			NewTextRun("Bad").WithSizePt(-1),
		})

	_, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err == nil {
		t.Fatalf("expected size validation error")
	}
	if !strings.Contains(err.Error(), "size must be >= 0") {
		t.Fatalf("unexpected error: %v", err)
	}
}
