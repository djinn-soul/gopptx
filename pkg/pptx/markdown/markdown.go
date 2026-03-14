package markdown

import (
	"errors"
	"regexp"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

var numberedListPattern = regexp.MustCompile(`^\d+\.\s+(.+)$`)

const defaultInlineRunsCapacity = 4

const (
	numberedListMatchCount = 2
	boldMarkerLength       = 2
)

type parsedMarkdownBullet struct {
	text  string
	style elements.ParagraphStyle
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
	if strings.TrimSpace(markdown) == "" {
		return nil, errors.New("markdown content cannot be empty")
	}
	return parseMarkdownWithAST(markdown)
}

func parseBulletLine(line string) (parsedMarkdownBullet, bool) {
	for _, marker := range []string{"- ", "* ", "+ "} {
		if after, ok := strings.CutPrefix(line, marker); ok {
			return parsedMarkdownBullet{
				text:  strings.TrimSpace(after),
				style: elements.DefaultParagraphStyle(),
			}, true
		}
	}

	matches := numberedListPattern.FindStringSubmatch(line)
	if len(matches) == numberedListMatchCount {
		return parsedMarkdownBullet{
			text:  strings.TrimSpace(matches[1]),
			style: elements.DefaultParagraphStyle().WithNumbered(),
		}, true
	}
	return parsedMarkdownBullet{}, false
}

func appendMarkdownBullet(slide *elements.SlideContent, text string, style elements.ParagraphStyle) {
	runs, rich := parseInlineTextRuns(text)
	if rich {
		*slide = slide.AddBulletRunsWithStyle(runs, style)
		return
	}
	*slide = slide.AddBulletWithStyle(text, style)
}

func parseInlineTextRuns(text string) ([]elements.Run, bool) {
	input := strings.TrimSpace(text)
	if input == "" {
		return nil, false
	}

	runs := make([]elements.Run, 0, defaultInlineRunsCapacity)
	hasStyled := false
	for i := 0; i < len(input); {
		if run, nextI, handled := tryHandleRichRun(input, i); handled {
			runs = append(runs, run)
			hasStyled = true
			i = nextI
			continue
		}

		next := nextInlineMarkerOffset(input[i:])
		if next < 0 {
			runs = append(runs, elements.Run{Text: input[i:]})
			break
		}
		if next == 0 {
			runs = append(runs, elements.Run{Text: input[i : i+1]})
			i++
			continue
		}
		runs = append(runs, elements.Run{Text: input[i : i+next]})
		i += next
	}

	return elements.NormalizeRuns(runs), hasStyled
}

func tryHandleRichRun(input string, i int) (elements.Run, int, bool) {
	if input[i] == '`' {
		return handleCodeRun(input, i)
	}
	if strings.HasPrefix(input[i:], "**") {
		return handleBoldRun(input, i)
	}
	if input[i] == '*' {
		return handleItalicRun(input, i)
	}
	return elements.Run{}, i, false
}

func handleCodeRun(input string, i int) (elements.Run, int, bool) {
	closeIdx := strings.Index(input[i+1:], "`")
	if closeIdx >= 0 {
		end := i + 1 + closeIdx
		if end > i+1 {
			return elements.Run{Text: input[i+1 : end], Code: true}, end + 1, true
		}
	}
	return elements.Run{}, i, false
}

func handleBoldRun(input string, i int) (elements.Run, int, bool) {
	closeIdx := strings.Index(input[i+boldMarkerLength:], "**")
	if closeIdx >= 0 {
		end := i + boldMarkerLength + closeIdx
		if end > i+boldMarkerLength {
			return elements.Run{Text: input[i+boldMarkerLength : end], Bold: true}, end + boldMarkerLength, true
		}
	}
	return elements.Run{}, i, false
}

func handleItalicRun(input string, i int) (elements.Run, int, bool) {
	closeIdx := strings.Index(input[i+1:], "*")
	if closeIdx >= 0 {
		end := i + 1 + closeIdx
		if end > i+1 {
			return elements.Run{Text: input[i+1 : end], Italic: true}, end + 1, true
		}
	}
	return elements.Run{}, i, false
}

func nextInlineMarkerOffset(input string) int {
	next := -1
	for _, marker := range []string{"`", "*"} {
		idx := strings.Index(input, marker)
		if idx < 0 {
			continue
		}
		if next < 0 || idx < next {
			next = idx
		}
	}
	return next
}
