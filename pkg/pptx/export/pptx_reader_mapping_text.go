package export

import (
	"reflect"
	"strings"

	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func consumeBodyPlaceholderAsBullets(sc *elements.SlideContent, es editorcommon.Shape) bool {
	if appendParagraphBullets(sc, es.Paragraphs) {
		return true
	}
	bodyText := strings.TrimSpace(es.Text)
	if bodyText == "" {
		return false
	}
	parts := strings.Split(bodyText, "\n")
	style := editorParagraphToExportStyle(es.Paragraph)
	for _, line := range parts {
		line = strings.TrimSpace(strings.TrimPrefix(line, "•"))
		if line == "" {
			continue
		}
		sc.Bullets = append(sc.Bullets, line)
		sc.BulletRuns = append(sc.BulletRuns, nil)
		sc.BulletStyles = append(sc.BulletStyles, style)
	}
	return len(sc.Bullets) > 0
}

func appendParagraphBullets(sc *elements.SlideContent, paragraphs []editorcommon.ShapeTextParagraph) bool {
	if len(paragraphs) == 0 {
		return false
	}
	startCount := len(sc.Bullets)
	for _, paragraph := range paragraphs {
		runs := editorRunsToExportRuns(paragraph.Runs)
		textValue := strings.TrimSpace(elements.RunsToPlainText(runs))
		if textValue == "" {
			continue
		}
		style := editorParagraphToExportStyle(paragraph.Paragraph)
		sc.Bullets = append(sc.Bullets, textValue)
		if len(runs) > 0 {
			sc.BulletRuns = append(sc.BulletRuns, runs)
		} else {
			sc.BulletRuns = append(sc.BulletRuns, nil)
		}
		sc.BulletStyles = append(sc.BulletStyles, style)
	}
	return len(sc.Bullets) > startCount
}

func editorParagraphsToExportParagraphs(es editorcommon.Shape) []text.Paragraph {
	if len(es.Paragraphs) > 0 {
		out := make([]text.Paragraph, 0, len(es.Paragraphs))
		for _, paragraph := range es.Paragraphs {
			runs := editorRunsToExportRuns(paragraph.Runs)
			style := editorParagraphToExportStyle(paragraph.Paragraph)
			if len(runs) == 0 && reflect.DeepEqual(style, elements.DefaultParagraphStyle()) {
				continue
			}
			out = append(out, text.Paragraph{
				Runs:  runs,
				Style: style,
			})
		}
		if len(out) > 0 {
			return out
		}
	}
	if len(es.Runs) > 0 || es.Paragraph != nil {
		return []text.Paragraph{{
			Runs:  editorRunsToExportRuns(es.Runs),
			Style: editorParagraphToExportStyle(es.Paragraph),
		}}
	}
	return nil
}

func editorParagraphToExportStyle(paragraph *editorcommon.Paragraph) elements.ParagraphStyle {
	style := elements.DefaultParagraphStyle()
	if paragraph == nil {
		return style
	}
	if paragraph.Alignment != nil {
		style.Align = *paragraph.Alignment
	}
	if paragraph.SpaceBeforePts != nil {
		style.SpaceBeforePt = *paragraph.SpaceBeforePts / textSpacingPtsScale
	}
	if paragraph.SpaceAfterPts != nil {
		style.SpaceAfterPt = *paragraph.SpaceAfterPts / textSpacingPtsScale
	}
	if paragraph.LineSpacingPct != nil {
		style.LineSpacingPct = *paragraph.LineSpacingPct / textSpacingPctScale
	}
	if paragraph.BulletStyle != nil {
		style.BulletStyle = *paragraph.BulletStyle
	}
	if paragraph.BulletChar != nil {
		style.BulletChar = *paragraph.BulletChar
	}
	if paragraph.BulletColor != nil {
		style.BulletColor = *paragraph.BulletColor
	}
	if paragraph.BulletSizePct != nil {
		style.BulletSize = *paragraph.BulletSizePct
	}
	if paragraph.Level != nil {
		style.Level = *paragraph.Level
	}
	if paragraph.Indent != nil {
		style.LeftIndent = styling.Emu(int64(*paragraph.Indent))
	}
	if paragraph.Hanging != nil {
		style.HangingIndent = styling.Emu(int64(*paragraph.Hanging))
	}
	return style
}

//nolint:gocognit
func editorRunsToExportRuns(runs []editorcommon.TextRun) []elements.Run {
	if len(runs) == 0 {
		return nil
	}
	out := make([]elements.Run, 0, len(runs))
	for _, run := range runs {
		exportRun := elements.NewRun(run.Text)
		if run.Bold != nil {
			exportRun.Bold = *run.Bold
		}
		if run.Italic != nil {
			exportRun.Italic = *run.Italic
		}
		if run.Underline != nil {
			exportRun.Underline = *run.Underline
		}
		if run.Strikethrough != nil {
			exportRun.Strikethrough = *run.Strikethrough
		}
		if run.Subscript != nil {
			exportRun.Subscript = *run.Subscript
		}
		if run.Superscript != nil {
			exportRun.Superscript = *run.Superscript
		}
		if run.Color != nil {
			exportRun.Color = *run.Color
		}
		if run.Highlight != nil {
			exportRun.Highlight = *run.Highlight
		}
		if run.Font != nil {
			exportRun.Font = *run.Font
		}
		if run.SizePt != nil {
			exportRun.SizePt = *run.SizePt
		}
		if run.Code != nil {
			exportRun.Code = *run.Code
		}
		if run.AllCaps != nil {
			exportRun.AllCaps = *run.AllCaps
		}
		if run.SmallCaps != nil {
			exportRun.SmallCaps = *run.SmallCaps
		}
		if run.OutlineColor != nil {
			exportRun.OutlineColor = *run.OutlineColor
		}
		if run.OutlineWidthPt != nil {
			exportRun.OutlineWidthPt = *run.OutlineWidthPt
		}
		exportRun.Hyperlink = editorHyperlinkToExportHyperlink(run.Hyperlink)
		exportRun.HoverAction = editorHyperlinkToExportHyperlink(run.HoverAction)
		out = append(out, exportRun)
	}
	return elements.NormalizeRuns(out)
}
