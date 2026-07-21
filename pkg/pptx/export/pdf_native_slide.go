package export

import (
	"errors"
	"fmt"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

//nolint:mnd // Footer placement and colors match the native PPT slide template.
func renderNativePDFFooter(pdf *gopdf.GoPdf, footerText string, page pageSize) {
	pdf.SetTextColor(100, 100, 100)
	// Measure the rendered width instead of counting bytes; len() would treat a
	// CJK or accented footer as several times wider than it is.
	textW := measuredWidthWithMetrics(pdf, footerText, "")
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

func renderNativePDFSlideText(pdf *gopdf.GoPdf, slide elements.SlideContent, page pageSize) {
	if slide.Title != "" {
		renderPDFTitle(pdf, slide, page)
	}
	if len(slide.Bullets) > 0 {
		renderPDFBullets(pdf, slide, page)
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

// renderNativePDFSlideImages draws every picture on the slide and reports the
// ones that could not be rendered, rather than dropping them silently.
func renderNativePDFSlideImages(pdf *gopdf.GoPdf, slide elements.SlideContent) error {
	var errs []error
	for i, img := range slide.Images {
		if err := renderPDFImageWithEffects(pdf, img); err != nil {
			errs = append(errs, fmt.Errorf("image %d: %w", i+1, err))
		}
	}
	return errors.Join(errs...)
}

func renderNativePDFSlideTable(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	if slide.Table != nil {
		renderPDFTable(pdf, *slide.Table)
	}
	for _, t := range slide.Tables {
		renderPDFTable(pdf, t)
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
