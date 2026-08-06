package presentation

import (
	"testing"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// A slide whose only chart is a 3-D pie used to produce no chart part at all:
// the field was settable and counted, but the generator never looked at it.
func TestBuildChartPartsEmitsPie3D(t *testing.T) {
	pie3d := charts.NewPie3DChart([]string{"A", "B"}, []float64{1, 2})
	slide := elements.NewSlide("3D pie")
	slide.Pie3D = &pie3d

	parts := BuildChartParts([]elements.SlideContent{slide})
	if len(parts) != 1 {
		t.Fatalf("chart parts = %d, want 1", len(parts))
	}
	if parts[0].spec.Kind != pptxxml.ChartKindThreeDPie {
		t.Fatalf("chart kind = %q, want %q", parts[0].spec.Kind, pptxxml.ChartKindThreeDPie)
	}
	if slideChartCount(slide) != 1 {
		t.Fatalf("slideChartCount = %d for a slide carrying Pie3D, want 1", slideChartCount(slide))
	}
}
