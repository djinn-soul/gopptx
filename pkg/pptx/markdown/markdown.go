package markdown

import (
	"errors"
	"regexp"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

var numberedListPattern = regexp.MustCompile(`^\d+\.\s+(.+)$`)

type parsedMarkdownBullet struct {
	text  string
	style elements.TextParagraphStyle
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
	parser := newMarkdownParser(markdown)
	return parser.parse()
}

func parseBulletLine(line string) (parsedMarkdownBullet, bool) {
	for _, marker := range []string{"- ", "* ", "+ "} {
		if after, ok := strings.CutPrefix(line, marker); ok {
			return parsedMarkdownBullet{
				text:  strings.TrimSpace(after),
				style: elements.DefaultTextParagraphStyle(),
			}, true
		}
	}

	matches := numberedListPattern.FindStringSubmatch(line)
	if len(matches) == 2 {
		return parsedMarkdownBullet{
			text:  strings.TrimSpace(matches[1]),
			style: elements.DefaultTextParagraphStyle().WithNumbered(),
		}, true
	}
	return parsedMarkdownBullet{}, false
}

func appendMarkdownBullet(slide *elements.SlideContent, text string, style elements.TextParagraphStyle) {
	runs, rich := parseInlineTextRuns(text)
	if rich {
		*slide = slide.AddBulletRunsWithStyle(runs, style)
		return
	}
	*slide = slide.AddBulletWithStyle(text, style)
}

func parseInlineTextRuns(text string) ([]elements.TextRun, bool) {
	input := strings.TrimSpace(text)
	if input == "" {
		return nil, false
	}

	runs := make([]elements.TextRun, 0, 4)
	hasStyled := false
	for i := 0; i < len(input); {
		if input[i] == '`' {
			close := strings.Index(input[i+1:], "`")
			if close >= 0 {
				end := i + 1 + close
				if end > i+1 {
					runs = append(runs, elements.TextRun{Text: input[i+1 : end], Code: true})
					hasStyled = true
				}
				i = end + 1
				continue
			}
			runs = append(runs, elements.TextRun{Text: input[i : i+1]})
			i++
			continue
		}

		if strings.HasPrefix(input[i:], "**") {
			close := strings.Index(input[i+2:], "**")
			if close >= 0 {
				end := i + 2 + close
				if end > i+2 {
					runs = append(runs, elements.TextRun{Text: input[i+2 : end], Bold: true})
					hasStyled = true
				}
				i = end + 2
				continue
			}
			runs = append(runs, elements.TextRun{Text: input[i : i+1]})
			i++
			continue
		}

		if input[i] == '*' {
			close := strings.Index(input[i+1:], "*")
			if close >= 0 {
				end := i + 1 + close
				if end > i+1 {
					runs = append(runs, elements.TextRun{Text: input[i+1 : end], Italic: true})
					hasStyled = true
				}
				i = end + 1
				continue
			}
			runs = append(runs, elements.TextRun{Text: input[i : i+1]})
			i++
			continue
		}

		next := nextInlineMarkerOffset(input[i:])
		if next < 0 {
			runs = append(runs, elements.TextRun{Text: input[i:]})
			break
		}
		if next == 0 {
			runs = append(runs, elements.TextRun{Text: input[i : i+1]})
			i++
			continue
		}
		runs = append(runs, elements.TextRun{Text: input[i : i+next]})
		i += next
	}

	return elements.NormalizeTextRuns(runs), hasStyled
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
