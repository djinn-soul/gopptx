package markdown

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

var markdownTableSeparatorPattern = regexp.MustCompile(`^:?-{3,}:?$`)

const (
	blockquotePartsCapacity    = 2
	fencedCodeLinesCapacity    = 16
	markdownTableRowsCapacity  = 4
	defaultTableWidthEMU       = 8230200
	markdownTableBaseYEMU      = 1600200
	markdownTableBulletLineEMU = 228600
	markdownTableGapEMU        = 152400
	markdownTableMaxYEMU       = 2743200
)

func isMarkdownFenceStart(trimmedLine string) bool {
	return strings.HasPrefix(trimmedLine, "```")
}

func (p *markdownParser) consumeBlockquote(startLine int) error {
	if err := p.ensureCurrent(startLine); err != nil {
		return err
	}

	parts := make([]string, 0, blockquotePartsCapacity)
	for p.index < len(p.lines) {
		trimmed := strings.TrimSpace(p.lines[p.index])
		if !strings.HasPrefix(trimmed, ">") {
			break
		}
		content := strings.TrimSpace(strings.TrimPrefix(trimmed, ">"))
		if content != "" {
			parts = append(parts, content)
		}
		p.index++
	}
	if len(parts) == 0 {
		return nil
	}

	for _, part := range parts {
		runs, _ := parseInlineTextRuns(part)
		para := elements.NewParagraph()
		para.Runs = runs
		*p.current = p.current.AddNoteParagraph(para)
	}
	return nil
}

func (p *markdownParser) consumeFencedBlock(startLine int) error {
	opening := strings.TrimSpace(p.lines[p.index])
	lang := strings.TrimSpace(strings.TrimPrefix(opening, "```"))

	p.index++
	codeLines := make([]string, 0, fencedCodeLinesCapacity)
	for p.index < len(p.lines) {
		trimmed := strings.TrimSpace(p.lines[p.index])
		if strings.HasPrefix(trimmed, "```") {
			break
		}
		codeLines = append(codeLines, strings.TrimRight(p.lines[p.index], "\r"))
		p.index++
	}
	if p.index >= len(p.lines) {
		return fmt.Errorf("line %d: unterminated fenced code block", startLine)
	}
	code := strings.TrimRight(strings.Join(codeLines, "\n"), "\n")
	if err := p.ensureCurrent(startLine); err != nil {
		return err
	}

	if strings.EqualFold(strings.TrimSpace(lang), "mermaid") {
		if err := addMermaidPlaceholder(p.current, code, startLine); err != nil {
			return err
		}
		p.index++
		return nil
	}

	addCodeBlock(p.current, lang, code)
	p.index++
	return nil
}

func (p *markdownParser) isTableStart() bool {
	if p.index+1 >= len(p.lines) {
		return false
	}
	header := strings.TrimSpace(p.lines[p.index])
	separator := strings.TrimSpace(p.lines[p.index+1])
	if !looksLikeMarkdownTableRow(header) {
		return false
	}
	return isMarkdownTableSeparatorRow(separator)
}

func (p *markdownParser) consumeTable(startLine int) error {
	if err := p.ensureCurrent(startLine); err != nil {
		return err
	}
	if p.current.Table != nil {
		return fmt.Errorf("line %d: multiple tables on one slide are not supported", startLine)
	}

	header, ok := parseMarkdownTableRow(strings.TrimSpace(p.lines[p.index]))
	if !ok {
		return fmt.Errorf("line %d: invalid markdown table header", startLine)
	}
	p.index += 2 // skip header + separator

	rows := make([][]string, 0, markdownTableRowsCapacity)
	rows = append(rows, header)
	for p.index < len(p.lines) {
		trimmed := strings.TrimSpace(p.lines[p.index])
		if !looksLikeMarkdownTableRow(trimmed) {
			break
		}
		row, rowOK := parseMarkdownTableRow(trimmed)
		if !rowOK {
			return fmt.Errorf("line %d: invalid markdown table row", p.index+1)
		}
		rows = append(rows, row)
		p.index++
	}
	if len(rows) == 0 || len(rows[0]) == 0 {
		return fmt.Errorf("line %d: markdown table must have at least one column", startLine)
	}

	columnCount := len(rows[0])
	for i, row := range rows {
		if len(row) != columnCount {
			return fmt.Errorf(
				"line %d: markdown table row has %d columns; expected %d",
				startLine+i,
				len(row),
				columnCount,
			)
		}
	}

	columnWidth := styling.Emu(int64(defaultTableWidthEMU / columnCount))
	columnWidths := make([]styling.Length, columnCount)
	for i := range columnWidths {
		columnWidths[i] = columnWidth
	}
	table := tables.NewTable(columnWidths)

	headerRow := make([]tables.TableCell, 0, columnCount)
	for _, text := range rows[0] {
		headerRow = append(headerRow, tables.NewTableCell(text).WithBold(true).WithBackgroundColor("4472C4"))
	}
	table = table.AddStyledRow(headerRow)
	for _, row := range rows[1:] {
		table = table.AddRow(row)
	}
	*p.current = p.current.WithTable(table)
	return nil
}

func positionMarkdownTable(table tables.Table, slide elements.SlideContent) tables.Table {
	if table.Y.Emu() != markdownTableBaseYEMU {
		return table
	}
	if len(slide.Bullets) == 0 {
		return table
	}
	newY := markdownTableBaseYEMU + int64(len(slide.Bullets))*markdownTableBulletLineEMU + markdownTableGapEMU
	if newY > markdownTableMaxYEMU {
		newY = markdownTableMaxYEMU
	}
	return table.Position(table.X, styling.Emu(newY))
}

func looksLikeMarkdownTableRow(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	if !strings.Contains(trimmed, "|") {
		return false
	}
	return strings.HasPrefix(trimmed, "|") || strings.HasSuffix(trimmed, "|")
}

func parseMarkdownTableRow(line string) ([]string, bool) {
	trimmed := strings.TrimSpace(line)
	if !looksLikeMarkdownTableRow(trimmed) {
		return nil, false
	}
	trimmed = strings.TrimPrefix(trimmed, "|")
	trimmed = strings.TrimSuffix(trimmed, "|")
	parts := strings.Split(trimmed, "|")
	if len(parts) == 0 {
		return nil, false
	}
	row := make([]string, 0, len(parts))
	for _, part := range parts {
		row = append(row, strings.TrimSpace(part))
	}
	return row, true
}

func isMarkdownTableSeparatorRow(line string) bool {
	parts, ok := parseMarkdownTableRow(line)
	if !ok || len(parts) == 0 {
		return false
	}
	for _, part := range parts {
		if !markdownTableSeparatorPattern.MatchString(strings.TrimSpace(part)) {
			return false
		}
	}
	return true
}
