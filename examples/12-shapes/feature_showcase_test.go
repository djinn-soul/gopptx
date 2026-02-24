package main

import (
	"fmt"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestBuildShowcaseSlidesGridSlideHasSixItems(t *testing.T) {
	slides, err := buildShowcaseSlides()
	if err != nil {
		t.Fatalf("buildShowcaseSlides() error = %v", err)
	}

	gridSlide := findSlideByTitle(t, slides, "Layout Helpers (2x3 Grid)")
	if got := len(gridSlide.Shapes); got != 6 {
		t.Fatalf("grid slide shape count = %d, want 6", got)
	}

	for i := 1; i <= 6; i++ {
		want := fmt.Sprintf("Item %d", i)
		if gridSlide.Shapes[i-1].Text != want {
			t.Fatalf("grid slide shape %d text = %q, want %q", i, gridSlide.Shapes[i-1].Text, want)
		}
	}
}

func TestBuildShowcaseSlidesTableBordersConfigured(t *testing.T) {
	slides, err := buildShowcaseSlides()
	if err != nil {
		t.Fatalf("buildShowcaseSlides() error = %v", err)
	}

	tableSlide := findSlideByTitle(t, slides, "Table Styling + Borders")
	if tableSlide.Table == nil {
		t.Fatal("table slide has nil Table")
	}
	if len(tableSlide.Table.StyledRows) < 2 || len(tableSlide.Table.StyledRows[1]) < 1 {
		t.Fatal("table slide missing expected data row/cell")
	}

	cell := tableSlide.Table.StyledRows[1][0]
	if cell.BorderLeft == nil || cell.BorderLeft.Color != "AA0000" || cell.BorderLeft.Dash != pptx.TableBorderDashDash {
		t.Fatalf("left border = %#v, want color AA0000 dash %q", cell.BorderLeft, pptx.TableBorderDashDash)
	}
	if cell.BorderRight == nil || cell.BorderRight.Color != "00AA00" ||
		cell.BorderRight.Dash != pptx.TableBorderDashDot {
		t.Fatalf("right border = %#v, want color 00AA00 dash %q", cell.BorderRight, pptx.TableBorderDashDot)
	}
	if cell.BorderTop == nil || cell.BorderTop.Color != "0000AA" ||
		cell.BorderTop.Dash != pptx.TableBorderDashLongDash {
		t.Fatalf("top border = %#v, want color 0000AA dash %q", cell.BorderTop, pptx.TableBorderDashLongDash)
	}
	if cell.BorderBottom == nil || cell.BorderBottom.Color != "777777" ||
		cell.BorderBottom.Dash != pptx.TableBorderDashSolid {
		t.Fatalf("bottom border = %#v, want color 777777 dash %q", cell.BorderBottom, pptx.TableBorderDashSolid)
	}
}

func TestBuildShowcaseSlidesFlowchartHasDashedConnector(t *testing.T) {
	slides, err := buildShowcaseSlides()
	if err != nil {
		t.Fatalf("buildShowcaseSlides() error = %v", err)
	}

	flowSlide, ok := findSlideWithConnectorLabel(slides, "next")
	if !ok {
		t.Fatal(`flowchart slide with connector label "next" not found`)
	}
	if got := len(flowSlide.Shapes); got != 2 {
		t.Fatalf("flowchart slide shape count = %d, want 2", got)
	}
	if got := len(flowSlide.Connectors); got != 1 {
		t.Fatalf("flowchart slide connector count = %d, want 1", got)
	}

	conn := flowSlide.Connectors[0]
	if conn.Label != "next" {
		t.Fatalf("connector label = %q, want %q", conn.Label, "next")
	}
	if conn.Line.Dash != pptx.LineDashDashDot {
		t.Fatalf("connector dash = %q, want %q", conn.Line.Dash, pptx.LineDashDashDot)
	}
}

func findSlideByTitle(t *testing.T, slides []pptx.SlideContent, title string) pptx.SlideContent {
	t.Helper()
	for _, slide := range slides {
		if slide.Title == title {
			return slide
		}
	}
	t.Fatalf("slide with title %q not found", title)
	return pptx.SlideContent{}
}

func findSlideWithConnectorLabel(slides []pptx.SlideContent, label string) (pptx.SlideContent, bool) {
	for _, slide := range slides {
		for _, connector := range slide.Connectors {
			if connector.Label == label {
				return slide, true
			}
		}
	}
	return pptx.SlideContent{}, false
}
