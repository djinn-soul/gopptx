package export

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestThreeDPieChartDetectionAndReconstruction(t *testing.T) {
	const raw = `<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart">` +
		`<c:chart><c:view3D/><c:plotArea><c:pie3DChart/></c:plotArea></c:chart></c:chartSpace>`
	if got := detectChartKind(raw); got != chartKindPie3D {
		t.Fatalf("detectChartKind() = %q, want %q", got, chartKindPie3D)
	}

	slide := elements.SlideContent{}
	applyParsedChart(&slide, parsedChart{
		Kind:  chartKindPie3D,
		Title: "Traffic share",
		X:     10,
		Y:     20,
		CX:    30,
		CY:    40,
		Series: []parsedChartSeries{{
			Name:       "Traffic",
			Categories: []string{"Desktop", "Mobile"},
			Values:     []float64{48, 52},
		}},
	})
	if slide.Pie3D == nil {
		t.Fatal("3D pie chart was not reconstructed into SlideContent.Pie3D")
	}
	if slide.Pie != nil {
		t.Fatal("3D pie chart was incorrectly reconstructed as a 2D pie")
	}
	if got := slide.Pie3D.ToChartSpec().Kind; got != chartKindPie3D {
		t.Fatalf("reconstructed chart kind = %q, want %q", got, chartKindPie3D)
	}
}
