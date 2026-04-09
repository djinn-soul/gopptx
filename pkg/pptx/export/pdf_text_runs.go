//nolint:mnd // Run-token wrapping uses small fixed preallocation sizes for perf.
package export

import (
	"math"
	"strings"

	"github.com/signintech/gopdf"
)

type pdfStyledRun struct {
	Text           string
	Bold           bool
	Italic         bool
	Color          [3]uint8
	FontHint       string
	SizePt         int
	HasHighlight   bool
	HighlightColor [3]uint8
	HasOutline     bool
	OutlineColor   [3]uint8
	OutlineWidthPt float64
}

func splitStyledRunsForWrap(runs []pdfStyledRun) []pdfStyledRun {
	out := make([]pdfStyledRun, 0, len(runs)*3)
	for _, r := range runs {
		out = append(out, splitOneStyledRun(r)...)
	}
	return out
}

func splitOneStyledRun(run pdfStyledRun) []pdfStyledRun {
	if run.Text == "" {
		return nil
	}
	tokens := make([]pdfStyledRun, 0, 4)
	var b strings.Builder
	flush := func() {
		if b.Len() == 0 {
			return
		}
		t := run
		t.Text = b.String()
		tokens = append(tokens, t)
		b.Reset()
	}
	for _, r := range run.Text {
		ch := string(r)
		if ch == "\n" || ch == " " || ch == "\t" {
			flush()
			t := run
			t.Text = ch
			tokens = append(tokens, t)
			continue
		}
		b.WriteRune(r)
	}
	flush()
	return tokens
}

func wrapStyledRuns(pdf *gopdf.GoPdf, runs []pdfStyledRun, maxWidth float64, tabStops []float64) [][]pdfStyledRun {
	tokens := splitStyledRunsForWrap(runs)
	lines := make([][]pdfStyledRun, 0, 4)
	line := make([]pdfStyledRun, 0, 6)
	lineW := 0.0

	pushLine := func() {
		lines = append(lines, trimLeadingSpaceRuns(line))
		line = make([]pdfStyledRun, 0, 6)
		lineW = 0
	}

	for _, tok := range tokens {
		if tok.Text == "\n" {
			pushLine()
			continue
		}
		w := measureStyledRunAdvance(pdf, tok, lineW, tabStops)
		if len(line) > 0 && lineW+w > maxWidth && !isWhitespaceOnlyRun(tok, false) {
			pushLine()
		}
		if len(line) == 0 && isWhitespaceOnlyRun(tok, true) {
			continue
		}
		line = append(line, tok)
		lineW += w
	}
	if len(line) > 0 || len(lines) == 0 {
		lines = append(lines, trimLeadingSpaceRuns(line))
	}
	return lines
}

func trimLeadingSpaceRuns(runs []pdfStyledRun) []pdfStyledRun {
	if len(runs) == 0 {
		return runs
	}
	idx := 0
	for idx < len(runs) && isWhitespaceOnlyRun(runs[idx], true) {
		idx++
	}
	if idx >= len(runs) {
		return []pdfStyledRun{}
	}
	return runs[idx:]
}

func measureStyledRunWidth(pdf *gopdf.GoPdf, run pdfStyledRun) float64 {
	if run.Text == "" {
		return 0
	}
	size := run.SizePt
	if size <= 0 {
		size = defaultFontSize
	}
	setPDFTextFontWithHint(pdf, size, run.Bold, run.Italic, run.FontHint)
	return measuredWidthWithMetrics(pdf, run.Text, run.FontHint)
}

func measureStyledRunAdvance(pdf *gopdf.GoPdf, run pdfStyledRun, cursorOffset float64, tabStops []float64) float64 {
	if run.Text == "\t" {
		return nextPDFTabAdvance(cursorOffset, tabStops)
	}
	return measureStyledRunWidth(pdf, run)
}

func renderStyledLine(pdf *gopdf.GoPdf, line []pdfStyledRun, x, y float64, opts pdfTextRenderOptions) {
	cursorX := x
	for _, run := range line {
		if run.Text == "" {
			continue
		}
		advance := measureStyledRunAdvance(pdf, run, cursorX-x, opts.TabStops)
		if run.Text == "\t" {
			cursorX += advance
			continue
		}
		size := run.SizePt
		if size <= 0 {
			size = defaultFontSize
		}
		lineHeight := opts.LineHeight
		if lineHeight <= 0 {
			lineHeight = pdfLineHeight(size)
		}
		if run.HasHighlight {
			renderPDFStyledRunHighlight(pdf, run, cursorX, y, advance, lineHeight)
		}
		renderPDFStyledRunText(pdf, run, cursorX, y)
		cursorX += advance
	}
}

func renderPDFStyledRunText(pdf *gopdf.GoPdf, run pdfStyledRun, x, y float64) {
	if run.HasOutline {
		renderPDFStyledRunOutline(pdf, run, x, y)
	}
	setPDFTextFontWithHint(pdf, runSizePt(run), run.Bold, run.Italic, run.FontHint)
	pdf.SetTextColor(run.Color[0], run.Color[1], run.Color[2])
	pdf.SetX(x)
	pdf.SetY(y + fontBaselineShift(run.FontHint, runSizePt(run)))
	_ = pdf.Cell(nil, run.Text)
}

func renderPDFStyledRunOutline(pdf *gopdf.GoPdf, run pdfStyledRun, x, y float64) {
	offset := pdfOutlineOffset(run.OutlineWidthPt)
	outlineRun := run
	outlineRun.HasOutline = false
	outlineRun.Color = run.OutlineColor
	for _, delta := range [][2]float64{
		{-offset, 0},
		{offset, 0},
		{0, -offset},
		{0, offset},
	} {
		setPDFTextFontWithHint(pdf, runSizePt(outlineRun), outlineRun.Bold, outlineRun.Italic, outlineRun.FontHint)
		pdf.SetTextColor(outlineRun.Color[0], outlineRun.Color[1], outlineRun.Color[2])
		pdf.SetX(x + delta[0])
		pdf.SetY(y + delta[1] + fontBaselineShift(outlineRun.FontHint, runSizePt(outlineRun)))
		_ = pdf.Cell(nil, outlineRun.Text)
	}
}

func renderPDFStyledRunHighlight(
	pdf *gopdf.GoPdf,
	run pdfStyledRun,
	x, y, width, lineHeight float64,
) {
	if width <= 0 {
		return
	}
	rectHeight := math.Max(lineHeight*0.72, 3)
	rectY := y + (lineHeight-rectHeight)/2
	pdf.SetFillColor(run.HighlightColor[0], run.HighlightColor[1], run.HighlightColor[2])
	pdf.RectFromUpperLeftWithStyle(x, rectY, width, rectHeight, "F")
}

func runSizePt(run pdfStyledRun) int {
	if run.SizePt > 0 {
		return run.SizePt
	}
	return defaultFontSize
}

func pdfOutlineOffset(widthPt float64) float64 {
	if widthPt <= 0 {
		return 0.45
	}
	return math.Min(math.Max(widthPt*0.35, 0.35), 1.4)
}

func isWhitespaceOnlyRun(run pdfStyledRun, trimLeading bool) bool {
	if trimLeading && run.Text == "\t" {
		return false
	}
	return strings.TrimSpace(run.Text) == ""
}
