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

func renderPDFBullets(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	yPos := 60.0
	if slide.Title != "" {
		yPos = 108.0
	}
	maxY := slideHeightPt - 24
	if slide.Table != nil {
		tableTop := emuToPt(slide.Table.Y.Emu())
		if tableTop > yPos {
			maxY = tableTop - 12
		}
	}
	baseX := 54.0
	maxWidth := slideWidthPt - 108
	for i, bullet := range slide.Bullets {
		style := bulletStyleForIndex(slide, i)
		runs := bulletRunsForIndex(slide, i, bullet)
		levelIndent := float64(style.Level * 14)
		fontSize := defaultFontSize
		if sz := firstRunSize(runs); sz > 0 {
			fontSize = sz
		}
		bold, italic := runTextStyle(runs, slide)
		prefix := bulletPrefix(style, i)
		renderedText := renderRunsPlain(runs)
		if prefix != "" {
			renderedText = prefix + " " + renderedText
		}
		availableWidth := maxWidth - levelIndent
		if availableWidth < 80 {
			availableWidth = 80
		}
		fontSize = fitPDFTextSize(pdf, renderedText, fontSize, minTextAutoFitSize, bold, italic, availableWidth)
		setPDFTextFont(pdf, fontSize, bold, italic)
		lines := wrapPDFText(pdf, renderedText, availableWidth)
		lineHeight := math.Max(float64(fontSize)+2, 14)
		pdf.SetTextColor(runTextColor(runs))
		for li, line := range lines {
			if yPos+lineHeight > maxY {
				setPDFTextFont(pdf, defaultFontSize, false, false)
				return
			}
			lineX := baseX + levelIndent
			if li > 0 && prefix != "" {
				lineX += 12
			}
			pdf.SetX(lineX)
			pdf.SetY(yPos)
			_ = pdf.Cell(nil, line)
			yPos += lineHeight
		}
		yPos += 2
	}
	setPDFTextFont(pdf, defaultFontSize, false, false)
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

func bulletPrefix(style text.ParagraphStyle, idx int) string {
	switch text.NormalizeBulletStyle(style.BulletStyle) {
	case text.BulletStyleNone:
		return ""
	case text.BulletStyleNumber:
		return fmt.Sprintf("%d.", idx+1)
	case text.BulletStyleLetterLower:
		return fmt.Sprintf("%c.", 'a'+(idx%26))
	case text.BulletStyleLetterUpper:
		return fmt.Sprintf("%c.", 'A'+(idx%26))
	case text.BulletStyleRomanLower:
		return strings.ToLower(romanNumeral(idx + 1))
	case text.BulletStyleRomanUpper:
		return romanNumeral(idx + 1)
	case text.BulletStyleCustom:
		if style.BulletChar != "" {
			return style.BulletChar
		}
		return "•"
	default:
		return "•"
	}
}

func romanNumeral(n int) string {
	if n <= 0 {
		return "I"
	}
	table := []struct {
		value int
		sym   string
	}{
		{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
		{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
		{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
	}
	var out strings.Builder
	for _, entry := range table {
		for n >= entry.value {
			out.WriteString(entry.sym)
			n -= entry.value
		}
	}
	return out.String()
}

func setPDFTextFont(pdf *gopdf.GoPdf, size int, bold bool, italic bool) {
	style := ""
	if bold {
		style += "B"
	}
	if italic {
		style += "I"
	}
	if size <= 0 {
		size = defaultFontSize
	}
	_ = pdf.SetFont("sans", style, size)
}
