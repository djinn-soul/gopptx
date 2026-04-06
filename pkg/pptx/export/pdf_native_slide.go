package export

import (
	"fmt"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func renderNativePDFFooter(pdf *gopdf.GoPdf, footerText string) {
	pdf.SetTextColor(100, 100, 100)
	pdf.SetX((slideWidthPt - float64(len(footerText))*4.5) / 2)
	pdf.SetY(slideHeightPt - 15)
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

func renderNativePDFSlideText(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	if slide.Title != "" {
		renderPDFTitle(pdf, slide)
	}
	if len(slide.Bullets) > 0 {
		renderPDFBullets(pdf, slide)
	}
}

func renderNativePDFSlideShapes(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	for _, shape := range slide.Shapes {
		renderPDFShape(pdf, shape)
	}
	for _, connector := range slide.Connectors {
		renderPDFConnector(pdf, connector)
	}
}

func renderNativePDFSlideAssets(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	for _, img := range slide.Images {
		_ = renderPDFImageWithEffects(pdf, img)
	}
	if slide.Table != nil {
		renderPDFTable(pdf, *slide.Table)
	}
}

func renderNativePDFSlideNumber(pdf *gopdf.GoPdf, index, total int) {
	pdf.SetTextColor(150, 150, 150)
	slideNum := fmt.Sprintf("%d / %d", index, total)
	pdf.SetX(slideWidthPt - 60)
	pdf.SetY(slideHeightPt - 15)
	_ = pdf.Cell(nil, slideNum)
}
