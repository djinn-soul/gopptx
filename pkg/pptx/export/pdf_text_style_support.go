package export

import (
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

const defaultPDFTabStepPt = 36.0

type pdfTextRenderOptions struct {
	LineHeight float64
	TabStops   []float64
}

func paragraphRenderedLineHeight(style text.ParagraphStyle, baseLineHeight float64) float64 {
	if style.LineSpacingPts > 0 {
		return float64(style.LineSpacingPts)
	}
	return baseLineHeight * paragraphLineSpacingFactor(style)
}

func paragraphTabStopsPt(style text.ParagraphStyle) []float64 {
	if len(style.TabStops) == 0 {
		return nil
	}
	out := make([]float64, 0, len(style.TabStops))
	for _, stop := range style.TabStops {
		out = append(out, emuToPt(stop.Emu()))
	}
	return out
}

func buildPDFStyledRuns(runs []text.Run, fittedSize int, defaultBold, defaultItalic bool) []pdfStyledRun {
	if len(runs) == 0 {
		return []pdfStyledRun{{Text: "", Color: [3]uint8{0, 0, 0}, SizePt: fittedSize}}
	}
	out := make([]pdfStyledRun, 0, len(runs))
	for _, run := range runs {
		out = append(out, pdfStyledRunFromTextRun(run, fittedSize, defaultBold, defaultItalic))
	}
	return out
}

func pdfStyledRunFromTextRun(run text.Run, fittedSize int, defaultBold, defaultItalic bool) pdfStyledRun {
	size := fittedSize
	if run.SizePt > 0 && run.SizePt < size {
		size = run.SizePt
	}
	color := [3]uint8{0, 0, 0}
	if run.Color != "" {
		r, g, b := hexToRGB(run.Color)
		color = [3]uint8{r, g, b}
	}
	fontHint := strings.TrimSpace(run.Font)
	if fontHint == "" {
		fontHint = inferCodeFontHint(run.Text)
	}
	styled := pdfStyledRun{
		Text:     run.Text,
		Bold:     run.Bold || defaultBold,
		Italic:   run.Italic || defaultItalic,
		Color:    color,
		FontHint: fontHint,
		Lang:     strings.TrimSpace(run.Lang),
		SizePt:   size,
	}
	if run.Highlight != "" {
		r, g, b := hexToRGB(run.Highlight)
		styled.HasHighlight = true
		styled.HighlightColor = [3]uint8{r, g, b}
	}
	if run.OutlineColor != "" {
		r, g, b := hexToRGB(run.OutlineColor)
		styled.HasOutline = true
		styled.OutlineColor = [3]uint8{r, g, b}
		styled.OutlineWidthPt = run.OutlineWidthPt
	}
	return styled
}

func nextPDFTabAdvance(cursorOffset float64, tabStops []float64) float64 {
	for _, stop := range tabStops {
		if stop > cursorOffset+0.01 {
			return stop - cursorOffset
		}
	}
	return defaultPDFTabAdvance(cursorOffset)
}

func defaultPDFTabAdvance(cursorOffset float64) float64 {
	remainder := math.Mod(cursorOffset, defaultPDFTabStepPt)
	if math.Abs(remainder) < 0.01 {
		return defaultPDFTabStepPt
	}
	return defaultPDFTabStepPt - remainder
}
