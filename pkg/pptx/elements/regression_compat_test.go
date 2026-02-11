package elements

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

func TestCalculateShapeIDs_IncludesChartSlot(t *testing.T) {
	withChart := NewSlide("Chart+Shape")
	withChart.Chart = &charts.BarChart{}
	withChart = withChart.AddShape(NewShape(ShapeTypeRectangle, 0, 0, 100, 100))

	withoutChart := NewSlide("ShapeOnly")
	withoutChart = withoutChart.AddShape(NewShape(ShapeTypeRectangle, 0, 0, 100, 100))

	idsWithChart := CalculateShapeIDs(withChart)
	idsWithoutChart := CalculateShapeIDs(withoutChart)

	if len(idsWithChart) != 1 || len(idsWithoutChart) != 1 {
		t.Fatalf("expected one shape id in each case, got withChart=%v withoutChart=%v", idsWithChart, idsWithoutChart)
	}
	if idsWithChart[0] != idsWithoutChart[0]+1 {
		t.Fatalf("expected chart to reserve one shape slot: withChart=%v withoutChart=%v", idsWithChart, idsWithoutChart)
	}
}

func TestPlaceholderMethods_AcceptLegacyAndTypedSignatures(t *testing.T) {
	slide := NewSlide("Placeholder API")
	img := NewImage("test.png", 1, 2, 3, 4)
	table := tables.NewTable([]int64{1000})
	chart := charts.NewBarChart([]string{"A"}, []float64{1})

	slide = slide.WithPlaceholderText(1, "legacy text")
	slide = slide.WithPlaceholderText(2, "title", "typed text")
	slide = slide.WithPlaceholderImage(3, img)
	slide = slide.WithPlaceholderImage(4, "pic", img)
	slide = slide.WithPlaceholderTable(5, table)
	slide = slide.WithPlaceholderTable(6, "tbl", table)
	slide = slide.WithPlaceholderChart(7, chart)
	slide = slide.WithPlaceholderChart(8, "chart", chart)

	if got := len(slide.PlaceholderOverrides); got != 8 {
		t.Fatalf("expected 8 placeholder overrides, got %d", got)
	}
	if slide.PlaceholderOverrides[0].Type != "" || slide.PlaceholderOverrides[0].Text != "legacy text" {
		t.Fatalf("legacy text signature not handled correctly: %#v", slide.PlaceholderOverrides[0])
	}
	if slide.PlaceholderOverrides[1].Type != "title" || slide.PlaceholderOverrides[1].Text != "typed text" {
		t.Fatalf("typed text signature not handled correctly: %#v", slide.PlaceholderOverrides[1])
	}
	if slide.PlaceholderOverrides[2].Image == nil || slide.PlaceholderOverrides[2].Type != "" {
		t.Fatalf("legacy image signature not handled correctly: %#v", slide.PlaceholderOverrides[2])
	}
	if slide.PlaceholderOverrides[3].Image == nil || slide.PlaceholderOverrides[3].Type != "pic" {
		t.Fatalf("typed image signature not handled correctly: %#v", slide.PlaceholderOverrides[3])
	}
	if slide.PlaceholderOverrides[4].Table == nil || slide.PlaceholderOverrides[4].Type != "" {
		t.Fatalf("legacy table signature not handled correctly: %#v", slide.PlaceholderOverrides[4])
	}
	if slide.PlaceholderOverrides[5].Table == nil || slide.PlaceholderOverrides[5].Type != "tbl" {
		t.Fatalf("typed table signature not handled correctly: %#v", slide.PlaceholderOverrides[5])
	}
	if slide.PlaceholderOverrides[6].Chart == nil || slide.PlaceholderOverrides[6].Type != "" {
		t.Fatalf("legacy chart signature not handled correctly: %#v", slide.PlaceholderOverrides[6])
	}
	if slide.PlaceholderOverrides[7].Chart == nil || slide.PlaceholderOverrides[7].Type != "chart" {
		t.Fatalf("typed chart signature not handled correctly: %#v", slide.PlaceholderOverrides[7])
	}
}

func TestTransitionOptions_RejectUnsupportedDirection(t *testing.T) {
	if err := (TransitionOptions{Type: TransitionRandomBars, Direction: TransitionDirRight}).Validate(); err == nil {
		t.Fatalf("expected direction validation error for randomBar transition")
	}

	if err := (TransitionOptions{Type: TransitionClock, Direction: TransitionDirRight}).Validate(); err == nil {
		t.Fatalf("expected direction validation error for wheel transition")
	}

	if err := (TransitionOptions{Type: TransitionPush, Direction: TransitionDirRight}).Validate(); err != nil {
		t.Fatalf("expected push/right to be valid, got: %v", err)
	}
}
