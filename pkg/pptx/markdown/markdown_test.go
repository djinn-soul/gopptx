package markdown

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
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
	if slides[1].BulletStyles[0].BulletStyle != elements.BulletStyleNumber {
		t.Fatalf("expected numbered style for markdown ordered list, got %#v", slides[1].BulletStyles[0])
	}
}

func TestSlidesFromMarkdown_FailsWhenContentPrecedesHeading(t *testing.T) {
	_, err := SlidesFromMarkdown("- orphan bullet")
	if err == nil {
		t.Fatalf("expected error for content before first heading")
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

func TestSlidesFromMarkdown_StrikethroughAndTaskList(t *testing.T) {
	input := `# Checklist
- [x] ~~Done~~ item
- [ ] Pending item
`
	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}
	if len(slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(slides))
	}
	if len(slides[0].BulletRuns) != 2 {
		t.Fatalf("expected 2 bullet run entries, got %d", len(slides[0].BulletRuns))
	}
	if !hasRun(slides[0].BulletRuns[0], "[x] ", false, false, false) {
		t.Fatalf("expected checked task marker run, got %#v", slides[0].BulletRuns[0])
	}
	if !hasStrikeRun(slides[0].BulletRuns[0], "Done") {
		t.Fatalf("expected strikethrough run in first task item, got %#v", slides[0].BulletRuns[0])
	}
	if !hasRun(slides[0].BulletRuns[1], "[ ] ", false, false, false) {
		t.Fatalf("expected unchecked task marker run, got %#v", slides[0].BulletRuns[1])
	}
}

func TestSlidesFromMarkdown_NestedListKeepsLevels(t *testing.T) {
	input := `# Plan
- parent
  - child
    - grandchild
`
	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}
	if len(slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(slides))
	}
	if len(slides[0].BulletStyles) != 3 {
		t.Fatalf("expected 3 bullet styles, got %d", len(slides[0].BulletStyles))
	}
	if slides[0].BulletStyles[0].Level != 0 {
		t.Fatalf("expected first bullet level 0, got %d", slides[0].BulletStyles[0].Level)
	}
	if slides[0].BulletStyles[1].Level != 1 {
		t.Fatalf("expected second bullet level 1, got %d", slides[0].BulletStyles[1].Level)
	}
	if slides[0].BulletStyles[2].Level != 2 {
		t.Fatalf("expected third bullet level 2, got %d", slides[0].BulletStyles[2].Level)
	}
}

func TestSlidesFromMarkdown_ImageEmbedsSlideImage(t *testing.T) {
	input := `# Images
![Architecture](https://example.com/arch.png)
`
	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}
	if len(slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(slides))
	}
	if len(slides[0].Images) != 1 {
		t.Fatalf("expected 1 embedded image, got %d", len(slides[0].Images))
	}
	if slides[0].Images[0].SourceURL != "https://example.com/arch.png" {
		t.Fatalf("expected source URL image, got %#v", slides[0].Images[0])
	}
	if slides[0].Images[0].AltText != "Architecture" {
		t.Fatalf("expected alt text to propagate, got %q", slides[0].Images[0].AltText)
	}
}

func TestSlidesFromMarkdown_DataURIImage(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString(testutil.TinyPNG())
	input := "# Inline\n![Pixel](data:image/png;base64," + encoded + ")\n"

	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}
	if len(slides) != 1 || len(slides[0].Images) != 1 {
		t.Fatalf("expected one image from data URI, got slides=%d images=%d", len(slides), len(slides[0].Images))
	}
	image := slides[0].Images[0]
	if len(image.Data) == 0 || image.Format != "png" {
		t.Fatalf("expected decoded png image data, got format=%q bytes=%d", image.Format, len(image.Data))
	}
}

func TestSlidesFromMarkdownFile_ResolvesRelativeImagePath(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "photo.png"), testutil.TinyPNG(), 0o600); err != nil {
		t.Fatalf("write image fixture: %v", err)
	}
	markdownPath := filepath.Join(tmpDir, "deck.md")
	content := "# Relative\n![Shot](photo.png)\n"
	if err := os.WriteFile(markdownPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write markdown fixture: %v", err)
	}

	slides, err := SlidesFromMarkdownFile(markdownPath)
	if err != nil {
		t.Fatalf("SlidesFromMarkdownFile returned error: %v", err)
	}
	if len(slides) != 1 || len(slides[0].Images) != 1 {
		t.Fatalf("expected one relative image, got slides=%d images=%d", len(slides), len(slides[0].Images))
	}

	expectedPath := filepath.Clean(filepath.Join(tmpDir, "photo.png"))
	if got := filepath.Clean(slides[0].Images[0].Path); got != expectedPath {
		t.Fatalf("expected resolved image path %q, got %q", expectedPath, got)
	}
}

func hasRun(runs []elements.Run, text string, bold bool, italic bool, code bool) bool {
	for _, run := range runs {
		if run.Text == text && run.Bold == bold && run.Italic == italic && run.Code == code {
			return true
		}
	}
	return false
}

func hasStrikeRun(runs []elements.Run, text string) bool {
	for _, run := range runs {
		if run.Text == text && run.Strikethrough == "sngStrike" {
			return true
		}
	}
	return false
}
