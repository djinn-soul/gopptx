//nolint:mnd // Run-token wrapping uses small fixed preallocation sizes for perf.
package export

import (
	"strings"

	"github.com/signintech/gopdf"
)

type pdfStyledRun struct {
	Text     string
	Bold     bool
	Italic   bool
	Color    [3]uint8
	FontHint string
	SizePt   int
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

func wrapStyledRuns(pdf *gopdf.GoPdf, runs []pdfStyledRun, maxWidth float64) [][]pdfStyledRun {
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
		w := measureStyledRunWidth(pdf, tok)
		if len(line) > 0 && lineW+w > maxWidth && strings.TrimSpace(tok.Text) != "" {
			pushLine()
		}
		if len(line) == 0 && strings.TrimSpace(tok.Text) == "" {
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
	for idx < len(runs) && strings.TrimSpace(runs[idx].Text) == "" {
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

func measureStyledLineWidth(pdf *gopdf.GoPdf, line []pdfStyledRun) float64 {
	total := 0.0
	for _, run := range line {
		total += measureStyledRunWidth(pdf, run)
	}
	return total
}

func renderStyledLine(pdf *gopdf.GoPdf, line []pdfStyledRun, x, y float64) {
	cursorX := x
	for _, run := range line {
		if run.Text == "" {
			continue
		}
		size := run.SizePt
		if size <= 0 {
			size = defaultFontSize
		}
		setPDFTextFontWithHint(pdf, size, run.Bold, run.Italic, run.FontHint)
		pdf.SetTextColor(run.Color[0], run.Color[1], run.Color[2])
		pdf.SetX(cursorX)
		pdf.SetY(y + fontBaselineShift(run.FontHint, size))
		_ = pdf.Cell(nil, run.Text)
		cursorX += measureStyledRunWidth(pdf, run)
	}
}
