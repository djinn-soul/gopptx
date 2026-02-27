package pptx

import (
	"testing"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestCompat_Charts(t *testing.T) {
	cats := []string{"A"}
	vals := []float64{1}
	_ = NewBarChart(cats, vals)
	_ = NewBarHorizontalChart(cats, vals)
	_ = NewBarStackedChart(cats, vals)
	_ = NewBarStacked100Chart(cats, vals)
	_ = NewLineChart(cats, vals)
	_ = NewLineMarkersChart(cats, vals)
	_ = NewLineStackedChart(cats, vals)
	_ = NewScatterChart([]float64{1}, []float64{2})
	_ = NewAreaChart(cats, vals)
	_ = NewAreaStackedChart(cats, vals)
	_ = NewAreaStacked100Chart(cats, vals)
	_ = NewPieChart(cats, vals)
	_ = NewDoughnutChart(cats, vals)
	_ = NewBubbleChart([]float64{1}, []float64{2}, []float64{3})
	_ = NewRadarChart(cats, vals)
	_ = NewRadarFilledChart(cats, vals)
	_ = NewStockHLCChart(cats, vals, vals, vals)
	_ = NewStockOHLCChart(cats, vals, vals, vals, vals)
	_ = NewComboChart(cats, nil, nil)
}

func TestCompat_Styling(t *testing.T) {
	_ = Inches(1)
	_ = InchesToEMU(1)
	_ = Centimeters(1)
	_ = CMToEMU(1)
	_ = Points(1)
	_ = PointsToEMU(1)
	_ = FontSize(12)
	_ = AllThemes()
}

func TestCompat_Text(t *testing.T) {
	_ = NewParagraph()
	_ = NewTextParagraph()
	_ = NewTextParagraphStyle()
	_ = DefaultParagraphStyle()
	_ = DefaultTextParagraphStyle()
	_ = NormalizeTextAlign("ctr")
	_ = NormalizeRuns(nil)
	_ = NormalizeTextRuns(nil)
}

func TestCompat_Shapes(t *testing.T) {
	_ = NewShapeLine("FF0000", Inches(1))
	_ = NewShapeGradientStop(0, "FF0000")
	_ = NewShapeGradientFill("linear", nil)
	_ = NewTextFrame()
	_ = NewTextBox("text", 0, 0, 1, 1)
	_ = NewWedgeRRectCallout(0, 0, 1, 1)
	_ = NewStar7(0, 0, 1)
	_ = NewStar10(0, 0, 1)
	_ = NewBadge("text", 0, 0, "FF0000")
	
	_ = NewRectangle(0, 0, 1, 1)
	_ = NewEllipse(0, 0, 1, 1)
}

func TestCompat_Connectors(t *testing.T) {
	_ = NewConnector("rect", Inches(0), Inches(0), Inches(1), Inches(1))
	_ = NewStraightConnector(Inches(0), Inches(0), Inches(1), Inches(1))
	_ = NewElbowConnector(Inches(0), Inches(0), Inches(1), Inches(1))
	_ = NewCurvedConnector(Inches(0), Inches(0), Inches(1), Inches(1))
	c := NewElbowConnector(Inches(0), Inches(0), Inches(1), Inches(1))
	_ = ConnectStartAuto(c, 1)
	_ = ConnectEndAuto(c, 2)
	_ = AutoReroute(c, nil)
}

func TestCompat_Slide(t *testing.T) {
	_ = NewNotesMaster()
	_ = NewMaster()
	_ = NewSolidBackground("FF0000")
	_ = NewGradientBackground(shapes.ShapeGradientFill{})
	_ = NewPictureBackground(shapes.Image{})
}
