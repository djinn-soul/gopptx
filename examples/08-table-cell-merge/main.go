// examples/08-table-cell-merge/main.go demonstrates table cell spanning (merge).
//
// Shows column spanning (ColSpan) and row spanning (RowSpan) using the TableCell
// fluent API to create merged header cells and merged body cells.
//
// Run with: go run ./examples/08-table-cell-merge/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir  = "examples/output"
	outputFile = "08_table_cell_merge.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// Slide 1: Column spanning — a header cell that spans 2 columns
	tableColSpan := pptx.NewTable([]styling.Length{
		styling.Inches(3),
		styling.Inches(3),
		styling.Inches(3),
	})
	// First row: a cell spanning columns 1-2 (needs a continuation placeholder), plus column 3
	tableColSpan = tableColSpan.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Spans 2 Columns").WithBold(true).WithColSpan(2),
		pptx.NewTableCell(""), // continuation placeholder for the ColSpan above
		pptx.NewTableCell("Regular Header").WithBold(true),
	})
	tableColSpan = tableColSpan.AddRow([]string{"Col A", "Col B", "Col C"})
	tableColSpan = tableColSpan.AddRow([]string{"Data 1", "Data 2", "Data 3"})

	slideColSpan := pptx.NewSlide("Column Span (ColSpan)").WithTable(tableColSpan)

	// Slide 2: Row spanning — a cell in column 2 that spans 2 rows downward
	tableRowSpan := pptx.NewTable([]styling.Length{
		styling.Inches(3),
		styling.Inches(3),
		styling.Inches(3),
	})
	tableRowSpan = tableRowSpan.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Row A, Col 1").WithBold(true),
		pptx.NewTableCell("Spans 2 Rows").WithBold(true).WithRowSpan(2),
		pptx.NewTableCell("Row A, Col 3").WithBold(true),
	})
	// Second row: column 2 is a continuation cell (empty, merged)
	tableRowSpan = tableRowSpan.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Row B, Col 1"),
		pptx.NewTableCell(""), // continuation of the RowSpan above
		pptx.NewTableCell("Row B, Col 3"),
	})
	tableRowSpan = tableRowSpan.AddRow([]string{"Row C, Col 1", "Row C, Col 2", "Row C, Col 3"})

	slideRowSpan := pptx.NewSlide("Row Span (RowSpan)").WithTable(tableRowSpan)

	// Slide 3: Combined col-span and row-span in one table
	tableCombined := pptx.NewTable([]styling.Length{
		styling.Inches(2.25),
		styling.Inches(2.25),
		styling.Inches(2.25),
		styling.Inches(2.25),
	})
	// Row 1: merged header spanning all 4 columns (+ 3 continuation placeholders)
	tableCombined = tableCombined.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Full-Width Merged Header").
			WithBold(true).
			WithBackgroundColor("4472C4").
			WithColSpan(4),
		pptx.NewTableCell(""), // continuation placeholder
		pptx.NewTableCell(""), // continuation placeholder
		pptx.NewTableCell(""), // continuation placeholder
	})
	// Row 2: two cells each spanning 2 columns (each needs one continuation placeholder)
	tableCombined = tableCombined.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Group A").WithBold(true).WithColSpan(2).WithBackgroundColor("DCE6F1"),
		pptx.NewTableCell(""), // continuation for Group A
		pptx.NewTableCell("Group B").WithBold(true).WithColSpan(2).WithBackgroundColor("DCE6F1"),
		pptx.NewTableCell(""), // continuation for Group B
	})
	// Row 3: first cell spans 2 rows, remaining cells are normal
	tableCombined = tableCombined.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Tall Cell").WithRowSpan(2),
		pptx.NewTableCell("A2"),
		pptx.NewTableCell("B1"),
		pptx.NewTableCell("B2"),
	})
	// Row 4: first cell is the row-span continuation (empty)
	tableCombined = tableCombined.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell(""), // continuation of RowSpan
		pptx.NewTableCell("A3"),
		pptx.NewTableCell("B3"),
		pptx.NewTableCell("B4"),
	})

	slideCombined := pptx.NewSlide("Combined ColSpan and RowSpan").WithTable(tableCombined)

	slides := []pptx.SlideContent{slideColSpan, slideRowSpan, slideCombined}

	data, err := pptx.CreateWithSlides("Table Cell Merge Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err = os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
