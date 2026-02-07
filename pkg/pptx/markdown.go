package pptx

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var numberedListPattern = regexp.MustCompile(`^\d+\.\s+(.+)$`)

// SlidesFromMarkdown converts a markdown document into slide content.
//
// Supported syntax:
// - "# Title" starts a new slide
// - "-", "*", "+" bullet lines become bullet points
// - numbered lines like "1. item" become bullet points
// - "---" ends the current slide
func SlidesFromMarkdown(markdown string) ([]SlideContent, error) {
	if strings.TrimSpace(markdown) == "" {
		return nil, fmt.Errorf("markdown content cannot be empty")
	}

	scanner := bufio.NewScanner(strings.NewReader(markdown))
	var slides []SlideContent
	var current *SlideContent
	lineNumber := 0

	flushCurrent := func() {
		if current != nil {
			slides = append(slides, *current)
			current = nil
		}
	}

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}
		if line == "---" {
			if current == nil {
				return nil, fmt.Errorf("line %d: slide separator found before any slide", lineNumber)
			}
			flushCurrent()
			continue
		}
		if strings.HasPrefix(line, "# ") {
			title := strings.TrimSpace(strings.TrimPrefix(line, "# "))
			if title == "" {
				return nil, fmt.Errorf("line %d: slide title cannot be empty", lineNumber)
			}
			flushCurrent()
			slide := NewSlide(title)
			current = &slide
			continue
		}
		if strings.HasPrefix(line, "## ") {
			if current == nil {
				return nil, fmt.Errorf("line %d: content found before first slide title", lineNumber)
			}
			appendMarkdownBullet(current, strings.TrimSpace(strings.TrimPrefix(line, "## ")))
			continue
		}

		bullet := parseBulletLine(line)
		if bullet == "" {
			bullet = line
		}
		if current == nil {
			return nil, fmt.Errorf("line %d: content found before first slide title", lineNumber)
		}
		appendMarkdownBullet(current, bullet)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	flushCurrent()

	if len(slides) == 0 {
		return nil, fmt.Errorf("markdown did not produce any slides")
	}
	return slides, nil
}

func parseBulletLine(line string) string {
	for _, marker := range []string{"- ", "* ", "+ "} {
		if strings.HasPrefix(line, marker) {
			return strings.TrimSpace(strings.TrimPrefix(line, marker))
		}
	}

	matches := numberedListPattern.FindStringSubmatch(line)
	if len(matches) == 2 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func appendMarkdownBullet(slide *SlideContent, text string) {
	runs, rich := parseInlineTextRuns(text)
	if rich {
		*slide = slide.AddBulletRuns(runs)
		return
	}
	*slide = slide.AddBullet(text)
}

func parseInlineTextRuns(text string) ([]TextRun, bool) {
	input := strings.TrimSpace(text)
	if input == "" {
		return nil, false
	}

	runs := make([]TextRun, 0, 4)
	hasStyled := false
	for i := 0; i < len(input); {
		if input[i] == '`' {
			close := strings.Index(input[i+1:], "`")
			if close >= 0 {
				end := i + 1 + close
				if end > i+1 {
					runs = append(runs, TextRun{Text: input[i+1 : end], Code: true})
					hasStyled = true
				}
				i = end + 1
				continue
			}
			runs = append(runs, TextRun{Text: input[i : i+1]})
			i++
			continue
		}

		if strings.HasPrefix(input[i:], "**") {
			close := strings.Index(input[i+2:], "**")
			if close >= 0 {
				end := i + 2 + close
				if end > i+2 {
					runs = append(runs, TextRun{Text: input[i+2 : end], Bold: true})
					hasStyled = true
				}
				i = end + 2
				continue
			}
			runs = append(runs, TextRun{Text: input[i : i+1]})
			i++
			continue
		}

		if input[i] == '*' {
			close := strings.Index(input[i+1:], "*")
			if close >= 0 {
				end := i + 1 + close
				if end > i+1 {
					runs = append(runs, TextRun{Text: input[i+1 : end], Italic: true})
					hasStyled = true
				}
				i = end + 1
				continue
			}
			runs = append(runs, TextRun{Text: input[i : i+1]})
			i++
			continue
		}

		next := nextInlineMarkerOffset(input[i:])
		if next < 0 {
			runs = append(runs, TextRun{Text: input[i:]})
			break
		}
		if next == 0 {
			runs = append(runs, TextRun{Text: input[i : i+1]})
			i++
			continue
		}
		runs = append(runs, TextRun{Text: input[i : i+next]})
		i += next
	}

	return normalizeTextRuns(runs), hasStyled
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
