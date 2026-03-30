// examples/67-table-api demonstrates table creation and cell formatting.
//
// Shows NewTable, AddRow (plain text), AddStyledRow (TableCell), cell background
// color, bold/italic text, cell margins, text alignment, and border configuration.
//
// Run with: go run ./examples/67-table-api/main.go
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
	outputFile = "67_table_api.pptx"
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

	slides := []pptx.SlideContent{
		buildPlainTableSlide(),
		buildStyledTableSlide(),
		buildMixedFormattingSlide(),
		buildMarginsTableSlide(),
	}

	outputPath := filepath.Join(outputDir, outputFile)
	data, err := pptx.CreateWithSlides("Table API Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func buildPlainTableSlide() pptx.SlideContent {
	tablePlain := pptx.NewTable([]styling.Length{
		styling.Inches(3),
		styling.Inches(3),
		styling.Inches(2.5),
	}).
		AddRow([]string{"Name", "Role", "Location"}).
		AddRow([]string{"Alice", "Engineer", "New York"}).
		AddRow([]string{"Bob", "Designer", "London"}).
		AddRow([]string{"Carol", "Manager", "Tokyo"})

	return pptx.NewSlide("Plain Text Table").WithTable(tablePlain)
}

func buildStyledTableSlide() pptx.SlideContent {
	tableStyled := pptx.NewTable([]styling.Length{
		styling.Inches(3),
		styling.Inches(2.5),
		styling.Inches(2),
	}).AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Product").WithBold(true).WithBackgroundColor("4472C4"),
		pptx.NewTableCell("Units").WithBold(true).WithBackgroundColor("4472C4"),
		pptx.NewTableCell("Price").WithBold(true).WithBackgroundColor("4472C4"),
	}).AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Widget A"),
		pptx.NewTableCell("150"),
		pptx.NewTableCell("$12.50"),
	}).AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Widget B").WithBackgroundColor("DCE6F1"),
		pptx.NewTableCell("200").WithBackgroundColor("DCE6F1"),
		pptx.NewTableCell("$9.99").WithBackgroundColor("DCE6F1"),
	}).AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Widget C"),
		pptx.NewTableCell("75"),
		pptx.NewTableCell("$24.00"),
	})

	return pptx.NewSlide("Styled Table – Bold Headers & Colors").WithTable(tableStyled)
}

func buildMixedFormattingSlide() pptx.SlideContent {
	tableMixed := pptx.NewTable([]styling.Length{
		styling.Inches(2.5),
		styling.Inches(4),
		styling.Inches(2),
	}).AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Name").WithBold(true),
		pptx.NewTableCell("Notes").WithBold(true),
		pptx.NewTableCell("Status").WithBold(true),
	}).AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Alpha"),
		pptx.NewTableCell("First entry in the system"),
		pptx.NewTableCell("Active"),
	}).AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Beta"),
		pptx.NewTableCell("Second entry, deprecated soon"),
		pptx.NewTableCell("Pending"),
	}).AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Gamma").WithBold(true),
		pptx.NewTableCell("Archived"),
		pptx.NewTableCell("Inactive").WithBackgroundColor("FFC7C7"),
	})

	return pptx.NewSlide("Mixed Cell Formatting").WithTable(tableMixed)
}

func buildMarginsTableSlide() pptx.SlideContent {
	tableMargins := pptx.NewTable([]styling.Length{
		styling.Inches(3.5),
		styling.Inches(3.5),
	}).AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Wide Margins Cell").
			WithMarginsPt(12).
			WithBold(true).
			WithBackgroundColor("E2EFDA"),
		pptx.NewTableCell("Tight Margins Cell").
			WithMarginsPt(2).
			WithBackgroundColor("FDE9D9"),
	}).AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("Font 10pt").
			WithSizePt(10),
		pptx.NewTableCell("Font 18pt").
			WithSizePt(18),
	})

	return pptx.NewSlide("Cell Margins & Font Size").WithTable(tableMargins)
}
