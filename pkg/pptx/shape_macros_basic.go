package pptx

import "github.com/djinn-soul/gopptx/pkg/pptx/shapes"

func NewRectangle(x, y, w, h float64) Shape { return shapes.NewRectangle(x, y, w, h) }

func NewEllipse(x, y, w, h float64) Shape { return shapes.NewEllipse(x, y, w, h) }

func NewTextBox(text string, x, y, w, h float64) Shape { return shapes.NewTextBox(text, x, y, w, h) }

func NewRoundedRectangle(x, y, w, h float64) Shape { return shapes.NewRoundedRectangle(x, y, w, h) }

func NewTriangle(x, y, w, h float64) Shape { return shapes.NewTriangle(x, y, w, h) }

func NewRightTriangle(x, y, w, h float64) Shape { return shapes.NewRightTriangle(x, y, w, h) }

func NewDiamond(x, y, w, h float64) Shape { return shapes.NewDiamond(x, y, w, h) }

func NewPentagon(x, y, w, h float64) Shape { return shapes.NewPentagon(x, y, w, h) }

func NewHexagon(x, y, w, h float64) Shape { return shapes.NewHexagon(x, y, w, h) }

func NewParallelogram(x, y, w, h float64) Shape { return shapes.NewParallelogram(x, y, w, h) }

func NewFlowChartProcess(x, y, w, h float64) Shape { return shapes.NewFlowChartProcess(x, y, w, h) }

func NewFlowChartDecision(x, y, w, h float64) Shape { return shapes.NewFlowChartDecision(x, y, w, h) }

func NewFlowChartTerminator(x, y, w, h float64) Shape {
	return shapes.NewFlowChartTerminator(x, y, w, h)
}

func NewRightArrow(x, y, w, h float64) Shape { return shapes.NewRightArrow(x, y, w, h) }

func NewLeftArrow(x, y, w, h float64) Shape { return shapes.NewLeftArrow(x, y, w, h) }

func NewUpArrow(x, y, w, h float64) Shape { return shapes.NewUpArrow(x, y, w, h) }

func NewDownArrow(x, y, w, h float64) Shape { return shapes.NewDownArrow(x, y, w, h) }

func NewLeftRightArrow(x, y, w, h float64) Shape { return shapes.NewLeftRightArrow(x, y, w, h) }

func NewUpDownArrow(x, y, w, h float64) Shape { return shapes.NewUpDownArrow(x, y, w, h) }

func NewQuadArrow(x, y, w, h float64) Shape { return shapes.NewQuadArrow(x, y, w, h) }

func NewBentArrow(x, y, w, h float64) Shape { return shapes.NewBentArrow(x, y, w, h) }

func NewUturnArrow(x, y, w, h float64) Shape { return shapes.NewUturnArrow(x, y, w, h) }

func NewCircularArrow(x, y, w, h float64) Shape { return shapes.NewCircularArrow(x, y, w, h) }

func NewChevron(x, y, w, h float64) Shape { return shapes.NewChevron(x, y, w, h) }

func NewWedgeRectCallout(x, y, w, h float64) Shape { return shapes.NewWedgeRectCallout(x, y, w, h) }

func NewWedgeRRectCallout(x, y, w, h float64) Shape { return shapes.NewWedgeRRectCallout(x, y, w, h) }

func NewWedgeEllipseCallout(x, y, w, h float64) Shape {
	return shapes.NewWedgeEllipseCallout(x, y, w, h)
}

func NewCloudCallout(x, y, w, h float64) Shape { return shapes.NewCloudCallout(x, y, w, h) }

func NewCloud(x, y, w, h float64) Shape { return shapes.NewCloud(x, y, w, h) }

func NewCircle(x, y, diameter float64) Shape { return shapes.NewCircle(x, y, diameter) }

func NewStar(x, y, size float64) Shape { return shapes.NewStar(x, y, size) }

func NewHeart(x, y, size float64) Shape { return shapes.NewHeart(x, y, size) }

func NewStar4(x, y, size float64) Shape { return shapes.NewStar4(x, y, size) }

func NewStar6(x, y, size float64) Shape { return shapes.NewStar6(x, y, size) }

func NewStar7(x, y, size float64) Shape { return shapes.NewStar7(x, y, size) }

func NewStar8(x, y, size float64) Shape { return shapes.NewStar8(x, y, size) }

func NewStar10(x, y, size float64) Shape { return shapes.NewStar10(x, y, size) }

func NewStar12(x, y, size float64) Shape { return shapes.NewStar12(x, y, size) }

func NewStar16(x, y, size float64) Shape { return shapes.NewStar16(x, y, size) }

func NewStar24(x, y, size float64) Shape { return shapes.NewStar24(x, y, size) }

func NewStar32(x, y, size float64) Shape { return shapes.NewStar32(x, y, size) }

func NewRibbon(x, y, w, h float64) Shape { return shapes.NewRibbon(x, y, w, h) }

func NewWave(x, y, w, h float64) Shape { return shapes.NewWave(x, y, w, h) }

func NewSeal(x, y, size float64) Shape { return shapes.NewSeal(x, y, size) }

func NewActionButtonHome(x, y, size float64) Shape { return shapes.NewActionButtonHome(x, y, size) }

func NewActionButtonHelp(x, y, size float64) Shape { return shapes.NewActionButtonHelp(x, y, size) }

func NewActionButtonInformation(x, y, size float64) Shape {
	return shapes.NewActionButtonInformation(x, y, size)
}

func NewActionButtonBack(x, y, size float64) Shape { return shapes.NewActionButtonBack(x, y, size) }

func NewActionButtonForward(x, y, size float64) Shape {
	return shapes.NewActionButtonForward(x, y, size)
}

func NewFlowChartDocument(x, y, w, h float64) Shape { return shapes.NewFlowChartDocument(x, y, w, h) }

func NewFlowChartData(x, y, w, h float64) Shape { return shapes.NewFlowChartData(x, y, w, h) }

func NewOctagon(x, y, w, h float64) Shape { return shapes.NewOctagon(x, y, w, h) }

func NewTrapezoid(x, y, w, h float64) Shape { return shapes.NewTrapezoid(x, y, w, h) }

func NewCube(x, y, w, h float64) Shape { return shapes.NewCube(x, y, w, h) }
