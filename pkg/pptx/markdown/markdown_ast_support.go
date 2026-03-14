package markdown

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
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

func (p *markdownASTParser) nextEmbeddedImageFrame() (styling.Length, styling.Length, styling.Length, styling.Length) {
	y := embeddedImageStartYInch + float64(p.imagePlacementCount)*embeddedImageGapInch
	return styling.Inches(embeddedImageXInches),
		styling.Inches(y),
		styling.Inches(embeddedImageWidthInches),
		styling.Inches(embeddedImageHeightInches)
}

func segmentText(lines *text.Segments, source []byte) string {
	if lines == nil || lines.Len() == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for i := 0; i < lines.Len(); i++ {
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
