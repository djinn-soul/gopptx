//nolint:mnd // Native PDF rendering relies on slide/layout geometry constants in points.
package export

import (
	"errors"
	"fmt"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// PDF coordinate constants.
// PPTX slides default to 9144000×6858000 EMU (10×7.5 inches).
// PDF uses 72 points/inch, so the slide page is 720×540 points.
const (
	ptPerInch     = 72.0
	slideWidthPt  = 720.0 // 10 inches, the 4:3 default
	slideHeightPt = 540.0 // 7.5 inches

	// defaultFontSize is the size PowerPoint gives a run that states none. It is
	// 18pt: measured off PowerPoint's own render of a text box whose a:rPr
	// carries no sz, which came out at 17.95pt. It was 14pt here, which made
	// every unsized run render a fifth too small.
	defaultFontSize = 18

	// Layout constants.
	//
	// PowerPoint's roundRect preset rounds by adj/100000 of the shorter side,
	// and its default adj is 16667. 0.1 drew a noticeably tighter corner.
	defaultRadiusFactor = 0.16667
	minStrokeWidth      = 0.5

	// nearZeroEpsilon is used for floating-point comparisons where a value
	// is considered effectively zero (e.g. rotation angle, cursor boundary).
	nearZeroEpsilon = 0.01

	// Default text-frame insets from the OOXML a:bodyPr defaults: lIns and rIns
	// are 91440 EMU (0.1in = 7.2pt), tIns and bIns are 45720 EMU (0.05in =
	// 3.6pt). PowerPoint applies these whenever a shape omits them, so text
	// starts 7.2pt inside its box rather than flush against the edge.
	defaultTextInsetLRPt = 7.2
	defaultTextInsetTBPt = 3.6
)

func emuToPt(emu int64) float64 {
	return (float64(emu) / emuPerInch) * ptPerInch
}

// pageSize is the PDF page geometry in points. It mirrors the deck's own
// <p:sldSz>, so a 16:9 deck renders onto a 960x540pt page rather than being
// cropped into the 4:3 default.
type pageSize struct {
	WidthPt  float64
	HeightPt float64
}

// defaultPageSize is the 4:3 page used when the deck's size is unknown.
func defaultPageSize() pageSize {
	return pageSize{WidthPt: slideWidthPt, HeightPt: slideHeightPt}
}

// pageSizeFromEMU converts a slide size in EMUs to a PDF page size, falling
// back to 4:3 for missing or nonsensical dimensions.
func pageSizeFromEMU(widthEMU, heightEMU int64) pageSize {
	if widthEMU <= 0 || heightEMU <= 0 {
		return defaultPageSize()
	}
	return pageSize{WidthPt: emuToPt(widthEMU), HeightPt: emuToPt(heightEMU)}
}

// optionsPageSize derives the page geometry for in-memory slides, which carry
// no size of their own, from PDFOptions.
func optionsPageSize(opts PDFOptions) pageSize {
	return pageSizeFromEMU(opts.SlideSize.Width, opts.SlideSize.Height)
}

// pdfViaNative renders slides directly to PDF using gopdf drawing primitives.
func pdfViaNative(
	_ string,
	slides []elements.SlideContent,
	orders []slidePaintOrder,
	outputPath string,
	opts PDFOptions,
	page pageSize,
) error {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{W: page.WidthPt, H: page.HeightPt},
	})
	// The document's font registry lives as long as the document does; dropping
	// it here keeps a long-lived process from holding one per export.
	defer releaseDocumentFonts(pdf)
	if err := configureNativePDFFont(pdf, opts); err != nil {
		return err
	}

	totalVisible := 0
	for _, slide := range slides {
		if !slide.Hidden {
			totalVisible++
		}
	}
	visibleIndex := 0
	var renderErrs []error
	for i, slide := range slides {
		if slide.Hidden {
			continue
		}
		visibleIndex++
		order := newSlidePaintOrder()
		if i < len(orders) {
			order = orders[i]
		}
		if err := renderNativePDFSlide(pdf, slide, order, visibleIndex, totalVisible, page); err != nil {
			renderErrs = append(renderErrs, err)
		}
	}

	if err := pdf.WritePdf(outputPath); err != nil {
		return err
	}
	// The PDF was written, but report anything that silently failed to render
	// so callers do not ship a deck with missing pictures.
	return errors.Join(renderErrs...)
}

// renderNativePDFSlide paints one slide: the background, then every element of
// the shape tree in the order the tree lists them, then the slide furniture.
//
// Everything between the background and the furniture goes through one paint
// list. PowerPoint has a single back-to-front order over placeholders,
// pictures, shapes, connectors, tables, charts and diagrams alike, so painting
// them as fixed layers put a table over a shape the deck had placed in front of
// it, and a chart over that table.
func renderNativePDFSlide(
	pdf *gopdf.GoPdf,
	slide elements.SlideContent,
	order slidePaintOrder,
	index, total int,
	page pageSize,
) error {
	pdf.AddPage()

	var errs []error
	if err := renderPDFBackground(pdf, slide.Background, page); err != nil {
		errs = append(errs, fmt.Errorf("slide %d background: %w", index, err))
	}
	for _, err := range buildSlidePaintList(pdf, slide, order, page).drawAll() {
		errs = append(errs, fmt.Errorf("slide %d: %w", index, err))
	}

	if slide.ShowSlideNumber {
		renderNativePDFSlideNumber(pdf, index, total, page)
	}
	if slide.FooterText != "" {
		renderNativePDFFooter(pdf, slide.FooterText, page)
	}
	if len(slide.PlaceholderOverrides) > 0 {
		renderNativePDFPlaceholderOverrides(pdf, slide)
	}
	return errors.Join(errs...)
}

func renderPDFTitle(pdf *gopdf.GoPdf, slide elements.SlideContent, page pageSize) {
	titleSize := slide.TitleSize
	if titleSize <= 0 {
		switch elements.NormalizeSlideLayout(slide.Layout) {
		case elements.SlideLayoutTitleAndContent:
			titleSize = 32
		default:
			titleSize = 44
		}
	}
	// Only apply the layout-based size cap when we used a layout default.
	// When TitleSize was explicitly read from the PPTX we honour it as-is.
	if slide.TitleSize <= 0 {
		titleMax := 44
		if elements.NormalizeSlideLayout(slide.Layout) == elements.SlideLayoutTitleAndContent {
			titleMax = 32
		}
		if titleSize > titleMax {
			titleSize = titleMax
		}
	}
	titleBoxX := 54.0
	titleBoxY := 44.0
	titleBoxW := page.WidthPt - 108
	titleBoxH := 72.0
	if b := slide.TitleBoundsEMU; b[2] > 0 || b[3] > 0 {
		titleBoxX = emuToPt(b[0])
		titleBoxY = emuToPt(b[1])
		titleBoxW = emuToPt(b[2])
		titleBoxH = emuToPt(b[3])
	} else if elements.NormalizeSlideLayout(slide.Layout) == elements.SlideLayoutCenteredTitle {
		// Default centered-title box: narrower, vertically centered per Office template defaults.
		// These values approximate the default ctrTitle placeholder from a standard Office theme
		// (x≈685800 EMU, y≈2130425 EMU, cx≈7772400 EMU, cy≈1470025 EMU on a 9144000×6858000 slide).
		titleBoxX = 54.0
		titleBoxY = 167.0
		titleBoxW = page.WidthPt - 108
		titleBoxH = 116.0
	}
	// The bounds above are the placeholder box; text sits inside its insets.
	titleBoxX += defaultTextInsetLRPt
	titleBoxY += defaultTextInsetTBPt
	titleBoxW -= 2 * defaultTextInsetLRPt
	titleBoxH -= 2 * defaultTextInsetTBPt
	titleSize = fitPDFTitleSize(
		pdf,
		slide.Title,
		titleSize,
		slide.TitleBold,
		slide.TitleItalic,
		titleBoxW,
		titleBoxH,
		slide.TitleFont,
	)
	setPDFTextFontWithHint(pdf, titleSize, slide.TitleBold, slide.TitleItalic, slide.TitleFont)
	if slide.TitleColor != "" {
		pdf.SetTextColor(hexToRGB(slide.TitleColor))
	} else {
		pdf.SetTextColor(0, 0, 0)
	}
	lines := wrapPDFTextWithMetrics(pdf, slide.Title, titleBoxW)
	// A title placeholder inherits the master's titleStyle, which spaces its
	// lines at 90%. That tightens the line box, and with it the first baseline.
	lineH := pdfLineHeight(titleSize) * titleDefaultLineSpacingFactor
	totalTextH := float64(len(lines)) * lineH
	yPos := titleBoxY + max(0, (titleBoxH-totalTextH)/2)
	for _, line := range lines {
		if yPos+lineH > titleBoxY+titleBoxH {
			break
		}
		pdf.SetX(alignedTextX(pdf, line, titleBoxX, titleBoxW, slide.TitleAlign))
		pdf.SetY(yPos + fontBaselineShiftInLineBox(pdf, slide.TitleFont, titleSize, lineH))
		_ = pdf.Cell(nil, line)
		yPos += lineH
	}
	setPDFTextFontWithHint(pdf, defaultFontSize, false, false, "")
}

func fitPDFTitleSize(
	pdf *gopdf.GoPdf,
	text string,
	initialSize int,
	bold bool,
	italic bool,
	maxWidth float64,
	maxHeight float64,
	fontHint string,
) int {
	size := max(14, min(initialSize, 44))
	for size > 14 {
		setPDFTextFontWithHint(pdf, size, bold, italic, fontHint)
		if fitPDFTextToBoxWithMetrics(
			pdf, text, size, 14, bold, italic, maxWidth, maxHeight, fontHint,
		) == size {
			return size
		}
		size--
	}
	return size
}

func alignedTextX(
	pdf *gopdf.GoPdf,
	text string,
	boxX float64,
	boxW float64,
	align string,
) float64 {
	textW := measuredWidth(pdf, text)
	switch elements.NormalizeTextAlign(align) {
	case elements.TextAlignCenter:
		return boxX + max((boxW-textW)/2, 0)
	case elements.TextAlignRight:
		return boxX + max(boxW-textW, 0)
	default:
		return boxX
	}
}

func renderPDFShape(pdf *gopdf.GoPdf, s shapes.Shape) {
	x, y, w, h := getShapeBounds(s)

	gradientRendered := renderPDFShapeGradient(pdf, s, x, y, w, h)
	setPDFShapeFill(pdf, s, gradientRendered)
	hasStroke := setPDFShapeStroke(pdf, s)
	hasFill := s.Fill != nil || ((s.GradientFill != nil && len(s.GradientFill.Stops) > 0) && !gradientRendered)
	style := drawStyle(hasFill, hasStroke)

	rotated := s.RotationDeg != nil && *s.RotationDeg != 0
	if rotated {
		pdf.Rotate(float64(*s.RotationDeg), x+w/2, y+h/2)
	}

	renderPDFShapeEffects(pdf, s, x, y, w, h, hasFill)
	softEdgesApplied := applyPDFShapeSoftEdges(pdf, s)
	if style != "" {
		drawPDFGeometry(pdf, s, x, y, w, h, style)
	}
	if softEdgesApplied {
		pdf.ClearTransparency()
	}

	// The outline is done; text decorations drawn below are lines too and must
	// not inherit the shape's dash.
	clearPDFLineDash(pdf)

	// Text stays inside the rotation: a shape's caption turns with the shape in
	// PowerPoint. Resetting before drawing it left rotated captions horizontal,
	// which pushed a quadrant chart's y-axis label out over the plot area.
	if s.Text != "" {
		renderPDFShapeText(pdf, s, x, y, w, h)
	}

	if rotated {
		pdf.RotateReset()
	}

	// Reset colors. The dash pattern is document-wide state in gopdf, so a
	// dashed shape would otherwise leave every later stroke dashed too.
	pdf.SetFillColor(255, 255, 255)
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(1)
	clearPDFLineDash(pdf)
}
