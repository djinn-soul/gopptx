//nolint:mnd // Shape text box fallback math uses fixed small offsets.
package export

import (
	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

const (
	defaultShapeTextPaddingPt = 4.0
	shapeTextMinBoxHeightPt   = 10.0
)

func renderPDFShapeText(pdf *gopdf.GoPdf, s shapes.Shape, x, y, w, h float64) {
	boxX, boxY, boxW, boxH, anchor := shapeTextBox(s, x, y, w, h)
	if boxW <= 2 || boxH <= 2 {
		return
	}
	fontSize := fitPDFTextToBox(
		pdf,
		s.Text,
		defaultFontSize,
		minTextAutoFitSize,
		false,
		false,
		boxW,
		boxH,
	)
	setPDFTextFont(pdf, fontSize, false, false)
	lines := wrapPDFText(pdf, s.Text, boxW)
	lineH := pdfLineHeight(fontSize)
	textBlockH := lineH * float64(len(lines))
	startY := shapeTextStartY(anchor, boxY, boxH, textBlockH)

	pdf.SetTextColor(0, 0, 0)
	yPos := startY
	for _, line := range lines {
		if yPos+lineH > boxY+boxH+0.5 {
			break
		}
		pdf.SetX(boxX)
		pdf.SetY(yPos)
		_ = pdf.Cell(nil, line)
		yPos += lineH
	}
	setPDFTextFont(pdf, defaultFontSize, false, false)
}

func shapeTextBox(
	s shapes.Shape,
	x, y, w, h float64,
) (float64, float64, float64, float64, shapes.TextFrameAnchor) {
	left := defaultShapeTextPaddingPt
	right := defaultShapeTextPaddingPt
	top := defaultShapeTextPaddingPt
	bottom := defaultShapeTextPaddingPt
	anchor := shapes.TextAnchorMiddle
	if s.TextFrame != nil {
		left = emuToPt(s.TextFrame.MarginLeft.Emu())
		right = emuToPt(s.TextFrame.MarginRight.Emu())
		top = emuToPt(s.TextFrame.MarginTop.Emu())
		bottom = emuToPt(s.TextFrame.MarginBottom.Emu())
		if s.TextFrame.Anchor != "" {
			anchor = s.TextFrame.Anchor
		}
	}
	boxX := x + left
	boxY := y + top
	boxW := w - left - right
	boxH := h - top - bottom
	if boxW <= 0 {
		boxW = w - 2
		boxX = x + 1
	}
	if boxH <= 0 {
		boxH = shapeTextMinBoxHeightPt
		boxY = y + 1
	}
	return boxX, boxY, boxW, boxH, anchor
}

func shapeTextStartY(anchor shapes.TextFrameAnchor, boxY, boxH, textBlockH float64) float64 {
	switch anchor {
	case shapes.TextAnchorTop:
		return boxY
	case shapes.TextAnchorBottom:
		return boxY + max(boxH-textBlockH, 0)
	default:
		return boxY + max((boxH-textBlockH)/2, 0)
	}
}
