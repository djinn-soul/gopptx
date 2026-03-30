package export

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

func TestSlidesFromPPTX_PreservesConnectorsAndFrameAccessibility(t *testing.T) {
	left := shapes.NewShape(shapes.ShapeTypeRectangle, styling.Inches(1), styling.Inches(2), styling.Inches(1.5), styling.Inches(0.8)).
		WithText("Left")
	right := shapes.NewShape(shapes.ShapeTypeRectangle, styling.Inches(5), styling.Inches(2), styling.Inches(1.5), styling.Inches(0.8)).
		WithText("Right")
	connector := shapes.NewStraightConnector(styling.Inches(2.5), styling.Inches(2.4), styling.Inches(5), styling.Inches(2.4)).
		ConnectStart(1, shapes.ConnectionSiteRight).
		ConnectEnd(2, shapes.ConnectionSiteLeft).
		WithAltText("Connector Alt")
	table := tables.NewTable([]styling.Length{styling.Inches(1.5), styling.Inches(1.5)}).
		WithAltText("Table Alt").
		Position(styling.Inches(1), styling.Inches(4)).
		Size(styling.Inches(3), styling.Inches(1.2)).
		AddRow([]string{"H1", "H2"}).
		AddRow([]string{"A", "B"})
	chart := charts.NewBarChart([]string{"Q1", "Q2"}, []float64{10, 20}).
		WithTitle("Revenue").
		WithAltText("Chart Alt").
		Position(styling.Inches(4.8), styling.Inches(3.4)).
		Size(styling.Inches(3.4), styling.Inches(2.2))

	slide := elements.NewSlide("Reader Structures").
		WithBlankLayout().
		AddShape(left).
		AddShape(right).
		AddConnector(connector).
		WithTable(table)
	slide.Chart = &chart

	data, err := pptx.CreateWithSlides("Reader Structures", []elements.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides failed: %v", err)
	}
	deckPath := filepath.Join(t.TempDir(), "reader_structures.pptx")
	if err := os.WriteFile(deckPath, data, 0o600); err != nil {
		t.Fatalf("write deck: %v", err)
	}

	_, got, err := SlidesFromPPTX(deckPath)
	if err != nil {
		t.Fatalf("SlidesFromPPTX failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(got))
	}
	s := got[0]
	if len(s.Shapes) != 2 {
		t.Fatalf("expected only the two anchor shapes, got %d", len(s.Shapes))
	}
	if len(s.Connectors) != 1 {
		t.Fatalf("expected 1 connector, got %d", len(s.Connectors))
	}
	readConnector := s.Connectors[0]
	if readConnector.AltText != "Connector Alt" {
		t.Fatalf("expected connector alt text, got %q", readConnector.AltText)
	}
	if readConnector.StartShapeIndex != 1 || readConnector.EndShapeIndex != 2 {
		t.Fatalf("expected connector anchors to map back to shape indices, got %+v", readConnector)
	}
	if readConnector.StartSite != shapes.ConnectionSiteRight || readConnector.EndSite != shapes.ConnectionSiteLeft {
		t.Fatalf("expected connector sites right/left, got %+v", readConnector)
	}
	if s.Table == nil || s.Table.AltText != "Table Alt" {
		t.Fatalf("expected table alt text to round-trip, got %+v", s.Table)
	}
	if s.Chart == nil || s.Chart.AltText != "Chart Alt" {
		t.Fatalf("expected chart alt text to round-trip, got %+v", s.Chart)
	}
}
