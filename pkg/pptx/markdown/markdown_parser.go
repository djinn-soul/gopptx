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
		lines: strings.Split(markdown, "\n"),
		//nolint:mnd // Initial capacity for slides
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

		if handled, err := p.handleSpecialBlocks(trimmed, lineNumber); err != nil {
			return nil, err
		} else if handled {
			continue
		}

		if handled, err := p.handleStructuralLine(trimmed, lineNumber); err != nil {
			return nil, err
		} else if handled {
			continue
		}

		if err := p.handleContentLine(trimmed, lineNumber); err != nil {
			return nil, err
		}
	}

	p.flushCurrent()
	if len(p.slides) == 0 {
		return nil, errors.New("markdown did not produce any slides")
	}
	return p.slides, nil
}

func (p *markdownParser) handleSpecialBlocks(trimmed string, lineNumber int) (bool, error) {
	if isMarkdownFenceStart(trimmed) {
		return true, p.consumeFencedBlock(lineNumber)
	}
	if p.isTableStart() {
		return true, p.consumeTable(lineNumber)
	}
	if strings.HasPrefix(trimmed, ">") {
		return true, p.consumeBlockquote(lineNumber)
	}
	return false, nil
}

func (p *markdownParser) handleStructuralLine(trimmed string, lineNumber int) (bool, error) {
	if trimmed == "---" {
		if p.current == nil {
			return false, fmt.Errorf("line %d: slide separator found before any slide", lineNumber)
		}
		p.flushCurrent()
		p.continuationTitle = p.lastTitle
		p.index++
		return true, nil
	}
	if after, ok := strings.CutPrefix(trimmed, "# "); ok {
		title := strings.TrimSpace(after)
		if title == "" {
			return false, fmt.Errorf("line %d: slide title cannot be empty", lineNumber)
		}
		p.flushCurrent()
		slide := elements.NewSlide(title)
		p.current = &slide
		p.lastTitle = title
		p.continuationTitle = ""
		p.index++
		return true, nil
	}
	if strings.HasPrefix(trimmed, "## ") || strings.HasPrefix(trimmed, "### ") {
		if err := p.ensureCurrent(lineNumber); err != nil {
			return false, err
		}
		subheading := strings.TrimSpace(trimMarkdownHeadingPrefix(trimmed))
		if subheading != "" {
			runs := []elements.Run{elements.NewRun(subheading).WithBold(true)}
			*p.current = p.current.AddBulletRunsWithStyle(runs, elements.DefaultParagraphStyle())
		}
		p.index++
		return true, nil
	}
	return false, nil
}

func (p *markdownParser) handleContentLine(trimmed string, lineNumber int) error {
	if err := p.ensureCurrent(lineNumber); err != nil {
		return err
	}
	bullet, ok := parseBulletLine(trimmed)
	if !ok {
		bullet = parsedMarkdownBullet{
			text:  trimmed,
			style: elements.DefaultParagraphStyle(),
		}
	}
	appendMarkdownBullet(p.current, bullet.text, bullet.style)
	p.index++
	return nil
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
