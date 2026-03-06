package markdown

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/mermaid"
)

const (
	codeHeaderFontSizePt = 13
	codeBodyFontSizePt   = 14
)

func addCodeBlock(slide *elements.SlideContent, lang string, code string) {
	if slide == nil {
		return
	}
	style := elements.DefaultParagraphStyle().WithNoBullet()
	normalizedLang := strings.ToLower(strings.TrimSpace(lang))
	if normalizedLang == "" {
		normalizedLang = "text"
	}
	header := []elements.Run{
		elements.NewRun("[" + strings.ToUpper(normalizedLang) + "]").
			WithBold(true).
			WithColor("1F4E78").
			WithSizePt(codeHeaderFontSizePt),
	}
	*slide = slide.AddBulletRunsWithStyle(header, style)

	// Use token-level syntax highlighting
	highlighter := newSyntaxHighlighter(normalizedLang)
	tokens := highlighter.Tokenize(code)

	// Group tokens by lines and render with coloring
	addTokenizedCodeLines(slide, tokens, style)
}

func addTokenizedCodeLines(slide *elements.SlideContent, tokens []Token, style elements.ParagraphStyle) {
	var currentLine []Token

	for _, tok := range tokens {
		if tok.Text == "\n" {
			// Render current line with token-level coloring
			renderTokenizedLine(slide, currentLine, style)
			currentLine = nil
		} else {
			currentLine = append(currentLine, tok)
		}
	}

	// Render remaining tokens
	if len(currentLine) > 0 {
		renderTokenizedLine(slide, currentLine, style)
	}
}

func renderTokenizedLine(slide *elements.SlideContent, tokens []Token, style elements.ParagraphStyle) {
	// Check if line is empty (only whitespace)
	isEmpty := true
	for _, tok := range tokens {
		if tok.Type != TokenWhitespace || strings.TrimSpace(tok.Text) != "" {
			isEmpty = false
			break
		}
	}

	if isEmpty || len(tokens) == 0 {
		*slide = slide.AddBulletWithStyle(" ", style)
		return
	}

	// Convert tokens to runs with their token-level colors (from Chroma/Solarized)
	runs := make([]elements.Run, 0, len(tokens))
	for _, tok := range tokens {
		// Use token's color from Chroma lexer, fallback to Solarized mapping
		color := tok.Color
		if color == "" {
			color = GetColor(tok.Type)
		}
		run := elements.NewRun(tok.Text).
			WithCode(true).
			WithColor(color).
			WithSizePt(codeBodyFontSizePt)
		runs = append(runs, run)
	}

	*slide = slide.AddBulletRunsWithStyle(runs, style)
}

func addMermaidPlaceholder(slide *elements.SlideContent, code string, lineNumber int) error {
	diagram, err := mermaid.CreateDiagram(code)
	if err != nil {
		return fmt.Errorf("line %d: %w", lineNumber, err)
	}

	for _, s := range diagram.Shapes {
		*slide = slide.AddShape(s)
	}
	for _, c := range diagram.Connectors {
		*slide = slide.AddConnector(c)
	}
	return nil
}
