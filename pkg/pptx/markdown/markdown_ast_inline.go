package markdown

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

type inlineStyleState struct {
	bold          bool
	italic        bool
	strikethrough bool
}

type markdownImage struct {
	Alt  string
	Dest string
}

const markdownTableRowsCapacity = 4

type markdownRunLinkResolver func(destination string) (action.Hyperlink, bool)

func extractInlineRuns(
	node ast.Node,
	source []byte,
	state inlineStyleState,
	resolveLink markdownRunLinkResolver,
) []elements.Run {
	runs := make([]elements.Run, 0, defaultInlineRunsCapacity)
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		runs = append(runs, extractInlineRunsFromNode(child, source, state, resolveLink)...)
	}
	return elements.NormalizeRuns(runs)
}

func extractInlineRunsFromNode(
	node ast.Node,
	source []byte,
	state inlineStyleState,
	resolveLink markdownRunLinkResolver,
) []elements.Run {
	switch typed := node.(type) {
	case *ast.Text:
		text := string(typed.Segment.Value(source))
		if typed.HardLineBreak() || typed.SoftLineBreak() {
			text += " "
		}
		return []elements.Run{styledTextRun(text, state)}
	case *ast.String:
		return []elements.Run{styledTextRun(string(typed.Value), state)}
	case *ast.CodeSpan:
		return []elements.Run{styledCodeRun(string(typed.Text(source)), state)}
	case *ast.Emphasis:
		next := state
		if typed.Level == 2 {
			next.bold = true
		} else {
			next.italic = true
		}
		return extractInlineRuns(typed, source, next, resolveLink)
	case *extast.Strikethrough:
		next := state
		next.strikethrough = true
		return extractInlineRuns(typed, source, next, resolveLink)
	case *ast.Link:
		return applyMarkdownLinkRuns(extractInlineRuns(typed, source, state, resolveLink), string(typed.Destination), resolveLink)
	case *ast.AutoLink:
		run := styledTextRun(string(typed.Label(source)), state)
		return applyMarkdownLinkRuns([]elements.Run{run}, string(typed.URL(source)), resolveLink)
	case *ast.Image:
		return nil
	default:
		return extractInlineRuns(node, source, state, resolveLink)
	}
}

func styledTextRun(text string, state inlineStyleState) elements.Run {
	return elements.NewRun(text).
		WithBold(state.bold).
		WithItalic(state.italic).
		WithStrikethrough(state.strikethrough)
}

func styledCodeRun(text string, state inlineStyleState) elements.Run {
	return elements.NewRun(text).
		WithCode(true).
		WithBold(state.bold).
		WithItalic(state.italic).
		WithStrikethrough(state.strikethrough)
}

func extractPlainText(node ast.Node, source []byte) string {
	runs := extractInlineRuns(node, source, inlineStyleState{}, nil)
	return strings.TrimSpace(elements.RunsToPlainText(runs))
}

func extractTableRows(table *extast.Table, source []byte) [][]string {
	rows := make([][]string, 0, markdownTableRowsCapacity)
	for row := table.FirstChild(); row != nil; row = row.NextSibling() {
		cells := make([]string, 0, 4)
		for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
			cells = append(cells, extractPlainText(cell, source))
		}
		if len(cells) > 0 {
			rows = append(rows, cells)
		}
	}
	return rows
}

func taskCheckboxState(node ast.Node) (bool, bool) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		found, checked := taskCheckboxState(child)
		if found {
			return true, checked
		}
		checkbox, ok := child.(*extast.TaskCheckBox)
		if ok {
			return true, checkbox.IsChecked
		}
	}
	return false, false
}

func collectParagraphImages(node ast.Node, source []byte) ([]markdownImage, bool) {
	images := make([]markdownImage, 0, 2)
	onlyImages := true
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		image, ok := child.(*ast.Image)
		if !ok {
			if textNode, textOK := child.(*ast.Text); textOK {
				if strings.TrimSpace(string(textNode.Segment.Value(source))) == "" {
					continue
				}
			}
			onlyImages = false
			continue
		}
		images = append(images, markdownImage{
			Alt:  strings.TrimSpace(extractPlainText(image, source)),
			Dest: strings.TrimSpace(string(image.Destination)),
		})
	}
	return images, onlyImages
}

func forceBoldRuns(runs []elements.Run) []elements.Run {
	out := make([]elements.Run, 0, len(runs))
	for _, run := range runs {
		out = append(out, run.WithBold(true))
	}
	return out
}
