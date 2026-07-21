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

// Line-item column keys used in the template rows.
const (
	keyName  = "name"
	keyNum   = "num"
	keyQty   = "qty"
	keyPrice = "price"
)

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

	// Step 1: build a visual .pptx template with {{variable}} tokens.
	slides := buildTemplateSlides()

	tmpTemplate, err := os.CreateTemp("", "gopptx-16-invoice-template-*.pptx")
	if err != nil {
		return fmt.Errorf("create temp template path: %w", err)
	}
	templatePath := tmpTemplate.Name()
	if err = tmpTemplate.Close(); err != nil {
		return fmt.Errorf("close temp template file: %w", err)
	}
	defer func() { _ = os.Remove(templatePath) }()

	if err := pptx.WriteFile(templatePath, "Invoice Template", slides); err != nil {
		return fmt.Errorf("save template: %w", err)
	}
	log.Println("Template prepared:", templatePath)

	// ── Step 2: render the template with data ──────────────────────────────────
	result, err := tplx.Render(templatePath, tplx.Context{
		"company":       "Acme Corp",
		"date":          "2026-03-01",
		"invoice_id":    "INV-2026-0042",
		"total":         "$3,600",
		"notes_present": true,
		"notes":         "Thank you for your business!",
		"line_items": []tplx.Row{
			{keyNum: "1", keyName: "Widget A", keyQty: "10", keyPrice: "$1,200"},
			{keyNum: "2", keyName: "Widget B", keyQty: "4", keyPrice: "$800"},
			{keyNum: "3", keyName: "Consulting", keyQty: "8h", keyPrice: "$1,600"},
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

func buildTemplateSlides() []pptx.SlideContent {
	coverSlide := pptx.NewSlide("").WithBlankLayout().
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(0), pptx.Inches(0), pptx.Inches(13.33), pptx.Inches(7.5)).
				WithFill(pptx.NewShapeFill("F5F9FD")),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(0.7), pptx.Inches(0.6), pptx.Inches(12), pptx.Inches(2)).
				WithGradientFill(
					pptx.NewShapeGradientFill(
						pptx.ShapeGradientTypeLinear,
						[]pptx.ShapeGradientStop{
							pptx.NewShapeGradientStop(0, "123A63"),
							pptx.NewShapeGradientStop(100, "2A6FA2"),
						},
					).WithLinearAngle(25),
				).
				WithLine(pptx.NewShapeLine("123A63", pptx.Points(1.2))).
				WithText("INVOICE • {{company}}"),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(0.9), pptx.Inches(3.1), pptx.Inches(3.8), pptx.Inches(1.2)).
				WithFill(pptx.NewShapeFill("DCE9F8")).
				WithText("Date\n{{date}}"),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(4.9), pptx.Inches(3.1), pptx.Inches(3.8), pptx.Inches(1.2)).
				WithFill(pptx.NewShapeFill("DBF1EE")).
				WithText("Invoice #\n{{invoice_id}}"),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(8.9), pptx.Inches(3.1), pptx.Inches(3.8), pptx.Inches(1.2)).
				WithFill(pptx.NewShapeFill("FEEBD7")).
				WithText("Total\n{{total}}"),
		)

	colWidths := []styling.Length{
		styling.Centimeters(1.2),
		styling.Centimeters(7.3),
		styling.Centimeters(2.6),
		styling.Centimeters(3.4),
	}
	tbl := pptx.NewTable(colWidths).
		AddStyledRow([]pptx.TableCell{
			pptx.NewTableCell("#").WithBold(true).WithBackgroundColor("DCE9F8"),
			pptx.NewTableCell("Item").WithBold(true).WithBackgroundColor("DCE9F8"),
			pptx.NewTableCell("Qty").WithBold(true).WithBackgroundColor("DCE9F8"),
			pptx.NewTableCell("Price").WithBold(true).WithBackgroundColor("DCE9F8"),
		}).
		AddStyledRow([]pptx.TableCell{
			pptx.NewTableCell("{{#each line_items}}{{.num}}"),
			pptx.NewTableCell("{{.name}}"),
			pptx.NewTableCell("{{.qty}}"),
			pptx.NewTableCell("{{.price}}{{/each}}"),
		})

	itemSlide := pptx.NewSlide("").WithBlankLayout().
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(0), pptx.Inches(0), pptx.Inches(13.33), pptx.Inches(1.1)).
				WithFill(pptx.NewShapeFill("123A63")).
				WithText("Line Items"),
		).
		WithTable(tbl)

	notesSlide := pptx.NewSlide("").WithBlankLayout().
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(0.9), pptx.Inches(0.7), pptx.Inches(11.6), pptx.Inches(1.5)).
				WithFill(pptx.NewShapeFill("EAF2FC")).
				WithLine(pptx.NewShapeLine("8AAED4", pptx.Points(1.1))).
				WithText("Client Notes"),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(0.9), pptx.Inches(2.6), pptx.Inches(11.6), pptx.Inches(2.2)).
				WithFill(pptx.NewShapeFill("FFFFFF")).
				WithLine(pptx.NewShapeLine("B9CBDD", pptx.Points(1.0))).
				WithText("{{#if notes_present}}{{notes}}{{/if}}"),
		).
		AddBullet("Payment due within 30 days.").
		AddBullet("Account team: {{company}} Finance")

	return []pptx.SlideContent{coverSlide, itemSlide, notesSlide}
}
