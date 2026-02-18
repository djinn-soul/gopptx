package pptx

import (
	"path/filepath"
	"testing"
)

func TestShapeShowcasePublicAPI(t *testing.T) {
	builder := NewPresentationBuilder("Shape Showcase Public API")
	builder.WithSlideSize(SlideSize16x9())

	groups := []struct {
		name   string
		shapes []Shape
	}{
		{
			"Basic Shapes",
			[]Shape{
				NewRectangle(0.5, 1.5, 1.5, 1.2).WithText("Rectangle"),
				NewRoundedRectangle(2.2, 1.5, 1.5, 1.2).WithText("RoundedRect"),
				NewEllipse(3.9, 1.5, 1.5, 1.2).WithText("Ellipse"),
				NewTriangle(5.6, 1.5, 1.5, 1.2).WithText("Triangle"),
				NewRightTriangle(7.3, 1.5, 1.5, 1.2).WithText("RightTriangle"),
				NewDiamond(9.0, 1.5, 1.5, 1.2).WithText("Diamond"),
				NewPentagon(0.5, 3.1, 1.5, 1.2).WithText("Pentagon"),
				NewHexagon(2.2, 3.1, 1.5, 1.2).WithText("Hexagon"),
				NewOctagon(3.9, 3.1, 1.5, 1.2).WithText("Octagon"),
				NewParallelogram(5.6, 3.1, 1.5, 1.2).WithText("Parallelogram"),
				NewTrapezoid(7.3, 3.1, 1.5, 1.2).WithText("Trapezoid"),
				NewCube(9.0, 3.1, 1.5, 1.2).WithText("Cube"),
				NewCircle(0.5, 4.7, 1.2).WithText("Circle"),
				NewCloud(2.2, 4.7, 1.5, 1.2).WithText("Cloud"),
				NewHeart(3.9, 4.7, 1.2).WithText("Heart"),
			},
		},
		{
			"Arrows",
			[]Shape{
				NewRightArrow(0.5, 1.5, 1.5, 1.2).WithText("Right"),
				NewLeftArrow(2.2, 1.5, 1.5, 1.2).WithText("Left"),
				NewUpArrow(3.9, 1.5, 1.5, 1.2).WithText("Up"),
				NewDownArrow(5.6, 1.5, 1.5, 1.2).WithText("Down"),
				NewLeftRightArrow(7.3, 1.5, 1.5, 1.2).WithText("LeftRight"),
				NewUpDownArrow(9.0, 1.5, 1.5, 1.2).WithText("UpDown"),
				NewQuadArrow(0.5, 3.1, 1.5, 1.2).WithText("Quad"),
				NewBentArrow(2.2, 3.1, 1.5, 1.2).WithText("Bent"),
				NewUturnArrow(3.9, 3.1, 1.5, 1.2).WithText("Uturn"),
				NewCircularArrow(5.6, 3.1, 1.5, 1.2).WithText("Circular"),
				NewChevron(7.3, 3.1, 1.5, 1.2).WithText("Chevron"),
			},
		},
		{
			"Callouts",
			[]Shape{
				NewWedgeRectCallout(0.5, 1.5, 1.5, 1.2).WithText("Rect Callout"),
				NewWedgeEllipseCallout(2.2, 1.5, 1.5, 1.2).WithText("Ellipse Callout"),
				NewCloudCallout(3.9, 1.5, 1.5, 1.2).WithText("Cloud Callout"),
			},
		},
		{
			"Flow Chart",
			[]Shape{
				NewFlowChartProcess(0.5, 1.5, 1.5, 1.2).WithText("Process"),
				NewFlowChartDecision(2.2, 1.5, 1.5, 1.2).WithText("Decision"),
				NewFlowChartTerminator(3.9, 1.5, 1.5, 1.2).WithText("Terminator"),
				NewFlowChartDocument(5.6, 1.5, 1.5, 1.2).WithText("Document"),
				NewFlowChartData(7.3, 1.5, 1.5, 1.2).WithText("Data"),
				NewFlowChartPredefinedProcess(9.0, 1.5, 1.5, 1.2).WithText("Predefined"),
				NewFlowChartInternalStorage(0.5, 3.1, 1.5, 1.2).WithText("Internal Storage"),
				NewFlowChartManualInput(2.2, 3.1, 1.5, 1.2).WithText("Manual Input"),
				NewFlowChartManualOperation(3.9, 3.1, 1.5, 1.2).WithText("Manual Operation"),
				NewFlowChartConnector(5.6, 3.1, 1.5, 1.2).WithText("Connector"),
				NewFlowChartOffpageConnector(7.3, 3.1, 1.5, 1.2).WithText("Offpage"),
			},
		},
		{
			"Stars & Banners",
			[]Shape{
				NewStar4(0.5, 1.5, 1.2).WithText("Star 4"),
				NewStar(2.2, 1.5, 1.2).WithText("Star 5"),
				NewStar6(3.9, 1.5, 1.2).WithText("Star 6"),
				NewStar8(5.6, 1.5, 1.2).WithText("Star 8"),
				NewStar12(7.3, 1.5, 1.2).WithText("Star 12"),
				NewStar16(9.0, 1.5, 1.2).WithText("Star 16"),
				NewStar24(0.5, 3.1, 1.2).WithText("Star 24"),
				NewStar32(2.2, 3.1, 1.2).WithText("Star 32"),
				NewRibbon(3.9, 3.1, 1.5, 1.2).WithText("Ribbon"),
				NewWave(5.6, 3.1, 1.5, 1.2).WithText("Wave"),
				NewSeal(7.3, 3.1, 1.2).WithText("Seal"),
			},
		},
		{
			"Action Buttons",
			[]Shape{
				NewActionButtonHome(0.5, 1.5, 1.0).WithText("Home"),
				NewActionButtonHelp(2.0, 1.5, 1.0).WithText("Help"),
				NewActionButtonInformation(3.5, 1.5, 1.0).WithText("Info"),
				NewActionButtonBack(5.0, 1.5, 1.0).WithText("Back"),
				NewActionButtonForward(6.5, 1.5, 1.0).WithText("Forward"),
			},
		},
	}

	for _, group := range groups {
		builder.AddShapesSlide(group.name, group.shapes...)
	}

	outPath := filepath.Join(t.TempDir(), "showcase_shapes.pptx")
	err := builder.WriteToFile(outPath)
	if err != nil {
		t.Fatalf("Failed to write showcase: %v", err)
	}
}
