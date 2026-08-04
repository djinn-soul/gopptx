package export

import (
	"errors"
	"fmt"
	"sort"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
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

func renderNativePDFSlideText(pdf *gopdf.GoPdf, slide elements.SlideContent, page pageSize) {
	if slide.Title != "" {
		renderPDFTitle(pdf, slide, page)
	}
	if len(slide.Bullets) > 0 {
		renderPDFBullets(pdf, slide, page)
	}
}

// renderNativePDFSlidePictureLayer draws the slide's pictures and shapes in
// shape-tree order, then its connectors.
//
// Pictures and shapes share one z-order in PowerPoint, so they cannot be drawn
// as separate layers. Elements read from a PPTX carry their tree position in
// ZOrder; those built in memory leave it at zero and keep their slice order,
// which the stable sort preserves.
func renderNativePDFSlidePictureLayer(pdf *gopdf.GoPdf, slide elements.SlideContent) error {
	type paintable struct {
		zOrder int
		image  *shapes.Image
		shape  *shapes.Shape
	}

	items := make([]paintable, 0, len(slide.Images)+len(slide.Shapes))
	for i := range slide.Images {
		items = append(items, paintable{zOrder: slide.Images[i].ZOrder, image: &slide.Images[i]})
	}
	for i := range slide.Shapes {
		items = append(items, paintable{zOrder: slide.Shapes[i].ZOrder, shape: &slide.Shapes[i]})
	}
	sort.SliceStable(items, func(a, b int) bool { return items[a].zOrder < items[b].zOrder })

	var errs []error
	imageNumber := 0
	for _, item := range items {
		if item.image != nil {
			imageNumber++
			if err := renderPDFImageWithEffects(pdf, *item.image); err != nil {
				errs = append(errs, fmt.Errorf("image %d: %w", imageNumber, err))
			}
			continue
		}
		renderPDFShape(pdf, *item.shape)
	}

	for _, connector := range slide.Connectors {
		renderPDFConnector(pdf, connector)
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
