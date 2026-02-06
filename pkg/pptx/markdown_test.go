package pptx

import "testing"

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
