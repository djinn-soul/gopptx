package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestSlidesFromMarkdown_Basic(t *testing.T) {
	input := `# Intro
- One
- Two

# Plan
1. Build
2. Validate
`

	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}

	if len(slides) != 2 {
		t.Fatalf("expected 2 slides, got %d", len(slides))
	}
	if slides[0].Title != "Intro" {
		t.Fatalf("expected first title Intro, got %q", slides[0].Title)
	}
	if len(slides[0].Bullets) != 2 || slides[0].Bullets[0] != "One" {
		t.Fatalf("unexpected bullets for first slide: %#v", slides[0].Bullets)
	}
	if slides[1].Bullets[1] != "Validate" {
		t.Fatalf("expected numbered bullet parsing, got %#v", slides[1].Bullets)
	}
}

func TestSlidesFromMarkdown_FailsWhenContentPrecedesHeading(t *testing.T) {
	_, err := SlidesFromMarkdown("- orphan bullet")
	if err == nil {
		t.Fatalf("expected error for content before first heading")
	}
}

func TestCreateWithMarkdownOutput(t *testing.T) {
	input := `# Intro
- Hello`
	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}

	data, err := CreateWithSlides("Deck", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected non-empty PPTX output")
	}
}

func TestSlidesFromMarkdown_InlineRichTextRuns(t *testing.T) {
	input := `# Intro
- **Bold** and *Italic* and ` + "`code`"

	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}
	if len(slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(slides))
	}
	if len(slides[0].Bullets) != 1 {
		t.Fatalf("expected 1 bullet, got %d", len(slides[0].Bullets))
	}
	if slides[0].Bullets[0] != "Bold and Italic and code" {
		t.Fatalf("unexpected plain bullet text: %q", slides[0].Bullets[0])
	}
	if len(slides[0].BulletRuns) != 1 {
		t.Fatalf("expected 1 bullet-run row, got %d", len(slides[0].BulletRuns))
	}

	runs := slides[0].BulletRuns[0]
	if !hasRun(runs, "Bold", true, false, false) {
		t.Fatalf("expected bold run in parsed markdown: %#v", runs)
	}
	if !hasRun(runs, "Italic", false, true, false) {
		t.Fatalf("expected italic run in parsed markdown: %#v", runs)
	}
	if !hasRun(runs, "code", false, false, true) {
		t.Fatalf("expected code run in parsed markdown: %#v", runs)
	}
}

func TestCreateWithMarkdownInlineRunsEmbedsRunProperties(t *testing.T) {
	input := `# Intro
- **Bold** and *Italic* and ` + "`code`"

	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}

	data, err := CreateWithSlides("Deck", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")

	checks := []string{
		`<a:rPr lang="en-US" sz="2800" b="1" i="0" u="none" dirty="0">`,
		`<a:rPr lang="en-US" sz="2800" b="0" i="1" u="none" dirty="0">`,
		`<a:latin typeface="Consolas"/>`,
		`<a:t>Bold</a:t>`,
		`<a:t>Italic</a:t>`,
		`<a:t>code</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func hasRun(runs []TextRun, text string, bold bool, italic bool, code bool) bool {
	for _, run := range runs {
		if run.Text == text && run.Bold == bold && run.Italic == italic && run.Code == code {
			return true
		}
	}
	return false
}
