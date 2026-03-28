// examples/07-tables/main.go demonstrates table creation and styling.
//
// Shows how to build tables with plain text rows, styled rows using TableCell,
// background colors, alignment, and border customization.
//
// Run with: go run ./examples/07-tables/main.go
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
	outputFile = "07_tables.pptx"
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

	// Slide 1: Plain text rows
	tablePlain := pptx.NewTable([]styling.Length{
		styling.Inches(3),
		styling.Inches(3),
		styling.Inches(3),
	})
	tablePlain = tablePlain.AddRow([]string{"Header 1", "Header 2", "Header 3"})
	tablePlain = tablePlain.AddRow([]string{"Row 1, Col 1", "Row 1, Col 2", "Row 1, Col 3"})
	tablePlain = tablePlain.AddRow([]string{"Row 2, Col 1", "Row 2, Col 2", "Row 2, Col 3"})
	tablePlain = tablePlain.AddRow([]string{"Row 3, Col 1", "Row 3, Col 2", "Row 3, Col 3"})

	slidePlain := pptx.NewSlide("Plain Text Table").WithTable(tablePlain)

	// Slide 2: Styled rows with bold headers and background colors
	tableStyled := pptx.NewTable([]styling.Length{
		styling.Inches(3),
		styling.Inches(3),
		styling.Inches(3),
	})
	tableStyled = tableStyled.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Product").WithBold(true).WithBackgroundColor("4472C4"),
		pptx.NewTableCell("Quantity").WithBold(true).WithBackgroundColor("4472C4"),
		pptx.NewTableCell("Price").WithBold(true).WithBackgroundColor("4472C4"),
	})
	tableStyled = tableStyled.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Widget A"),
		pptx.NewTableCell("150"),
		pptx.NewTableCell("$12.50"),
	})
	tableStyled = tableStyled.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Widget B").WithBackgroundColor("DCE6F1"),
		pptx.NewTableCell("200").WithBackgroundColor("DCE6F1"),
		pptx.NewTableCell("$9.99").WithBackgroundColor("DCE6F1"),
	})
	tableStyled = tableStyled.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Widget C"),
		pptx.NewTableCell("75"),
		pptx.NewTableCell("$24.00"),
	})

	slideStyled := pptx.NewSlide("Styled Table with Headers").WithTable(tableStyled)

	// Slide 3: Mixed bold, italic, and alignment
	tableMixed := pptx.NewTable([]styling.Length{
		styling.Inches(2),
		styling.Inches(4),
		styling.Inches(3),
	})
	tableMixed = tableMixed.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Name").WithBold(true),
		pptx.NewTableCell("Description").WithBold(true),
		pptx.NewTableCell("Status").WithBold(true),
	})
	tableMixed = tableMixed.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Alpha"),
		pptx.NewTableCell("First entry in the list."),
		pptx.NewTableCell("Active"),
	})
	tableMixed = tableMixed.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Beta"),
		pptx.NewTableCell("Second entry, slightly longer description."),
		pptx.NewTableCell("Pending"),
	})
	tableMixed = tableMixed.AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Gamma"),
		pptx.NewTableCell("Deprecated entry."),
		pptx.NewTableCell("Inactive"),
	})

	slideMixed := pptx.NewSlide("Mixed Cell Formatting").WithTable(tableMixed)

	slides := []pptx.SlideContent{slidePlain, slideStyled, slideMixed}

	data, err := pptx.CreateWithSlides("Tables Demo", slides)
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
