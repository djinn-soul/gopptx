//nolint:mnd // Text/background renderer uses fixed PPT layout constants and spacing defaults.
package export

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

const (
	defaultBulletPrefix = "•"
)

//nolint:gocognit // Background rendering handles multiple media/fill modes with explicit branches.
func renderPDFBackground(pdf *gopdf.GoPdf, bg *elements.SlideBackground) error {
	if bg == nil {
		pdf.SetFillColor(255, 255, 255)
		pdf.RectFromUpperLeftWithStyle(0, 0, slideWidthPt, slideHeightPt, "F")
		return nil
	}

	switch bg.Type {
	case elements.SlideBackgroundSolid:
		if bg.SolidFill != nil && bg.SolidFill.Color != "" {
			pdf.SetFillColor(hexToRGB(bg.SolidFill.Color))
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.RectFromUpperLeftWithStyle(0, 0, slideWidthPt, slideHeightPt, "F")
	case elements.SlideBackgroundGradient:
		if !renderPDFGradientBackground(pdf, bg.GradientFill) {
			pdf.SetFillColor(255, 255, 255)
			if bg.GradientFill != nil && len(bg.GradientFill.Stops) > 0 {
				pdf.SetFillColor(hexToRGB(bg.GradientFill.Stops[0].Color))
			}
			pdf.RectFromUpperLeftWithStyle(0, 0, slideWidthPt, slideHeightPt, "F")
		}
	case elements.SlideBackgroundPicture:
		pdf.SetFillColor(255, 255, 255)
		pdf.RectFromUpperLeftWithStyle(0, 0, slideWidthPt, slideHeightPt, "F")
		if bg.PictureFill == nil {
			return nil
		}
		data := bg.PictureFill.Data
		if len(data) == 0 && bg.PictureFill.Path != "" {
			fileData, err := os.ReadFile(bg.PictureFill.Path)
			if err != nil {
				return fmt.Errorf("read picture background: %w", err)
			}
			data = fileData
		}
		if len(data) == 0 {
			return nil
		}
		holder, err := gopdf.ImageHolderByBytes(data)
		if err != nil {
			return fmt.Errorf("load picture background: %w", err)
		}
		return pdf.ImageByHolder(holder, 0, 0, &gopdf.Rect{W: slideWidthPt, H: slideHeightPt})
	default:
		pdf.SetFillColor(255, 255, 255)
		pdf.RectFromUpperLeftWithStyle(0, 0, slideWidthPt, slideHeightPt, "F")
	}
	return nil
}

//nolint:funlen,gocognit // Bullet rendering applies paragraph spacing/indents/line-fit rules in one ordered pass.
func renderPDFBullets(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	yPos := 60.0
	if slide.Title != "" {
		yPos = 108.0
	}
	maxY := slideHeightPt - 24
	baseX := 54.0
	maxWidth := slideWidthPt - 108
	if b := slide.ContentBoundsEMU; b[2] > 0 || b[3] > 0 {
		baseX = emuToPt(b[0])
		yPos = emuToPt(b[1])
		maxWidth = emuToPt(b[2])
		maxY = yPos + emuToPt(b[3])
	}
	if slide.Table != nil {
		tableTop := emuToPt(slide.Table.Y.Emu())
		if tableTop > yPos {
			maxY = tableTop - 12
		}
	}

	// Apply vertical alignment when ContentVAlign is "ctr" (middle) or "b" (bottom).
	switch slide.ContentVAlign {
	case "ctr":
		totalH := measureBulletsHeight(pdf, slide, maxWidth, maxY-yPos)
		yPos += math.Max(0, (maxY-yPos-totalH)/2)
	case "b":
		totalH := measureBulletsHeight(pdf, slide, maxWidth, maxY-yPos)
		yPos += math.Max(0, maxY-yPos-totalH)
	}

	prevSpaceAfter := 0.0
	for i, bullet := range slide.Bullets {
		style := bulletStyleForIndex(slide, i)
		runs := bulletRunsForIndex(slide, i, bullet)
		levelIndent := float64(style.Level * 14)
		leftIndent := emuToPt(style.LeftIndent.Emu())
		rightIndent := emuToPt(style.RightIndent.Emu())
		hangingIndent := emuToPt(style.HangingIndent.Emu())
		fontSize := defaultFontSize
		if sz := firstRunSize(runs); sz > 0 {
			fontSize = sz
		}
		bold, italic := runTextStyle(runs, slide)
		fontHint := firstRunFont(runs)
		prefix := bulletPrefix(style, i)
		if prefix == "" {
			// SlidesFromPPTX can lose explicit bullet-style metadata; preserve bullet intent.
			prefix = defaultBulletPrefix
		}
		renderedText := renderRunsPlain(runs)
		if strings.TrimSpace(fontHint) == "" {
			fontHint = inferCodeFontHint(renderedText)
		}
		yPos += paragraphStartGap(i, prevSpaceAfter, style)
		availableWidth := maxWidth - levelIndent - leftIndent - rightIndent
		if availableWidth < 80 {
			availableWidth = 80
		}
		fontSize = fitPDFTextToBoxWithMetrics(
			pdf,
			renderedText,
			fontSize,
			minTextAutoFitSize,
			bold,
			italic,
			availableWidth,
			maxY-yPos,
			fontHint,
		)
		setPDFTextFontWithHint(pdf, fontSize, bold, italic, fontHint)
		styledRuns := buildBulletStyledRuns(runs, slide, fontSize)
		lines := wrapStyledRuns(pdf, styledRuns, availableWidth)
		lineHeight := math.Max(pdfLineHeight(fontSize)*paragraphLineSpacingFactor(style), 12)
		for li, line := range lines {
			if yPos+lineHeight > maxY {
				setPDFTextFontWithHint(pdf, defaultFontSize, false, false, "")
				return
			}
			lineX := baseX + levelIndent + leftIndent

			if li == 0 && prefix != "" {
				prefixRuns := buildBulletPrefixRuns(prefix, style, slide, fontSize, fontHint, runs)
				prefixX := lineX - hangingIndent
				if hangingIndent == 0 {
					// Fallback gap if no hanging indent is specified.
					prefixX = lineX - 14
				}
				renderStyledLine(pdf, prefixRuns, prefixX, yPos)
			}

			align := elements.NormalizeTextAlign(style.Align)
			if align == elements.TextAlignCenter || align == elements.TextAlignRight {
				lineText := styledLinePlain(line)
				lineX = alignedTextX(pdf, lineText, baseX+levelIndent+leftIndent, availableWidth, style.Align, fontHint)
			}
			renderStyledLine(pdf, line, lineX, yPos)
			yPos += lineHeight
		}
		prevSpaceAfter = paragraphAfterGap(style)
	}
	setPDFTextFontWithHint(pdf, defaultFontSize, false, false, "")
}

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
		levelIndent := float64(style.Level * 14)
		leftIndent := emuToPt(style.LeftIndent.Emu())
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
		availableWidth := maxWidth - levelIndent - leftIndent - rightIndent
		if availableWidth < 80 {
			availableWidth = 80
		}
		fontSize = fitPDFTextToBoxWithMetrics(pdf, renderedText, fontSize, minTextAutoFitSize, bold, italic, availableWidth, availH-total, fontHint)
		lineHeight := math.Max(pdfLineHeight(fontSize)*paragraphLineSpacingFactor(style), 12)
		styledRuns := buildBulletStyledRuns(runs, slide, fontSize)
		lines := wrapStyledRuns(pdf, styledRuns, availableWidth)
		total += float64(len(lines)) * lineHeight
		prevSpaceAfter = paragraphAfterGap(style)
		}
		return total
}
