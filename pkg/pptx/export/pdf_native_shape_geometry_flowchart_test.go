package export

import (
	"math"
	"testing"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func assertPointsInsideBox(t *testing.T, name string, points []gopdf.Point, x, y, w, h float64) {
	t.Helper()
	if len(points) < 3 {
		t.Fatalf("%s: got %d points, want a closed outline", name, len(points))
	}
	const slack = 0.001
	for i, p := range points {
		if p.X < x-slack || p.X > x+w+slack || p.Y < y-slack || p.Y > y+h+slack {
			t.Errorf("%s: point %d at %.2f,%.2f is outside the %vx%v box", name, i, p.X, p.Y, w, h)
		}
	}
}

func TestFlowchartOutlinesStayInsideTheirBox(t *testing.T) {
	const x, y, w, h = 10.0, 20.0, 200.0, 100.0
	cases := map[string][]gopdf.Point{
		"diamond":       diamondPoints(x, y, w, h),
		"parallelogram": flowParallelogramPoints(x, y, w, h),
		"offpage":       offpageConnectorPoints(x, y, w, h),
		"card":          clippedCornerPoints(x, y, w, h),
		"document":      flowDocumentPoints(x, y, w, h),
		"delay":         flowDelayPoints(x, y, w, h),
		"cylinder":      cylinderPoints(x, y, w, h),
		"stored data":   flowStoredDataPoints(x, y, w, h),
		"cross":         crossPoints(x, y, w, h),
	}
	for name, points := range cases {
		assertPointsInsideBox(t, name, points, x, y, w, h)
	}
}

func TestCrossIsSymmetric(t *testing.T) {
	const x, y, w, h = 0.0, 0.0, 100.0, 100.0
	points := crossPoints(x, y, w, h)
	if len(points) != 12 {
		t.Fatalf("got %d points, want the 12 corners of a plus sign", len(points))
	}
	var sumX, sumY float64
	for _, p := range points {
		sumX += p.X
		sumY += p.Y
	}
	if math.Abs(sumX/float64(len(points))-w/2) > 0.001 {
		t.Errorf("centroid x = %.3f, want %v", sumX/float64(len(points)), w/2)
	}
	if math.Abs(sumY/float64(len(points))-h/2) > 0.001 {
		t.Errorf("centroid y = %.3f, want %v", sumY/float64(len(points)), h/2)
	}
}

func TestFlowchartPresetsAreRecognised(t *testing.T) {
	pdf := newTestPDF(t)

	recognised := []string{
		shapes.ShapeTypeFlowChartProcess,
		shapes.ShapeTypeFlowChartDecision,
		shapes.ShapeTypeFlowChartTerminator,
		shapes.ShapeTypeFlowChartInputOutput,
		shapes.ShapeTypeFlowChartPreparation,
		shapes.ShapeTypeFlowChartConnector,
		shapes.ShapeTypeFlowChartDocument,
		shapes.ShapeTypeFlowChartDelay,
		shapes.ShapeTypeFlowChartMagneticDisk,
		shapes.ShapeTypeFlowChartStoredData,
		shapes.ShapeTypeFlowChartOffpageConnector,
		shapes.ShapeTypeFlowChartCard,
		shapes.ShapeTypeFlowChartManualInput,
		shapes.ShapeTypeFlowChartManualOperation,
		shapes.ShapeTypeFlowChartMerge,
		shapes.ShapeTypeFlowChartExtract,
		shapes.ShapeTypeFlowChartCollate,
		shapes.ShapeTypeMathPlus,
		shapes.ShapeTypeMathMinus,
	}
	for _, shapeType := range recognised {
		if !drawPDFExtendedGeometry(pdf, flipState{unflippedShape: true}, shapeType, 0, 0, 100, 50, "D") {
			t.Errorf("%s fell through to the rectangle fallback", shapeType)
		}
	}

	if drawPDFExtendedGeometry(pdf, flipState{unflippedShape: true}, "notAPresetAtAll", 0, 0, 100, 50, "D") {
		t.Error("an unknown preset claimed to be drawn")
	}
}

func TestDecisionAndProcessDrawDifferentOutlines(t *testing.T) {
	// The whole point of the table: two presets that used to be the same box.
	diamond := diamondPoints(0, 0, 100, 100)
	if len(diamond) != 4 {
		t.Fatalf("got %d points for a decision, want 4", len(diamond))
	}
	if diamond[0].X != 50 || diamond[0].Y != 0 {
		t.Errorf("decision starts at %.0f,%.0f, want its top vertex at 50,0", diamond[0].X, diamond[0].Y)
	}
}
