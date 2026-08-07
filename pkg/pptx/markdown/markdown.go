package markdown

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// defaultInlineRunsCapacity is the run capacity the AST walker preallocates for
// one line of inline content.
const defaultInlineRunsCapacity = 4

// ParseOptions controls markdown-to-slide conversion behavior.
type ParseOptions struct {
	BaseDir string
}

// SlidesFromMarkdown converts a markdown document into slide content.
//
// Supported syntax:
// - "# Title" starts a new slide
// - "-", "*", "+" bullet lines become bullet points
// - numbered lines like "1. item" become bullet points
// - "---" ends the current slide
// - GFM tables are mapped to native table elements
// - fenced code blocks are rendered as no-bullet code paragraphs
// - fenced mermaid blocks are converted to placeholder shapes
// - blockquotes are parsed into slide speaker notes.
func SlidesFromMarkdown(markdown string) ([]elements.SlideContent, error) {
	return SlidesFromMarkdownWithOptions(markdown, ParseOptions{})
}

// SlidesFromMarkdownWithOptions converts markdown with parser options.
func SlidesFromMarkdownWithOptions(markdown string, options ParseOptions) ([]elements.SlideContent, error) {
	if strings.TrimSpace(markdown) == "" {
		return nil, errors.New("markdown content cannot be empty")
	}
	return parseMarkdownWithAST(markdown, options)
}

// SlidesFromMarkdownFile reads markdown from disk and parses it with BaseDir set
// to the markdown file's directory for relative image resolution.
func SlidesFromMarkdownFile(markdownPath string) ([]elements.SlideContent, error) {
	content, err := os.ReadFile(filepath.Clean(markdownPath))
	if err != nil {
		return nil, err
	}
	options := ParseOptions{BaseDir: filepath.Dir(filepath.Clean(markdownPath))}
	return SlidesFromMarkdownWithOptions(string(content), options)
}

// parseInlineTextRuns splits a line on the inline markdown markers. The parser
// itself now lives with the run model, so the plain generator path can reach it
// too.
func parseInlineTextRuns(line string) ([]elements.Run, bool) {
	return elements.ParseInlineRuns(line)
}
