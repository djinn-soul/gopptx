package markdown

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

type markdownASTParser struct {
	source              []byte
	lineOffsets         []int
	options             ParseOptions
	slides              []elements.SlideContent
	current             *elements.SlideContent
	lastTitle           string
	continuationTitle   string
	imagePlacementCount int
}

const (
	initialASTSlideCapacity = 8
)

func parseMarkdownWithAST(markdownContent string, options ParseOptions) ([]elements.SlideContent, error) {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	source := []byte(markdownContent)
	doc := md.Parser().Parse(text.NewReader(source))

	parser := &markdownASTParser{
		source:      source,
		lineOffsets: buildLineOffsets(source),
		options:     options,
		slides:      make([]elements.SlideContent, 0, initialASTSlideCapacity),
	}
	if err := parser.consumeDocument(doc); err != nil {
		return nil, err
	}
	parser.flushCurrent()
	if len(parser.slides) == 0 {
		return nil, errors.New("markdown did not produce any slides")
	}
	return parser.slides, nil
}

func (p *markdownASTParser) consumeDocument(doc ast.Node) error {
	for node := doc.FirstChild(); node != nil; node = node.NextSibling() {
		if err := p.consumeBlock(node, 0); err != nil {
			return err
		}
	}
	return nil
}

func (p *markdownASTParser) consumeBlock(node ast.Node, listDepth int) error {
	switch n := node.(type) {
	case *ast.Heading:
		return p.consumeHeading(n)
	case *ast.ThematicBreak:
		return p.consumeThematicBreak(n)
	case *ast.Paragraph:
		return p.consumeParagraph(n, listDepth, false, false, false)
	case *ast.TextBlock:
		return p.consumeParagraph(n, listDepth, false, false, false)
	case *ast.List:
		return p.consumeList(n, listDepth)
	case *ast.FencedCodeBlock:
		return p.consumeFencedCodeBlock(n)
	case *ast.CodeBlock:
		return p.consumeCodeBlock(n)
	case *extast.Table:
		return p.consumeTable(n)
	case *ast.Blockquote:
		return p.consumeBlockquote(n)
	default:
		return nil
	}
}

func (p *markdownASTParser) consumeHeading(node *ast.Heading) error {
	line := p.nodeLine(node)
	if node.Level == 1 {
		title := strings.TrimSpace(extractPlainText(node, p.source))
		if title == "" {
			return fmt.Errorf("line %d: slide title cannot be empty", line)
		}
		p.flushCurrent()
		slide := elements.NewSlide(title)
		p.current = &slide
		p.lastTitle = title
		p.continuationTitle = ""
		return nil
	}

	if err := p.ensureCurrent(line); err != nil {
		return err
	}
	runs := extractInlineRuns(node, p.source, inlineStyleState{}, p.resolveRunHyperlink)
	if len(runs) == 0 {
		return nil
	}
	runs = forceBoldRuns(runs)
	*p.current = p.current.AddBulletRunsWithStyle(runs, elements.DefaultParagraphStyle())
	return nil
}

func (p *markdownASTParser) consumeThematicBreak(node *ast.ThematicBreak) error {
	line := p.nodeLine(node)
	if p.current == nil {
		return fmt.Errorf("line %d: slide separator found before any slide", line)
	}
	p.flushCurrent()
	p.continuationTitle = p.lastTitle
	return nil
}

func (p *markdownASTParser) consumeParagraph(
	node ast.Node,
	listDepth int,
	isOrdered bool,
	isTaskList bool,
	isTaskChecked bool,
) error {
	line := p.nodeLine(node)
	if err := p.ensureCurrent(line); err != nil {
		return err
	}

	images, onlyImages := collectParagraphImages(node, p.source)
	if len(images) > 0 {
		for _, image := range images {
			if err := p.addMarkdownImage(image); err != nil {
				return fmt.Errorf("line %d: %w", line, err)
			}
		}
		if onlyImages {
			return nil
		}
	}

	runs := extractInlineRuns(node, p.source, inlineStyleState{}, p.resolveRunHyperlink)
	if len(runs) == 0 {
		return nil
	}
	if isTaskList {
		prefix := "[ ] "
		if isTaskChecked {
			prefix = "[x] "
		}
		runs = append([]elements.Run{elements.NewRun(prefix)}, runs...)
	}

	style := elements.DefaultParagraphStyle().WithLevel(clampBulletLevel(listDepth))
	if isOrdered {
		style = style.WithNumbered()
	}
	*p.current = p.current.AddBulletRunsWithStyle(elements.NormalizeRuns(runs), style)
	return nil
}

//nolint:gocognit // List parsing keeps markdown edge-case handling explicit for deterministic slide output.
func (p *markdownASTParser) consumeList(node *ast.List, depth int) error {
	for item := node.FirstChild(); item != nil; item = item.NextSibling() {
		if item.Kind() != ast.KindListItem {
			continue
		}
		isTaskList, isTaskChecked := taskCheckboxState(item)
		for child := item.FirstChild(); child != nil; child = child.NextSibling() {
			switch typed := child.(type) {
			case *ast.Paragraph:
				if err := p.consumeParagraph(typed, depth, node.IsOrdered(), isTaskList, isTaskChecked); err != nil {
					return err
				}
			case *ast.TextBlock:
				if err := p.consumeParagraph(typed, depth, node.IsOrdered(), isTaskList, isTaskChecked); err != nil {
					return err
				}
			case *ast.List:
				if err := p.consumeList(typed, depth+1); err != nil {
					return err
				}
			case *ast.FencedCodeBlock:
				if err := p.consumeFencedCodeBlock(typed); err != nil {
					return err
				}
			case *ast.CodeBlock:
				if err := p.consumeCodeBlock(typed); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *markdownASTParser) consumeFencedCodeBlock(node *ast.FencedCodeBlock) error {
	line := p.nodeLine(node)
	if err := p.ensureCurrent(line); err != nil {
		return err
	}
	lang := strings.TrimSpace(string(node.Language(p.source)))
	code := strings.TrimRight(segmentText(node.Lines(), p.source), "\n")
	if strings.EqualFold(lang, "mermaid") {
		return addMermaidPlaceholder(p.current, code, line)
	}
	addCodeBlock(p.current, lang, code)
	return nil
}

func (p *markdownASTParser) consumeCodeBlock(node *ast.CodeBlock) error {
	line := p.nodeLine(node)
	if err := p.ensureCurrent(line); err != nil {
		return err
	}
	code := strings.TrimRight(segmentText(node.Lines(), p.source), "\n")
	addCodeBlock(p.current, "text", code)
	return nil
}

func (p *markdownASTParser) consumeTable(node *extast.Table) error {
	line := p.nodeLine(node)
	if err := p.ensureCurrent(line); err != nil {
		return err
	}
	if p.current.Table != nil {
		return fmt.Errorf("line %d: multiple tables on one slide are not supported", line)
	}

	rows := extractTableRows(node, p.source)
	if len(rows) == 0 || len(rows[0]) == 0 {
		return fmt.Errorf("line %d: markdown table must have at least one column", line)
	}
	columnCount := len(rows[0])
	for i := range rows {
		if len(rows[i]) != columnCount {
			return fmt.Errorf("line %d: markdown table row has inconsistent columns", line)
		}
	}

	columnWidth := styling.Emu(int64(defaultTableWidthEMU / columnCount))
	columnWidths := make([]styling.Length, columnCount)
	for i := range columnWidths {
		columnWidths[i] = columnWidth
	}
	table := tables.NewTable(columnWidths)
	header := make([]tables.TableCell, 0, columnCount)
	for _, cell := range rows[0] {
		header = append(header, tables.NewTableCell(cell).WithBold(true).WithBackgroundColor("4472C4"))
	}
	table = table.AddStyledRow(header)
	for _, row := range rows[1:] {
		table = table.AddRow(row)
	}
	*p.current = p.current.WithTable(table)
	return nil
}

func (p *markdownASTParser) consumeBlockquote(node *ast.Blockquote) error {
	line := p.nodeLine(node)
	if err := p.ensureCurrent(line); err != nil {
		return err
	}
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		lines := segmentText(child.Lines(), p.source)
		if strings.TrimSpace(lines) == "" {
			continue
		}
		for raw := range strings.SplitSeq(lines, "\n") {
			trimmed := strings.TrimSpace(raw)
			if trimmed == "" {
				continue
			}
			runs := extractInlineRunsFromMarkdownText(trimmed)
			if len(runs) == 0 {
				continue
			}
			note := elements.NewParagraph()
			note.Runs = runs
			*p.current = p.current.AddNoteParagraph(note)
		}
	}
	return nil
}
