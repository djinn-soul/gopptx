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
		`<a:rPr lang="en-US" sz="1800" b="0" i="0" u="none" dirty="0"><a:latin typeface="Consolas"/></a:rPr>`,
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

func TestCreateWithSlidesEmbedsTextEnhancementRunFormatting(t *testing.T) {
	slide := NewSlide("Enhancements").
		AddBulletRuns([]TextRun{
			NewTextRun("Strike").WithStrikethrough(true).WithHighlight("FFFF00"),
			NewTextRun("H2O").WithSubscript(true),
			NewTextRun("x2").WithSuperscript(true),
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
		`<a:rPr lang="en-US" sz="1800" b="0" i="0" u="none" strike="sngStrike" dirty="0"><a:highlight><a:srgbClr val="FFFF00"/></a:highlight></a:rPr>`,
		`baseline="-25000"`,
		`baseline="30000"`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidTextRunHighlight(t *testing.T) {
	slide := NewSlide("Rich").
		AddBulletRuns([]TextRun{
			NewTextRun("Bad").WithHighlight("not-hex"),
		})

	_, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err == nil {
		t.Fatalf("expected highlight validation error")
	}
	if !strings.Contains(err.Error(), "highlight must be 6-digit RGB hex") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsConflictingBaselineModes(t *testing.T) {
	slide := NewSlide("Rich").
		AddBulletRuns([]TextRun{
			{
				Text:        "Bad",
				Subscript:   true,
				Superscript: true,
			},
		})

	_, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err == nil {
		t.Fatalf("expected conflicting baseline validation error")
	}
	if !strings.Contains(err.Error(), "cannot be both subscript and superscript") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTextSizePresetsMatchPptRsPrelude(t *testing.T) {
	if TextSizeTitle != 44 || TextSizeSubtitle != 32 || TextSizeHeading != 28 {
		t.Fatalf("unexpected title/subtitle/heading presets")
	}
	if TextSizeBody != 18 || TextSizeSmall != 14 || TextSizeCaption != 12 {
		t.Fatalf("unexpected body/small/caption presets")
	}
	if TextSizeCode != 14 || TextSizeLarge != 36 || TextSizeXLarge != 48 {
		t.Fatalf("unexpected code/large/xlarge presets")
	}
}
