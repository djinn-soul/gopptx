package presentation

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/markdown"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

// A code block renders as an ordinary text shape: a filled rectangle whose
// paragraphs are the highlighted lines. Doing it here rather than in the XML
// layer keeps the highlighter — and its lexer dependency — on the authoring
// side, where the Markdown importer already lives.

// codeBlockPadding is the inset between the block's box and its text.
const codeBlockPadding = styling.Length(91440) // 0.1 inch

// codeBlockShapes converts each code block on the slide into a shape.
func codeBlockShapes(slide elements.SlideContent) []shapes.Shape {
	if len(slide.CodeBlocks) == 0 {
		return nil
	}
	out := make([]shapes.Shape, 0, len(slide.CodeBlocks))
	for _, block := range slide.CodeBlocks {
		out = append(out, codeBlockShape(block))
	}
	return out
}

func codeBlockShape(block elements.CodeBlock) shapes.Shape {
	shape := shapes.NewShape("rect", block.X, block.Y, block.CX, block.CY).
		WithFill(shapes.NewShapeFill(block.BackgroundOrDefault())).
		WithTextMargins(codeBlockPadding, codeBlockPadding, codeBlockPadding, codeBlockPadding).
		WithVerticalAnchor(shapes.TextAnchorTop).
		WithTextWrap(shapes.TextWrapNone).
		WithAltText(block.AltText)
	shape.TextParagraphs = codeBlockParagraphs(block)
	return shape
}

func codeBlockParagraphs(block elements.CodeBlock) []text.Paragraph {
	lines := markdown.HighlightedLines(
		block.Language,
		block.Code,
		block.FontSizeOrDefault(),
		block.FontOrDefault(),
	)

	paragraphs := make([]text.Paragraph, 0, len(lines)+1)
	if block.ShowLanguageLabel {
		label := text.NewRun(markdown.LanguageLabel(block.Language)).
			WithBold(true).
			WithSizePt(block.FontSizeOrDefault()).
			WithFont(block.FontOrDefault())
		paragraphs = append(paragraphs, codeParagraph([]text.Run{label}))
	}
	for _, runs := range lines {
		paragraphs = append(paragraphs, codeParagraph(runs))
	}
	return paragraphs
}

// codeParagraph builds one unbulleted line. An empty line still needs a
// paragraph, or the listing's blank lines collapse.
func codeParagraph(runs []text.Run) text.Paragraph {
	paragraph := text.NewParagraph()
	paragraph.Style = paragraph.Style.WithNoBullet()
	paragraph.Runs = runs
	return paragraph
}
