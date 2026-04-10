package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

const (
	minVerticalLineHeightPt = 12.0 // minimum line height for vertical text, in points
	splitRunsCapMultiplier  = 4    // capacity multiplier when splitting runs into single-rune tokens
)

func renderPDFShapeParagraphTextVertical(pdf *gopdf.GoPdf, s shapes.Shape, x, y, w, h float64) {
	boxX, boxY, boxW, boxH, anchor := shapeTextBox(s, x, y, w, h)
	if boxW <= 2 || boxH <= 2 {
		return
	}
	boxX, boxY, boxW, boxH, restoreOrientation := beginShapeTextOrientation(
		pdf, s.TextFrame, boxX, boxY, boxW, boxH, x, y, w, h,
	)
	defer restoreOrientation()
	paragraphs := normalizedShapeParagraphs(s)
	fontSize := fitPDFShapeParagraphText(pdf, paragraphs, boxW, boxH)
	lines := make([]verticalStyledLine, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		style := text.NormalizeParagraphStyle(paragraph.Style)
		runs := buildShapeParagraphStyledRuns(paragraph.Runs, fontSize)
		lineHeight := math.Max(
			paragraphRenderedLineHeight(style, maxStyledRunsLineHeight(runs)),
			minVerticalLineHeightPt,
		)
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
	boxX, boxY, boxW, boxH, restoreOrientation := beginShapeTextOrientation(
		pdf, s.TextFrame, boxX, boxY, boxW, boxH, x, y, w, h,
	)
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
		lineHeight: math.Max(pdfLineHeight(fontSize), minVerticalLineHeightPt),
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
		colWidth = minVerticalLineHeightPt
	}
	columnCount := countVerticalColumnsNeeded(lines, boxH)
	usedWidth := float64(columnCount) * colWidth
	startX := boxX + boxW - colWidth
	switch anchor {
	case shapes.TextAnchorTop:
		// startX already set to rightmost column — top anchor is the default layout
	case shapes.TextAnchorBottom:
		startX = boxX + usedWidth - colWidth
	case shapes.TextAnchorMiddle:
		startX = boxX + (boxW+usedWidth)/2 - colWidth
	}
	cursorX := min(startX, boxX+boxW-colWidth)
	for _, line := range lines {
		runes := splitStyledRunsForVertical(line.runs)
		var done bool
		cursorX, done = renderVerticalLineRunes(pdf, runes, line, cursorX, colWidth, boxX, boxY, boxH)
		if done {
			break
		}
		cursorX -= colWidth
		if cursorX < boxX-nearZeroEpsilon {
			break
		}
	}
	setPDFTextFontWithHint(pdf, defaultFontSize, false, false, "")
}

// renderVerticalLineRunes renders one line's rune tokens vertically, wrapping to the next
// column when the box height is exceeded. It returns the updated cursorX and whether the
// box boundary was exceeded (caller should stop).
func renderVerticalLineRunes(
	pdf *gopdf.GoPdf,
	runes []pdfStyledRun,
	line verticalStyledLine,
	cursorX, colWidth, boxX, boxY, boxH float64,
) (float64, bool) {
	cursorY := boxY
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
		if cursorX < boxX-nearZeroEpsilon {
			return cursorX, true
		}
		advance := max(measureStyledRunWidth(pdf, run), colWidth/2)
		if run.HasHighlight {
			renderPDFStyledRunHighlight(pdf, run, cursorX, cursorY, advance, line.lineHeight)
		}
		renderPDFStyledRunText(pdf, run, cursorX, cursorY)
		cursorY += line.lineHeight
	}
	return cursorX, false
}

func splitStyledRunsForVertical(runs []pdfStyledRun) []pdfStyledRun {
	out := make([]pdfStyledRun, 0, len(runs)*splitRunsCapMultiplier)
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
	return max(maxWidth+2, minVerticalLineHeightPt)
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
