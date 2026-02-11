package markdown

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestSlidesFromMarkdown_GFMTable(t *testing.T) {
	input := `# Data
| Feature | Status |
|---------|--------|
| Tables | Done |
| Mermaid | Done |
`
	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}
	if len(slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(slides))
	}
	if slides[0].Table == nil {
		t.Fatalf("expected table on slide")
	}
	if len(slides[0].Table.Rows) != 3 {
		t.Fatalf("expected header + 2 rows, got %#v", slides[0].Table.Rows)
	}
	if len(slides[0].Table.StyledRows) != 1 || !slides[0].Table.StyledRows[0][0].Bold {
		t.Fatalf("expected styled header row, got %#v", slides[0].Table.StyledRows)
	}
}

func TestSlidesFromMarkdown_CodeFence(t *testing.T) {
	input := `# Code
` + "```rust" + `
fn main() {
    println!("hi");
}
` + "```"

	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}
	if len(slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(slides))
	}
	if len(slides[0].Bullets) < 2 {
		t.Fatalf("expected code block bullets, got %#v", slides[0].Bullets)
	}
	if slides[0].Bullets[0] != "[RUST]" {
		t.Fatalf("expected language header bullet, got %q", slides[0].Bullets[0])
	}
	if len(slides[0].BulletRuns) < 2 || len(slides[0].BulletRuns[1]) == 0 || !slides[0].BulletRuns[1][0].Code {
		t.Fatalf("expected code runs in fenced block, got %#v", slides[0].BulletRuns)
	}
	if slides[0].BulletStyles[0].BulletStyle != elements.BulletStyleNone {
		t.Fatalf("expected no-bullet style for code block, got %#v", slides[0].BulletStyles[0])
	}
}

func TestSlidesFromMarkdown_MermaidBlock(t *testing.T) {
	input := `# Diagram
` + "```mermaid" + `
flowchart LR
    A --> B --> C
` + "```"

	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}
	if len(slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(slides))
	}
	if len(slides[0].Shapes) != 1 {
		t.Fatalf("expected one mermaid placeholder shape, got %d", len(slides[0].Shapes))
	}
	if !strings.Contains(slides[0].Shapes[0].Text, "Flowchart") {
		t.Fatalf("expected flowchart placeholder text, got %q", slides[0].Shapes[0].Text)
	}
}

func TestSlidesFromMarkdown_RejectsUnsupportedMermaid(t *testing.T) {
	input := `# Diagram
` + "```mermaid" + `
unknownDiagram
    A --> B
` + "```"
	_, err := SlidesFromMarkdown(input)
	if err == nil {
		t.Fatalf("expected unsupported mermaid error")
	}
	if !strings.Contains(err.Error(), "unsupported mermaid") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSlidesFromMarkdown_SpeakerNotes(t *testing.T) {
	input := `# Notes
- Agenda

> First note line.
> Second note line.`

	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}
	if len(slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(slides))
	}
	if slides[0].Notes != "First note line.\nSecond note line." {
		t.Fatalf("unexpected notes value: %q", slides[0].Notes)
	}
}

func TestSlidesFromMarkdown_Md2PptDemoFixture(t *testing.T) {
	// Only tests parsing, not generation to avoid cycle
	content := "# Slide 1\n- Bullet"
	slides, err := SlidesFromMarkdown(content)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}
	if len(slides) == 0 {
		t.Fatalf("expected slides")
	}
}
