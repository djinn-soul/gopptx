//nolint:mnd // Native PDF rendering relies on slide/layout geometry constants in points.
package export

import (
	"errors"
	"fmt"
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// PDF coordinate constants.
// PPTX slides default to 9144000×6858000 EMU (10×7.5 inches).
// PDF uses 72 points/inch, so the slide page is 720×540 points.
const (
	ptPerInch       = 72.0
	slideWidthPt    = 720.0 // 10 inches
	slideHeightPt   = 540.0 // 7.5 inches
	defaultFontSize = 14

	// Layout constants.
	defaultRadiusFactor = 0.1
	minStrokeWidth      = 0.5
)

func emuToPt(emu int64) float64 {
	return (float64(emu) / emuPerInch) * ptPerInch
}

// pdfViaNative renders slides directly to PDF using gopdf drawing primitives.
func pdfViaNative(_ string, slides []elements.SlideContent, outputPath string, opts PDFOptions) error {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{W: slideWidthPt, H: slideHeightPt},
	})
	if err := configureNativePDFFont(pdf, opts); err != nil {
		return err
	}

	for i, slide := range slides {
		renderNativePDFSlide(pdf, slide, i+1, len(slides))
	}

	return pdf.WritePdf(outputPath)
}

func configureNativePDFFont(pdf *gopdf.GoPdf, opts PDFOptions) error {
	sansAlias := ""
	if tryNativePDFFonts(pdf, opts.NativeFontPaths, fontFamilySans) {
		sansAlias = fontFamilySans
	} else if tryNativePDFFonts(pdf, systemFontPathsForFamily(fontFamilySans), fontFamilySans) {
		sansAlias = fontFamilySans
	}
	if sansAlias == "" {
		return errors.New("no system TTF font found; install Arial or DejaVu Sans, or specify NativeFontPaths")
	}
	serifAlias := ""
	if tryNativePDFFonts(pdf, systemFontPathsForFamily(fontFamilySerif), fontFamilySerif) {
		serifAlias = fontFamilySerif
	}
	monoAlias := ""
	if tryNativePDFFonts(pdf, systemFontPathsForFamily(fontFamilyMono), fontFamilyMono) {
		monoAlias = fontFamilyMono
	}
	setPDFFontAliases(sansAlias, serifAlias, monoAlias)
	return nil
}

func tryNativePDFFonts(pdf *gopdf.GoPdf, fontPaths []string, alias string) bool {
	for _, path := range fontPaths {
		if err := pdf.AddTTFFont(alias, path); err != nil {
			continue
		}
		if err := pdf.SetFont(alias, "", defaultFontSize); err == nil {
			return true
		}
	}
	return false
}

func renderNativePDFSlide(pdf *gopdf.GoPdf, slide elements.SlideContent, index, total int) {
	pdf.AddPage()
	_ = renderPDFBackground(pdf, slide.Background)
	renderNativePDFSlideText(pdf, slide)
	renderNativePDFSlideShapes(pdf, slide)
	renderNativePDFSlideAssets(pdf, slide)
	renderNativePDFSlideNumber(pdf, index, total)
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
		_ = renderPDFImage(pdf, img)
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

func renderPDFTitle(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	titleSize := slide.TitleSize
	if titleSize <= 0 {
		titleSize = defaultFontSize
	}
	titleMax := 44
	if elements.NormalizeSlideLayout(slide.Layout) == elements.SlideLayoutTitleAndContent {
		titleMax = 32
	}
	if titleSize > titleMax {
		titleSize = titleMax
	}
	titleBoxX := 54.0
	titleBoxY := 44.0
	titleBoxW := slideWidthPt - 108
	titleBoxH := 72.0
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
	lines := wrapPDFTextWithMetrics(pdf, slide.Title, titleBoxW, slide.TitleFont)
	lineH := pdfLineHeight(titleSize)
	yPos := titleBoxY
	for _, line := range lines {
		if yPos+lineH > titleBoxY+titleBoxH {
			break
		}
		pdf.SetX(alignedTextX(pdf, line, titleBoxX, titleBoxW, slide.TitleAlign, slide.TitleFont))
		pdf.SetY(yPos + fontBaselineShift(slide.TitleFont, titleSize))
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
	fontHint string,
) float64 {
	textW := measuredWidthWithMetrics(pdf, text, fontHint)
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

	if rotated {
		pdf.RotateReset()
	}

	if s.Text != "" {
		renderPDFShapeText(pdf, s, x, y, w, h)
	}

	// Reset colors
	pdf.SetFillColor(255, 255, 255)
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(1)
}

func getShapeBounds(s shapes.Shape) (float64, float64, float64, float64) {
	return emuToPt(int64(s.X)), emuToPt(int64(s.Y)), emuToPt(int64(s.CX)), emuToPt(int64(s.CY))
}

func setPDFShapeFill(pdf *gopdf.GoPdf, s shapes.Shape, gradientRendered bool) {
	if s.Fill != nil && s.Fill.Color != "" {
		pdf.SetFillColor(hexToRGB(s.Fill.Color))
	} else if !gradientRendered && s.GradientFill != nil && len(s.GradientFill.Stops) > 0 {
		pdf.SetFillColor(hexToRGB(s.GradientFill.Stops[0].Color))
	}
}

func setPDFShapeStroke(pdf *gopdf.GoPdf, s shapes.Shape) bool {
	if s.Line != nil && s.Line.Width > 0 {
		strokeW := emuToPt(int64(s.Line.Width))
		if strokeW < minStrokeWidth {
			strokeW = minStrokeWidth
		}
		pdf.SetLineWidth(strokeW)
		pdf.SetStrokeColor(hexToRGB(s.Line.Color))
		return true
	}
	return false
}

func drawPDFGeometry(pdf *gopdf.GoPdf, s shapes.Shape, x, y, w, h float64, style string) {
	switch s.Type {
	case shapes.ShapeTypeRectangle:
		pdf.RectFromUpperLeftWithStyle(x, y, w, h, style)
	case shapes.ShapeTypeRoundedRectangle:
		radius := math.Min(w, h) * defaultRadiusFactor
		_ = pdf.Rectangle(x, y, x+w, y+h, style, radius, 0)
	case shapes.ShapeTypePie, shapes.ShapeTypePieWedge, shapes.ShapeTypeChord:
		drawPieShape(pdf, s, x, y, w, h, style)
	case shapes.ShapeTypeEllipse:
		pdf.Oval(x, y, x+w, y+h)
	case shapes.ShapeTypeTriangle:
		pdf.Polygon([]gopdf.Point{{X: x + w/2, Y: y}, {X: x, Y: y + h}, {X: x + w, Y: y + h}}, style)
	case shapes.ShapeTypeRightArrow:
		pdf.Polygon(rightArrowPoints(x, y, w, h), style)
	case shapes.ShapeTypeLeftArrow:
		pdf.Polygon(leftArrowPoints(x, y, w, h), style)
	default:
		pdf.RectFromUpperLeftWithStyle(x, y, w, h, style)
	}
}

// renderPDFImage embeds a raster image into the PDF at the given EMU position.
func renderPDFImage(pdf *gopdf.GoPdf, img shapes.Image) error {
	return renderPDFImageWithEffects(pdf, img)
}
