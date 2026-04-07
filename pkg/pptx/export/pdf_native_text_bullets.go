//nolint:mnd // Bullet layout helpers use fixed typographic constants inherited from PPT rendering rules.
package export

import (
	"math"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func bulletStyleForIndex(slide elements.SlideContent, idx int) text.ParagraphStyle {
	if idx < len(slide.BulletStyles) {
		return slide.BulletStyles[idx]
	}
	return slide.DefaultBulletStyle
}

func bulletRunsForIndex(slide elements.SlideContent, idx int, fallbackText string) []elements.Run {
	if idx < len(slide.BulletRuns) && len(slide.BulletRuns[idx]) > 0 {
		return slide.BulletRuns[idx]
	}
	return []elements.Run{elements.NewRun(fallbackText)}
}

func firstRunSize(runs []elements.Run) int {
	if len(runs) == 0 {
		return 0
	}
	return runs[0].SizePt
}

func firstRunFont(runs []elements.Run) string {
	if len(runs) == 0 {
		return ""
	}
	return runs[0].Font
}

func runTextColor(runs []elements.Run) (uint8, uint8, uint8) {
	if len(runs) > 0 && runs[0].Color != "" {
		return hexToRGB(runs[0].Color)
	}
	return 60, 60, 60
}

func runTextStyle(runs []elements.Run, slide elements.SlideContent) (bool, bool) {
	if len(runs) > 0 {
		return runs[0].Bold || slide.ContentBold, runs[0].Italic || slide.ContentItalic
	}
	return slide.ContentBold, slide.ContentItalic
}

func renderRunsPlain(runs []elements.Run) string {
	var out strings.Builder
	for _, run := range runs {
		out.WriteString(run.Text)
	}
	return out.String()
}

func styledLinePlain(line []pdfStyledRun) string {
	var b strings.Builder
	for _, run := range line {
		b.WriteString(run.Text)
	}
	return b.String()
}

func buildBulletStyledRuns(runs []elements.Run, slide elements.SlideContent, fittedSize int) []pdfStyledRun {
	out := make([]pdfStyledRun, 0, len(runs))
	for _, run := range runs {
		size := fittedSize
		if run.SizePt > 0 && run.SizePt < size {
			size = run.SizePt
		}
		cr, cg, cb := runTextColor(runs)
		if run.Color != "" {
			cr, cg, cb = hexToRGB(run.Color)
		}
		out = append(out, pdfStyledRun{
			Text:     run.Text,
			Bold:     run.Bold || slide.ContentBold,
			Italic:   run.Italic || slide.ContentItalic,
			Color:    [3]uint8{cr, cg, cb},
			FontHint: run.Font,
			SizePt:   size,
		})
	}
	if len(out) == 0 {
		cr, cg, cb := runTextColor(runs)
		out = append(out, pdfStyledRun{
			Text:   "",
			Color:  [3]uint8{cr, cg, cb},
			SizePt: fittedSize,
		})
	}
	return out
}

func buildBulletPrefixRuns(
	prefix string,
	style text.ParagraphStyle,
	slide elements.SlideContent,
	fittedSize int,
	fontHint string,
	runs []elements.Run,
) []pdfStyledRun {
	if prefix == "" {
		return nil
	}
	cr, cg, cb := runTextColor(runs)
	if style.BulletColor != "" {
		cr, cg, cb = hexToRGB(style.BulletColor)
	}
	return []pdfStyledRun{
		{
			Text:     prefix,
			Bold:     slide.ContentBold,
			Italic:   slide.ContentItalic,
			Color:    [3]uint8{cr, cg, cb},
			FontHint: fontHint,
			SizePt:   fittedSize,
		},
	}
}

// measureBulletsHeight computes the total rendered height of all bullets
// to support vertical alignment pre-positioning.
func measureBulletsHeight(pdf *gopdf.GoPdf, slide elements.SlideContent, maxWidth, availH float64) float64 {
	total := 0.0
	prevSpaceAfter := 0.0
	for i, bullet := range slide.Bullets {
		style := bulletStyleForIndex(slide, i)
		runs := bulletRunsForIndex(slide, i, bullet)
		leftIndent, _ := resolveIndent(style)
		rightIndent := emuToPt(style.RightIndent.Emu())
		fontSize := defaultFontSize
		if sz := firstRunSize(runs); sz > 0 {
			fontSize = sz
		}
		bold, italic := runTextStyle(runs, slide)
		fontHint := firstRunFont(runs)
		renderedText := renderRunsPlain(runs)
		gap := paragraphStartGap(i, prevSpaceAfter, style)
		total += gap
		availableWidth := maxWidth - leftIndent - rightIndent
		if availableWidth < 80 {
			availableWidth = 80
		}
		fontSize = fitPDFTextToBoxWithMetrics(
			pdf, renderedText, fontSize, minTextAutoFitSize, bold, italic, availableWidth, availH-total, fontHint,
		)
		lineHeight := math.Max(pdfLineHeight(fontSize)*paragraphLineSpacingFactor(style), 12)
		styledRuns := buildBulletStyledRuns(runs, slide, fontSize)
		lines := wrapStyledRuns(pdf, styledRuns, availableWidth)
		total += float64(len(lines)) * lineHeight
		prevSpaceAfter = paragraphAfterGap(style)
	}
	return total
}

func resolveIndent(style text.ParagraphStyle) (float64, float64) {
	leftIndent := emuToPt(style.LeftIndent.Emu())
	hangingIndent := emuToPt(style.HangingIndent.Emu())
	if leftIndent == 0 {
		leftIndent = 27.0 + float64(style.Level)*31.5
	}
	if hangingIndent == 0 {
		hangingIndent = -18.0
	}
	return leftIndent, hangingIndent
}
