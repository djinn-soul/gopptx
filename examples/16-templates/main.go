// examples/16-templates/main.go demonstrates the tplx Jinja-style template engine.
//
// It generates a simple invoice PPTX template using the gopptx builder,
// then renders it with tplx.Render, expanding a table of line items.
//
// Run with: go run ./examples/16-templates/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tplx"
)

const outputDir = "examples/output"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return err
	}

	// ── Step 1: build the .pptx template with {{variable}} tokens ─────────────
	// In a real project you'd author this once in PowerPoint; here we generate
	// it programmatically so the demo is fully self-contained.

	coverSlide := pptx.NewSlide("Invoice: {{company}}").
		WithCenteredTitleLayout().
		AddBullet("Date: {{date}}").
		AddBullet("Invoice #: {{invoice_id}}").
		AddBullet("Total: {{total}}")

	// A table slide.  The second row is the loop template.
	// Because NewTable requires column widths we use the styling helpers.
	colWidths := []styling.Length{
		styling.Centimeters(1),
		styling.Centimeters(6),
		styling.Centimeters(2),
		styling.Centimeters(3),
	}
	tbl := pptx.NewTable(colWidths).AddStyledRow([]pptx.TableCell{
		pptx.NewTableCell("#"),
		pptx.NewTableCell("Item"),
		pptx.NewTableCell("Qty"),
		pptx.NewTableCell("Price"),
	}).AddStyledRow([]pptx.TableCell{
		// The {{#each}} / {{/each}} tokens wrap the row-template cells.
		pptx.NewTableCell("{{#each line_items}}{{.num}}"),
		pptx.NewTableCell("{{.name}}"),
		pptx.NewTableCell("{{.qty}}"),
		pptx.NewTableCell("{{.price}}{{/each}}"),
	})

	itemSlide := pptx.NewSlide("Line Items").WithTable(tbl)

	notesSlide := pptx.NewSlide("Notes").
		AddBullet("{{#if has_notes}}{{notes}}{{/if}}").
		AddBullet("Payment due within 30 days.")

	templatePath := filepath.Join(outputDir, "16_invoice_template.pptx")
	slides := []pptx.SlideContent{coverSlide, itemSlide, notesSlide}
	if err := pptx.WriteFile(templatePath, "Invoice Template", slides); err != nil {
		return fmt.Errorf("save template: %w", err)
	}
	log.Println("Template written:", templatePath)

	// ── Step 2: render the template with data ──────────────────────────────────
	result, err := tplx.Render(templatePath, tplx.Context{
		"company":    "Acme Corp",
		"date":       "2026-03-01",
		"invoice_id": "INV-2026-0042",
		"total":      "$3,600",
		"has_notes":  true,
		"notes":      "Thank you for your business!",
		"line_items": []tplx.Row{
			{"num": "1", "name": "Widget A", "qty": "10", "price": "$1,200"},
			{"num": "2", "name": "Widget B", "qty": "4", "price": "$800"},
			{"num": "3", "name": "Consulting", "qty": "8h", "price": "$1,600"},
		},
	})
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	outPath := filepath.Join(outputDir, "16_template_invoice.pptx")
	if err = result.Save(outPath); err != nil {
		return fmt.Errorf("save rendered: %w", err)
	}
	log.Println("Rendered invoice written:", outPath)
	log.Println("Open in PowerPoint/LibreOffice to verify all tokens are resolved.")
	return nil
}
