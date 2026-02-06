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
			current.Bullets = append(current.Bullets, strings.TrimSpace(strings.TrimPrefix(line, "## ")))
			continue
		}

		bullet := parseBulletLine(line)
		if bullet == "" {
			bullet = line
		}
		if current == nil {
			return nil, fmt.Errorf("line %d: content found before first slide title", lineNumber)
		}
		current.Bullets = append(current.Bullets, bullet)
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
