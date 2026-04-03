package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

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
	if s.Line == nil || s.Line.Width <= 0 {
		return false
	}
	strokeWidth := emuToPt(int64(s.Line.Width))
	if strokeWidth < minStrokeWidth {
		strokeWidth = minStrokeWidth
	}
	pdf.SetLineWidth(strokeWidth)
	pdf.SetStrokeColor(hexToRGB(s.Line.Color))
	return true
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
		pdf.Polygon(
			[]gopdf.Point{{X: x + w/2, Y: y}, {X: x, Y: y + h}, {X: x + w, Y: y + h}},
			style,
		)
	case shapes.ShapeTypeRightArrow:
		pdf.Polygon(rightArrowPoints(x, y, w, h), style)
	case shapes.ShapeTypeLeftArrow:
		pdf.Polygon(leftArrowPoints(x, y, w, h), style)
	default:
		pdf.RectFromUpperLeftWithStyle(x, y, w, h, style)
	}
}
