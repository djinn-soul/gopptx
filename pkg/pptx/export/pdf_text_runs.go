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
	Lang           string
	SizePt         int
	HasHighlight   bool
	HighlightColor [3]uint8
	HasOutline     bool
	OutlineColor   [3]uint8
	OutlineWidthPt float64
	Underline      string // OOXML u value: "sng", "dbl", "dotted", …
	Strikethrough  string // OOXML strike value: "sngStrike", "dblStrike"
	Subscript      bool
	Superscript    bool
	SmallCaps      bool
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

// breakLongStyledRun splits a token that cannot fit its box on any line.
//
// It always returns at least one piece, so the caller can treat the last one as
// the remainder to keep filling. A single rune wider than the box still gets
// its own line rather than looping forever.
func breakLongStyledRun(
	pdf *gopdf.GoPdf,
	run pdfStyledRun,
	maxWidth float64,
	tabStops []float64,
) []pdfStyledRun {
	pieces := make([]pdfStyledRun, 0, 2)
	var b strings.Builder
	flush := func() {
		piece := run
		piece.Text = b.String()
		pieces = append(pieces, piece)
		b.Reset()
	}
	for _, r := range run.Text {
		candidate := run
		candidate.Text = b.String() + string(r)
		if b.Len() == 0 || measureStyledRunAdvance(pdf, candidate, 0, tabStops) <= maxWidth {
			b.WriteRune(r)
			continue
		}
		flush()
		b.WriteRune(r)
	}
	flush()
	return pieces
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
			w = measureStyledRunAdvance(pdf, tok, lineW, tabStops)
		}
		if len(line) == 0 && isWhitespaceOnlyRun(tok, true) {
			continue
		}
		// A single token wider than the box has nowhere to wrap. The plain-text
		// path splits it mid-word; without the same here, a long label such as
		// [Subroutine] ran straight out of its node instead of breaking the way
		// PowerPoint does.
		if len(line) == 0 && w > maxWidth {
			pieces := breakLongStyledRun(pdf, tok, maxWidth, tabStops)
			for _, piece := range pieces[:len(pieces)-1] {
				lines = append(lines, []pdfStyledRun{piece})
			}
			last := pieces[len(pieces)-1]
			line = append(line, last)
			lineW = measureStyledRunAdvance(pdf, last, 0, tabStops)
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
	// The run's typeface has to be embedded before its size is worked out: the
	// script and small-caps factors come from the font's own metrics, and an
	// unregistered family would answer with the fallback's.
	ensureStyledRunFonts(pdf, run)
	setPDFTextFontWithHintAndLang(
		pdf, effectiveRunSizePt(pdf, run), run.Bold, run.Italic, run.FontHint, run.Lang,
	)
	return measuredWidth(pdf, run.Text)
}

// ensureStyledRunFonts embeds the typeface each run names, so that every metric
// read afterwards belongs to the font the text is actually drawn in.
func ensureStyledRunFonts(pdf *gopdf.GoPdf, runs ...pdfStyledRun) {
	for _, run := range runs {
		if run.FontHint != "" {
			ensureNamedPDFFont(pdf, run.FontHint)
		}
	}
}

func measureStyledRunAdvance(pdf *gopdf.GoPdf, run pdfStyledRun, cursorOffset float64, tabStops []float64) float64 {
	if run.Text == "\t" {
		return nextPDFTabAdvance(cursorOffset, tabStops)
	}
	return measureStyledRunWidth(pdf, run)
}

func renderStyledLine(pdf *gopdf.GoPdf, line []pdfStyledRun, x, y float64, opts pdfTextRenderOptions) {
	cursorX := x
	// Every run's font must be embedded before the shared baseline is derived
	// from their metrics.
	ensureStyledRunFonts(pdf, line...)
	baseline := styledLineBaselineOffset(pdf, line)
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
		renderPDFStyledRunTextOnBaseline(pdf, run, cursorX, y, baseline)
		cursorX += advance
	}
}

// renderPDFStyledRunText draws a run that sits alone on its line, on the
// baseline its own font asks for.
func renderPDFStyledRunText(pdf *gopdf.GoPdf, run pdfStyledRun, x, y float64) {
	ensureStyledRunFonts(pdf, run)
	renderPDFStyledRunTextOnBaseline(pdf, run, x, y, runBaselineOffset(pdf, run))
}

// renderPDFStyledRunTextOnBaseline draws a run with its baseline at
// y+baseline, whatever its own font's metrics would have put it. Every run of a
// line shares one baseline, so a line mixing Georgia with Verdana does not step
// up and down mid-word.
func renderPDFStyledRunTextOnBaseline(pdf *gopdf.GoPdf, run pdfStyledRun, x, y, baseline float64) {
	if run.HasOutline {
		renderPDFStyledRunOutline(pdf, run, x, y, baseline)
	}
	size := effectiveRunSizePt(pdf, run)
	// A sub- or superscript run is drawn shifted off the shared baseline.
	runBaseline := baseline + scriptBaselineOffsetPt(pdf, run)
	setPDFTextFontWithHintAndLang(pdf, size, run.Bold, run.Italic, run.FontHint, run.Lang)
	pdf.SetTextColor(run.Color[0], run.Color[1], run.Color[2])
	pdf.SetX(x)
	// Cell() drops the baseline by the font's own typoAscender below the Y it is
	// given, so that much is taken back out of the requested baseline.
	pdf.SetY(y + runBaseline - gopdfBaselineDropPt(pdf, run))
	_ = pdf.Cell(nil, run.Text)
	renderPDFStyledRunDecorations(pdf, run, x, y+runBaseline)
}

func renderPDFStyledRunOutline(pdf *gopdf.GoPdf, run pdfStyledRun, x, y, baseline float64) {
	offset := pdfOutlineOffset(run.OutlineWidthPt)
	outlineRun := run
	outlineRun.HasOutline = false
	outlineRun.Color = run.OutlineColor
	size := effectiveRunSizePt(pdf, outlineRun)
	y += baseline + scriptBaselineOffsetPt(pdf, outlineRun) - gopdfBaselineDropPt(pdf, outlineRun)
	for _, delta := range [][2]float64{
		{-offset, 0},
		{offset, 0},
		{0, -offset},
		{0, offset},
	} {
		setPDFTextFontWithHintAndLang(
			pdf,
			size,
			outlineRun.Bold,
			outlineRun.Italic,
			outlineRun.FontHint,
			outlineRun.Lang,
		)
		pdf.SetTextColor(outlineRun.Color[0], outlineRun.Color[1], outlineRun.Color[2])
		pdf.SetX(x + delta[0])
		pdf.SetY(y + delta[1])
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
