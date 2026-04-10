package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func renderPDFShapeParagraphTextVertical(pdf *gopdf.GoPdf, s shapes.Shape, x, y, w, h float64) {
	boxX, boxY, boxW, boxH, anchor := shapeTextBox(s, x, y, w, h)
	if boxW <= 2 || boxH <= 2 {
		return
	}
	boxX, boxY, boxW, boxH, restoreOrientation := beginShapeTextOrientation(pdf, s.TextFrame, boxX, boxY, boxW, boxH, x, y, w, h)
	defer restoreOrientation()
	paragraphs := normalizedShapeParagraphs(s)
	fontSize := fitPDFShapeParagraphText(pdf, paragraphs, boxW, boxH)
	lines := make([]verticalStyledLine, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		style := text.NormalizeParagraphStyle(paragraph.Style)
		runs := buildShapeParagraphStyledRuns(paragraph.Runs, fontSize)
		lineHeight := math.Max(paragraphRenderedLineHeight(style, maxStyledRunsLineHeight(runs)), 12)
		lines = append(lines, verticalStyledLine{
			runs:       runs,
			lineHeight: lineHeight,
		})
	}
	renderVerticalStyledLines(pdf, lines, boxX, boxY, boxW, boxH, anchor)
}

func renderPDFShapePlainTextVertical(pdf *gopdf.GoPdf, s shapes.Shape, x, y, w, h float64) {
	boxX, boxY, boxW, boxH, anchor := shapeTextBox(s, x, y, w, h)
	if boxW <= 2 || boxH <= 2 {
		return
	}
	boxX, boxY, boxW, boxH, restoreOrientation := beginShapeTextOrientation(pdf, s.TextFrame, boxX, boxY, boxW, boxH, x, y, w, h)
	defer restoreOrientation()

	fontHint := inferCodeFontHint(s.Text)
	fontSize := fitPDFTextToBoxWithMetrics(
		pdf,
		s.Text,
		defaultFontSize,
		minTextAutoFitSize,
		false,
		false,
		boxW,
		boxH,
		fontHint,
	)
	runs := buildPDFStyledRuns([]text.Run{text.NewRun(s.Text).WithFont(fontHint)}, fontSize, false, false)
	renderVerticalStyledLines(pdf, []verticalStyledLine{{
		runs:       runs,
		lineHeight: math.Max(pdfLineHeight(fontSize), 12),
	}}, boxX, boxY, boxW, boxH, anchor)
}

type verticalStyledLine struct {
	runs       []pdfStyledRun
	lineHeight float64
}

func renderVerticalStyledLines(
	pdf *gopdf.GoPdf,
	lines []verticalStyledLine,
	boxX, boxY, boxW, boxH float64,
	anchor shapes.TextFrameAnchor,
) {
	if len(lines) == 0 {
		return
	}
	colWidth := maxVerticalColumnWidth(pdf, lines)
	if colWidth <= 0 {
		colWidth = 12
	}
	columnCount := countVerticalColumnsNeeded(lines, boxH)
	usedWidth := float64(columnCount) * colWidth
	startX := boxX + boxW - colWidth
	switch anchor {
	case shapes.TextAnchorBottom:
		startX = boxX + usedWidth - colWidth
	case shapes.TextAnchorMiddle:
		startX = boxX + (boxW+usedWidth)/2 - colWidth
	}
	cursorX := min(startX, boxX+boxW-colWidth)
	for _, line := range lines {
		cursorY := boxY
		runes := splitStyledRunsForVertical(line.runs)
		for _, run := range runes {
			if run.Text == "\n" {
				cursorX -= colWidth
				cursorY = boxY
				continue
			}
			if run.Text == "\t" {
				cursorY += line.lineHeight * 2
				continue
			}
			if cursorY+line.lineHeight > boxY+boxH {
				cursorX -= colWidth
				cursorY = boxY
			}
			if cursorX < boxX-0.01 {
				return
			}
			advance := max(measureStyledRunWidth(pdf, run), colWidth*0.5)
			if run.HasHighlight {
				renderPDFStyledRunHighlight(pdf, run, cursorX, cursorY, advance, line.lineHeight)
			}
			renderPDFStyledRunText(pdf, run, cursorX, cursorY)
			cursorY += line.lineHeight
		}
		cursorX -= colWidth
		if cursorX < boxX-0.01 {
			return
		}
	}
	setPDFTextFontWithHint(pdf, defaultFontSize, false, false, "")
}

func splitStyledRunsForVertical(runs []pdfStyledRun) []pdfStyledRun {
	out := make([]pdfStyledRun, 0, len(runs)*4)
	for _, run := range runs {
		for _, r := range run.Text {
			tok := run
			tok.Text = string(r)
			out = append(out, tok)
		}
	}
	return out
}

func maxVerticalColumnWidth(pdf *gopdf.GoPdf, lines []verticalStyledLine) float64 {
	maxWidth := 0.0
	for _, line := range lines {
		for _, run := range splitStyledRunsForVertical(line.runs) {
			if run.Text == "\n" || run.Text == "\t" {
				continue
			}
			if w := measureStyledRunWidth(pdf, run); w > maxWidth {
				maxWidth = w
			}
		}
	}
	return max(maxWidth+2, 12)
}

func countVerticalColumnsNeeded(lines []verticalStyledLine, boxH float64) int {
	total := 0
	for _, line := range lines {
		if line.lineHeight <= 0 {
			continue
		}
		rowsPerCol := max(int(boxH/line.lineHeight), 1)
		runes := 0
		for _, run := range splitStyledRunsForVertical(line.runs) {
			if run.Text == "\n" {
				total += max((runes+rowsPerCol-1)/rowsPerCol, 1)
				runes = 0
				continue
			}
			if run.Text == "\t" {
				runes += 2
				continue
			}
			runes++
		}
		total += max((runes+rowsPerCol-1)/rowsPerCol, 1)
	}
	return max(total, 1)
}
