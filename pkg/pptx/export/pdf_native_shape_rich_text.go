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
	tabStops   []float64
}

func renderPDFShapeParagraphText(pdf *gopdf.GoPdf, s shapes.Shape, x, y, w, h float64) {
	if s.TextFrame != nil && isVerticalShapeText(s.TextFrame.Orientation) {
		renderPDFShapeParagraphTextVertical(pdf, s, x, y, w, h)
		return
	}
	boxX, boxY, boxW, boxH, anchor := shapeTextBox(s, x, y, w, h)
	if boxW <= 2 || boxH <= 2 {
		return
	}
	boxX, boxY, boxW, boxH, restoreOrientation := beginShapeTextOrientation(
		pdf, s.TextFrame, boxX, boxY, boxW, boxH, x, y, w, h,
	)
	defer restoreOrientation()
	paragraphs := normalizedShapeParagraphs(s)
	fontSize := shapeParagraphNaturalSize(paragraphs)
	if shapeTextShrinksToFit(s) {
		fontSize = fitPDFShapeParagraphText(pdf, paragraphs, boxW, boxH)
	}
	layout, totalHeight := layoutShapeParagraphs(pdf, paragraphs, boxW, fontSize)
	startY := shapeTextStartY(anchor, boxY, boxH, totalHeight)

	pdf.SetTextColor(0, 0, 0)
	// A shape that does not shrink its text lets that text spill past the box,
	// exactly as PowerPoint does; clipping it there dropped the tail of any
	// caption longer than its shape.
	clipToBox := shapeTextShrinksToFit(s)
	yPos := startY
	for _, line := range layout {
		if clipToBox && yPos+line.lineHeight > boxY+boxH+0.5 {
			break
		}
		lineX := boxX + line.xOffset
		if elements.NormalizeTextAlign(line.align) == elements.TextAlignCenter ||
			elements.NormalizeTextAlign(line.align) == elements.TextAlignRight {
			lineText := styledLinePlain(line.runs)
			lineX = alignedTextX(pdf, lineText, boxX+line.xOffset, line.availWidth, line.align)
		}
		renderStyledLine(pdf, line.runs, lineX, yPos, pdfTextRenderOptions{
			LineHeight: line.lineHeight,
			TabStops:   line.tabStops,
		})
		yPos += line.lineHeight
	}
	setPDFTextFontWithHint(pdf, defaultFontSize, false, false, "")
}

// shapeTextShrinksToFit reports whether the shape asked PowerPoint to shrink its
// text on overflow.
//
// OOXML's default is noAutofit: PowerPoint renders text at its stated size and
// lets it spill out of the shape. The renderer used to shrink unconditionally,
// so a caption in a short box came out a third of its real size while PowerPoint
// drew it full size across two lines.
func shapeTextShrinksToFit(s shapes.Shape) bool {
	return s.TextFrame != nil && s.TextFrame.AutoFit == shapes.TextAutoFitNormal
}

// shapeParagraphNaturalSize is the size the runs ask for, with no fitting.
func shapeParagraphNaturalSize(paragraphs []text.Paragraph) int {
	maxSize := defaultFontSize
	for _, paragraph := range paragraphs {
		for _, run := range paragraph.Runs {
			if run.SizePt > maxSize {
				maxSize = run.SizePt
			}
		}
	}
	return maxSize
}

func fitPDFShapeParagraphText(
	pdf *gopdf.GoPdf,
	paragraphs []text.Paragraph,
	boxW, boxH float64,
) int {
	maxSize := shapeParagraphNaturalSize(paragraphs)
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
		style := shapeParagraphStyle(paragraph)
		totalHeight += paragraphStartGap(idx, prevSpaceAfter, style, shapeTextSpacing())
		levelIndent := float64(style.Level * 14)
		leftIndent := emuToPt(style.LeftIndent.Emu())
		rightIndent := emuToPt(style.RightIndent.Emu())
		hangingIndent := emuToPt(style.HangingIndent.Emu())
		tabStops := paragraphTabStopsPt(style)
		availWidth := boxW - levelIndent - leftIndent - rightIndent
		if availWidth < 40 {
			availWidth = 40
		}
		runs := buildShapeParagraphStyledRuns(paragraph.Runs, fittedSize)
		prefixRuns := buildShapeParagraphPrefixRuns(style, idx, fittedSize, runs)
		wrapped := wrapStyledRuns(pdf, runs, availWidth, tabStops)
		lineHeight := paragraphRenderedLineHeight(style, maxStyledRunsLineHeight(runs), shapeTextSpacing())
		if lineHeight < 12 {
			lineHeight = 12
		}
		for lineIdx, line := range wrapped {
			xOffset := levelIndent + leftIndent

			if lineIdx == 0 && len(prefixRuns) > 0 {
				// With a hanging indent the bullet sits in the hang. Without one
				// it goes at the paragraph's own left edge: PowerPoint keeps the
				// bullet inside the shape, so it must not be pulled out to the
				// left of the text box.
				prefixX := xOffset + hangingIndent
				lines = append(lines, shapeParagraphLayoutLine{
					runs:       prefixRuns,
					xOffset:    prefixX,
					lineHeight: 0,
					align:      elements.TextAlignLeft,
					availWidth: availWidth,
					tabStops:   nil,
				})
			}

			lines = append(lines, shapeParagraphLayoutLine{
				runs:       line,
				xOffset:    xOffset,
				lineHeight: lineHeight,
				align:      style.Align,
				availWidth: availWidth,
				tabStops:   tabStops,
			})
			totalHeight += lineHeight
		}
		prevSpaceAfter = paragraphAfterGap(style)
		totalHeight += prevSpaceAfter
	}
	return lines, totalHeight
}

// shapeParagraphStyle normalizes a paragraph belonging to shape text.
//
// NormalizeParagraphStyle defaults an unset bullet style to "bullet", which is
// what a body placeholder wants but not what a shape wants: PowerPoint draws no
// bullet on shape text unless the paragraph asks for one. Left alone, every
// plain text box came out with a bullet stamped over its first letter.
func shapeParagraphStyle(paragraph text.Paragraph) text.ParagraphStyle {
	style := elements.NormalizeParagraphStyle(paragraph.Style)
	if strings.TrimSpace(paragraph.Style.BulletStyle) == "" {
		style.BulletStyle = text.BulletStyleNone
	}
	return style
}

func normalizedShapeParagraphs(s shapes.Shape) []text.Paragraph {
	if len(s.TextParagraphs) > 0 {
		return s.TextParagraphs
	}
	return []text.Paragraph{{Runs: []text.Run{text.NewRun(s.Text)}}}
}

func buildShapeParagraphStyledRuns(runs []text.Run, fittedSize int) []pdfStyledRun {
	return buildPDFStyledRuns(runs, fittedSize, false, false)
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
