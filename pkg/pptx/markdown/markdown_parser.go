package markdown

import (
	"errors"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type markdownParser struct {
	lines             []string
	index             int
	slides            []elements.SlideContent
	current           *elements.SlideContent
	lastTitle         string
	continuationTitle string
}

func newMarkdownParser(markdown string) *markdownParser {
	return &markdownParser{
		lines:  strings.Split(markdown, "\n"),
		slides: make([]elements.SlideContent, 0, 8),
	}
}

func (p *markdownParser) parse() ([]elements.SlideContent, error) {
	for p.index < len(p.lines) {
		lineNumber := p.index + 1
		trimmed := strings.TrimSpace(p.lines[p.index])

		if trimmed == "" {
			p.index++
			continue
		}
		if isMarkdownFenceStart(trimmed) {
			if err := p.consumeFencedBlock(lineNumber); err != nil {
				return nil, err
			}
			continue
		}
		if p.isTableStart() {
			if err := p.consumeTable(lineNumber); err != nil {
				return nil, err
			}
			continue
		}
		if strings.HasPrefix(trimmed, ">") {
			if err := p.consumeBlockquote(lineNumber); err != nil {
				return nil, err
			}
			continue
		}
		if trimmed == "---" {
			if p.current == nil {
				return nil, fmt.Errorf("line %d: slide separator found before any slide", lineNumber)
			}
			p.flushCurrent()
			p.continuationTitle = p.lastTitle
			p.index++
			continue
		}
		if after, ok := strings.CutPrefix(trimmed, "# "); ok {
			title := strings.TrimSpace(after)
			if title == "" {
				return nil, fmt.Errorf("line %d: slide title cannot be empty", lineNumber)
			}
			p.flushCurrent()
			slide := elements.NewSlide(title)
			p.current = &slide
			p.lastTitle = title
			p.continuationTitle = ""
			p.index++
			continue
		}
		if strings.HasPrefix(trimmed, "## ") || strings.HasPrefix(trimmed, "### ") {
			if err := p.ensureCurrent(lineNumber); err != nil {
				return nil, err
			}
			subheading := strings.TrimSpace(trimMarkdownHeadingPrefix(trimmed))
			if subheading != "" {
				runs := []elements.TextRun{elements.NewTextRun(subheading).WithBold(true)}
				*p.current = p.current.AddBulletRunsWithStyle(runs, elements.DefaultTextParagraphStyle())
			}
			p.index++
			continue
		}

		if err := p.ensureCurrent(lineNumber); err != nil {
			return nil, err
		}
		bullet, ok := parseBulletLine(trimmed)
		if !ok {
			bullet = parsedMarkdownBullet{
				text:  trimmed,
				style: elements.DefaultTextParagraphStyle(),
			}
		}
		appendMarkdownBullet(p.current, bullet.text, bullet.style)
		p.index++
	}

	p.flushCurrent()
	if len(p.slides) == 0 {
		return nil, errors.New("markdown did not produce any slides")
	}
	return p.slides, nil
}

func (p *markdownParser) flushCurrent() {
	if p.current == nil {
		return
	}
	p.slides = append(p.slides, *p.current)
	p.current = nil
}

func (p *markdownParser) ensureCurrent(lineNumber int) error {
	if p.current != nil {
		return nil
	}
	if strings.TrimSpace(p.continuationTitle) != "" {
		title := p.continuationTitle + " (continued)"
		slide := elements.NewSlide(title)
		p.current = &slide
		p.lastTitle = title
		p.continuationTitle = ""
		return nil
	}
	return fmt.Errorf("line %d: content found before first slide title", lineNumber)
}

func trimMarkdownHeadingPrefix(line string) string {
	return strings.TrimLeft(line, "# ")
}
