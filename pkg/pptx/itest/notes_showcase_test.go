package itest

import (
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/markdown"
)

func TestSpeakerNotesShowcase(t *testing.T) {
	builder := pptx.NewPresentationBuilder("Speaker Notes Showcase")

	// 1. Simple String Notes (legacy/compatibility)
	s1 := pptx.NewSlide("Simple Notes")
	s1 = s1.WithNotes("This is a simple note string.")
	builder.AddSlide(s1)

	// 2. Rich Text Notes (Bold, Italic) via API
	s2 := pptx.NewSlide("Rich API Notes")
	para1 := pptx.NewTextParagraph()
	para1 = para1.AddRun(pptx.NewTextRun("This is "))
	para1 = para1.AddRun(pptx.NewTextRun("bold").WithBold(true))
	para1 = para1.AddRun(pptx.NewTextRun(" text."))
	s2 = s2.AddNoteParagraph(para1)

	para2 := pptx.NewTextParagraph()
	para2 = para2.AddRun(pptx.NewTextRun("This is "))
	para2 = para2.AddRun(pptx.NewTextRun("italic").WithItalic(true))
	para2 = para2.AddRun(pptx.NewTextRun(" text."))
	s2 = s2.AddNoteParagraph(para2)
	builder.AddSlide(s2)

	// 3. Markdown Parsed Notes
	mdContent := `
# Markdown Slide
> This is a note from markdown.
> It has **bold** and *italic* text.
> - And a bullet list? (Future work)
`
	slides, err := markdown.SlidesFromMarkdown(mdContent)
	if err != nil {
		t.Fatalf("failed to parse markdown: %s", err)
	}
	for _, s := range slides {
		builder.AddSlide(s)
	}

	outPath := filepath.Join(t.TempDir(), "notes_showcase.pptx")
	if err := builder.WriteToFile(outPath); err != nil {
		t.Fatalf("failed to write presentation: %s", err)
	}
}

