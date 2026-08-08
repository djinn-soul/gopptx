package elements

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

func TestCalculateShapeIDs_IncludesChartSlot(t *testing.T) {
	withChart := NewSlide("Chart+Shape")
	withChart.Chart = &charts.BarChart{}
	withChart = withChart.AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 100, 100))

	withoutChart := NewSlide("ShapeOnly")
	withoutChart = withoutChart.AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 100, 100))

	idsWithChart := CalculateShapeIDs(withChart)
	idsWithoutChart := CalculateShapeIDs(withoutChart)

	if len(idsWithChart) != 1 || len(idsWithoutChart) != 1 {
		t.Fatalf("expected one shape id in each case, got withChart=%v withoutChart=%v", idsWithChart, idsWithoutChart)
	}
	if idsWithChart[0] != idsWithoutChart[0]+1 {
		t.Fatalf(
			"expected chart to reserve one shape slot: withChart=%v withoutChart=%v",
			idsWithChart,
			idsWithoutChart,
		)
	}
}

func TestPlaceholderMethods_AcceptTypedSignatures(t *testing.T) {
	slide := NewSlide("Placeholder API")
	img := shapes.NewImage("test.png", 1, 2, 3, 4)
	table := tables.NewTable([]styling.Length{styling.Emu(1000)})
	chart := charts.NewBarChart([]string{"A"}, []float64{1})

	slide = slide.WithPlaceholderText(1, "body text")
	slide = slide.WithPlaceholderTextAs(2, "title", "typed text")
	slide = slide.WithPlaceholderImage(3, img)
	slide = slide.WithPlaceholderImageAs(4, "pic", img)
	slide = slide.WithPlaceholderTable(5, table)
	slide = slide.WithPlaceholderTableAs(6, "tbl", table)
	slide = slide.WithPlaceholderChart(7, chart)
	slide = slide.WithPlaceholderChartAs(8, "chart", chart)

	if got := len(slide.PlaceholderOverrides); got != 8 {
		t.Fatalf("expected 8 placeholder overrides, got %d", got)
	}
	if slide.PlaceholderOverrides[0].Type != "body" || slide.PlaceholderOverrides[0].Text != "body text" {
		t.Fatalf("default text signature not handled correctly: %#v", slide.PlaceholderOverrides[0])
	}
	if slide.PlaceholderOverrides[1].Type != "title" || slide.PlaceholderOverrides[1].Text != "typed text" {
		t.Fatalf("text-as signature not handled correctly: %#v", slide.PlaceholderOverrides[1])
	}
	if slide.PlaceholderOverrides[2].Image == nil || slide.PlaceholderOverrides[2].Type != "pic" {
		t.Fatalf("default image signature not handled correctly: %#v", slide.PlaceholderOverrides[2])
	}
	if slide.PlaceholderOverrides[3].Image == nil || slide.PlaceholderOverrides[3].Type != "pic" {
		t.Fatalf("image-as signature not handled correctly: %#v", slide.PlaceholderOverrides[3])
	}
	if slide.PlaceholderOverrides[4].Table == nil || slide.PlaceholderOverrides[4].Type != "body" {
		t.Fatalf("default table signature not handled correctly: %#v", slide.PlaceholderOverrides[4])
	}
	if slide.PlaceholderOverrides[5].Table == nil || slide.PlaceholderOverrides[5].Type != "tbl" {
		t.Fatalf("table-as signature not handled correctly: %#v", slide.PlaceholderOverrides[5])
	}
	if slide.PlaceholderOverrides[6].Chart == nil || slide.PlaceholderOverrides[6].Type != "body" {
		t.Fatalf("default chart signature not handled correctly: %#v", slide.PlaceholderOverrides[6])
	}
	if slide.PlaceholderOverrides[7].Chart == nil || slide.PlaceholderOverrides[7].Type != "chart" {
		t.Fatalf("chart-as signature not handled correctly: %#v", slide.PlaceholderOverrides[7])
	}
}

func TestTransitionOptions_RejectUnsupportedDirection(t *testing.T) {
	if err := (transitions.TransitionOptions{Type: transitions.TransitionRandomBars, Direction: transitions.TransitionDirRight}).Validate(); err == nil {
		t.Fatalf("expected direction validation error for randomBar transition")
	}

	if err := (transitions.TransitionOptions{Type: transitions.TransitionClock, Direction: transitions.TransitionDirRight}).Validate(); err == nil {
		t.Fatalf("expected direction validation error for wheel transition")
	}

	if err := (transitions.TransitionOptions{Type: transitions.TransitionPush, Direction: transitions.TransitionDirRight}).Validate(); err != nil {
		t.Fatalf("expected push/right to be valid, got: %v", err)
	}
}

// TestCalculateShapeIDs_CountsOverflowTablesAndCharts guards the ID arithmetic
// against the overflow collections. The renderer emits a shape-tree entry for
// every table in Tables and every chart in Charts, so each one has to push the
// shapes that follow it along by one; counting only Table and the primary chart
// left an animation pointing at a table, a chart, or nothing at all.
func TestCalculateShapeIDs_CountsOverflowTablesAndCharts(t *testing.T) {
	base := NewSlide("ShapeOnly").
		AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 100, 100))

	withTables := NewSlide("Tables").
		AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 100, 100))
	withTables.Tables = []tables.Table{
		tables.NewTable([]styling.Length{styling.Inches(1)}).AddRow([]string{"a"}),
		tables.NewTable([]styling.Length{styling.Inches(1)}).AddRow([]string{"b"}),
	}

	withCharts := NewSlide("Charts").
		AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 100, 100))
	withCharts.Charts = []charts.ChartDefinition{
		charts.BarChart{Categories: []string{"a"}, Values: []float64{1}},
	}

	baseIDs := CalculateShapeIDs(base)
	tableIDs := CalculateShapeIDs(withTables)
	chartIDs := CalculateShapeIDs(withCharts)

	if len(baseIDs) != 1 || len(tableIDs) != 1 || len(chartIDs) != 1 {
		t.Fatalf("expected one shape id each: base=%v tables=%v charts=%v", baseIDs, tableIDs, chartIDs)
	}
	if tableIDs[0] != baseIDs[0]+2 {
		t.Errorf("two overflow tables moved the shape id to %d, want %d", tableIDs[0], baseIDs[0]+2)
	}
	if chartIDs[0] != baseIDs[0]+1 {
		t.Errorf("one overflow chart moved the shape id to %d, want %d", chartIDs[0], baseIDs[0]+1)
	}
}
