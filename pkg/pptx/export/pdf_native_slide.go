package export

import (
	"fmt"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

//nolint:mnd // Footer placement and colors match the native PPT slide template.
func renderNativePDFFooter(pdf *gopdf.GoPdf, footerText string, page pageSize) {
	pdf.SetTextColor(100, 100, 100)
	// Measure the rendered width instead of counting bytes; len() would treat a
	// CJK or accented footer as several times wider than it is.
	textW := measuredWidth(pdf, footerText)
	pdf.SetX(max((page.WidthPt-textW)/2, 0))
	pdf.SetY(page.HeightPt - 15)
	_ = pdf.Cell(nil, footerText)
	pdf.SetTextColor(0, 0, 0)
}

func renderNativePDFPlaceholderOverrides(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	for _, ph := range slide.PlaceholderOverrides {
		if ph.Text == "" || ph.Override == nil {
			continue
		}
		if ph.Override.X == nil || ph.Override.Y == nil || ph.Override.CX == nil || ph.Override.CY == nil {
			continue
		}
		x := emuToPt(ph.Override.X.Emu())
		y := emuToPt(ph.Override.Y.Emu())
		pdf.SetTextColor(0, 0, 0)
		pdf.SetX(x)
		pdf.SetY(y)
		_ = pdf.Cell(nil, ph.Text)
	}
}

//nolint:mnd // Slide number placement and colors match the native PPT slide template.
func renderNativePDFSlideNumber(pdf *gopdf.GoPdf, index, total int, page pageSize) {
	pdf.SetTextColor(150, 150, 150)
	slideNum := fmt.Sprintf("%d / %d", index, total)
	pdf.SetX(page.WidthPt - 60)
	pdf.SetY(page.HeightPt - 15)
	_ = pdf.Cell(nil, slideNum)
}
