package markdown

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

const (
	defaultTableWidthEMU       = 8230200
	markdownTableBaseYEMU      = 1600200
	markdownTableBulletLineEMU = 228600
	markdownTableGapEMU        = 152400
	markdownTableMaxYEMU       = 2743200
)

func (p *markdownASTParser) ensureCurrent(line int) error {
	if p.current != nil {
		return nil
	}
	if strings.TrimSpace(p.continuationTitle) != "" {
		title := p.continuationTitle + " (continued)"
		slide := elements.NewSlide(title)
		p.current = &slide
		p.lastTitle = title
		p.continuationTitle = ""
		p.imagePlacementCount = 0
		return nil
	}
	return fmt.Errorf("line %d: content found before first slide title", line)
}

func (p *markdownASTParser) flushCurrent() {
	if p.current == nil {
		return
	}
	if p.current.Table != nil {
		table := positionMarkdownTable(*p.current.Table, *p.current)
		p.current.Table = &table
	}
	p.slides = append(p.slides, *p.current)
	p.current = nil
	p.imagePlacementCount = 0
}

func (p *markdownASTParser) nodeLine(node ast.Node) int {
	lines := node.Lines()
	if lines != nil && lines.Len() > 0 {
		return lineFromOffset(lines.At(0).Start, p.lineOffsets)
	}
	return 1
}

func (p *markdownASTParser) addMarkdownImage(image markdownImage) error {
	placed, err := p.resolveMarkdownImage(image)
	if err != nil {
		return err
	}
	*p.current = p.current.AddImage(placed)
	p.imagePlacementCount++
	return nil
}

func segmentText(lines *text.Segments, source []byte) string {
	if lines == nil || lines.Len() == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for i := range lines.Len() {
		segment := lines.At(i)
		buffer.Write(segment.Value(source))
	}
	return buffer.String()
}

func buildLineOffsets(source []byte) []int {
	offsets := []int{0}
	for i, b := range source {
		if b == '\n' {
			offsets = append(offsets, i+1)
		}
	}
	return offsets
}

func lineFromOffset(offset int, lineOffsets []int) int {
	if len(lineOffsets) == 0 {
		return 1
	}
	idx := sort.Search(len(lineOffsets), func(i int) bool {
		return lineOffsets[i] > offset
	})
	if idx <= 0 {
		return 1
	}
	return idx
}

func clampBulletLevel(level int) int {
	if level < 0 {
		return 0
	}
	if level > elements.MaxBulletLevel {
		return elements.MaxBulletLevel
	}
	return level
}

func extractInlineRunsFromMarkdownText(textLine string) []elements.Run {
	runs, rich := parseInlineTextRuns(textLine)
	if !rich {
		return []elements.Run{elements.NewRun(textLine)}
	}
	return runs
}

func positionMarkdownTable(table tables.Table, slide elements.SlideContent) tables.Table {
	if table.Y.Emu() != markdownTableBaseYEMU {
		return table
	}
	if len(slide.Bullets) == 0 {
		return table
	}
	newY := markdownTableBaseYEMU + int64(len(slide.Bullets))*markdownTableBulletLineEMU + markdownTableGapEMU
	newY = min(newY, markdownTableMaxYEMU)
	return table.Position(table.X, styling.Emu(newY))
}
