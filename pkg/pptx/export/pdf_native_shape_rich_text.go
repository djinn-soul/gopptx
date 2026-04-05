//nolint:mnd // Rich text layout uses fixed typographic spacing constants.
package export

import (
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

type shapeParagraphLayoutLine struct {
	runs       []pdfStyledRun
	xOffset    float64
	lineHeight float64
	align      string
	availWidth float64
}

func renderPDFShapeParagraphText(pdf *gopdf.GoPdf, s shapes.Shape, x, y, w, h float64) {
	boxX, boxY, boxW, boxH, anchor := shapeTextBox(s, x, y, w, h)
	if boxW <= 2 || boxH <= 2 {
		return
	}
	paragraphs := normalizedShapeParagraphs(s)
	fontSize := fitPDFShapeParagraphText(pdf, paragraphs, boxW, boxH)
	layout, totalHeight := layoutShapeParagraphs(pdf, paragraphs, boxW, fontSize)
	startY := shapeTextStartY(anchor, boxY, boxH, totalHeight)

	pdf.SetTextColor(0, 0, 0)
	yPos := startY
	for _, line := range layout {
		if yPos+line.lineHeight > boxY+boxH+0.5 {
			break
		}
		lineX := boxX + line.xOffset
		if elements.NormalizeTextAlign(line.align) == elements.TextAlignCenter ||
			elements.NormalizeTextAlign(line.align) == elements.TextAlignRight {
			lineText := styledLinePlain(line.runs)
			lineX = alignedTextX(
				pdf,
				lineText,
				boxX+line.xOffset,
				line.availWidth,
				line.align,
				firstStyledFontHint(line.runs),
			)
		}
		renderStyledLine(pdf, line.runs, lineX, yPos)
		yPos += line.lineHeight
	}
	setPDFTextFontWithHint(pdf, defaultFontSize, false, false, "")
}

func fitPDFShapeParagraphText(
	pdf *gopdf.GoPdf,
	paragraphs []text.Paragraph,
	boxW, boxH float64,
) int {
	maxSize := defaultFontSize
	for _, paragraph := range paragraphs {
		for _, run := range paragraph.Runs {
			if run.SizePt > maxSize {
				maxSize = run.SizePt
			}
		}
	}
	low, high := minTextAutoFitSize, maxSize
	bestSize := minTextAutoFitSize
	for low <= high {
		mid := (low + high) / 2
		_, totalHeight := layoutShapeParagraphs(pdf, paragraphs, boxW, mid)
		if totalHeight <= boxH {
			bestSize = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return bestSize
}

func layoutShapeParagraphs(
	pdf *gopdf.GoPdf,
	paragraphs []text.Paragraph,
	boxW float64,
	fittedSize int,
) ([]shapeParagraphLayoutLine, float64) {
	lines := make([]shapeParagraphLayoutLine, 0, 8)
	totalHeight := 0.0
	prevSpaceAfter := 0.0
	for idx, paragraph := range paragraphs {
		style := elements.NormalizeParagraphStyle(paragraph.Style)
		totalHeight += paragraphStartGap(idx, prevSpaceAfter, style)
		levelIndent := float64(style.Level * 14)
		leftIndent := emuToPt(style.LeftIndent.Emu())
		rightIndent := emuToPt(style.RightIndent.Emu())
		hangingIndent := emuToPt(style.HangingIndent.Emu())
		availWidth := boxW - levelIndent - leftIndent - rightIndent
		if availWidth < 40 {
			availWidth = 40
		}
		runs := buildShapeParagraphStyledRuns(paragraph.Runs, fittedSize)
		prefixRuns := buildShapeParagraphPrefixRuns(style, idx, fittedSize, runs)
		wrapped := wrapStyledRuns(pdf, runs, availWidth)
		lineHeight := maxStyledRunsLineHeight(runs) * paragraphLineSpacingFactor(style)
		if lineHeight < 12 {
			lineHeight = 12
		}
		for lineIdx, line := range wrapped {
			xOffset := levelIndent + leftIndent

			if lineIdx == 0 && len(prefixRuns) > 0 {
				prefixX := xOffset - hangingIndent
				if hangingIndent == 0 {
					prefixX = xOffset - 14
				}
				lines = append(lines, shapeParagraphLayoutLine{
					runs:       prefixRuns,
					xOffset:    prefixX,
					lineHeight: 0,
					align:      elements.TextAlignLeft,
					availWidth: availWidth,
				})
			}

			lines = append(lines, shapeParagraphLayoutLine{
				runs:       line,
				xOffset:    xOffset,
				lineHeight: lineHeight,
				align:      style.Align,
				availWidth: availWidth,
			})
			totalHeight += lineHeight
		}
		prevSpaceAfter = paragraphAfterGap(style)
		totalHeight += prevSpaceAfter
	}
	return lines, totalHeight
}

func normalizedShapeParagraphs(s shapes.Shape) []text.Paragraph {
	if len(s.TextParagraphs) > 0 {
		return s.TextParagraphs
	}
	return []text.Paragraph{{Runs: []text.Run{text.NewRun(s.Text)}}}
}

func buildShapeParagraphStyledRuns(runs []text.Run, fittedSize int) []pdfStyledRun {
	if len(runs) == 0 {
		return []pdfStyledRun{{Text: "", Color: [3]uint8{0, 0, 0}, SizePt: fittedSize}}
	}
	out := make([]pdfStyledRun, 0, len(runs))
	for _, run := range runs {
		size := fittedSize
		if run.SizePt > 0 && run.SizePt < size {
			size = run.SizePt
		}
		color := [3]uint8{0, 0, 0}
		if run.Color != "" {
			r, g, b := hexToRGB(run.Color)
			color = [3]uint8{r, g, b}
		}
		fontHint := run.Font
		if strings.TrimSpace(fontHint) == "" {
			fontHint = inferCodeFontHint(run.Text)
		}
		out = append(out, pdfStyledRun{
			Text:     run.Text,
			Bold:     run.Bold,
			Italic:   run.Italic,
			Color:    color,
			FontHint: fontHint,
			SizePt:   size,
		})
	}
	return out
}

func buildShapeParagraphPrefixRuns(
	style text.ParagraphStyle,
	index int,
	fittedSize int,
	runs []pdfStyledRun,
) []pdfStyledRun {
	prefix := bulletPrefix(style, index)
	if prefix == "" {
		return nil
	}
	color := [3]uint8{0, 0, 0}
	if style.BulletColor != "" {
		r, g, b := hexToRGB(style.BulletColor)
		color = [3]uint8{r, g, b}
	} else if len(runs) > 0 {
		color = runs[0].Color
	}
	fontHint := ""
	if len(runs) > 0 {
		fontHint = runs[0].FontHint
	}
	return []pdfStyledRun{{
		Text:     prefix,
		Color:    color,
		FontHint: fontHint,
		SizePt:   fittedSize,
	}}
}

func maxStyledRunsLineHeight(runs []pdfStyledRun) float64 {
	maxHeight := pdfLineHeight(defaultFontSize)
	for _, run := range runs {
		size := run.SizePt
		if size <= 0 {
			size = defaultFontSize
		}
		height := pdfLineHeight(size)
		if height > maxHeight {
			maxHeight = height
		}
	}
	return maxHeight
}

func firstStyledFontHint(runs []pdfStyledRun) string {
	for _, run := range runs {
		if strings.TrimSpace(run.FontHint) != "" {
			return run.FontHint
		}
	}
	return ""
}
