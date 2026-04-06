//nolint:mnd // Text/background renderer uses fixed PPT layout constants and spacing defaults.
package export

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

const (
	defaultBulletPrefix = "•"
	vAlignCenter        = "ctr"
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
	baseX := 36.0 // matches the standard PPT "Title and Content" content placeholder left edge (457200 EMU)
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
	case vAlignCenter:
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
		leftIndent, hangingIndent := resolveIndent(style)
		rightIndent := emuToPt(style.RightIndent.Emu())
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
		availableWidth := maxWidth - leftIndent - rightIndent
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
			lineX := baseX + leftIndent

			if li == 0 && prefix != "" {
				prefixRuns := buildBulletPrefixRuns(prefix, style, slide, fontSize, fontHint, runs)
				prefixX := lineX + hangingIndent
				renderStyledLine(pdf, prefixRuns, prefixX, yPos)
			}

			align := elements.NormalizeTextAlign(style.Align)
			if align == elements.TextAlignCenter || align == elements.TextAlignRight {
				lineText := styledLinePlain(line)
				lineX = alignedTextX(pdf, lineText, baseX+leftIndent, availableWidth, style.Align, fontHint)
			}
			renderStyledLine(pdf, line, lineX, yPos)
			yPos += lineHeight
		}
		prevSpaceAfter = paragraphAfterGap(style)
	}
	setPDFTextFontWithHint(pdf, defaultFontSize, false, false, "")
}
